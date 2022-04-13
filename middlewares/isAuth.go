package middlewares

import (
	"strings"
	"torchizm/library-backend/config"
	"torchizm/library-backend/helpers"
	"torchizm/library-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BearerHeader struct {
	Token    string `validate:"required,min=7,max=500,sid"`
	Platform string `validate:"required,min=3,max=7,platform"`
}

func IsAuth(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()

	if len(headers["Authorization"]) < 7 {
		return helpers.BadResponse(ctx, "No valid token provided", nil)
	}

	bearerHeader := &BearerHeader{
		Token:    strings.Split(headers["Authorization"], "Bearer ")[1],
		Platform: headers["Platform"],
	}

	if err := helpers.ValidateStruct(bearerHeader); err != nil {
		return helpers.BadResponse(ctx, "Failed", err)
	}

	sessionCollection := config.Instance.Database.Collection("session")
	sessionFilter := bson.D{
		{Key: "sid", Value: bearerHeader.Token},
		{Key: "platform", Value: bearerHeader.Platform},
	}
	session := &models.Session{}

	if err := sessionCollection.FindOne(ctx.Context(), sessionFilter).Decode(&session); err != nil {
		return helpers.ServerResponse(ctx, "Identification failed", nil)
	}

	userCollection := config.Instance.Database.Collection("user")
	userFilter := bson.D{{Key: "_id", Value: session.Owner}}
	user := &models.User{}

	if err := userCollection.FindOne(ctx.Context(), userFilter).Decode(&user); err != nil {
		return helpers.BadResponse(ctx, "User not found", nil)
	}

	ctx.Locals("session", session)
	ctx.Locals("user", user)

	return ctx.Next()
}
