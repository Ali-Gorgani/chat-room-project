package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
)

func SetupChatRouter(chatHandler *handler.ChatHandler) *fiber.App {
	// Initialize the HTML engine for rendering views
	engine := html.New("./views", ".html")

	// Create a new Fiber app with custom config
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Serve static files from the "public" directory
	app.Static("/", "./public", fiber.Static{
		Index: "home.html", // Serve index.html by default
	})

	// WebSocket routes with middleware
	app.Use("/ws", middleware.WSMiddleware())
	app.Get("/ws", chatHandler.ServeWS)

	return app
}
