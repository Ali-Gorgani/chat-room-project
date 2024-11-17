package user

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/repository/user"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
)

type UsersService struct {
	c user.IClient
}

func NewUserService(logger *logger.Logger, config *configs.Config) *UsersService {
	c, err := user.NewClient(logger, config)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to establish connection with UserService: %v", err))
		return nil
	}
	return &UsersService{
		c: c,
	}
}

func (s *UsersService) GetUserByUsername(ctx context.Context, req domain.Auth) (domain.Auth, error) {
	dtoReq := MapDomainUserToDtoGetUserReq(req)
	dtoRes, err := s.c.GetUserByUsername(ctx, dtoReq)
	if err != nil {
		return domain.Auth{}, err
	}
	return MapDtoUserResToDomainUser(dtoRes), nil
}
