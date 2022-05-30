package feed

import (
	"time"
	"torchizm/library-backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Author    models.User        `json:"author" bson:"author"`
	Content   string             `json:"content" bson:"content" validate:"required,min=5"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt"`
}
