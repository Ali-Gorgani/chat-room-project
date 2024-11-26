package handler

import (
	"log"
	"sync"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/usecase"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ws"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ChatHandler struct {
	usecase *usecase.ChatUseCase
	client  *redis.Client
	hub     *ws.Hub
}

func NewChatHandler(usecase *usecase.ChatUseCase, client *redis.Client, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
		client:  client,
		hub:     hub,
	}
}

// CreateRoom godoc
// @Summary Create a new chat room
// @Description Create a new chat room with the given ID and name
// @Tags chat
// @Accept json
// @Produce json
// @Param CreateRoomRequest body CreateRoomRequest true "Create Room Request"
// @Success 201 {object} CreateRoomRequest
// @Failure 400 {object} map[string]interface{}
// @Router /ws/create-room [post]
func (h *ChatHandler) CreateRoom(ctx *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errors.NewError(errors.ErrorBadRequest, err)
	}

	// Initialize a room with updated Clients structure
	h.hub.Rooms[req.ID] = &ws.Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string][]*ws.Client), // Updated to match new struct
	}

	return ctx.Status(fiber.StatusCreated).JSON(req)
}

var config = websocket.Config{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Origins:         []string{"https://localhost:3002"},
}

func (h *ChatHandler) JoinRoom(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomId")
	if roomID == "" {
		log.Println("Room ID is missing")
		return fiber.NewError(fiber.StatusBadRequest, "Room ID is required")
	}

	var req JoinRoomRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Println("Error parsing query parameters:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	if websocket.IsWebSocketUpgrade(ctx) {
		return websocket.New(func(conn *websocket.Conn) {
			client := &ws.Client{
				Conn:     conn,
				Message:  make(chan *ws.Message, 10),
				ID:       req.UserID,
				RoomID:   roomID,
				Username: req.Username,
			}

			// Register the client
			h.hub.Register <- client

			// Handle message reading and writing
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				client.ReadMessage(h.hub)
			}()

			go func() {
				defer wg.Done()
				client.WriteMessage()
			}()

			wg.Wait()
		}, config)(ctx)
	}

	return fiber.ErrUpgradeRequired
}

// GetRooms godoc
// @Summary Get all chat rooms
// @Description Retrieve a list of all chat rooms
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {array} RoomRes
// @Router /ws/get-rooms [get]
func (h *ChatHandler) GetRooms(ctx *fiber.Ctx) error {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(rooms)
}

// GetClients godoc
// @Summary Get clients in a chat room
// @Description Retrieve a list of clients in the specified chat room
// @Tags chat
// @Accept json
// @Produce json
// @Param roomId path string true "Room ID"
// @Success 200 {array} ClientRes
// @Failure 404 {object} map[string]interface{}
// @Router /ws/get-clients/{roomId} [get]
func (h *ChatHandler) GetClients(ctx *fiber.Ctx) error {
	var clients []ClientRes
	roomID := ctx.Params("roomId")

	// Check if the room exists
	room, ok := h.hub.Rooms[roomID]
	if !ok {
		clients = make([]ClientRes, 0)
		return ctx.Status(fiber.StatusOK).JSON(clients)
	}

	// Aggregate all connections for each user ID
	for _, clientList := range room.Clients {
		for _, client := range clientList {
			clients = append(clients, ClientRes{
				ID:       client.ID,
				Username: client.Username,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(clients)
}
