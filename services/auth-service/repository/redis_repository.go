package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"github.com/go-redis/redis/v8"
)

type AuthRepositoryWithRedis struct {
	client *redis.Client
	logger *logger.Logger
}

func NewAuthRepositoryWithRedis(client *redis.Client, logger *logger.Logger) ports.IAuthRepository {
	return &AuthRepositoryWithRedis{
		client: client,
		logger: logger,
	}
}

// StoreToken is a method to store a token in Redis with auto-expiration
func (r *AuthRepositoryWithRedis) CreateToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	key := fmt.Sprintf("sessionID:%s", auth.Claims.SessionID)
	refreshTokenKey := fmt.Sprintf("refreshToken:%s", auth.RefreshToken) // Mapping key

	// Use a pipeline to execute multiple commands in a single round-trip
	pipe := r.client.TxPipeline()

	// Set the token information as fields in the Redis hash
	pipe.HSet(ctx, key, map[string]interface{}{
		"refresh_token": auth.RefreshToken,
		"user_id":       auth.User.ID,
		"expires_at":    auth.RefreshTokenExpiresAt,
		"is_revoked":    auth.RefreshTokenIsRevoked,
	})

	// Store the mapping from refresh token to session ID
	pipe.Set(ctx, refreshTokenKey, auth.Claims.SessionID, time.Until(auth.RefreshTokenExpiresAt))

	// Set the key to expire after 1 day
	expiration := 15 * time.Minute
	pipe.Expire(ctx, key, expiration)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to execute pipeline for storing token: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}

	return auth, nil
}

// GetToken is a method to retrieve a token by id from Redis
func (r *AuthRepositoryWithRedis) GetTokenByID(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	key := fmt.Sprintf("sessionID:%s", auth.ID)

	// Retrieve the fields for the token from the Redis hash
	sessionData, err := r.client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		r.logger.Warn(fmt.Sprintf("token not found in Redis: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorNotFound, err)
	}
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to retrieve token from Redis: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}

	if len(sessionData) == 0 {
		r.logger.Warn(fmt.Sprintf("token not found in Redis"))
		return domain.Auth{}, errors.NewError(errors.ErrorNotFound, fmt.Errorf("token not found in Redis"))
	}

	// Parse the session data
	refreshToken := sessionData["refresh_token"]

	userID, err := strconv.Atoi(sessionData["user_id"])
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to parse user_id: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}

	expiresAt, err := time.Parse(time.RFC3339, sessionData["expires_at"])
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to parse expires_at: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}

	isRevoked, err := strconv.ParseBool(sessionData["is_revoked"])
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to parse is_revoked: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}

	foundAuth := domain.Auth{
		ID:                    auth.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: expiresAt,
		RefreshTokenIsRevoked: isRevoked,
		User: domain.User{
			ID: uint(userID),
		},
		Claims: domain.Claims{
			ID:        uint(userID),
			SessionID: auth.ID,
			IssuedAt:  time.Now(),
			ExpiresAt: expiresAt,
		},
	}
	return foundAuth, nil
}

// GetToken is a method to retrieve a token by refresh token from Redis
func (r *AuthRepositoryWithRedis) GetTokenByRefreshToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	refreshTokenKey := fmt.Sprintf("refreshToken:%s", auth.RefreshToken)

	// Step 1: Retrieve the sessionID using the refresh token
	sessionID, err := r.client.Get(ctx, refreshTokenKey).Result()
	if err == redis.Nil {
		r.logger.Warn(fmt.Sprintf("refresh token not found in Redis: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorNotFound, err)
	} else if err != nil {
		r.logger.Error(fmt.Sprintf("failed to get sessionID by refresh token: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}

	if sessionID == "" {
		r.logger.Warn(fmt.Sprintf("sessionID not found in Redis"))
		return domain.Auth{}, errors.NewError(errors.ErrorNotFound, fmt.Errorf("sessionID not found in Redis"))
	}

	auth.ID = sessionID

	// Step 2: Retrieve the token information using the sessionID
	auth, err = r.GetTokenByID(ctx, auth)
	if err != nil {
		return domain.Auth{}, err
	}

	return auth, nil
}

// DeleteToken is a method to delete a token from Redis
func (r *AuthRepositoryWithRedis) DeleteToken(ctx context.Context, auth domain.Auth) error {
	auth, err := r.GetTokenByID(ctx, auth)
	if err != nil {
		return err
	}

	// Use a pipeline to execute multiple commands in a single round-trip
	pipe := r.client.TxPipeline()

	// Step 1: delete the refresh token mapping to sessionID
	refreshTokenKey := fmt.Sprintf("refreshToken:%s", auth.RefreshToken)
	pipe.Del(ctx, refreshTokenKey)

	// Step 2: delete the sessionID hash
	key := fmt.Sprintf("sessionID:%s", auth.ID)
	pipe.Del(ctx, key)

	// Execute the pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			// Handle specific cases where keys are missing
			r.logger.Warn(fmt.Sprintf("key not found for deletion in Redis: %v", err))
			return errors.NewError(errors.ErrorNotFound, err)
		}
		// General error handling for failed pipeline execution
		r.logger.Error(fmt.Sprintf("failed to execute pipeline for deleting token: %v", err))
		return errors.NewError(errors.ErrorInternal, err)
	}

	return nil
}

// RevokeToken is a method to revoke a token in Redis
func (r *AuthRepositoryWithRedis) RevokeToken(ctx context.Context, auth domain.Auth) error {
	key := fmt.Sprintf("sessionID:%s", auth.Claims.SessionID)

	// Set is_revoked to true for the given token
	err := r.client.HSet(ctx, key, "is_revoked", true).Err()
	if err != nil {
		if err == redis.Nil {
			r.logger.Warn(fmt.Sprintf("token not found for revocation in Redis: %v", err))
			return errors.NewError(errors.ErrorNotFound, err)
		}
		r.logger.Error(fmt.Sprintf("failed to revoke token in Redis: %v", err))
		return errors.NewError(errors.ErrorInternal, err)
	}
	return nil
}
