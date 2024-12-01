package router

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,https://localhost:3002", // Comma-separated origins as a single string
		AllowCredentials: true,                                                                 // Allow cookies and credentials
		AllowMethods:     "GET,POST,PUT,DELETE",                                                // Specify allowed HTTP methods
		AllowHeaders:     "Content-Type,Authorization",                                         // Specify allowed headers
	}))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Serve static assets (CSS, JS, etc.) from the "public" directory
	app.Static("/assets", "./public")

	// Route to render signup.html
	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("signup", fiber.Map{
			"Title": "Signup",
		})
	})

	// Route to render login.html
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "Login",
		})
	})

	// Route to render rooms.html
	app.Get("/rooms", func(c *fiber.Ctx) error {
		return c.Render("rooms", fiber.Map{
			"Title": "Chat Rooms",
		})
	})

	// Route to render chat.html
	app.Get("/chat", func(c *fiber.Ctx) error {
		return c.Render("chat", fiber.Map{
			"Title": "Chat Room",
		})
	})

	// WebSocket routes
	app.Post("/ws/create-room", chatHandler.CreateRoom)
	app.Get("/ws/join-room/:roomId", chatHandler.JoinRoom)
	app.Get("/ws/get-rooms", chatHandler.GetRooms)
	app.Get("/ws/get-clients/:roomId", chatHandler.GetClients)

	return app
}
