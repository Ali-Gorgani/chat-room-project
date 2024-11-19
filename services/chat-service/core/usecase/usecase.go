package usecase

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
)

type ChatUseCase struct {
	chatRepository ports.IChatRepository
	logger         *logger.Logger
	config         *configs.Config
}

func NewChatUseCase(chatRepository ports.IChatRepository, logger *logger.Logger, config *configs.Config) *ChatUseCase {
	return &ChatUseCase{
		chatRepository: chatRepository,
		logger:         logger,
		config:         config,
	}
}
