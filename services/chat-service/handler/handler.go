package handler

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/usecase"
	myWebsocket "github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/websocket"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	usecase *usecase.ChatUseCase
	client  *redis.Client
	manager *myWebsocket.Manager
}

func NewChatHandler(usecase *usecase.ChatUseCase, client *redis.Client, manager *myWebsocket.Manager) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
		client:  client,
		manager: manager,
	}
}

func (h *ChatHandler) ServeWS(ctx *fiber.Ctx) error {
	h.manager.ServeWS(ctx)
	return nil
}
