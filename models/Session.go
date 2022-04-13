package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Owner     primitive.ObjectID `json:"owner" bson:"owner"`
	Jwt       string             `json:"sid" bson:"sid"`
	Platform  string             `json:"platform" bson:"platform" validate:"required,platform"`
	Expires   time.Time          `json:"expires" bson:"expires"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
