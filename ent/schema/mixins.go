package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

type OwnerMixin struct {
	mixin.Schema
}

func (OwnerMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_by").Immutable(),
	}
}

func (OwnerMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_by"),
	}
}
