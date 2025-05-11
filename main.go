package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	const PORT = ":3000"
	if err := app.Listen(PORT); err != nil {
		panic(err)
	}
}
