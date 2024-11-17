package auth

import (
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/repository/auth"
)

func MapDomainHashPasswordReqToDtoHashPasswordReq(req domain.User) auth.HashPasswordReq {
	return auth.HashPasswordReq{
		Password: req.Password,
	}
}

func MapDtoHashPasswordResToDomainHashPasswordRes(res auth.HashPasswordRes) domain.User {
	return domain.User{
		Password: res.HashedPassword,
	}
}

func MapDomainVerifyTokenReqToDtoVerifyTokenReq(req domain.Auth) auth.VerifyTokenReq {
	return auth.VerifyTokenReq{
		Token: req.AccessToken,
	}
}

func MapDtoVerifyTokenResToDomainVerifyTokenRes(res auth.VerifyTokenRes) domain.User {
	return domain.User{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Role: domain.Role{
			Name: res.Role,
		},
	}
}
