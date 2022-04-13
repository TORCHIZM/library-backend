package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivateParams struct {
	User primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty" validate:"required"`
	Code int                `json:"code" bson:"code" validate:"required,min=100000,max=999999,number"`
}
