package grpc

import (
	"context"
	"fmt"
	"net"

	grpchandler "github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/grpc-handler"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/pkg/user"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	client *ent.Client
	logger *logger.Logger
	config *configs.Config
}

func NewGRPCServer(client *ent.Client, logger *logger.Logger, config *configs.Config) *GRPCServer {
	srv := grpc.NewServer()
	userHandler := grpchandler.NewUserHandler(client, logger, config)
	user.RegisterUsersServiceServer(srv, userHandler)

	return &GRPCServer{
		server: srv,
		client: client,
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
				listener, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.config.GRPC.UserPort))
				if err != nil {
					srv.logger.Fatal(fmt.Sprintf("Failed to listen on port %d: %v", srv.config.GRPC.UserPort, err))
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
