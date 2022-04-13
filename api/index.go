package api

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Username string `json:"username" validate:"required,min=6,max=32,string"`
	Password string `json:"password" validate:"required,min=6,max=32,password"`
	Email    string `json:"email" validate:"required,email"`
}

func Index(c *fiber.Ctx) error {
	return c.JSON(c.App().Stack())
}
