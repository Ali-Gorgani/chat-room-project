package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.String("content").
			NotEmpty(),
		field.String("room_id").
			NotEmpty(),
		field.String("username").
			NotEmpty(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return nil
}
