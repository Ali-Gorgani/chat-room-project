package auth

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/repository/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
)

type AuthService struct {
	c      auth.IClient
	logger *logger.Logger
}

func NewAuthService(logger *logger.Logger, config *configs.Config) *AuthService {
	c, err := auth.NewClient(logger, config)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to establish connection with AuthService: %v", err))
		return nil
	}
	return &AuthService{
		c:      c,
		logger: logger,
	}
}

func (s *AuthService) HashPassword(ctx context.Context, req domain.User) (domain.User, error) {
	dtoReq := MapDomainHashPasswordReqToDtoHashPasswordReq(req)
	dtoRes, err := s.c.HashPassword(ctx, dtoReq)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to call HashPassword: %v", err))
		return domain.User{}, err
	}
	return MapDtoHashPasswordResToDomainHashPasswordRes(dtoRes), nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req domain.Auth) (domain.User, error) {
	dtoReq := MapDomainVerifyTokenReqToDtoVerifyTokenReq(req)
	dtoRes, err := s.c.VerifyToken(ctx, dtoReq)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to call VerifyToken: %v", err))
		return domain.User{}, err
	}
	return MapDtoVerifyTokenResToDomainVerifyTokenRes(dtoRes), nil
}
