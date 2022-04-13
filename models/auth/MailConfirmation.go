package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MailConfirmation struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User      primitive.ObjectID `json:"User" valid:"required" bson:"User"`
	Code      int                `json:"Code" valid:"required" bson:"Code"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
