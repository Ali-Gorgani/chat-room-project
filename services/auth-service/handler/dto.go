package handler

import (
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	ID                    uint         `json:"id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type RevokeTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func LoginRequestToDomainAuth(req LoginRequest) domain.Auth {
	return domain.Auth{
		User: domain.User{
			Username: req.Username,
			Password: req.Password,
		},
	}
}

func DomainAuthToLoginResponse(auth domain.Auth) LoginResponse {
	return LoginResponse{
		ID:                    auth.ID,
		AccessToken:           auth.AccessToken,
		AccessTokenExpiresAt:  auth.AccessTokenExpiresAt,
		RefreshToken:          auth.RefreshToken,
		RefreshTokenExpiresAt: auth.RefreshTokenExpiresAt,
		User: UserResponse{
			ID:       auth.User.ID,
			Username: auth.User.Username,
			Email:    auth.User.Email,
			Role:     auth.User.Role,
		},
	}
}

func RefreshTokenRequestToDomainAuth(req RefreshTokenRequest) domain.Auth {
	return domain.Auth{
		RefreshToken: req.RefreshToken,
	}
}

func DomainAuthToRefreshTokenResponse(auth domain.Auth) RefreshTokenResponse {
	return RefreshTokenResponse{
		AccessToken:          auth.AccessToken,
		AccessTokenExpiresAt: auth.AccessTokenExpiresAt,
	}
}

func RevokeTokenRequestToDomainAuth(req RevokeTokenRequest) domain.Auth {
	return domain.Auth{
		RefreshToken: req.RefreshToken,
	}
}
