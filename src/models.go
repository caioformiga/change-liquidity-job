package src

import "time"

type LiquidityPoolConfigs struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	UserID      string    `json:"userID,omitempty" bson:"userID,omitempty"`
	BaseAmount  float64   `json:"baseAmount" bson:"baseAmount"`
	QuoteAmount float64   `json:"quoteAmount" bson:"quoteAmount"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
