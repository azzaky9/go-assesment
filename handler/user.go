package handler

import (
	"go-task/database"
	"go-task/models"
	"go-task/utils"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	type RegisterUserInput struct {
		Username  string `json:"username" validate:"required,min=3,max=32"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8,max=16"`
		FirstName string `json:"firstName" validate:"required,min=3,max=8"`
		LastName  string `json:"lastName" `
		Role      string `json:"role" validate:"required,oneof=admin user"`
	}

	var input RegisterUserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	log.Printf("Incoming input: %+v", input)

	if errs := utils.ValidationHandler(input); errs != nil {
		return errs
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server errror",
		})
	}

	user := models.Users{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		FirstName: input.FirstName,
		LastName:  &input.LastName,
		Role:      (*models.Role)(&input.Role),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User registered.",
	})
}

func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username" validate:"max=32"`
		Email    string `json:"email" validate:"max=32"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid body request",
		})
	}

	log.Printf("Incoming input: %+v", input)

	if errs := utils.ValidationHandler(input); errs != nil {
		return errs
	}

	var user models.Users
	query := database.DB
	if input.Username != "" {
		query = query.Where("username = ?", input.Username)
	} else if input.Email != "" {
		query = query.Where("email = ?", input.Email)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Login must be include one of credential email or username",
		})
	}

	if err := query.First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	credential := utils.JwtCredentialStruct{
		Id:       user.Id,
		Username: user.Username,
		Role:     *user.Role,
	}

	token := utils.CreateJWT(credential)
	if token == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	cookieOpts := fiber.Cookie{
		Name:     "_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(utils.JwtExpire),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	}

	c.Cookie(&cookieOpts)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login success",
		"data":    token,
	})
}
