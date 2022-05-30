package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type      string             `json:"type" bson:"type" validate:"required,commentType"`
	Object    primitive.ObjectID `json:"object" bson:"object" validate:"required"`
	Name      string             `json:"username" bson:"username" validate:"required,min=6,max=32,string"`
	Author    string             `json:"password" bson:"password" validate:"required,min=6,max=32,password"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt"`
}
