package service

// TaskRequest - запрос на создание задачи
type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

// UpdateTaskRequest - запрос на обновление задачи
type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required,oneof=new in_progress done"`
}