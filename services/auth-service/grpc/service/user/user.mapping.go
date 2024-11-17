package user

import (
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/repository/user"
)

func MapDomainUserToDtoGetUserReq(req domain.Auth) user.GetUserReq {
	return user.GetUserReq{
		Username: req.User.Username,
	}
}

func MapDtoUserResToDomainUser(res user.UserRes) domain.Auth {
	return domain.Auth{
		User: domain.User{
			ID:       res.ID,
			Username: res.Username,
			Password: res.Password,
			Email:    res.Email,
			Role:     res.Role.Name,
		},
	}
}
