package websocket

import "github.com/gofiber/websocket/v2"

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string
	RoomID   string
	Username string
}

type Message struct {
	Content   string `json:"content"`
	RoomID    string `json:"room_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}
