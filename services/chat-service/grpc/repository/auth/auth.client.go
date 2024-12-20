package auth

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/grpc/pkg/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client interface for AuthService
type IClient interface {
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
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", config.GRPC.AuthHost, config.GRPC.AuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to establish connection with AuthService: %v", err))
		return nil, err
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
