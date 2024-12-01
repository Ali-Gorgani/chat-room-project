package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/handler"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func SetupUserRouter(handler *handler.UserHandler) *fiber.App {
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,https://localhost:3002", // Comma-separated origins as a single string
		AllowCredentials: true,                                                                 // Allow cookies and credentials
		AllowMethods:     "GET,POST,PUT,DELETE",                                                // Specify allowed HTTP methods
		AllowHeaders:     "Content-Type,Authorization",                                         // Specify allowed headers
	}))

	app.Post("/users", handler.CreateUser)
	app.Get("/users/:id", middleware.AuthMiddleware(), handler.FindUserByID)
	app.Put("/users/:id", middleware.AuthMiddleware(), handler.UpdateUser)
	app.Delete("/users/:id", middleware.AuthMiddleware(), handler.DeleteUser)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
