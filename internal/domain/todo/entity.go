package todo

import "time"

// Todo represents a single task tracked by the API.
type Todo struct {
	ID        TodoID    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
