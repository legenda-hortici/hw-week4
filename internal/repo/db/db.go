package db

import (
	"context"
	"database/sql"
	"fmt"
	"restapi/internal/config"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

// Запросы
const (
	insertTaskQuery  = "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id"
	gatTaskQuery     = "SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1"
	getAllTasksQuery = "SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at LIMIT $1 OFFSET $2"
	deleteTaskQuery  = "DELETE FROM tasks WHERE id = $1"
	updateTaskQuery  = "UPDATE tasks SET title = $2, description = $3, status = $4, updated_at = $5 WHERE id = $1"
)

// Таймаут
const (
	timeout = 5 * time.Second
)

type DBrepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	CreateTask(ctx context.Context, task Task) (int64, error)
	GetTask(ctx context.Context, id int64) (*Task, error)
	GetAllTasks(ctx context.Context, limit int, offset int) ([]*Task, error)
	DeleteTask(ctx context.Context, id int64) error
	UpdateTask(ctx context.Context, id int64, task UpdateTask) error
}

// NewRepo создает новый репозиторий
func NewRepo(ctx context.Context, log *zap.SugaredLogger, cfg config.AppConfig) (Repository, error) {
	// Составляем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
		cfg.Database.PoolMaxConns,
		cfg.Database.PoolMaxConnLifetime.String(),
		cfg.Database.PoolMaxConnIdleTime.String(),
	)

	// Парсим строку подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to parse PostgreSQL config"))
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Устанавливаем режим кэширования
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаем подключение
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to create PostgreSQL connection pool"))
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	// Пингуем бд для проверки
	err = pool.Ping(ctx)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to ping database"))
		return nil, errors.Wrap(err, "failed to ping database")
	}

	log.Info("Connected to database")

	return &DBrepository{
		pool: pool,
	}, nil
}

// CreateTask создает новую задачу
func (r *DBrepository) CreateTask(ctx context.Context, task Task) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var id int64
	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description, task.Status).Scan(&id)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to create task"))
		return -1, errors.Wrap(err, "failed to create task")
	}

	return id, nil
}

// GetTask возвращает задачу по id
func (r *DBrepository) GetTask(ctx context.Context, id int64) (*Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var task Task
	if err := r.pool.QueryRow(ctx, gatTaskQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Error(errors.Wrap(err, "task not found"))
			return nil, errors.Wrap(err, "task not found")
		}
		log.Error(errors.Wrap(err, "failed to scan task"))
		return nil, errors.Wrap(err, "failed to get task")
	}

	return &task, nil
}

// GetAllTasks возвращает все задачи
func (r *DBrepository) GetAllTasks(ctx context.Context, limit int, offset int) ([]*Task, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var tasks []*Task
	rows, err := r.pool.Query(ctx, getAllTasksQuery, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error(errors.Wrap(err, "no tasks found"))
			return nil, errors.Wrap(err, "no tasks found")
		}
		log.Error(errors.Wrap(err, "failed to get all tasks"))
		return nil, errors.Wrap(err, "failed to get all tasks")
	}

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			log.Error(errors.Wrap(err, "failed to scan task"))
			return nil, errors.Wrap(err, "failed to scan task")
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// DeleteTask удаляет задачу
func (r *DBrepository) DeleteTask(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.pool.Exec(ctx, deleteTaskQuery, id)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to delete task"))
		return errors.Wrap(err, "failed to delete task")
	}

	return nil
}

// UpdateTask обновляет задачу
func (r *DBrepository) UpdateTask(ctx context.Context, id int64, task UpdateTask) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.pool.Exec(ctx, updateTaskQuery, id, task.Title, task.Description, task.Status, time.Now())
	if err != nil {
		log.Error(errors.Wrap(err, "failed to update task"))
		return errors.Wrap(err, "failed to update task")
	}

	return nil
}
