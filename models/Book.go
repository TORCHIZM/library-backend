package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"username" bson:"username" validate:"required,min=6,max=32,string"`
	Author     string             `json:"password" bson:"password" validate:"required,min=6,max=32,password"`
	Email      string             `json:"email" bson:"email" validate:"required,min=5,max=96,email"`
	FullName   string             `json:"fullName" bson:"fullName" validate:"required,min=6,max=48,stringWithSpace"`
	Image      string             `json:"profileImage" bson:"profileImage" validate:"imagelink"`
	Comments   []Comment
	Quotations []Quotation
	CreatedAt  time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
}
