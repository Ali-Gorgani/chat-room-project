package domain

import "time"

// User represents a user in the chat
type User struct {
	ID       int
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

// Room represents a chat room.
type Room struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Messages []Message `json:"messages"`
}

// Message represents a chat message in a room.
type Message struct {
	ID      string    `json:"id"`
	RoomID  string    `json:"room_id"`
	UserID  string    `json:"user_id"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}
