package handlers

import (
	"drw6/internal/drw6"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	ctx       fiber.Ctx
	drw6      *drw6.Drw6
	allowhost string
}

func New(
	_ctx fiber.Ctx,
	_drw6 *drw6.Drw6,
	_allowhost string,
) *Handler {
	return &Handler{
		ctx:       _ctx,
		drw6:      _drw6,
		allowhost: _allowhost,
	}
}
