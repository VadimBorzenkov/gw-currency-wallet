package handlers

import (
	"context"
	"time"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestTimeout = 1 * time.Second

// Register регистрирует нового пользователя
// @Summary Register new user
// @Description Создает нового пользователя с предоставленными данными
// @Tags Users
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
		h.logger.Errorf("Invalid input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	id, err := h.service.RegisterUser(&user)
	if err != nil {
		if err.Error() == "username already exists" {
			h.logger.Errorf("Username already exists")
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
// @Summary Authorization user
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

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	user, err := h.service.AuthenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		h.logger.Error("Invalid password")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password"})
	}

	deviceID := ctx.IP()
	accessToken, err := h.service.NewJWT(user.ID, user.Email, deviceID, uuid.New().String())
	if err != nil {
		h.logger.Error("Failed to generate access token.")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to generate access token."})
	}

	refreshToken, err := h.service.NewJWT(user.ID, user.Email, deviceID, uuid.New().String())
	if err != nil {
		h.logger.Error("Failed to generate refresh token.")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to generate refresh token."})
	}

	now := time.Now()
	expire := now.Add(time.Duration(h.service.AccessTTL() * time.Minute))

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		DeviceID:  deviceID,
		CreatedAt: now,
		ExpiresAt: expire,
	}

	err = h.service.CreateRefreshTokenModel(ctxWithTimeout, &refreshTokenModel)
	if err != nil {
		h.logger.Error("Failed to save refresh token.")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to save refresh token."})
	}

	return ctx.JSON(fiber.Map{"token": accessToken})
}
