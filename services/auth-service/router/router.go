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

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001,https://localhost:3002",     // Comma-separated origins as a single string
		AllowCredentials: true,                         // Allow cookies and credentials
		AllowMethods:     "GET,POST,PUT,DELETE",        // Specify allowed HTTP methods
		AllowHeaders:     "Content-Type,Authorization", // Specify allowed headers
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
