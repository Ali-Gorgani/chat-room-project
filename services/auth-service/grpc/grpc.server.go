package grpc

import (
	"context"
	"fmt"
	"net"

	grpchandler "github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/grpc-handler"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc/pkg/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	logger *logger.Logger
	config *configs.Config
}

func NewGRPCServer(grpcHandler *grpchandler.AuthHandler, logger *logger.Logger, config *configs.Config) *GRPCServer {
	srv := grpc.NewServer()
	auth.RegisterAuthServiceServer(srv, grpcHandler)

	return &GRPCServer{
		server: srv,
		logger: logger,
		config: config,
	}
}

func (srv *GRPCServer) SetupGRPCServer(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			srv.logger.Info("Starting gRPC server")

			// Start the gRPC server in a separate goroutine
			go func() {
				listener, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.config.GRPC.AuthPort))
				if err != nil {
					srv.logger.Fatal(fmt.Sprintf("Failed to listen on port %d: %v", srv.config.GRPC.AuthPort, err))
				}
				if err := srv.server.Serve(listener); err != nil {
					srv.logger.Fatal(fmt.Sprintf("Failed to serve gRPC: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.logger.Info("Shutting down gRPC server")
			srv.server.GracefulStop()
			return nil
		},
	})
}
