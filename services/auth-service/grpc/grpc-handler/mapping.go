package grpchandler

import (
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/pkg/auth"
)

func MapProtoHashPasswordReqToDomainAuth(in *auth.HashPasswordReq) domain.Auth {
	return domain.Auth{
		User: domain.User{
			Password: in.Password,
		},
	}
}

func MapDomainAuthToProtoHashPasswordRes(res domain.Auth) *auth.HashPasswordRes {
	return &auth.HashPasswordRes{
		HashedPassword: res.User.Password,
	}
}

func MapProtoVerifyTokenReqToDomainAuth(in *auth.VerifyTokenReq) domain.Auth {
	return domain.Auth{
		AccessToken: in.Token,
	}
}

func MapDomainAuthToProtoVerifyTokenRes(res domain.Auth) *auth.VerifyTokenRes {
	return &auth.VerifyTokenRes{
		Id:       int32(res.User.ID),
		Username: res.User.Username,
		Email:    res.User.Email,
		Role:     res.User.Role,
	}
}
