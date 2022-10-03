package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Expense holds the schema definition for the Expense entity.
type Expense struct {
	ent.Schema
}

func (Expense) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
		OwnerMixin{},
	}
}

// Fields of the Expense.
func (Expense) Fields() []ent.Field {
	return []ent.Field{
		field.Int("amount").
			Positive(),
		field.String("category").Default("UNCATEGORIZED"),
	}
}

// Edges of the Expense.
func (Expense) Edges() []ent.Edge {
	return nil
}
