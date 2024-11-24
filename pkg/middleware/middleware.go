package middleware

import (
	"strings"

	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	claimsKey           = "claims"
)

func AuthMiddleware(tokenManager utils.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(authorizationHeader)
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing Authorization header")
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if token == authHeader {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token format")
		}

		claims, err := tokenManager.ParseJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token: " + err.Error()})
		}

		c.Locals(claimsKey, claims)

		return c.Next()
	}
}
