package book

type ForgotPasswordConfirmParams struct {
	Name     int    `json:"code"     validate:"required,min=100000,max=999999,number"`
	Password string `json:"password" validate:"required,min=6,max=32,password"`
}
