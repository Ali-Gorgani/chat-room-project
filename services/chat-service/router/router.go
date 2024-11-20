package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupChatRouter(chatHandler *handler.ChatHandler) *fiber.App {
	app := fiber.New()

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Chat routes
	app.Get("/protected", middleware.AuthMiddleware(), chatHandler.ProtectedEndpoint)

	// Serve static files
	app.Static("/", "./views", fiber.Static{
		Index: "index.html", // Serve index.html by default
	})

	return app
}
