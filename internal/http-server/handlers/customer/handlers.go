package handlers

import (
	"errors"
	"log/slog"

	"github.com/voznikaetnepriyazn/Customer-service/internal/http-server/middleware/logger"
	uuidparam "github.com/voznikaetnepriyazn/Customer-service/internal/http-server/middleware/uuid"
	response "github.com/voznikaetnepriyazn/Customer-service/internal/lib/api/response"
	decodejson "github.com/voznikaetnepriyazn/Customer-service/internal/lib/decode"
	"github.com/voznikaetnepriyazn/Customer-service/internal/lib/logger/sl"
	valid "github.com/voznikaetnepriyazn/Customer-service/internal/lib/validate"
	"github.com/voznikaetnepriyazn/Customer-service/internal/models/customer"
	"github.com/voznikaetnepriyazn/Customer-service/internal/storage"

	"github.com/gofiber/fiber/v3"
)

type Response struct {
	URL string `json:"url" validate:"required, url"`
}

type RequestFullStruct struct {
	Customer customer.Customer
}

type Request struct {
	response.Response
	URL string `json:"url" validate:"required, url"`
}

type Crud interface {
	NewAdd(log *slog.Logger, adder storage.CustomerService) fiber.Handler
	NewDelete(log *slog.Logger, deleter storage.CustomerService) fiber.Handler
	NewGetAll(log *slog.Logger, get storage.CustomerService) fiber.Handler
	NewGetById(log *slog.Logger, get storage.CustomerService) fiber.Handler
	NewUpdate(log *slog.Logger, update storage.CustomerService) fiber.Handler
	NewIsOrderCreated(log *slog.Logger, ord storage.CustomerService) fiber.Handler
}

func NewAdd(log *slog.Logger, adder storage.CustomerService) fiber.Handler {
	return func(c fiber.Ctx) error {
		const op = "handlers.url.add.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req RequestFullStruct

		if !decodejson.DecodeJSON(c, &req, log) {
			return nil
		}

		if !valid.Validate(c, &req, log) {
			return nil
		}

		//проверка на уже существующее значение
		id, err := adder.AddURL(req.Customer)
		if errors.Is(err, storage.ErrUrlExist) {
			log.Info("url already exists", slog.Any("url", req.Customer))

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "url already exists",
			})

			return nil
		}

		//прочие ошибки
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to add url",
			})

			return nil
		}

		log.Info("url added", slog.Any("id", id))

		responseOK(c)

		return nil
	}
}

func responseOK(c fiber.Ctx) error {
	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "OK",
	})

	return nil
}

func NewDelete(log *slog.Logger, deleter storage.CustomerService) fiber.Handler {
	return func(c fiber.Ctx) error {
		const op = "handlers.url.delete.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(c, "id")
		if ok {
			return nil
		}

		err := deleter.DeleteURL(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", uuidParam)

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "not found",
			})

			return nil
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})

			return nil
		}

		log.Info("deleted url", slog.Any("deleted", uuidParam))

		responseOK(c)

		return nil
	}
}

func NewGetAll(log *slog.Logger, get storage.CustomerService) fiber.Handler {
	return func(c fiber.Ctx) error {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		resURL, err := get.GetAllURL()
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("urls not found")

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "not found",
			})

			return nil
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})

			return nil
		}

		log.Info("got urls", slog.Any("urls", resURL))

		c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"urls": resURL,
		})

		return nil
	}
}

func NewGetById(log *slog.Logger, get storage.CustomerService) fiber.Handler {
	return func(c fiber.Ctx) error {
		const op = "handlers.url.getById.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		uuidParam, ok := uuidparam.UUIDFromCtx(c, "id")
		if ok {
			return nil
		}
		resURL, err := get.GetByIdURL(uuidParam)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", slog.String("id", uuidParam.String()))

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "not found",
			})

			return nil
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})

			return nil
		}

		log.Info("got url", slog.Any("url", resURL))

		c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"url": resURL,
		})

		return nil
	}
}

func NewUpdate(log *slog.Logger, update storage.CustomerService) fiber.Handler {
	return func(c fiber.Ctx) error {
		const op = "handlers.url.update.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req RequestFullStruct

		if !decodejson.DecodeJSON(c, &req, log) {
			return nil
		}

		if !valid.Validate(c, &req, log) {
			return nil
		}

		err := update.UpdateURL(req.Customer)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "id", req)

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "not found",
			})

			return nil
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})

			return nil
		}

		log.Info("updated url", slog.Any("url", req))

		responseOK(c)

		return nil
	}
}

func NewIscustomerCreated(log *slog.Logger, ord storage.CustomerService) fiber.Handler {
	return func(c fiber.Ctx) error {
		const op = "handlers.url.IsOrderCreated.New"

		log := logger.FromCtx(c)

		log.Info("handling request")

		var req Request

		if !decodejson.DecodeJSON(c, &req, log) {
			return nil
		}

		if !valid.Validate(c, &req, log) {
			return nil
		}

		uuidParam, ok := uuidparam.UUIDFromCtx(c, "id")
		if ok {
			return nil
		}

		resId, err := ord.IsCustomerCreatedURL(uuidParam)
		if resId == false {
			log.Info("url not found", "id", resId)

			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "not found",
			})

			return nil
		}

		if err != nil {
			log.Error("failed to check url", sl.Err(err))

			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})

			return nil
		}

		log.Info("url exist", slog.Any("url", uuidParam))

		responseOK(c)

		return nil
	}
}
