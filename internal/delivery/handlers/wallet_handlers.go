package handlers

import (
	"errors"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func extractUserIDFromToken(c *fiber.Ctx) (uint64, error) {
	// Извлекаем данные из локального контекста
	claims, ok := c.Locals("claims").(*utils.Claims)
	if !ok || claims.UserID == 0 {
		return 0, errors.New("invalid token claims")
	}
	return claims.UserID, nil
}

// GetBalance возвращает баланс пользователя.
// @Summary Get user balance
// @Description Получает текущий баланс пользователя.
// @Tags wallet
// @Accept json
// @Produce json
// @Success 200 {object} models.BalanceResponse
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /wallet/balance [get]
func (h *handler) GetBalance(ctx *fiber.Ctx) error {
	// Извлекаем ID пользователя из токена
	userID, err := extractUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Получаем баланс пользователя
	balance, err := h.service.GetBalance(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get balance",
		})
	}

	// Возвращаем ответ с балансом
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"balance": balance,
	})
}

// Deposit пополняет баланс пользователя.
// @Summary Deposit funds to user balance
// @Description Пополняет баланс пользователя на указанную сумму.
// @Tags wallet
// @Accept json
// @Produce json
// @Param deposit body models.DepositRequest true "Deposit request"
// @Success 200 {object} models.DepositResponse
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /wallet/deposit [post]
func (h *handler) Deposit(ctx *fiber.Ctx) error {
	var deposit models.DepositRequest
	if err := ctx.BodyParser(&deposit); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Извлекаем ID пользователя из токена
	userID, err := extractUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Пополнение счета
	newBalance, err := h.service.Deposit(userID, deposit.Amount, deposit.Currency)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid amount or currency",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Account topped up successfully",
		"new_balance": newBalance,
	})
}

// Withdraw выводит средства со счета.
// @Summary Withdraw funds from user balance
// @Description Выводит указанную сумму со счета пользователя.
// @Tags wallet
// @Accept json
// @Produce json
// @Param withdraw body models.WithdrawRequest true "Withdraw request"
// @Success 200 {object} models.WithdrawResponse
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /wallet/withdraw [post]
func (h *handler) Withdraw(ctx *fiber.Ctx) error {
	var withdraw models.WithdrawRequest
	if err := ctx.BodyParser(&withdraw); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Извлекаем ID пользователя из токена
	userID, err := extractUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Вывод средств
	newBalance, err := h.service.Withdraw(userID, withdraw.Amount, withdraw.Currency)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Insufficient funds or invalid amount",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Withdrawal successful",
		"new_balance": newBalance,
	})
}

// GetExchangeRates получает актуальные курсы валют.
// @Summary Get exchange rates
// @Description Получает актуальные курсы валют для различных валютных пар.
// @Tags exchange
// @Accept json
// @Produce json
// @Success 200 {object} models.RatesResponse
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /wallet/rates [get]
func (h *handler) GetExchangeRates(ctx *fiber.Ctx) error {
	rates, err := h.service.GetAllRates()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve exchange rates",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"rates": rates,
		"USD":   rates["USD"],
		"RUB":   rates["RUB"],
		"EUR":   rates["EUR"],
	})
}

// ExchangeCurrency - обработка обмена валют.
// @Summary Exchange currency
// @Description Обменивает одну валюту на другую по актуальному курсу.
// @Tags exchange
// @Accept json
// @Produce json
// @Param exchange body models.ExchangeRequest true "Exchange request"
// @Success 200 {object} models.ExchangeResponse
// @Failure 400 {object} models.ErrorResponse "Insufficient funds or invalid currencies"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /wallet/exchange [post]
func (h *handler) ExchangeCurrency(ctx *fiber.Ctx) error {
	var exchangeRequest models.ExchangeRequest

	// Парсим тело запроса
	if err := ctx.BodyParser(&exchangeRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Извлекаем ID пользователя из токена
	userID, err := extractUserIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Проверка наличия достаточного баланса
	userBalance, err := h.service.GetBalance(userID)
	if err != nil {
		h.logger.Errorf("Failed to retrieve user balance for user %d: %v", userID, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user balance",
		})
	}

	fromCurrencyBalance, ok := userBalance[exchangeRequest.FromCurrency]
	if !ok || fromCurrencyBalance < exchangeRequest.Amount {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient funds or invalid currencies",
		})
	}

	// Конвертация валюты
	exchangedAmount, err := h.service.ExchangeCurrency(exchangeRequest.FromCurrency, exchangeRequest.ToCurrency, exchangeRequest.Amount)
	if err != nil {
		h.logger.Errorf("Failed to convert currency: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Currency conversion failed",
		})
	}

	// Обновление баланса пользователя
	newBalance, err := h.service.UpdateUserBalance(userID, exchangeRequest.FromCurrency, exchangeRequest.ToCurrency, exchangeRequest.Amount, exchangedAmount)
	if err != nil {
		h.logger.Errorf("Failed to update user balance: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user balance",
		})
	}

	// Возвращаем успешный ответ
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Exchange successful",
		"exchanged_amount": exchangedAmount,
		"new_balance":      newBalance,
	})
}
