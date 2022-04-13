package routers

import (
	"torchizm/library-backend/api/auth"
	"torchizm/library-backend/api/user"
	"torchizm/library-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes Employe
func UserRoutes(ctx fiber.Router) {
	route := ctx.Group("/user")
	route.Get("/", user.GetAll)
	route.Post("/login", auth.Login)
	route.Post("/register", auth.Register)
	route.Post("/resend-confirmation", auth.ResendMail)
	route.Post("/activate", auth.ActivateAccount)
	route.Post("/logout", middlewares.IsAuth, auth.LogOut)
}
