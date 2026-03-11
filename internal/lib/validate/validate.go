package valid

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func Validate(c fiber.Ctx, req interface{}, log *slog.Logger) bool {
	if err := validate.Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		if !errors.As(err, &validateErr) {
			log.Error("unknown validation error", slog.Any("error", err))
			c.Status(500).JSON(fiber.Map{
				"error": "internal error"})
			return false
		}

		log.Warn("validation failed", slog.Any("errors", FormatValidationError(validateErr)))

		c.Status(400).JSON(fiber.Map{
			"error":   "validation failed",
			"details": FormatValidationError(validateErr),
		})
		return false
	}
	return true
}

func FormatValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Error()
	}
	return errors
}
