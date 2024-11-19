package repository

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
)

type ChatRepository struct {
	client *ent.Client
	logger *logger.Logger
}

func NewChatRepository(client *ent.Client, logger *logger.Logger) ports.IChatRepository {
	return &ChatRepository{
		client: client,
		logger: logger,
	}
}
