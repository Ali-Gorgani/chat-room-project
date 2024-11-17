package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent"
	entAuth "github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
)

type AuthRepository struct {
	client *ent.Client
	logger *logger.Logger
}

func NewAuthRepository(client *ent.Client, logger *logger.Logger) ports.IAuthRepository {
	return &AuthRepository{
		client: client,
		logger: logger,
	}
}

// CreateToken is a method to create a token
func (r *AuthRepository) CreateToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	token, err := r.client.Auth.
		Create().
		SetUserID(auth.User.ID).
		SetRefreshToken(auth.RefreshToken).
		SetExpiresAt(auth.RefreshTokenExpiresAt).
		Save(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			r.logger.Error(fmt.Sprintf("failed to create token: %v", err))
			return domain.Auth{}, errors.NewError(errors.ErrorConflict, fmt.Errorf("token already exists"))
		}
		r.logger.Error(fmt.Sprintf("failed to create token: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}
	auth.ID = uint(token.ID)
	return auth, nil
}

// GetToken is a method to get a token
func (r *AuthRepository) GetToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	token, err := r.client.Auth.
		Query().
		Where(entAuth.RefreshTokenEQ(auth.RefreshToken)).
		Only(ctx)
	if ent.IsNotFound(err) {
		r.logger.Warn(fmt.Sprintf("token not found: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorNotFound, err)
	}
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to retrieve token: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}
	auth = domain.Auth{
		ID:                    uint(token.ID),
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresAt: token.ExpiresAt,
		RefreshTokenIsRevoked: token.IsRevoked,
		User: domain.User{
			ID: token.UserID,
		},
		Claims: domain.Claims{
			ID:        token.UserID,
			IssuedAt:  token.CreatedAt,
			ExpiresAt: token.ExpiresAt,
		},
	}
	return auth, nil
}

// DeleteToken is a method to delete a token
func (r *AuthRepository) DeleteToken(ctx context.Context, auth domain.Auth) error {
	err := r.client.Auth.
		DeleteOneID(int(auth.ID)).
		Exec(ctx)
	if ent.IsNotFound(err) {
		r.logger.Warn(fmt.Sprintf("token not found for deletion: %v", err))
		return errors.NewError(errors.ErrorNotFound, err)
	}
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to delete token: %v", err))
		return errors.NewError(errors.ErrorInternal, err)
	}
	return nil
}

// RevokedToken is a method to revoke a token
func (r *AuthRepository) RevokedToken(ctx context.Context, auth domain.Auth) error {
	_, err := r.client.Auth.
		UpdateOneID(int(auth.ID)).
		SetIsRevoked(true).
		Save(ctx)
	if ent.IsNotFound(err) {
		r.logger.Warn(fmt.Sprintf("token not found for revocation: %v", err))
		return errors.NewError(errors.ErrorNotFound, err)
	}
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to revoke token: %v", err))
		return errors.NewError(errors.ErrorInternal, err)
	}
	return nil
}
