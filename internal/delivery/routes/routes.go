package routes

import (
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/delivery/handlers"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/middleware"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func RegistrationRoutes(app *fiber.App, h handlers.HandlerInterface, tokenManager utils.TokenManager) {
	// Middleware для CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Группа API
	api := app.Group("/api/v1")

	// Регистрация и авторизация пользователей
	api.Post("/register", h.RegisterUser)
	api.Post("/login", h.LoginUser)

	// Маршруты с авторизацией (используют JWT-токен)
	api.Get("/balance", middleware.AuthMiddleware(tokenManager), h.GetBalance)
	api.Post("/wallet/deposit", middleware.AuthMiddleware(tokenManager), h.Deposit)
	api.Post("/wallet/withdraw", middleware.AuthMiddleware(tokenManager), h.Withdraw)
	api.Get("/exchange/rates", middleware.AuthMiddleware(tokenManager), h.GetExchangeRates)
	api.Post("/exchange", middleware.AuthMiddleware(tokenManager), h.ExchangeCurrency)

	// Включаем Swagger-документацию
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
}
