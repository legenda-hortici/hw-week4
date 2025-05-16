package service

import (
	"restapi/internal/dto"
	"restapi/internal/repo"
	"restapi/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type service struct {
	log  *zap.SugaredLogger
	repo repo.Repository
}

type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTask(ctx *fiber.Ctx) error
	GetAllTasks(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
}

func NewService(log *zap.SugaredLogger, repo repo.Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	if err := ctx.BodyParser(&req); err != nil {
		s.log.Errorf("Invalid request body: %v", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if err := validator.Validate(ctx.Context(), req); err != nil {
		s.log.Errorf("Invalid request body: %v", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid request body")
	}

	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	id, err := s.repo.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Errorf("Error creating task: %v", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	responce := dto.Response{
		Status: "success",
		Data:   map[string]uuid.UUID{"id": id},
	}

	return ctx.Status(fiber.StatusCreated).JSON(responce)
}

func (s *service) GetTask(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid task id")
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	task, err := s.repo.GetTask(ctx.Context(), id)
	if err != nil {
		s.log.Errorf("Error getting task: %v", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	responce := dto.Response{
		Status: "success",
		Data:   task,
	}

	return ctx.Status(fiber.StatusOK).JSON(responce)
}

func (s *service) GetAllTasks(ctx *fiber.Ctx) error {
	tasks, err := s.repo.GetAllTasks(ctx.Context())
	if err != nil {
		s.log.Errorf("Error getting tasks: %v", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	responce := dto.Response{
		Status: "success",
		Data:   tasks,
	}

	return ctx.Status(fiber.StatusOK).JSON(responce)
}

func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid task id")
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	err = s.repo.DeleteTask(ctx.Context(), id)
	if err != nil {
		s.log.Errorf("Error deleting task: %v", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	responce := dto.Response{
		Status: "success",
	}

	return ctx.Status(fiber.StatusOK).JSON(responce)
}

func (s *service) UpdateTask(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid id")
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	var req UpdateTaskRequest

	if err := ctx.BodyParser(&req); err != nil {
		s.log.Errorf("Invalid request body: %v", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if err := validator.Validate(ctx.Context(), req); err != nil {
		s.log.Errorf("Invalid request body: %v", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid request body")
	}

	task := repo.UpdateTask{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	err = s.repo.UpdateTask(ctx.Context(), id, task)
	if err != nil {
		s.log.Errorf("Error updating task: %v", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	responce := dto.Response{
		Status: "success",
	}

	return ctx.Status(fiber.StatusOK).JSON(responce)
}
