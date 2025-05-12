package main

import (
	"go-task/database"
	"go-task/routes"
	"go-task/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

type User struct {
	Name string `validate:"required,min=5,max=20"` // Required field, min 5 char long max 20
	Age  int    `validate:"required"`              // Required field, and client needs to implement our 'teener' tag format which we'll see later
}

func main() {

	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(utils.GlobalErrorHandlerResp{
				Success: true,
				Message: err.Error(),
			})
		},
	})

	app.Use(cors.New())

	database.ConnectDB()

	routes.SetupRoutes(app)

	const PORT = ":3000"
	if err := app.Listen(PORT); err != nil {
		panic(err)
	}
}
