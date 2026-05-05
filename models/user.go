package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName     string             `json:"first_name" bson:"first_name"`
	LastName      string             `json:"last_name" bson:"last_name"`
	EmailID       string             `json:"email_id" bson:"email_id"`
	Password      string             `json:"password" bson:"password"`
	WalletAddress string             `json:"wallet_address,omitempty" bson:"wallet_address,omitempty"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}