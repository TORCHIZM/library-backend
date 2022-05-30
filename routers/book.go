package routers

import (
	"torchizm/library-backend/api/book"
	"torchizm/library-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func BookRoutes(ctx fiber.Router) {
	route := ctx.Group("/book")
	route.Post("/get-reading-books", middlewares.IsAuth, book.GetReadingBooks)
	route.Post("/new-reading", middlewares.IsAuth, book.NewReading)
}
