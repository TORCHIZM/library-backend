package routers

import (
	"torchizm/library-backend/api/feed"
	"torchizm/library-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FeedRoutes(ctx fiber.Router) {
	route := ctx.Group("/feed")
	route.Get("/", middlewares.IsAuth, feed.GetFeed)
	route.Post("/new-post", middlewares.IsAuth, feed.NewPost)
}
