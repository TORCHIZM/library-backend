package config

import (
	"context"
	"fmt"
	"time"
	"torchizm/library-backend/models"
	"torchizm/library-backend/utils"

	"github.com/firmanJS/fiber-with-mongo/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var Instance MongoInstance

func Connect() error {
	DatabaseConnection := config.Config("MONGO_HOST")
	DatabaseName := config.Config("MONGO_DB_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI(DatabaseConnection))

	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(DatabaseName)

	if err != nil {
		return err
	}

	Instance = MongoInstance{
		Client:   client,
		Database: db,
	}
	fmt.Println("-----------------------------------")
	fmt.Println("  Database connection established.")

	roleCollection := Instance.Database.Collection("role")

	roles := []string{"user", "partner", "admin"}
	for _, roleName := range roles {
		role := &models.Role{}
		userRoleFilter := bson.D{{Key: "rolename", Value: roleName}}
		if err := roleCollection.FindOne(ctx, userRoleFilter).Decode(&role); err != nil {
			roleCollection.InsertOne(ctx, &models.Role{
				RoleName:  roleName,
				CreatedAt: utils.MakeTimestamp(),
				UpdatedAt: utils.MakeTimestamp(),
			})

			fmt.Printf("  Role %s created\n", roleName)
		}
	}

	fmt.Println("-----------------------------------")
	return nil
}
