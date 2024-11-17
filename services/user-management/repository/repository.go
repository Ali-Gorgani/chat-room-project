package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent"
	entRole "github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent/role"
	entUser "github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent/user"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
)

type UserRepository struct {
	client *ent.Client
	logger *logger.Logger
}

func NewUserRepository(client *ent.Client, logger *logger.Logger) ports.IUserRepository {
	return &UserRepository{
		client: client,
		logger: logger,
	}
}

// CreateUserWithTransaction creates a user with role and profile assignment in a transaction
func (r *UserRepository) CreateUserWithTransaction(ctx context.Context, user domain.User) (domain.User, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to start transaction: %v", err))
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			r.logger.Error(fmt.Sprintf("transaction panic: %v", p))
		}
	}()

	// Create user
	createdUser, err := tx.User.Create().
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetEmail(user.Email).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return domain.User{}, errors.NewError(errors.ErrorConflict, fmt.Errorf("user already exists"))
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Create profile if provided
	var createdProfile *ent.Profile
	if user.Profile.FirstName != "" || user.Profile.LastName != "" || user.Profile.ProfilePicture != "" {
		createdProfile, err = tx.Profile.Create().
			SetFirstName(user.Profile.FirstName).
			SetLastName(user.Profile.LastName).
			SetProfilePicture(user.Profile.ProfilePicture).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return domain.User{}, errors.NewError(errors.ErrorInternal, err)
		}
	}

	// Assign role
	role, err := r.findOrCreateRole(ctx, tx, user.Role.Name)
	if err != nil {
		tx.Rollback()
		return domain.User{}, err
	}

	// Assign role and profile to the user
	_, err = tx.User.UpdateOneID(createdUser.ID).
		SetRole(role).
		SetProfile(createdProfile). // Assign profile to user
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	if err := tx.Commit(); err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	return domain.User{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Password: createdUser.Password,
		Email:    createdUser.Email,
		Role: domain.Role{
			ID:          role.ID,
			Name:        role.Name,
			Premissions: role.Permissions,
		},
		Profile: domain.Profile{
			ID:             createdProfile.ID,
			FirstName:      createdProfile.FirstName,
			LastName:       createdProfile.LastName,
			ProfilePicture: createdProfile.ProfilePicture,
		},
	}, nil
}

// FindUserByIDWithTransaction retrieves a user with their profile and role within a transaction
func (r *UserRepository) FindUserByIDWithTransaction(ctx context.Context, user domain.User) (domain.User, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to start transaction: %v", err))
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}
	defer tx.Rollback()

	// Query user with profile and role information
	foundUser, err := tx.User.Query().
		Where(entUser.IDEQ(user.ID)).
		WithProfile(). // Load the associated profile
		WithRole().    // Load the associated role information
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Map the result to domain.User including profile and role information
	return domain.User{
		ID:       foundUser.ID,
		Username: foundUser.Username,
		Password: foundUser.Password,
		Email:    foundUser.Email,
		Role: domain.Role{
			ID:          foundUser.Edges.Role.ID,
			Name:        foundUser.Edges.Role.Name,
			Premissions: foundUser.Edges.Role.Permissions,
		},
		Profile: domain.Profile{
			ID:             foundUser.Edges.Profile.ID,
			FirstName:      foundUser.Edges.Profile.FirstName,
			LastName:       foundUser.Edges.Profile.LastName,
			ProfilePicture: foundUser.Edges.Profile.ProfilePicture,
		},
	}, nil
}

// FindUserByUsernameWithTransaction retrieves a user with their profile and role within a transaction
func (r *UserRepository) FindUserByUsernameWithTransaction(ctx context.Context, user domain.User) (domain.User, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to start transaction: %v", err))
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}
	defer tx.Rollback()

	// Query user with profile and role information
	foundUser, err := tx.User.Query().
		Where(entUser.UsernameEQ(user.Username)).
		WithProfile(). // Load the associated profile
		WithRole().    // Load the associated role information
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Map the result to domain.User including profile and role information
	return domain.User{
		ID:       foundUser.ID,
		Username: foundUser.Username,
		Password: foundUser.Password,
		Email:    foundUser.Email,
		Role: domain.Role{
			ID:          foundUser.Edges.Role.ID,
			Name:        foundUser.Edges.Role.Name,
			Premissions: foundUser.Edges.Role.Permissions,
		},
		Profile: domain.Profile{
			ID:             foundUser.Edges.Profile.ID,
			FirstName:      foundUser.Edges.Profile.FirstName,
			LastName:       foundUser.Edges.Profile.LastName,
			ProfilePicture: foundUser.Edges.Profile.ProfilePicture,
		},
	}, nil
}

// UpdateUserWithTransaction updates a user with their profile and role within a transaction
func (r *UserRepository) UpdateUserWithTransaction(ctx context.Context, user domain.User) (domain.User, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to start transaction: %v", err))
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}
	defer tx.Rollback()

	// Update the profile if it exists
	var updatedProfile *ent.Profile
	if user.Profile.ID != 0 {
		updatedProfile, err = tx.Profile.UpdateOneID(user.Profile.ID).
			SetFirstName(user.Profile.FirstName).
			SetLastName(user.Profile.LastName).
			SetProfilePicture(user.Profile.ProfilePicture).
			Save(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
			}
			return domain.User{}, errors.NewError(errors.ErrorInternal, err)
		}
	}

	// Update the role if it exists
	var updatedRole *ent.Role
	if user.Role.ID != 0 {
		existingRole, err := tx.Role.Query().Where(entRole.NameEQ(user.Role.Name)).Only(ctx)
		if err == nil {
			updatedRole = existingRole
		} else if ent.IsNotFound(err) {
			updatedRole, err = tx.Role.UpdateOneID(user.Role.ID).
				SetName(user.Role.Name).
				SetPermissions(defaultPermissions(user.Role.Name)).
				Save(ctx)
			if err != nil {
				return domain.User{}, errors.NewError(errors.ErrorInternal, err)
			}
		} else {
			return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
		}
	}

	// Update the user
	updatedUser, err := tx.User.UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetEmail(user.Email).
		SetRole(updatedRole).       // Assign role to user
		ClearProfile().             // Clear the profile association
		SetProfile(updatedProfile). // Assign profile to user
		Save(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return domain.User{}, errors.NewError(errors.ErrorConflict, fmt.Errorf("username or email already exists"))
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	return domain.User{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Password: updatedUser.Password,
		Email:    updatedUser.Email,
		Role: domain.Role{
			ID:          updatedRole.ID,
			Name:        updatedRole.Name,
			Premissions: updatedRole.Permissions,
		},
		Profile: domain.Profile{
			ID:             updatedProfile.ID,
			FirstName:      updatedProfile.FirstName,
			LastName:       updatedProfile.LastName,
			ProfilePicture: updatedProfile.ProfilePicture,
		},
	}, nil
}

// DeleteUserWithTransaction deletes a user and their profile within a transaction
func (r *UserRepository) DeleteUserWithTransaction(ctx context.Context, user domain.User) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to start transaction: %v", err))
		return errors.NewError(errors.ErrorInternal, err)
	}
	defer tx.Rollback()

	// Delete the user's profile if it exists
	if user.Profile.ID != 0 {
		err = tx.Profile.DeleteOneID(user.Profile.ID).Exec(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return errors.NewError(errors.ErrorNotFound, err)
			}
			return errors.NewError(errors.ErrorInternal, err)
		}
	}

	// Delete the user
	err = tx.User.DeleteOneID(user.ID).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NewError(errors.ErrorNotFound, err)
		}
		return errors.NewError(errors.ErrorInternal, err)
	}

	return tx.Commit()
}

// Helper function to find or create a role by name within the transaction
func (r *UserRepository) findOrCreateRole(ctx context.Context, tx *ent.Tx, roleName string) (*ent.Role, error) {
	role, err := tx.Role.Query().Where(entRole.NameEQ(roleName)).Only(ctx)
	if err == nil {
		return role, nil
	}

	if ent.IsNotFound(err) {
		createdRole, err := tx.Role.Create().
			SetName(roleName).
			SetPermissions(defaultPermissions(roleName)).
			Save(ctx)
		if err != nil {
			return nil, errors.NewError(errors.ErrorInternal, err)
		}
		return createdRole, nil
	}

	return nil, errors.NewError(errors.ErrorInternal, err)
}

// Default permissions based on role name
func defaultPermissions(roleName string) []string {
	if roleName == "admin" {
		return []string{"read", "create", "update", "delete"}
	}
	return []string{"create", "update"}
}
