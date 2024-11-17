package user

import "github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/pkg/user"

func MapDtoGetUserReqToPbGetUserReq(req GetUserReq) *user.GetUserReq {
	return &user.GetUserReq{
		Username: req.Username,
	}
}

func MapPbGetUserResToDtoGetUserRes(res *user.UserRes) UserRes {
	return UserRes{
		ID:       uint(res.Id),
		Username: res.Username,
		Password: res.Password,
		Email:    res.Email,
		Role: Role{
			Name:        res.Role.Name,
			Premissions: res.Role.Premissions,
		},
	}
}
