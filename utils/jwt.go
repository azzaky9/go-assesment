package utils

import (
	"fmt"
	"go-task/config"
	"go-task/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCredentialStruct struct {
	Id       uint
	Username string
	Role     models.Role
}

var JwtExpire int64

func CreateJWT(identity JwtCredentialStruct) string {
	token := jwt.New(jwt.SigningMethodHS256)

	JwtExpire = time.Now().Add(time.Hour * 24).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = identity.Id
	claims["username"] = identity.Username
	claims["admin"] = identity.Role == models.Admin
	claims["exp"] = JwtExpire

	secretKey := config.GetEnv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "secret"
	}

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Error signing JWT token: %v", err)
		return ""
	}

	return t
}

func GetUserIDFromToken(c *fiber.Ctx) (uint, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userIdStr := fmt.Sprintf("%v", claims["id"])
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format: %w", err)
	}

	return uint(userId), nil
}

func IsAdmin(c *fiber.Ctx) bool {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	isAdmin := claims["admin"] == true
	return isAdmin
}
