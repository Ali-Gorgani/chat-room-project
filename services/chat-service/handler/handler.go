package handler

import (
	"time"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ChatHandler struct {
	usecase *usecase.ChatUseCase
	client *redis.Client
}

func NewChatHandler(usecase *usecase.ChatUseCase, client *redis.Client) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
		client: client,
	}
}

// ProtectedEndpoint godoc
// @Summary Protected Endpoint
// @Description Access this endpoint only with a valid JWT Bearer token.
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProtectedResponse "Successfully authorized"
// @Failure 400 {object} map[string]interface{} "Bad Request - Missing or invalid token"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Token expired or invalid"
// @Router /protected [get]
func (h *ChatHandler) ProtectedEndpoint(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	h.client.Set(ctx.Context(), "token", token, 24*time.Hour)
	h.client.HSet(ctx.Context(), "refresh_tokens", token, time.Now().Add(24*time.Hour).Unix())

	redisToken := h.client.Get(ctx.Context(), "token").Val()
	log.Infof("Token: %s", redisToken)

	// Use the token to perform actions, such as verifying claims
	return ctx.JSON(fiber.Map{
		"message": "You are authorized!",
		"token":   token,
	})
}
