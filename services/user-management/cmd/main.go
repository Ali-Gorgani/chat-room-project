package main

import (
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/core/usecase"
	_ "github.com/Ali-Gorgani/chat-room-project/services/user-management/docs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc"
	grpcHandler "github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/grpc-handler"
	grpcAuthRepository "github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/repository/auth"
	grpcAuthService "github.com/Ali-Gorgani/chat-room-project/services/user-management/grpc/service/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/repository"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/router"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/server"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/db"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// @title Chat Room User API
// @version 1.0
// @description API documentation for the Chat Room Authentication Service
// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "JWT Authorization header using the Bearer scheme. Example: \"Bearer {token}\""
func main() {
	app := fx.New(
		logger.Module,
		db.Module,
		configs.Module,
		fx.Provide(
			// http service
			handler.NewUserHandler,
			router.SetupUserRouter,
			fx.Annotate(
				repository.NewUserRepository,
				fx.As(new(ports.IUserRepository)),
			),
			usecase.NewUserUseCase,
			server.NewServer,

			// gRPC service
			grpcHandler.NewUserHandler,
			fx.Annotate(
				grpcAuthRepository.NewClient,
				fx.As(new(grpcAuthRepository.IClient)),
			),
			grpcAuthService.NewAuthService,
			grpc.NewGRPCServer,
		),
		fx.Invoke(func(
			lc fx.Lifecycle,
			app *fiber.App,
			logger *logger.Logger,
			config *configs.Config,
			srv *server.Server, // Inject the Fiber server
			grpcSrv *grpc.GRPCServer, // Inject the gRPC server
		) {
			// Set up the Fiber server
			srv.SetupUserServer(lc)

			// Set up the gRPC server
			grpcSrv.SetupGRPCServer(lc)
		}),
	)
	app.Run()
}
