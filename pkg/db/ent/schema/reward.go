package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/oracle-manager/pkg/db/mixin"
	"github.com/google/uuid"
)

// Reward holds the schema definition for the Reward entity.
type Reward struct {
	ent.Schema
}

func (Reward) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

// Fields of the Reward.
func (Reward) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("coin_type_id", uuid.UUID{}).Unique(),
		field.Uint64("daily_reward"),
	}
}
