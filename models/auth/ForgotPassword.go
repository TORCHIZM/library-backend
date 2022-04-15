package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForgotPassword struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	Code      int                `json:"code"  bson:"code"  validate:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
