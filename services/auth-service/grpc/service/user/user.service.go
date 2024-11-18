package user

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/repository/user"
)

type UsersService struct {
	c user.IClient
}

func NewUserService(c user.IClient) *UsersService {
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
