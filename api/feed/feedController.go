package feed

import (
	"fmt"
	"torchizm/library-backend/config"
	"torchizm/library-backend/helpers"
	"torchizm/library-backend/models"
	"torchizm/library-backend/models/feed"
	"torchizm/library-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetFeed(ctx *fiber.Ctx) error {
	feedCollection := config.Instance.Database.Collection("feed")

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{}},
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user"},
			{Key: "localField", Value: "author"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "author"},
		}}}

	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$author"},
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "createdAt", Value: -1}},
		},
	}

	feedResult, feedError := feedCollection.Aggregate(ctx.Context(), mongo.Pipeline{matchStage, lookupStage, unwindStage, sortStage})

	if feedError != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	posts := []feed.Post{}

	if err := feedResult.All(ctx.Context(), &posts); err != nil {
		return helpers.ServerResponse(ctx, "Error", err.Error())
	}

	return helpers.CrudResponse(ctx, "Feed", &posts)
}

func NewPost(ctx *fiber.Ctx) error {
	fmt.Println("slema test")
	params := &feed.NewPostParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, errors.Error(), errors)
	}

	if err := helpers.ValidateStruct(params); err != nil {
		fmt.Println("test")
		return helpers.ServerResponse(ctx, "Failed", err)
	}

	feedCollection := config.Instance.Database.Collection("feed")
	user := ctx.Locals("user").(*models.User)

	if _, err := feedCollection.InsertOne(ctx.Context(), &feed.PostBson{
		Content:   params.Content,
		Author:    user.ID,
		CreatedAt: utils.MakeTimestamp(),
		UpdatedAt: utils.MakeTimestamp(),
	}); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	return helpers.MsgResponse(ctx, "Success")
}
