package handler

import "github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/usecase"

type ChatHandler struct {
	usecase *usecase.ChatUseCase
}

func NewChatHandler(usecase *usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
	}
}
