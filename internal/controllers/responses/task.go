package responses

type Task struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Duration  string `json:"duration,omitempty"`
}
