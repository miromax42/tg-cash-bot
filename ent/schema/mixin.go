package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

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
