package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent"

	// entProfile "github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent/profile"
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
	// Start a transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("failed to start transaction: %v", err))
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}
	defer tx.Rollback()

	// Create the user
	createdUser, err := tx.User.Create().
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetEmail(user.Email).
		Save(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return domain.User{}, errors.NewError(errors.ErrorConflict, fmt.Errorf("user already exists"))
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Ensure role defaults to "user" if not provided
	roleName := user.Role.Name
	if roleName == "" {
		roleName = "user"
	}

	// Find or create the role
	role, err := tx.Role.Query().Where(entRole.NameEQ(entRole.Name(roleName))).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			role, err = tx.Role.Create().
				SetName(entRole.Name(roleName)).
				SetPermissions(defaultPermissions(roleName)).
				Save(ctx)
			if err != nil {
				return domain.User{}, errors.NewError(errors.ErrorInternal, err)
			}
		} else {
			return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
		}
	}

	// Update the user role
	_, err = tx.User.UpdateOneID(createdUser.ID).
		SetRole(role).
		Save(ctx)
	if err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Create the profile if provided
	if user.Profile.FirstName != "" || user.Profile.LastName != "" || user.Profile.ProfilePicture != "" {
		createdProfile, err := tx.Profile.Create().
			SetFirstName(user.Profile.FirstName).
			SetLastName(user.Profile.LastName).
			SetProfilePicture(user.Profile.ProfilePicture).
			Save(ctx)
		if err != nil {
			return domain.User{}, errors.NewError(errors.ErrorInternal, err)
		}

		// Update the user profile
		_, err = tx.User.UpdateOneID(createdUser.ID).
			SetProfile(createdProfile).
			Save(ctx)
		if err != nil {
			return domain.User{}, errors.NewError(errors.ErrorInternal, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// get user with profile and role information
	createdUser, err = r.client.User.Query().
		Where(entUser.IDEQ(createdUser.ID)).
		WithProfile(). // Load the associated profile
		WithRole().    // Load the associated role information
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Return the created user with profile and role
	return domain.User{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Password: createdUser.Password,
		Email:    createdUser.Email,
		Role: domain.Role{
			ID:          createdUser.Edges.Role.ID,
			Name:        string(createdUser.Edges.Role.Name),
			Premissions: createdUser.Edges.Role.Permissions,
		},
		Profile: func() domain.Profile {
			if createdUser.Edges.Profile != nil {
				return domain.Profile{
					ID:             createdUser.Edges.Profile.ID,
					FirstName:      createdUser.Edges.Profile.FirstName,
					LastName:       createdUser.Edges.Profile.LastName,
					ProfilePicture: createdUser.Edges.Profile.ProfilePicture,
				}
			}
			return domain.Profile{} // Return an empty profile if no profile exists
		}(),
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
			Name:        string(foundUser.Edges.Role.Name),
			Premissions: foundUser.Edges.Role.Permissions,
		},
		Profile: func() domain.Profile {
			if foundUser.Edges.Profile != nil {
				return domain.Profile{
					ID:             foundUser.Edges.Profile.ID,
					FirstName:      foundUser.Edges.Profile.FirstName,
					LastName:       foundUser.Edges.Profile.LastName,
					ProfilePicture: foundUser.Edges.Profile.ProfilePicture,
				}
			}
			return domain.Profile{} // Return an empty profile if no profile exists
		}(),
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
			Name:        string(foundUser.Edges.Role.Name),
			Premissions: foundUser.Edges.Role.Permissions,
		},
		Profile: func() domain.Profile {
			if foundUser.Edges.Profile != nil {
				return domain.Profile{
					ID:             foundUser.Edges.Profile.ID,
					FirstName:      foundUser.Edges.Profile.FirstName,
					LastName:       foundUser.Edges.Profile.LastName,
					ProfilePicture: foundUser.Edges.Profile.ProfilePicture,
				}
			}
			return domain.Profile{} // Return an empty profile if no profile exists
		}(),
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

	// get user with profile and role information
	updatedUser, err := tx.User.Query().
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

	// Update the user
	_, err = tx.User.UpdateOneID(updatedUser.ID).
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetEmail(user.Email).
		Save(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return domain.User{}, errors.NewError(errors.ErrorConflict, fmt.Errorf("user already exists"))
		}
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// Update the user role if provided
	if user.Role.Name != "" {
		role, err := tx.Role.Query().Where(entRole.NameEQ(entRole.Name(user.Role.Name))).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				role, err = tx.Role.Create().
					SetName(entRole.Name(user.Role.Name)).
					SetPermissions(defaultPermissions(user.Role.Name)).
					Save(ctx)
				if err != nil {
					return domain.User{}, errors.NewError(errors.ErrorInternal, err)
				}
			} else {
				return domain.User{}, errors.NewError(errors.ErrorNotFound, err)
			}
		}

		_, err = tx.User.UpdateOneID(updatedUser.ID).
			SetRole(role).
			Save(ctx)
		if err != nil {
			return domain.User{}, errors.NewError(errors.ErrorInternal, err)
		}
	}

	// Update the user profile if provided
	if user.Profile.FirstName != "" || user.Profile.LastName != "" || user.Profile.ProfilePicture != "" {
		if updatedUser.Edges.Profile != nil {
			_, err = tx.Profile.UpdateOneID(updatedUser.Edges.Profile.ID).
				SetFirstName(user.Profile.FirstName).
				SetLastName(user.Profile.LastName).
				SetProfilePicture(user.Profile.ProfilePicture).
				Save(ctx)
			if err != nil {
				return domain.User{}, errors.NewError(errors.ErrorInternal, err)
			}
		} else {
			createdProfile, err := tx.Profile.Create().
				SetFirstName(user.Profile.FirstName).
				SetLastName(user.Profile.LastName).
				SetProfilePicture(user.Profile.ProfilePicture).
				Save(ctx)
			if err != nil {
				return domain.User{}, errors.NewError(errors.ErrorInternal, err)
			}

			_, err = tx.User.UpdateOneID(updatedUser.ID).
				SetProfile(createdProfile).
				Save(ctx)
			if err != nil {
				return domain.User{}, errors.NewError(errors.ErrorInternal, err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return domain.User{}, errors.NewError(errors.ErrorInternal, err)
	}

	// get user with profile and role information
	updatedUser, err = r.client.User.Query().
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

	// Map the updated entities back to the domain model
	return domain.User{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Password: updatedUser.Password,
		Email:    updatedUser.Email,
		Role: domain.Role{
			ID:          updatedUser.Edges.Role.ID,
			Name:        string(updatedUser.Edges.Role.Name),
			Premissions: updatedUser.Edges.Role.Permissions,
		},
		Profile: func() domain.Profile {
			if updatedUser.Edges.Profile != nil {
				return domain.Profile{
					ID:             updatedUser.Edges.Profile.ID,
					FirstName:      updatedUser.Edges.Profile.FirstName,
					LastName:       updatedUser.Edges.Profile.LastName,
					ProfilePicture: updatedUser.Edges.Profile.ProfilePicture,
				}
			}
			return domain.Profile{}
		}(),
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

// Default permissions based on role name
func defaultPermissions(roleName string) []string {
	if roleName == "admin" {
		return []string{"read", "create", "update", "delete"}
	}
	return []string{"create", "update"}
}
