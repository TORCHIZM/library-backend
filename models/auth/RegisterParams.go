package auth

import "torchizm/library-backend/models"

type RegisterParams struct {
	Platform string `json:"platform,omitempty" validate:"required,platform"`
	models.User
}
