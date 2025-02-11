package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlayerID  string             `bson:"player_id" json:"player_id"`
	Nickname  string             `bson:"nickname" json:"nickname"`
	Level     int                `bson:"level" json:"level"`
	VipLevel  int                `bson:"vip_level" json:"vip_level"`
	Items     []Item             `bson:"items" json:"items"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Item struct {
	ItemID string `bson:"item_id" json:"item_id"`
	Amount int    `bson:"amount" json:"amount"`
}
