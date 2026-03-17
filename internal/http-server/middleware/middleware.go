package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v3"
)

// генерация ид запроса
func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		id := fmt.Sprintf("%d", time.Now().UnixNano())

		ctx := context.WithValue(c.Context(), "request_id", id)

		c.Locals("requestID", ctx)

		c.Set("X-Request-ID", id)

		return c.Next()
	}
}

// перехват паники - обертка запроса в дефер и рековер
func Recoverer(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				stackBuf := make([]byte, 4096)
				stackSize := runtime.Stack(stackBuf, false)
				stackTrace := string(stackBuf[:stackSize])

				log.Error("panic recovered",
					slog.Any("error", err),
					slog.String("stack_trace", stackTrace),
					slog.String("method", c.Method()),
					slog.String("path", c.Path()),
				)

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "internal server error",
				})
			}
		}()

		return c.Next()
	}
}

// получение ид из контекста
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}

	return ""
}
