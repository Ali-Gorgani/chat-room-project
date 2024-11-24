package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func SetupAuthRouter(handler *handler.AuthHandler) *fiber.App {
	app := fiber.New()

	// Customize CORS settings (optional)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://localhost:3002", // Specify allowed origins
		AllowMethods: "GET,POST,PUT,DELETE",   // Specify allowed methods
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

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
