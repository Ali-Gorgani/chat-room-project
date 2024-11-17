package ports

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
)

type IAuthRepository interface {
	CreateToken(ctx context.Context, auth domain.Auth) (domain.Auth, error)
	GetToken(ctx context.Context, auth domain.Auth) (domain.Auth, error)
	DeleteToken(ctx context.Context, auth domain.Auth) error
	RevokedToken(ctx context.Context, auth domain.Auth) error
}
