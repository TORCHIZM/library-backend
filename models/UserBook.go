package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBook struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User      primitive.ObjectID `json:"user" bson:"user" validate:"required"`
	Book      primitive.ObjectID `json:"book" bson:"book" validate:"required"`
	Status    string             `json:"status" bson:"status" validate:"required,userBookStatus"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt"`
}
