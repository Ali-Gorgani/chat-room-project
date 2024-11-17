package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupAuthRouter(handler *handler.UserHandler) *fiber.App {
	app := fiber.New()

	// Define Swagger route first to keep it unaffected by AuthMiddleware.
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Public routes
	app.Post("/refresh-token", handler.RefreshToken)
	app.Post("/login", handler.Login)

	// Routes protected by AuthMiddleware
	app.Post("/logout", middleware.AuthMiddleware(), handler.Logout)
	app.Post("/revoke-token", middleware.AuthMiddleware(), handler.RevokeToken)

	return app
}
