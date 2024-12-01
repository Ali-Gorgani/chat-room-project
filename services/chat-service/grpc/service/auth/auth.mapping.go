package auth

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/grpc/repository/auth"
	"strconv"
)

func MapDomainVerifyTokenReqToDtoVerifyTokenReq(req domain.Auth) auth.VerifyTokenReq {
	return auth.VerifyTokenReq{
		Token: req.AccessToken,
	}
}

func MapDtoVerifyTokenResToDomainVerifyTokenRes(res auth.VerifyTokenRes) domain.User {
	return domain.User{
		ID:       strconv.Itoa(res.ID),
		Username: res.Username,
		Email:    res.Email,
		Role: domain.Role{
			Name: res.Role,
		},
	}
}
