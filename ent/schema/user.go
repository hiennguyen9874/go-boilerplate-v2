package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name"),
		field.String("email").Unique(),
		field.String("password"),
		field.Bool("is_active").Default(true),
		field.Bool("is_super_user").Default(false),
		field.Bool("verified").Default(false),
		field.String("verification_code").Optional().Nillable(),
		field.String("password_reset_token").Optional().Nillable(),
		field.Time("password_reset_at").Optional().Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Item.Type),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
		// Or, mixin.CreateTime only for create_time
		// and mixin.UpdateTime only for update_time.
	}
}
