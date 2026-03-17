package logger

import (
	"Customer/internal/http-server/middleware"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func New(log *slog.Logger) fiber.Handler {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	return func(c fiber.Ctx) error {
		entry := log.With(
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("remote_addr", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
		)

		c.Locals("logger", entry)

		start := time.Now()

		entry.Info("request completed",
			slog.Int("status", c.Response().StatusCode()),
			slog.Int("bytes", len(c.Response().Body())),
			slog.Duration("duration", time.Since(start)),
		)

		return c.Next()
	}
}

func LogQuery(log *slog.Logger, op string) fiber.Handler {
	return func(c fiber.Ctx) error {
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(c.Context())),
		)

		c.Locals("logger", log)
		return c.Next()
	}
}

func FromCtx(c fiber.Ctx) *slog.Logger {
	if v := c.Locals("logger"); v != nil {
		if log, ok := v.(*slog.Logger); ok {
			return log
		}
	}
	return slog.Default()
}
