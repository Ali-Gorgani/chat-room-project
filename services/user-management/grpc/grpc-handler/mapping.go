package grpchandler

import (
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/pkg/user"
)

func MapProtoGetUserReqToDomainAuth(req *user.GetUserReq) domain.User {
	return domain.User{
		Username: req.Username,
	}
}

func MapDomainUserToProtoUserRes(res domain.User) *user.UserRes {
	return &user.UserRes{
		Id:       int32(res.ID),
		Username: res.Username,
		Password: res.Password,
		Email:    res.Email,
		Role: &user.Role{
			Name:        res.Role.Name,
			Premissions: res.Role.Premissions,
		},
	}
}
