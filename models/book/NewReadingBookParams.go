package book

import "go.mongodb.org/mongo-driver/bson/primitive"

type NewReadingBookParams struct {
	BookId primitive.ObjectID `json:"book" validate:"required"`
}
