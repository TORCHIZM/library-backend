package auth

type LogoutParams struct {
	Username string `json:"username,omitempty" validate:"required,min=6,max=32,string"`
	Platform string `json:"platform" validate:"required,platform"`
	Sid      string `json:"sid" validate:"required,min=200,max=500,sid"`
}
