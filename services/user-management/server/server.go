package server

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Server struct {
	app    *fiber.App
	logger *logger.Logger
	config *configs.Config
}

func NewServer(app *fiber.App, logger *logger.Logger, config *configs.Config) *Server {
	return &Server{
		app:    app,
		logger: logger,
		config: config,
	}
}
func (srv *Server) SetupUserServer(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			srv.logger.Info("Starting server")
			go srv.app.Listen(fmt.Sprintf(":%d", srv.config.Server.Port))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.logger.Info("Shutting down server")
			return srv.app.Shutdown()
		},
	})
}
