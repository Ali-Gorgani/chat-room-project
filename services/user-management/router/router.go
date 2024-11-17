package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupUserRouter(handler *handler.UserHandler) *fiber.App {
	app := fiber.New()

	app.Post("/users", handler.CreateUser)
	app.Get("/users/:id", middleware.AuthMiddleware(), handler.FindUserByID)
	app.Put("/users/:id", middleware.AuthMiddleware(), handler.UpdateUser)
	app.Delete("/users/:id", middleware.AuthMiddleware(), handler.DeleteUser)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
