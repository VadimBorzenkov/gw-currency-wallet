package app

import (
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/config"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/db"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/delivery/handlers"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/delivery/routes"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/grpc"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/repository"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/services"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/logger"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/migrator"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Run() {
	logger := logger.InitLogger()

	// Загрузка конфигурации
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Could not load config: %v", err)
	}

	// Инициализация подключения к базе данных
	dbase, err := db.Init(config)
	if err != nil {
		logger.Fatalf("Could not initialize DB connection: %v", err)
	}
	defer func() {
		if err := db.Close(dbase); err != nil {
			log.Errorf("Failed to close database: %v", err)
		}
	}()

	if err := migrator.RunDatabaseMigrations(dbase); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Создание объекта для работы с хранилищем
	repo := repository.NewRepository(dbase, logger)

	grpcClient, err := grpc.NewCurrencyClient()
	if err != nil {
		log.Fatalf("Failed to create grpc-client: %v", err)
	}

	tokenManager := utils.NewManager(config)

	service := services.NewService(repo, grpcClient, tokenManager, logger)

	handler := handlers.NewHandler(service, &tokenManager, logger)

	app := fiber.New()

	routes.RegistrationRoutes(app, handler, &tokenManager)

	logger.Infof("Starting server on port %s", config.Port)
	if err := app.Listen(":" + config.Port); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}

}
