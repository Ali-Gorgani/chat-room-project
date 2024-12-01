package handler

import (
	"log"

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

	createRoomRes, err := h.usecase.CreateRoom(ctx.Context(), CreateRoomReqToDomainChat(req))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainChatToRoomRes(createRoomRes)

	return ctx.Status(fiber.StatusCreated).JSON(res)
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
			h.usecase.JoinRoom(ctx.Context(), JoinRoomReqToDomainChat(roomID, req, conn))
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
	rooms, err := h.usecase.GetRooms(ctx.Context())
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainChatToGetRoomsRes(rooms)

	return ctx.Status(fiber.StatusOK).JSON(res)
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
	roomID := ctx.Params("roomId")

	getClientsRes, err := h.usecase.GetClients(ctx.Context(), GetClientsReqToDomainChat(roomID))
	if err != nil {
		apiErr := errors.FromError(err)
		return ctx.Status(apiErr.Status).JSON(apiErr)
	}
	res := DomainChatToGetClientsRes(getClientsRes)

	return ctx.Status(fiber.StatusOK).JSON(res)
}
