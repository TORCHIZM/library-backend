package book

import (
	"fmt"
	"torchizm/library-backend/config"
	"torchizm/library-backend/helpers"
	"torchizm/library-backend/models"
	"torchizm/library-backend/models/book"
	"torchizm/library-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetReadingBooks(ctx *fiber.Ctx) error {
	userBooksCollection := config.Instance.Database.Collection("userbooks")
	user := ctx.Locals("user").(*models.User)
	bookFilter := bson.D{
		{Key: "reader", Value: user.ID},
		{Key: "status", Value: "reading"},
	}

	books := &[]models.Book{}
	userBooksCollection.Find(ctx.Context(), bookFilter)
	userBooks, userBooksErr := userBooksCollection.Find(ctx.Context(), bookFilter)

	if userBooksErr != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	if err := userBooks.Decode(&books); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	return helpers.CrudResponse(ctx, "Books", &books)
}

func NewReading(ctx *fiber.Ctx) error {
	fmt.Println("ananı sikim")
	params := &book.NewReadingBookParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, errors.Error(), errors)
	}

	if err := helpers.ValidateStruct(params); err != nil {
		return helpers.ServerResponse(ctx, "Failed", err)
	}

	userBooksCollection := config.Instance.Database.Collection("userbooks")
	user := ctx.Locals("user").(*models.User)

	if _, err := userBooksCollection.InsertOne(ctx.Context(), &models.UserBook{
		Book:      params.BookId,
		User:      user.ID,
		Status:    "reading",
		CreatedAt: utils.MakeTimestamp(),
		UpdatedAt: utils.MakeTimestamp(),
	}); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	return helpers.MsgResponse(ctx, "Success")
}
