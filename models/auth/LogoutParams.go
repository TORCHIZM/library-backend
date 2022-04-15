package auth

type LogoutParams struct {
	Username string `json:"username,omitempty" validate:"required,min=6,max=32,string"`
	Sid      string `json:"sid" validate:"required,min=200,max=500,sid"`
	Platform string `json:"platform" validate:"required,platform"`
}
