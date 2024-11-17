package grpchandler

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/pkg/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
)

type AuthHandler struct {
	auth.AuthServiceServer
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(client *ent.Client, logger *logger.Logger, config *configs.Config) *AuthHandler {
	return &AuthHandler{
		authUseCase: usecase.NewAuthUseCase(client, logger, config),
	}
}

func (c *AuthHandler) HashPassword(ctx context.Context, in *auth.HashPasswordReq) (*auth.HashPasswordRes, error) {
	hashedPassword, err := c.authUseCase.HashPassword(ctx, MapProtoHashPasswordReqToDomainAuth(in))
	if err != nil {
		return nil, err
	}
	return MapDomainAuthToProtoHashPasswordRes(hashedPassword), nil
}

func (c *AuthHandler) VerifyToken(ctx context.Context, in *auth.VerifyTokenReq) (*auth.VerifyTokenRes, error) {
	claims, err := c.authUseCase.VerifyToken(ctx, MapProtoVerifyTokenReqToDomainAuth(in))
	if err != nil {
		return nil, err
	}
	return MapDomainAuthToProtoVerifyTokenRes(claims), nil
}
