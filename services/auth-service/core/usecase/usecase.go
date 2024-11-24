package usecase

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/service/user"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/hash"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/jwt"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	authRepository ports.IAuthRepository
	userService    *user.UsersService
	logger         *logger.Logger
	config         *configs.Config
}

func NewAuthUseCase(authRepository ports.IAuthRepository, userService *user.UsersService, logger *logger.Logger, config *configs.Config) *AuthUseCase {
	return &AuthUseCase{
		authRepository: authRepository,
		userService:    userService,
		logger:         logger,
		config:         config,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	foundAuth, err := a.userService.GetUserByUsername(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in getting user by username: %v", err))
		return domain.Auth{}, err
	}

	// check user password
	ok, err := hash.ComparePassword(foundAuth.User.Password, auth.User.Password, a.logger)
	if err != nil {
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}
	if !ok {
		return domain.Auth{}, errors.NewError(errors.ErrorUnauthorized, err)
	}

	// create access token
	auth.Claims = domain.Claims{
		ID:       foundAuth.User.ID,
		Username: foundAuth.User.Username,
		Email:    foundAuth.User.Email,
		Role:     foundAuth.User.Role,
		Duration: a.config.JWT.AccessTokenDuration,
	}

	sessionID := uuid.New().String()
	auth.ID = sessionID
	auth.Claims.SessionID = sessionID

	auth, err = a.CreateToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in creating access token: %v", err))
		return domain.Auth{}, err
	}

	// create refresh token
	auth.Claims.Duration = a.config.JWT.RefreshTokenDuration
	refreshToken, err := a.CreateToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in creating refresh token: %v", err))
		return domain.Auth{}, err
	}

	// save refresh token
	auth.RefreshToken = refreshToken.AccessToken
	auth.RefreshTokenExpiresAt = refreshToken.AccessTokenExpiresAt
	auth, err = a.authRepository.CreateToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in saving refresh token: %v", err))
		return domain.Auth{}, err
	}

	return auth, nil
}

func (a *AuthUseCase) Logout(ctx context.Context, auth domain.Auth) error {
	// get token from context
	contextToken, ok := ctx.Value("token").(string)
	if !ok {
		err := fmt.Errorf("error in getting token from context")
		a.logger.Error(err.Error())
		return errors.NewError(errors.ErrorBadRequest, err)
	}

	// verify access token
	auth.AccessToken = contextToken
	auth, err := a.VerifyToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in verifying token: %v", err))
		return err
	}

	// delete token by finding it via RefreshToken in the database
	err = a.authRepository.DeleteToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in deleting token: %v", err))
		return err
	}

	return nil
}

func (a *AuthUseCase) RefreshToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	// get refresh token from database
	auth, err := a.authRepository.GetTokenByRefreshToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in getting token from database: %v", err))
		return domain.Auth{}, err
	}

	// check refresh token is expired
	if auth.RefreshTokenIsRevoked {
		err := fmt.Errorf("refresh token is revoked")
		a.logger.Error(err.Error())
		return domain.Auth{}, errors.NewError(errors.ErrorUnauthorized, err)
	}

	// create access token
	auth.Claims.Duration = a.config.JWT.AccessTokenDuration
	auth, err = a.CreateToken(ctx, auth)
	if err != nil {
		return domain.Auth{}, err
	}

	return auth, nil
}

func (a *AuthUseCase) RevokeToken(ctx context.Context, auth domain.Auth) error {
	// get token from context
	contextToken, ok := ctx.Value("token").(string)
	if !ok {
		err := fmt.Errorf("error in getting token from context")
		a.logger.Error(err.Error())
		return errors.NewError(errors.ErrorBadRequest, err)
	}

	// verify access token
	auth.AccessToken = contextToken
	accessTokenClaims, err := a.VerifyToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in verifying token: %v", err))
		return err
	}

	// get refresh token from database
	auth, err = a.authRepository.GetTokenByRefreshToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in getting token from database: %v", err))
		return err
	}

	// compare access token and refresh token are for the same person
	if accessTokenClaims.Claims.ID != auth.Claims.ID {
		err := fmt.Errorf("access token and refresh token are not for the same person")
		a.logger.Error(err.Error())
		return errors.NewError(errors.ErrorForbidden, err)
	}

	// check refresh token is expired
	if auth.RefreshTokenIsRevoked {
		err := fmt.Errorf("refresh token is revoked")
		a.logger.Error(err.Error())
		return errors.NewError(errors.ErrorUnauthorized, err)
	}

	// revoke access token
	err = a.authRepository.RevokeToken(ctx, auth)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in revoking token: %v", err))
		return err
	}

	return nil
}

func (a *AuthUseCase) HashPassword(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	hashedPassword, err := hash.HashedPassword(auth.User.Password, a.logger)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in hashing password: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}
	auth.User.Password = hashedPassword
	return auth, nil
}

func (a *AuthUseCase) CreateToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	secretKey := a.config.JWT.SecretKey
	userClaims := jwt.UserClaims{
		ID:        auth.Claims.ID,
		SessionID: auth.Claims.SessionID,
		Username:  auth.Claims.Username,
		Email:     auth.Claims.Email,
		Role:      auth.Claims.Role,
		Duration:  auth.Claims.Duration,
	}
	claims, err := jwt.NewUserClaims(userClaims)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in creating token claims: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}
	accessToken, err := jwt.CreateToken(secretKey, claims)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in creating token claims: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorInternal, err)
	}
	auth = domain.Auth{
		ID:                   claims.RegisteredClaims.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
		User: domain.User{
			ID:       claims.ID,
			Username: claims.Username,
			Email:    claims.Email,
			Role:     claims.Role,
		},
		Claims: domain.Claims{
			ID:        claims.ID,
			Username:  claims.Username,
			Email:     claims.Email,
			Role:      claims.Role,
			Duration:  claims.Duration,
			SessionID: claims.RegisteredClaims.ID,
			IssuedAt:  claims.RegisteredClaims.IssuedAt.Time,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
		},
	}
	return auth, nil
}

func (a *AuthUseCase) VerifyToken(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	secretKey := a.config.JWT.SecretKey
	claims, err := jwt.VerifyToken(auth.AccessToken, secretKey)
	if err != nil {
		a.logger.Error(fmt.Sprintf("error in verifying token: %v", err))
		return domain.Auth{}, errors.NewError(errors.ErrorUnauthorized, err)
	}
	auth = domain.Auth{
		ID:                   claims.SessionID,
		AccessToken:          auth.AccessToken,
		AccessTokenExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
		User: domain.User{
			ID:       claims.ID,
			Username: claims.Username,
			Email:    claims.Email,
			Role:     claims.Role,
		},
		Claims: domain.Claims{
			ID:        claims.ID,
			Username:  claims.Username,
			Email:     claims.Email,
			Role:      claims.Role,
			Duration:  claims.Duration,
			SessionID: claims.RegisteredClaims.ID,
			IssuedAt:  claims.RegisteredClaims.IssuedAt.Time,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
		},
	}
	return auth, nil
}
