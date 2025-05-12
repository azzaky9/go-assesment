package middleware

import (
	"go-task/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Protected() func(*fiber.Ctx) error {
	secretKey := config.GetEnv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "s3cret"
	}
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secretKey)},
		ErrorHandler: jwtError,
		TokenLookup:  "cookie:_token",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"success": false, "message": "Missing or malformed JWT"})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"success": false, "message": "Invalid or expired JWT"})
	}
}
