package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupChatRouter(handler *handler.ChatHandler) *fiber.App {
	app := fiber.New()

	// Define Swagger route first to keep it unaffected by ChatMiddleware.
	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
