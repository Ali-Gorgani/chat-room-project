package schema

import "entgo.io/ent"

// Chat holds the schema definition for the Chat entity.
type Chat struct {
	ent.Schema
}

// Fields of the Chat.
func (Chat) Fields() []ent.Field {
	return nil
}

// Edges of the Chat.
func (Chat) Edges() []ent.Edge {
	return nil
}
