package ws

import "sync"

type Room struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	Clients map[string][]*Client `json:"clients"`
}

type Hub struct {
	Rooms        map[string]*Room
	Register     chan *Client
	Unregister   chan *Client
	Broadcast    chan *Message
	sync.RWMutex // Mutex to protect shared data
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			h.Lock()
			room, ok := h.Rooms[cl.RoomID]
			if !ok {
				// If the room does not exist, create it
				room = &Room{
					ID:      cl.RoomID,
					Name:    "Default Room Name", // Add default name or fetch it dynamically
					Clients: make(map[string][]*Client),
				}
				h.Rooms[cl.RoomID] = room
			}

			// Check if the user is already connected
			isNewUser := len(room.Clients[cl.ID]) == 0

			// Add the client to the slice
			room.Clients[cl.ID] = append(room.Clients[cl.ID], cl)
			h.Unlock()

			// Broadcast "joined the room" only for the first connection
			if isNewUser {
				h.Broadcast <- &Message{
					Content:  cl.Username + " has joined the room",
					RoomID:   cl.RoomID,
					Username: cl.Username,
				}
			}

		case cl := <-h.Unregister:
			h.Lock()
			if room, ok := h.Rooms[cl.RoomID]; ok {
				if clientList, ok := room.Clients[cl.ID]; ok {
					// Remove the specific client from the slice
					for i, c := range clientList {
						if c == cl {
							room.Clients[cl.ID] = append(clientList[:i], clientList[i+1:]...)
							break
						}
					}

					// If no more connections exist for the user ID, remove the entry
					if len(room.Clients[cl.ID]) == 0 {
						delete(room.Clients, cl.ID)

						// Broadcast "left the chat" when the user disconnects entirely
						h.Broadcast <- &Message{
							Content:  cl.Username + " has left the room",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					close(cl.Message) // Clean up the client's message channel
				}
			}
			h.Unlock()

		case m := <-h.Broadcast:
			h.RLock() // Use RLock for reading
			if room, ok := h.Rooms[m.RoomID]; ok {
				for _, clientList := range room.Clients {
					for _, cl := range clientList {
						cl.Message <- m
					}
				}
			}
			h.RUnlock()
		}
	}
}
