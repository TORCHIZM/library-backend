package auth

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"`
	jwt.RegisteredClaims
}
