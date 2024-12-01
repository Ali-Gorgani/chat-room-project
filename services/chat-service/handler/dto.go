package handler

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/core/domain"
	"github.com/gofiber/websocket/v2"
)

type CreateRoomRequest struct {
	ID   string `json:"id"`
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

func dtoJoinRoomReqToDomainChat(roomID string, req JoinRoomRequest, conn *websocket.Conn) domain.Chat {
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
