package auth

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/grpc/repository/auth"
)

type AuthService struct {
	c auth.IClient
}

func NewAuthService(c auth.IClient) *AuthService {
	return &AuthService{
		c: c,
	}
}

func (s *AuthService) VerifyToken(ctx context.Context, req domain.Auth) (domain.User, error) {
	dtoReq := MapDomainVerifyTokenReqToDtoVerifyTokenReq(req)
	dtoRes, err := s.c.VerifyToken(ctx, dtoReq)
	if err != nil {
		return domain.User{}, err
	}
	return MapDtoVerifyTokenResToDomainVerifyTokenRes(dtoRes), nil
}
