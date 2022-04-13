package auth

type LoginParams struct {
	Username string `json:"username,omitempty" validate:"required,min=6,max=32,string"`
	Password string `json:"password" validate:"required,min=6,max=32,password"`
	Platform string `json:"platform" validate:"required,platform"`
}
