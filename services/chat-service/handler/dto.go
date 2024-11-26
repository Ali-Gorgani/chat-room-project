package handler

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
