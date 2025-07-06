package requests

type CreateTask struct {
	Text string `json:"text" binding:"required" example:"task 1"`
}
