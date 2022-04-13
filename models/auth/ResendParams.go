package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResendParams struct {
	User primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty" validate:"required"`
}
