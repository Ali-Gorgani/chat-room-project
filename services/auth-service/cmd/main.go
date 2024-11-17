package main

import (
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/core/usecase"
	_ "github.com/Ali-Gorgani/chat-room-project/services/auth-service/docs" // Import Swagger docs
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/grpc"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/repository"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/router"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/server"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/db"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// @title Chat Room Auth API
// @version 1.0
// @description API documentation for the Chat Room Authentication Service
// @host localhost:3001
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
			handler.NewAuthHandler,
			router.SetupAuthRouter,
			server.NewServer,
			grpc.NewGRPCServer, // Provide the gRPC server
			fx.Annotate(
				repository.NewAuthRepository,
				fx.As(new(ports.IAuthRepository)),
			),
			usecase.NewAuthUseCase,
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
			srv.SetupAuthServer(lc)

			// Set up the gRPC server
			grpcSrv.SetupGRPCServer(lc)
		}),
	)
	app.Run()
}
