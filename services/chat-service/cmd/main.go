package main

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/usecase"
	_ "github.com/Ali-Gorgani/chat-room-project/services/chat-service/docs" // Import Swagger docs
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/repository"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/router"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/server"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/db"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// @title Chat Room Chat API
// @version 1.0
// @description API documentation for the Chat Room Chat Service
// @host localhost:3002
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
			handler.NewChatHandler,
			router.SetupChatRouter,
			fx.Annotate(
				repository.NewChatRepository,
				fx.As(new(ports.IChatRepository)),
			),
			usecase.NewChatUseCase,
			server.NewServer,
		),
		fx.Invoke(func(
			lc fx.Lifecycle,
			app *fiber.App,
			logger *logger.Logger,
			config *configs.Config,
			srv *server.Server, // Inject the Fiber server
		) {
			// Set up the Fiber server
			srv.SetupChatServer(lc)
		}),
	)
	app.Run()
}
