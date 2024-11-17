package grpchandler

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/pkg/user"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
)

type UserHandler struct {
	user.UsersServiceServer
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(client *ent.Client, logger *logger.Logger, config *configs.Config) *UserHandler {
	return &UserHandler{
		userUseCase: usecase.NewUserUseCase(client, logger, config),
	}
}

func (h *UserHandler) GetUserByUsername(ctx context.Context, req *user.GetUserReq) (*user.UserRes, error) {
	res, err := h.userUseCase.FindUserByUsername(ctx, MapProtoGetUserReqToDomainAuth(req))
	if err != nil {
		return &user.UserRes{}, err
	}
	return MapDomainUserToProtoUserRes(res), nil
}
