package api

import (
	"restapi/internal/api/middleware"
	"restapi/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, token string, log *zap.SugaredLogger) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Явно укажите разрешенные домены
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-SomeID",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api := app.Group("/v1", middleware.Autorization(token))
	{
		// Создание задачи
		api.Post("/tasks", r.Service.CreateTask)

		// Получение задачи
		api.Get("/tasks/:id", r.Service.GetTask)

		// Получение всех задач
		api.Get("/tasks", r.Service.GetAllTasks)

		// Удаление задачи
		api.Delete("/tasks/:id", r.Service.DeleteTask)

		// Обновление задачи
		api.Put("/tasks/:id", r.Service.UpdateTask)
	}

	return app
}
