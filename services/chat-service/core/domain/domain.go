package domain

import "github.com/gofiber/websocket/v2"

// User represents a user in the chat
type User struct {
	ID       string
	Username string
	Email    string
	Role     Role
}

type Role struct {
	Name        string
	Premissions []string
}

type Auth struct {
	AccessToken string
}

type Room struct {
	ID   string
	Name string
}

type Message struct {
	ID       int
	RoomID   string
	Username string
	Content  string
}

type Chat struct {
	Room    Room
	Message Message
	User    User
	Conn    *websocket.Conn
}
