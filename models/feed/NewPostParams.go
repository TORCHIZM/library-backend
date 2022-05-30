package feed

type NewPostParams struct {
	Content string `json:"content" validate:"required,min=5"`
}
