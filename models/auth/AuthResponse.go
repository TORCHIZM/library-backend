package auth

import "torchizm/library-backend/models"

type AuthResponse struct {
	User    *models.User    `json:"user,omitempty"`
	Session *models.Session `json:"session,omitempty"`
}
