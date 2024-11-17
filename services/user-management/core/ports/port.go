package ports

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
)

type IUserRepository interface {
	CreateUserWithTransaction(ctx context.Context, user domain.User) (domain.User, error)
	FindUserByIDWithTransaction(ctx context.Context, user domain.User) (domain.User, error)
	FindUserByUsernameWithTransaction(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUserWithTransaction(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUserWithTransaction(ctx context.Context, user domain.User) error
}
