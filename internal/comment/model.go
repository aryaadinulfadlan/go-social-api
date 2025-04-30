package comment

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,min=10,max=100"`
}
