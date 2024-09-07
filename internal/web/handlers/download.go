package handlers

import "github.com/gofiber/fiber/v3"

func Download(h *Handler) error {
	if h.drw6.State.IsLoad() {
		return h.ctx.JSON(
			JsonMessage{
				Code:  fiber.StatusBadRequest,
				Error: "already in loading",
			},
		)
	}
	go h.drw6.UpdateMust()
	return h.ctx.JSON(
		JsonMessage{
			Code:    fiber.StatusOK,
			Message: "loading started",
		},
	)
}
