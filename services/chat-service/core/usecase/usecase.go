package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/grpc/service/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ws"
	"github.com/google/uuid"
)

type ChatUseCase struct {
	chatRepository ports.IChatRepository
	authService    *auth.AuthService
	logger         *logger.Logger
	config         *configs.Config
	hub            *ws.Hub
}

func NewChatUseCase(chatRepository ports.IChatRepository, authService *auth.AuthService, logger *logger.Logger, config *configs.Config, hub *ws.Hub) *ChatUseCase {
	return &ChatUseCase{
		chatRepository: chatRepository,
		logger:         logger,
		config:         config,
		hub: hub,
	}
}

func (uc *ChatUseCase) CreateRoom(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	chat.Room.ID = uuid.New().String()
	createdRoom, err := uc.chatRepository.AddRoom(ctx, chat)
	if err != nil {
		uc.logger.Error(fmt.Sprintf("error creating room: %v", err))
		return domain.Chat{}, err
	}

	return createdRoom, nil
}

func (uc *ChatUseCase) JoinRoom(ctx context.Context, chat domain.Chat) error {
	client := &ws.Client{
		Conn:     chat.Conn,
		Message:  make(chan *ws.Message, 10),
		ID:       chat.User.ID,
		RoomID:   chat.Room.ID,
		Username: chat.User.Username,
	}

	// Register the client
	uc.hub.Register <- client

	// Handle message reading and writing
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		client.ReadMessage(uc.hub)
	}()

	go func() {
		defer wg.Done()
		client.WriteMessage()
	}()

	wg.Wait()
	return nil
}
