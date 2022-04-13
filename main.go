package main

import (
	"log"
	"os"

	"torchizm/library-backend/api"
	"torchizm/library-backend/config"
	"torchizm/library-backend/helpers"
	"torchizm/library-backend/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// connect to the database
	dbError := config.Connect()
	if dbError != nil {
		log.Fatal(dbError)
		return
	}

	app := fiber.New()

	// middlewares
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(logger.New())

	app.Get("/", api.Index)

	// routers setup
	routers.SetupRoutes(app)

	// handle 404
	app.Use(func(c *fiber.Ctx) error {
		return helpers.NotFoundResponse(c, nil)
	})

	helpers.RegisterCustomValidations()

	// get the port
	port := os.Getenv("PORT")

	// launch the app
	launchError := app.Listen(":" + port)
	if launchError != nil {
		panic(launchError)
	}
}
