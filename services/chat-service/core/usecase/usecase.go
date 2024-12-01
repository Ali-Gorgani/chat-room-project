package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/grpc/service/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/configs"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/errors"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/logger"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ws"
)

type ChatUseCase struct {
	chatRepository ports.IChatRepository
	// authService    *auth.AuthService
	logger *logger.Logger
	config *configs.Config
	hub    *ws.Hub
}

func NewChatUseCase(chatRepository ports.IChatRepository, authService *auth.AuthService, logger *logger.Logger, config *configs.Config, hub *ws.Hub) *ChatUseCase {
	return &ChatUseCase{
		chatRepository: chatRepository,
		logger:         logger,
		config:         config,
		hub:            hub,
	}
}

func (uc *ChatUseCase) CreateRoom(ctx context.Context, chat domain.Chat) (domain.Chat, error) {

	createdRoom, err := uc.chatRepository.AddRoom(ctx, chat)
	if err != nil {
		uc.logger.Error(fmt.Sprintf("error creating room: %v", err))
		return domain.Chat{}, err
	}

	// Initialize a room with updated Clients structure
	uc.hub.Lock()
	uc.hub.Rooms[createdRoom.Room.ID] = &ws.Room{
		ID:      createdRoom.Room.ID,
		Name:    createdRoom.Room.Name,
		Clients: make(map[string][]*ws.Client), // Updated to match new struct
	}
	uc.hub.Unlock()

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

func (uc *ChatUseCase) GetRooms(ctx context.Context) ([]domain.Chat, error) {
	rooms, err := uc.chatRepository.GetRooms(ctx)
	if err != nil {
		uc.logger.Error(fmt.Sprintf("error getting rooms: %v", err))
		return nil, err
	}

	uc.hub.Lock()
	for _, room := range rooms {
		uc.hub.Rooms[room.Room.ID] = &ws.Room{
			ID:      room.Room.ID,
			Name:    room.Room.Name,
			Clients: uc.hub.Rooms[room.Room.ID].Clients,
		}
	}
	uc.hub.Unlock()

	return rooms, nil
}

func (uc *ChatUseCase) GetClients(ctx context.Context, chat domain.Chat) ([]domain.Chat, error) {
	uc.hub.Lock()
	defer uc.hub.Unlock()

	room, ok := uc.hub.Rooms[chat.Room.ID]
	if !ok {
		return nil, errors.NewError(errors.ErrorNotFound, fmt.Errorf("room not found"))
	}

	var clients []domain.Chat
	for _, clientList := range room.Clients {
		for _, client := range clientList {
			clients = append(clients, domain.Chat{
				User: domain.User{
					ID:       client.ID,
					Username: client.Username,
				},
			})
		}
	}

	return clients, nil
}
