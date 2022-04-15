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

func Register(ctx *fiber.Ctx) error {
	user := &auth.RegisterParams{}

	if err := ctx.BodyParser(user); err != nil {
		return helpers.ServerResponse(ctx, err.Error(), err.Error())
	}

	if err := helpers.ValidateStruct(user); err != nil {
		return helpers.ServerResponse(ctx, "Failed", err)
	}

	userCollection := config.Instance.Database.Collection("user")
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.M{"username": user.Username},
			bson.M{"email": user.Email},
		}},
	}

	var userRecord bson.M

	if err := userCollection.FindOne(ctx.Context(), filter).Decode(&userRecord); err == nil {
		return helpers.ServerResponse(ctx, "Error", "User already exists")
	}

	hashedPassword, passwordErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if passwordErr != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	roleCollection := config.Instance.Database.Collection("role")
	role := &models.Role{}

	if err := roleCollection.FindOne(ctx.Context(), bson.M{"rolename": "user"}).Decode(&role); err != nil {
		return helpers.ServerResponse(ctx, "Error", "Server is not ready yet")
	}

	insertUser := &models.User{
		Active:       false,
		Username:     user.Username,
		FullName:     user.FullName,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		TrustLevel:   user.TrustLevel,
		Password:     string(hashedPassword),
		Role:         role.ID,
		DateOfBirth:  user.DateOfBirth,
		CreatedAt:    utils.MakeTimestamp(),
		UpdatedAt:    utils.MakeTimestamp(),
	}

	if result, errs := userCollection.InsertOne(ctx.Context(), insertUser); errs != nil {
		return helpers.ServerResponse(ctx, "Error", errs.Error())
	} else {
		filter := bson.D{{Key: "_id", Value: result.InsertedID}}
		createdRecord := userCollection.FindOne(ctx.Context(), filter)
		createduser := &models.User{}
		createdRecord.Decode(createduser)

		rand.Seed(time.Now().UnixNano())
		confirmationNumber := rand.Intn((999999 - 100000)) + 100000

		confirmationCollection := config.Instance.Database.Collection("confirmations")
		confirmation := &auth.MailConfirmation{}
		confirmation.User = createduser.ID
		confirmation.Code = confirmationNumber
		confirmation.CreatedAt = utils.MakeTimestamp()
		confirmation.UpdatedAt = utils.MakeTimestamp()

		if _, errs := confirmationCollection.InsertOne(ctx.Context(), confirmation); errs != nil {
			return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
		}

		var jwtKey = []byte(c.Config("SECRET_KEY"))
		expirationTime := time.Now().Add(time.Hour * 24)

		claims := &auth.Claims{
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Platform: user.Platform,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		sessionCollection := config.Instance.Database.Collection("session")
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

		if err != nil {
			return helpers.ServerResponse(ctx, "Sid token couldn't signed", nil)
		}

		session := &models.Session{
			Owner:     createduser.ID,
			Jwt:       token,
			Platform:  user.Platform,
			Expires:   expirationTime,
			CreatedAt: utils.MakeTimestamp(),
		}

		if _, err := sessionCollection.InsertOne(ctx.Context(), session); err != nil {
			return helpers.BadResponse(ctx, "An error has been occurred")
		}

		response := &auth.AuthResponse{
			User:    createduser,
			Session: session,
		}

		helpers.SendMail(user.Email, confirmationNumber)
		return helpers.CrudResponse(ctx, "Create", response)
	}
}

func ActivateAccount(ctx *fiber.Ctx) error {
	params := &auth.ActivateParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, "Error", errors.Error())
	}
	if errors := helpers.ValidateStruct(params); errors != nil {
		return helpers.ServerResponse(ctx, "Failed", errors)
	}

	userCollection := config.Instance.Database.Collection("user")
	userFilter := bson.D{{Key: "_id", Value: params.User}, {Key: "active", Value: false}}
	user := &models.User{}

	if err := userCollection.FindOne(ctx.Context(), userFilter).Decode(&user); err != nil {
		return helpers.NotFoundResponse(ctx, "User not found")
	}

	confirmationCollection := config.Instance.Database.Collection("confirmations")
	confirmation := &auth.MailConfirmation{}
	confirmationFilter := bson.D{{Key: "user", Value: user.ID}}

	if err := confirmationCollection.FindOne(ctx.Context(), confirmationFilter).Decode(&confirmation); err != nil {
		return helpers.NotFoundResponse(ctx, "Confirmation not found")
	}

	if confirmation.Code != params.Code {
		return helpers.BadResponse(ctx, "Wrong code")
	}

	created := confirmation.CreatedAt.Add(time.Minute * 60)
	if created.Sub(utils.MakeTimestamp()) < 0 {
		return helpers.BadResponse(ctx, "Code expired")
	}

	user.Active = true
	user.UpdatedAt = utils.MakeTimestamp()

	if _, err := userCollection.UpdateOne(ctx.Context(), userFilter, bson.M{
		"$set": user,
	}); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	if _, err := confirmationCollection.DeleteOne(ctx.Context(), confirmationFilter); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	return helpers.MsgResponse(ctx, "Account activated successfully")
}

func ResendMail(ctx *fiber.Ctx) error {
	params := &auth.ResendParams{}

	if errors := ctx.BodyParser(params); errors != nil {
		return helpers.ServerResponse(ctx, "Error", errors.Error())
	}

	if err := helpers.ValidateStruct(params); err != nil {
		return helpers.ServerResponse(ctx, "Failed", err)
	}

	userCollection := config.Instance.Database.Collection("user")
	userFilter := bson.D{{Key: "_id", Value: params.User}, {Key: "active", Value: false}}
	user := &models.User{}

	if err := userCollection.FindOne(ctx.Context(), userFilter).Decode(&user); err != nil {
		return helpers.NotFoundResponse(ctx, "User not found")
	}

	confirmationCollection := config.Instance.Database.Collection("confirmations")
	confirmation := &auth.MailConfirmation{}
	confirmationFilter := bson.D{{Key: "user", Value: user.ID}}

	if err := confirmationCollection.FindOne(ctx.Context(), confirmationFilter).Decode(&confirmation); err != nil {
		return helpers.NotFoundResponse(ctx, "Confirmation not found")
	}

	created := confirmation.CreatedAt.Add(time.Second * 60)
	if created.Sub(utils.MakeTimestamp()) > 0 {
		return helpers.BadResponse(ctx, "You have to wait 1 minute")
	}

	rand.Seed(time.Now().UnixNano())
	confirmationNumber := rand.Intn((999999 - 100000)) + 100000

	if _, err := confirmationCollection.UpdateOne(ctx.Context(), confirmationFilter, bson.M{
		"$set": &auth.MailConfirmation{
			User:      user.ID,
			Code:      confirmationNumber,
			CreatedAt: utils.MakeTimestamp(),
			UpdatedAt: utils.MakeTimestamp(),
		},
	}); err != nil {
		return helpers.ServerResponse(ctx, "Error", "An error has been occurred")
	}

	helpers.SendMail(user.Email, confirmationNumber)
	return helpers.MsgResponse(ctx, "Code sent")
}
