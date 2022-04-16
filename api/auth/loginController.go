package auth

import (
	"math/rand"
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
		return helpers.BadResponse(ctx, "Credentials not matching our records")
	}

	if !user.Active {
		return helpers.BadResponse(ctx, "Your account is not verified yet")
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
		return helpers.ServerResponse(ctx, "Error", "Sid token couldn't signed")
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
			return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
		}
	}

	response := &auth.AuthResponse{
		User:    user,
		Session: session,
	}

	return helpers.CrudResponse(ctx, "Create", response)
}

func LogOut(ctx *fiber.Ctx) error {
	sessionCollection := config.Instance.Database.Collection("session")
	session := ctx.Locals("session").(*models.Session)
	sessionFilter := bson.D{{Key: "_id", Value: session.ID}}

	if _, err := sessionCollection.DeleteOne(ctx.Context(), sessionFilter); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	return helpers.MsgResponse(ctx, "Logged out")
}

func ForgotPassword(ctx *fiber.Ctx) error {
	params := &auth.ForgotPasswordParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, errors.Error(), errors)
	}

	if err := helpers.ValidateStruct(params); err != nil {
		return helpers.ServerResponse(ctx, "Failed", err)
	}

	user := &models.User{}
	userCollection := config.Instance.Database.Collection("user")
	userFilter := bson.D{{Key: "email", Value: params.Email}}

	if err := userCollection.FindOne(ctx.Context(), userFilter).Decode(user); err != nil {
		return helpers.NotFoundResponse(ctx, "User not found")
	}

	forgotPasswordCollection := config.Instance.Database.Collection("passwordconfirmations")
	forgotPassword := &auth.ForgotPassword{}
	forgotPasswordFilter := bson.D{{Key: "email", Value: user.Email}}

	rand.Seed(time.Now().UnixNano())
	confirmationNumber := rand.Intn((999999 - 100000)) + 100000

	if err := forgotPasswordCollection.FindOne(ctx.Context(), forgotPasswordFilter).Decode(&forgotPassword); err != nil {
		if _, err := forgotPasswordCollection.InsertOne(ctx.Context(), &auth.ForgotPassword{
			Email:     user.Email,
			Code:      confirmationNumber,
			CreatedAt: utils.MakeTimestamp(),
			UpdatedAt: utils.MakeTimestamp(),
		}); err != nil {
			return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
		}

		helpers.SendForgotPassword(user.Email, confirmationNumber)
		return helpers.MsgResponse(ctx, "Code sent")
	}

	updated := forgotPassword.UpdatedAt.Add(time.Minute * 5)
	if updated.Sub(utils.MakeTimestamp()) > 0 {
		return helpers.BadResponse(ctx, "You have to wait 5 minute")
	}

	if _, err := forgotPasswordCollection.UpdateOne(ctx.Context(), forgotPasswordFilter, bson.M{
		"$set": &auth.ForgotPassword{
			Email:     user.Email,
			Code:      confirmationNumber,
			UpdatedAt: utils.MakeTimestamp(),
			CreatedAt: utils.MakeTimestamp(),
		},
	}); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	helpers.SendForgotPassword(user.Email, confirmationNumber)
	return helpers.MsgResponse(ctx, "Code sent")
}

func ForgotPasswordConfirm(ctx *fiber.Ctx) error {
	params := &auth.ForgotPasswordConfirmParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, "Error", errors.Error())
	}

	if errors := helpers.ValidateStruct(params); errors != nil {
		return helpers.ServerResponse(ctx, "Failed", errors)
	}

	forgotPasswordCollection := config.Instance.Database.Collection("passwordconfirmations")
	forgotPassword := &auth.ForgotPassword{}
	forgotPasswordFilter := bson.D{{Key: "code", Value: params.Code}}

	if err := forgotPasswordCollection.FindOne(ctx.Context(), forgotPasswordFilter).Decode(&forgotPassword); err != nil {
		return helpers.NotFoundResponse(ctx, "Confirmation not found")
	}

	if forgotPassword.Code != params.Code {
		return helpers.BadResponse(ctx, "Wrong code")
	}

	updated := forgotPassword.UpdatedAt.Add(time.Minute * 60)
	if updated.Sub(utils.MakeTimestamp()) < 0 {
		return helpers.BadResponse(ctx, "Code expired")
	}

	userCollection := config.Instance.Database.Collection("user")
	userFilter := bson.D{{Key: "email", Value: forgotPassword.Email}}
	user := &models.User{}

	if err := userCollection.FindOne(ctx.Context(), userFilter).Decode(&user); err != nil {
		return helpers.NotFoundResponse(ctx, "User not found")
	}

	hashedPassword, passwordErr := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if passwordErr != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	if string(hashedPassword) == user.Password {
		return helpers.ServerResponse(ctx, "Error", "Your password cannot be the same as the old password")
	}

	user.Password = string(hashedPassword)

	if _, err := userCollection.UpdateOne(ctx.Context(), userFilter, bson.M{
		"$set": user,
	}); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	if _, err := forgotPasswordCollection.DeleteOne(ctx.Context(), forgotPasswordFilter); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	return helpers.MsgResponse(ctx, "Password changed")
}
