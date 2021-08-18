package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wallets are tagged to a single user in a 1-to-1 relationship
type Wallet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty"`
	Balance   int64              `bson:"balance,omitempty"` // assumes a user cannot have a negative balance
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
	IsActive  *bool              `bson:"isActive,omitempty"`
}
