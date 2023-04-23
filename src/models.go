package src

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LiquidityPoolConfigs struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ConfigID    string             `json:"configID" bson:"configID"`
	BaseAmount  float64            `json:"baseAmount" bson:"baseAmount"`
	QuoteAmount float64            `json:"quoteAmount" bson:"quoteAmount"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
