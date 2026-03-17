package uuidparam

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func UUIDParam(paramName string, log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		alias := c.Params("paramName")
		if alias == "" {
			slog.Info("id is empty")

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request",
			})
			return nil
		}

		id, err := uuid.Parse(alias)
		if err != nil {
			slog.Info("invalid uuid format", slog.String("id", alias), slog.Any("error", err))
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid id format",
			})
			return nil
		}

		c.Locals("uuid_"+paramName, id)
		return c.Next()
	}
}

func UUIDFromCtx(c fiber.Ctx, paramName string) (uuid.UUID, bool) {
	if v := c.Locals("uuid_" + paramName); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return id, true
		}
	}
	return uuid.Nil, false
}
