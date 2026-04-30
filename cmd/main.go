package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/voznikaetnepriyazn/Customer-service/internal/config"
	handlers "github.com/voznikaetnepriyazn/Customer-service/internal/http-server/handlers/customer"
	"github.com/voznikaetnepriyazn/Customer-service/internal/lib/logger/sl"
	"github.com/voznikaetnepriyazn/Customer-service/internal/storage"
	"github.com/voznikaetnepriyazn/Customer-service/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error(".env file not found", sl.Err(err))
	}

	cfg := config.MustLoad()

	logger := setUpLogger(cfg.Env)

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		slog.SetDefault(logger)
	} else {
		slog.SetDefault(logger)
	}

	logger.Info("starting customer service", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	db, err := postgresql.New(cfg.DB.DSN())
	if err != nil {
		slog.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	customerService := storage.CustomerService(db)

	app := fiber.New(fiber.Config{
		AppName:       "Customer Service v1.0",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})

	registerRouter(app, logger, customerService)

	addr := cfg.HTTPServer.Address
	if addr == "" {
		addr = ":8081"
	}

	logger.Info("starting server", slog.String("address", addr))
	if err := app.Listen(addr); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}

func registerRouter(app *fiber.App, log *slog.Logger, service storage.CustomerService) {
	api := app.Group("/api/v1")

	api.Post("/customer", handlers.NewAdd(log, service))
	api.Get("/customer/:id", handlers.NewGetById(log, service))
	api.Get("/customer", handlers.NewGetAll(log, service))
	api.Put("/customer/:id/fullName", handlers.NewUpdate(log, service))
	api.Delete("customer/:id", handlers.NewDelete(log, service))
	api.Get("/customer/:id", handlers.NewIscustomerCreated(log, service))
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
