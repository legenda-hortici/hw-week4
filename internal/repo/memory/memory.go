package memory

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

const (
	statusNew        = "new"
	statusDone       = "done"
	statusInProgress = "in_progress"
)

type repository struct {
	mu   sync.RWMutex
	Task map[uuid.UUID]*Task
}

type Repository interface {
	CreateTask(ctx context.Context, task Task) (uuid.UUID, error)
	GetTask(ctx context.Context, id uuid.UUID) (*Task, error)
	GetAllTasks(ctx context.Context) ([]*Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
	UpdateTask(ctx context.Context, id uuid.UUID, task UpdateTask) error
}

func NewRepo(ctx context.Context) Repository {
	return &repository{
		Task: make(map[uuid.UUID]*Task),
	}
}

func checkStatus(status string) string {
	if status == statusDone || status == statusInProgress {
		return status
	}

	return statusNew
}

func (r *repository) CreateTask(ctx context.Context, task Task) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return uuid.Nil, errors.Wrap(ctx.Err(), "failed to create task")
	default:
		if task.Title == "" {
			return uuid.Nil, errors.Wrap(errors.New("title is required"), "Title is required")
		}

		newTask := &Task{
			ID:          uuid.New(),
			Title:       task.Title,
			Description: task.Description,
			Status:      statusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		r.Task[newTask.ID] = newTask
		return newTask.ID, nil
	}
}

func (r *repository) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "failed to get task")
	default:
		task, ok := r.Task[id]
		if !ok {
			return nil, errors.Wrap(errors.New("task not found"), "Task not found")
		}

		return task, nil
	}
}

func (r *repository) GetAllTasks(ctx context.Context) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*Task, 0, len(r.Task))

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "failed to get tasks")
	default:
		if len(r.Task) == 0 {
			return []*Task{}, errors.Wrap(errors.New("tasks not found"), "Tasks not found")
		}

		for _, task := range r.Task {
			tasks = append(tasks, task)
		}
		return tasks, nil
	}
}

func (r *repository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed to delete task")
	default:
		delete(r.Task, id)
		return nil
	}
}

func (r *repository) UpdateTask(ctx context.Context, id uuid.UUID, task UpdateTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed to update task")
	default:
		if task.Title == "" {
			return errors.Wrap(errors.New("title is required"), "Title is required")
		}

		if task.Description == "" {
			return errors.Wrap(errors.New("description is required"), "Description is required")
		}

		if task.Status == "" {
			return errors.Wrap(errors.New("status is required"), "Status is required")
		}

		status := checkStatus(task.Status)
		task.Status = status
		task.UpdatedAt = time.Now()

		newTask := &Task{
			ID:          id,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   r.Task[id].CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}

		r.Task[id] = newTask
		return nil
	}
}
