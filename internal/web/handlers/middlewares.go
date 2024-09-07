package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

const (
	Host = "Host"
)

func AllowHost(h *Handler) error {
	if host, ok := h.ctx.GetReqHeaders()[Host]; ok {
		if strings.Contains(host[0], h.allowhost) {
			return h.ctx.Next()
		}
	}
	return fiber.ErrBadRequest
}
