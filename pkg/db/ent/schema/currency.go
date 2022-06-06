package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/oracle-manager/pkg/db/mixin"
	"github.com/google/uuid"
)

// Currency holds the schema definition for the Currency entity.
type Currency struct {
	ent.Schema
}

func (Currency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

// Fields of the Currency.
func (Currency) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("coin_type_id", uuid.UUID{}).Unique(),
		field.Uint64("price_vs_usdt"),
	}
}
