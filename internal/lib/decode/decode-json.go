package decodejson

import (
	"Customer/internal/lib/logger/sl"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func DecodeJSON(c fiber.Ctx, req interface{}, log *slog.Logger) bool {
	if err := c.Bind().Body(req); err != nil {
		slog.Error("failed to decode request body", sl.Err(err))

		err := c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to decode request",
		})
		if err != nil {
			slog.Error("failed to send error response", sl.Err(err))
		}

		return false
	}

	slog.Info("request body decoded", slog.Any("request", req))
	return true
}
