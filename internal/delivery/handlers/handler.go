package handlers

import (
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/services"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// HandlerInterface определяет интерфейс для обработчиков.
type HandlerInterface interface {
	RegisterUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error

	GetBalance(ctx *fiber.Ctx) error
	Deposit(ctx *fiber.Ctx) error
	Withdraw(ctx *fiber.Ctx) error
	GetExchangeRates(ctx *fiber.Ctx) error
	ExchangeCurrency(ctx *fiber.Ctx) error
}

type handler struct {
	service      services.Service
	tokenManager utils.TokenManager
	logger       *logrus.Logger
}

// NewHandler создаёт новый экземпляр обработчиков.
func NewHandler(service services.Service, tokenManager utils.TokenManager, logger *logrus.Logger) HandlerInterface {
	return &handler{
		service:      service,
		tokenManager: tokenManager,
		logger:       logger,
	}
}
