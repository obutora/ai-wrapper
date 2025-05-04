package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Dot struct {
	ent.Schema
}

func (Dot) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid"),
		field.String("title"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Annotations(entgql.OrderField("UPDATED_AT")),
	}
}

func (Dot) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("uid", "updated_at"),
	}
}

func (Dot) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.RelayConnection(),
		entgql.MultiOrder(),
	}
}
