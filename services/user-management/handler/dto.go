package handler

import "github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	ProfilePicture string `json:"profilePicture"`
}

type UpdateUserRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	ProfilePicture string `json:"profilePicture"`
}

func CreateUserRequestToDomainUser(req CreateUserRequest) domain.User {
	return domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role: domain.Role{
			Name: req.Role,
		},
		Profile: domain.Profile{
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			ProfilePicture: req.ProfilePicture,
		},
	}
}

func UpdateUserRequestToDomainUser(req UpdateUserRequest, userID int) domain.User {
	return domain.User{
		ID:       userID,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role: domain.Role{
			Name: req.Role,
		},
		Profile: domain.Profile{
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			ProfilePicture: req.ProfilePicture,
		},
	}
}

func DomainUserToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.Profile.FirstName,
		LastName:  user.Profile.LastName,
	}
}
