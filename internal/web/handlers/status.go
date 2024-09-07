package handlers

import (
	"github.com/gofiber/fiber/v3"
)

type StatusMessage struct {
	Code    int    `json:"code"`
	Loading bool   `json:"loading"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func Status(h *Handler) error {
	var (
		ErrFunc = func() string {
			if err := h.drw6.State.Error(); err != nil {
				return err.Error()
			}
			return ""
		}
	)

	return h.ctx.JSON(
		StatusMessage{
			Code:    fiber.StatusOK,
			Loading: h.drw6.State.IsLoad(),
			Message: h.drw6.State.Message(),
			Error:   ErrFunc(),
		},
	)
}
