package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/pkg/user"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client interface for UserService
type IClient interface {
	GetUserByUsername(ctx context.Context, req GetUserReq) (UserRes, error)
}

// Client struct for managing connection
type Client struct {
	c      user.UsersServiceClient // gRPC client
	logger *logger.Logger
}

// NewClient creates a new gRPC client for AuthService
func NewClient(logger *logger.Logger, config *configs.Config) (IClient, error) {
	// Establish gRPC connection with the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// UserService connection
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", config.GRPC.UserHost, config.GRPC.UserPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to UserService: %v", err)
	}
	client := user.NewUsersServiceClient(conn)

	return &Client{
		c:      client,
		logger: logger,
	}, nil
}

func (c *Client) GetUserByUsername(ctx context.Context, req GetUserReq) (UserRes, error) {
	res, err := c.c.GetUserByUsername(ctx, MapDtoGetUserReqToPbGetUserReq(req))
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to call GetUserByUsername: %v", err))
		return UserRes{}, err
	}
	return MapPbGetUserResToDtoGetUserRes(res), nil
}
