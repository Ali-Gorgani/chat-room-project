package repository

import (
	"context"
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/ports"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent"
	EntMessage "github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/message"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/errors"
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

func (r *ChatRepository) AddRoom(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	room := chat.Room
	createdRoom, err := r.client.Room.Create().
		SetID(room.ID).
		SetName(room.Name).
		Save(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("error creating room: %v", err))
		return domain.Chat{}, errors.NewError(errors.ErrorInternal, err)
	}

	res := domain.Chat{
		Room: domain.Room{
			ID:   createdRoom.ID,
			Name: createdRoom.Name,
		},
	}

	return res, nil
}

func (r *ChatRepository) GetRooms(ctx context.Context) ([]domain.Chat, error) {
	rooms, err := r.client.Room.Query().All(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("error getting rooms: %v", err))
		return nil, errors.NewError(errors.ErrorInternal, err)
	}

	var res []domain.Chat
	for _, room := range rooms {
		res = append(res, domain.Chat{
			Room: domain.Room{
				ID:   room.ID,
				Name: room.Name,
			},
		})
	}

	return res, nil
}

func (r *ChatRepository) AddMessage(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	message := chat.Message
	createdMessage, err := r.client.Message.Create().
		SetRoomID(message.RoomID).
		SetUsername(message.Username).
		SetContent(message.Content).
		Save(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("error creating message: %v", err))
		return domain.Chat{}, errors.NewError(errors.ErrorInternal, err)
	}

	res := domain.Chat{
		Message: domain.Message{
			ID:       createdMessage.ID,
			RoomID:   createdMessage.RoomID,
			Username: createdMessage.Username,
			Content:  createdMessage.Content,
		},
	}

	return res, nil
}

func (r *ChatRepository) GetMessagesByRoomID(ctx context.Context, chat domain.Chat) ([]domain.Chat, error) {
	messages, err := r.client.Message.Query().
		Where(EntMessage.RoomIDEQ(chat.Message.RoomID)).
		All(ctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("error getting messages: %v", err))
		return nil, errors.NewError(errors.ErrorInternal, err)
	}

	var res []domain.Chat
	for _, message := range messages {
		res = append(res, domain.Chat{
			Message: domain.Message{
				ID:       message.ID,
				RoomID:   message.RoomID,
				Username: message.Username,
				Content:  message.Content,
			},
		})
	}

	return res, nil
}
