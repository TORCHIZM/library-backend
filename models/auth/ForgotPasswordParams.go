package auth

type ForgotPasswordParams struct {
	Email string `json:"email,omitempty" bson:"email,omitempty" validate:"required,email,max=255,min=5"`
}
