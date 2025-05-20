package dto

import (
	"github.com/gofiber/fiber/v2"
)

// Коды ошибок
const (
	FieldBadFormat     = "FIELD_BADFORMAT"
	FieldIncorrect     = "FIELD_INCORRECT"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	FieldNotFound      = "FIELD_NOT_FOUND"
	InternalError      = "Service is currently unavailable. Please try again later."
)

// Response - структура ответа
type Response struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

// Error - структура ошибки
type Error struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

// BadResponseError - возвращает ошибку
func BadResponseError(ctx *fiber.Ctx, code, desc string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: code,
			Desc: desc,
		},
	})
}

// InternalServerError - возвращает ошибку
func InternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: ServiceUnavailable,
			Desc: InternalError,
		},
	})
}
