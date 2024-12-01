package ports

import (
	"context"

	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
)

type IChatRepository interface {
	AddRoom(ctx context.Context, chat domain.Chat) (domain.Chat, error)
	GetRooms(ctx context.Context) ([]domain.Chat, error)
	AddMessage(ctx context.Context, message domain.Chat) (domain.Chat, error)
	GetMessagesByRoomID(ctx context.Context, chat domain.Chat) ([]domain.Chat, error)
}
