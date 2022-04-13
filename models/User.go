package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Active       bool               `json:"active" bson:"active"`
	Username     string             `json:"username" bson:"username" validate:"required,min=6,max=32,string"`
	Password     string             `json:"password" bson:"password" validate:"required,min=6,max=32,password"`
	Email        string             `json:"email" bson:"email" validate:"required,min=5,max=96,email"`
	FullName     string             `json:"fullName" bson:"fullName" validate:"required,min=6,max=48,stringWithSpace"`
	ProfileImage string             `json:"profileImage" bson:"profileImage" validate:"imagelink"`
	TrustLevel   int64              `json:"trustLevel,omitempty" bson:"trustLevel" validate:"number"`
	Role         primitive.ObjectID `json:"role,omitempty" bson:"role"`
	DateOfBirth  time.Time          `json:"dateOfBirth" bson:"dateOfBirth" validate:"required"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt,omitempty" bson:"updatedAt"`
}
