package handler

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	connections = make(map[string]*websocket.Conn) // Active connections
	mu          sync.Mutex                         // Synchronize access
)

func WebSocketHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleWebSocket(c *websocket.Conn) {
	userID := c.Params("userID") // Extract userID from params
	mu.Lock()
	connections[userID] = c
	mu.Unlock()
	defer func() {
		mu.Lock()
		delete(connections, userID)
		mu.Unlock()
		c.Close()
	}()

	fmt.Printf("User %s connected\n", userID)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}
		fmt.Printf("Message from user %s: %s\n", userID, string(msg))

		// Echo message back to user
		if err := c.WriteMessage(websocket.TextMessage, []byte("Message received")); err != nil {
			fmt.Printf("Error writing message: %v\n", err)
			break
		}
	}
}
