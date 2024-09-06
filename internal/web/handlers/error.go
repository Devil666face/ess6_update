package handlers

import (
	"github.com/gofiber/fiber/v3"
)

type Error struct {
	code    int
	message string
}

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.JSON(Error{
		code:    code,
		message: err.Error(),
	})
}
