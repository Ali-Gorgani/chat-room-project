package grpchandler

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/pkg/user"
)

type UserHandler struct {
	user.UnimplementedUsersServiceServer
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetUserByUsername(ctx context.Context, req *user.GetUserReq) (*user.UserRes, error) {
	res, err := h.userUseCase.FindUserByUsername(ctx, MapProtoGetUserReqToDomainAuth(req))
	if err != nil {
		return &user.UserRes{}, err
	}
	return MapDomainUserToProtoUserRes(res), nil
}
