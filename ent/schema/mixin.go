package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		index.Fields("created_by").
			Annotations(entsql.IndexType("HASH")),
	}
}

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Immutable().
			Default(time.Now()).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		field.Time("update_time").
			UpdateDefault(time.Now).
			Default(time.Now()).
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
	}
}
