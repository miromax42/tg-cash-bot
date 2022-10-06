package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
)

// PersonalSettings holds the schema definition for the PersonalSettings entity.
type PersonalSettings struct {
	ent.Schema
}

func (PersonalSettings) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the PersonalSettings.
func (PersonalSettings) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Immutable(),
		field.Enum("currency").
			Optional().
			GoType(currency.Token(0)).
			Annotations(&entsql.Annotation{
				Default: currency.TokenRUB.String(),
			}),
	}
}

// Edges of the PersonalSettings.
func (PersonalSettings) Edges() []ent.Edge {
	return nil
}