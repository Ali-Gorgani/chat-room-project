package ports

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
)

type IAuthRepository interface {
	CreateToken(ctx context.Context, auth domain.Auth) (domain.Auth, error)
	GetTokenByID(ctx context.Context, auth domain.Auth) (domain.Auth, error)
	GetTokenByRefreshToken(ctx context.Context, auth domain.Auth) (domain.Auth, error)
	DeleteToken(ctx context.Context, auth domain.Auth) error
	RevokeToken(ctx context.Context, auth domain.Auth) error
}
