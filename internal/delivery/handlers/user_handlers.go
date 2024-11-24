package handlers

import (
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
	"github.com/gofiber/fiber/v2"
)

// Register регистрирует нового пользователя
// @Summary Register new user
// @Description Создает нового пользователя с предоставленными данными
// @Tags auth
// @Accept json
// @Produce json
// @Param register body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/v1/register [post]
func (h *handler) RegisterUser(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	id, err := h.service.RegisterUser(&user)
	if err != nil {
		if err.Error() == "username already exists" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username already exists"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user_id": id,
	})
}

// LoginUser авторизует пользователя
// @Summary Авторизация пользователя
// @Description Авторизация пользователя с возвратом JWT-токена для дальнейших запросов
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.ErrorResponse "Invalid username or password"
// @Router /api/v1/login [post]
func (h *handler) LoginUser(ctx *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&credentials); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.service.AuthenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	token, err := h.tokenManager.NewJWT(uint64(user.ID), user.Email, ctx.IP(), "tokenID")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.JSON(fiber.Map{"token": token})
}
