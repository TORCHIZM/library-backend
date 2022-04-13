package auth

import (
	"time"
	"torchizm/library-backend/config"
	"torchizm/library-backend/helpers"
	"torchizm/library-backend/models"
	"torchizm/library-backend/models/auth"
	"torchizm/library-backend/utils"

	c "github.com/firmanJS/fiber-with-mongo/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *fiber.Ctx) error {
	params := &auth.LoginParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, errors.Error(), errors)
	}

	if err := helpers.ValidateStruct(params); err != nil {
		return helpers.ServerResponse(ctx, "Failed", err)
	}

	userCollection := config.Instance.Database.Collection("user")
	filter := bson.D{{Key: "username", Value: params.Username}}

	user := &models.User{}

	if err := userCollection.FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return helpers.NotFoundResponse(ctx, "User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return helpers.NotFoundResponse(ctx, "Credentials not matching our records")
	}

	if !user.Active {
		return helpers.BadResponse(ctx, "Your account is not verified yet", nil)
	}

	session := &models.Session{}

	var jwtKey = []byte(c.Config("SECRET_KEY"))
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &auth.Claims{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Platform: params.Platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	sessionCollection := config.Instance.Database.Collection("session")

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	if err != nil {
		return helpers.CrudResponse(ctx, "Sid token couldn't signed", nil)
	}

	if err := sessionCollection.FindOne(ctx.Context(), bson.D{
		{Key: "owner", Value: user.ID},
		{Key: "platform", Value: params.Platform},
	}).Decode(&session); err != nil {
		session := &models.Session{
			Owner:     user.ID,
			Jwt:       token,
			Platform:  params.Platform,
			Expires:   expirationTime,
			CreatedAt: utils.MakeTimestamp(),
		}

		if _, err := sessionCollection.InsertOne(ctx.Context(), session); err != nil {
			return helpers.BadResponse(ctx, "An error has been occurred", nil)
		}

		return helpers.CrudResponse(ctx, "Create", session)
	}

	return helpers.CrudResponse(ctx, "Create", session)
}

func LogOut(ctx *fiber.Ctx) error {
	sessionCollection := config.Instance.Database.Collection("session")
	session := ctx.Locals("session").(*models.Session)
	sessionFilter := bson.D{{Key: "_id", Value: session.ID}}

	if _, err := sessionCollection.DeleteOne(ctx.Context(), sessionFilter); err != nil {
		return helpers.ServerResponse(ctx, "An error has been occurred", nil)
	}

	return helpers.MsgResponse(ctx, "Logged out", nil)
}
