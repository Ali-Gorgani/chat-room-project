package usecase

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/service/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
)

type UserUseCase struct {
	userRepository ports.IUserRepository
	authService    *auth.AuthService
	logger         *logger.Logger
}

func NewUserUseCase(userRepository ports.IUserRepository, authService *auth.AuthService, logger *logger.Logger) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		authService:    authService,
		logger:         logger,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	// hash password with auth service
	hashedPassword, err := u.authService.HashPassword(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}
	user.Password = hashedPassword.Password

	createdUser, err := u.userRepository.CreateUserWithTransaction(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	return createdUser, nil
}

func (u *UserUseCase) FindUserByID(ctx context.Context, user domain.User) (domain.User, error) {
	// get token from context
	contextToken, ok := ctx.Value("token").(string)
	if !ok {
		err := fmt.Errorf("error in getting token from context")
		u.logger.Error(err.Error())
		return domain.User{}, errors.NewError(errors.ErrorBadRequest, err)
	}

	// verify token with auth service and get user claims
	_, err := u.authService.VerifyToken(ctx, domain.Auth{AccessToken: contextToken})
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	existingUser, err := u.userRepository.FindUserByIDWithTransaction(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	if !u.hasPermission(existingUser, "read") {
		return domain.User{}, errors.NewError(errors.ErrorForbidden, fmt.Errorf("user does not have permission to read user"))
	}

	return existingUser, nil
}

func (u *UserUseCase) FindUserByUsername(ctx context.Context, user domain.User) (domain.User, error) {
	existingUser, err := u.userRepository.FindUserByUsernameWithTransaction(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	return existingUser, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	// get token from context
	contextToken, ok := ctx.Value("token").(string)
	if !ok {
		err := fmt.Errorf("error in getting token from context")
		u.logger.Error(err.Error())
		return domain.User{}, errors.NewError(errors.ErrorBadRequest, err)
	}

	// verify token with auth service and get user claims
	userClaims, err := u.authService.VerifyToken(ctx, domain.Auth{AccessToken: contextToken})
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	if userClaims.ID != user.ID && userClaims.Role.Name != "admin" {
		return domain.User{}, errors.NewError(errors.ErrorForbidden, fmt.Errorf("user does not have permission to update user"))
	}

	existingUser, err := u.userRepository.FindUserByIDWithTransaction(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	if !u.hasPermission(existingUser, "update") {
		return domain.User{}, errors.NewError(errors.ErrorForbidden, fmt.Errorf("user does not have permission to update user"))
	}

	user.ID = existingUser.ID
	user.Role.ID = existingUser.Role.ID
	user.Profile.ID = existingUser.Profile.ID

	// hash password with auth service
	hashedPassword, err := u.authService.HashPassword(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}
	user.Password = hashedPassword.Password

	updatedUser, err := u.userRepository.UpdateUserWithTransaction(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return domain.User{}, err
	}

	return updatedUser, nil
}

func (u *UserUseCase) DeleteUser(ctx context.Context, user domain.User) error {
	// get token from context
	contextToken, ok := ctx.Value("token").(string)
	if !ok {
		err := fmt.Errorf("error in getting token from context")
		u.logger.Error(err.Error())
		return errors.NewError(errors.ErrorBadRequest, err)
	}

	// verify token with auth service and get user claims
	claims, err := u.authService.VerifyToken(ctx, domain.Auth{AccessToken: contextToken})
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	userClaims, err := u.userRepository.FindUserByIDWithTransaction(ctx, claims)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	existingUser, err := u.userRepository.FindUserByIDWithTransaction(ctx, user)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	if !u.hasPermission(userClaims, "delete") {
		return errors.NewError(errors.ErrorForbidden, fmt.Errorf("user does not have permission to delete user"))
	}

	err = u.userRepository.DeleteUserWithTransaction(ctx, existingUser)
	if err != nil {
		u.logger.Error(err.Error())
		return err
	}

	return nil
}

func (u *UserUseCase) hasPermission(user domain.User, requiredPermission string) bool {
	// Step 1: Check if the user has a role
	if user.Role.Name == "" {
		u.logger.Error(fmt.Sprintf("User %d has no role assigned", user.ID))
		return false
	}

	// Step 2: Check if the user's role has the required permission
	rolePermissions := user.Role.Premissions // Assuming the Role has a `Permissions` field that is a list of strings

	// Step 3: Loop through the role's permissions and check if the required permission exists
	for _, permission := range rolePermissions {
		if permission == requiredPermission {
			u.logger.Info(fmt.Sprintf("User %d has required permission: %s", user.ID, requiredPermission))
			return true
		}
	}

	// If the required permission wasn't found
	u.logger.Info(fmt.Sprintf("User %d does not have the required permission: %s", user.ID, requiredPermission))
	return false
}
