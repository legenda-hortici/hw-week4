package db

import (
	"time"
)

// Task - задача
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateTask - обновленная задача
type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UpdatedAt   time.Time
}
