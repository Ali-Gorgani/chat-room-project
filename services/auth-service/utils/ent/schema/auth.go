package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Auth holds the schema definition for the Auth entity.
type Auth struct {
	ent.Schema
}

// Fields of the Auth.
func (Auth) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Unique(),
		field.Uint("user_id").
			Positive(),
		field.String("refresh_token"),
		field.Bool("is_revoked").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("expires_at"),
	}
}

// Edges of the Auth.
func (Auth) Edges() []ent.Edge {
	return nil
}
