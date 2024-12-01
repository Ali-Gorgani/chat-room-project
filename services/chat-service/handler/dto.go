package handler

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/gofiber/websocket/v2"
)

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type JoinRoomRequest struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func CreateRoomReqToDomainChat(req CreateRoomRequest) domain.Chat {
	return domain.Chat{
		Room: domain.Room{
			Name: req.Name,
		},
	}
}

func DomainChatToRoomRes(chat domain.Chat) RoomRes {
	return RoomRes{
		ID:   chat.Room.ID,
		Name: chat.Room.Name,
	}
}

func JoinRoomReqToDomainChat(roomID string, req JoinRoomRequest, conn *websocket.Conn) domain.Chat {
	return domain.Chat{
		Room: domain.Room{
			ID: roomID,
		},
		User: domain.User{
			ID:       req.UserID,
			Username: req.Username,
		},
		Conn: conn,
	}
}

func DomainChatToGetRoomsRes(chat []domain.Chat) []RoomRes {
	var res []RoomRes
	for _, c := range chat {
		res = append(res, RoomRes{
			ID:   c.Room.ID,
			Name: c.Room.Name,
		})
	}
	return res
}

func GetClientsReqToDomainChat(roomID string) domain.Chat {
	return domain.Chat{
		Room: domain.Room{
			ID: roomID,
		},
	}
}

func DomainChatToGetClientsRes(chat []domain.Chat) []ClientRes {
	var res []ClientRes
	for _, c := range chat {
		res = append(res, ClientRes{
			ID:       c.User.ID,
			Username: c.User.Username,
		})
	}
	return res
}
