package schema

import "entgo.io/ent"

// ChatRoom holds the schema definition for the ChatRoom entity.
type ChatRoom struct {
	ent.Schema
}

// Fields of the ChatRoom.
func (ChatRoom) Fields() []ent.Field {
	return nil
}

// Edges of the ChatRoom.
func (ChatRoom) Edges() []ent.Edge {
	return nil
}
