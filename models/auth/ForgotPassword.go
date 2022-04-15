package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForgotPassword struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"Email" validate:"required" bson:"Email"`
	Code      int                `json:"Code" validate:"required" bson:"Code"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
