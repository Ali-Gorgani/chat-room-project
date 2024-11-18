package domain

import "time"

// Room represents a chat room.
type Room struct {
	ID   string 
	Name string 
}

// Message represents a message in a chat room.
type Message struct {
	ID        string    
	RoomID    string    
	SenderID  string    
	Content   string    
	CreatedAt time.Time 
}

// UserRoomAccess represents a user's access to a chat room.
type UserRoomAccess struct {
	UserID string 
	RoomID string 
}
