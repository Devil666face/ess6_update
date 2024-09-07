package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type JsonMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ErrorHandler(c fiber.Ctx, err error) error {
	fmt.Println("error", err)
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.JSON(JsonMessage{
		Code:  code,
		Error: err.Error(),
	})
}
