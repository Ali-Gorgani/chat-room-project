package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/pkg/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client interface for AuthService
type IClient interface {
	HashPassword(ctx context.Context, req HashPasswordReq) (HashPasswordRes, error)
	VerifyToken(ctx context.Context, req VerifyTokenReq) (VerifyTokenRes, error)
}

// Client struct for managing connection
type Client struct {
	c      auth.AuthServiceClient // gRPC client
	logger *logger.Logger
}

// NewClient creates a new gRPC client for AuthService
func NewClient(logger *logger.Logger, config *configs.Config) (IClient, error) {
	// Establish gRPC connection with the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// AuthService connection
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", config.GRPC.AuthHost, config.GRPC.AuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to UserService: %v", err)
	}
	client := auth.NewAuthServiceClient(conn)

	return &Client{
		c:      client,
		logger: logger,
	}, nil
}

func (c *Client) HashPassword(ctx context.Context, req HashPasswordReq) (HashPasswordRes, error) {
	res, err := c.c.HashPassword(ctx, MapDtoHashPasswordReqToPbHashPasswordReq(req))
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to call HashPassword: %v", err))
		return HashPasswordRes{}, err
	}
	return MapPbHashPasswordResToDtoHashPasswordRes(res), nil
}

func (c *Client) VerifyToken(ctx context.Context, req VerifyTokenReq) (VerifyTokenRes, error) {
	res, err := c.c.VerifyToken(ctx, MapDtoVerifyTokenReqToPbVerifyTokenReq(req))
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to call VerifyToken: %v", err))
		return VerifyTokenRes{}, err
	}
	return MapPbVerifyTokenResToDtoVerifyTokenRes(res), nil
}
