package grpchandler

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/pkg/auth"
)

type AuthHandler struct {
	auth.AuthServiceServer
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
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
