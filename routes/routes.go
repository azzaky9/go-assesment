package routes

import (
	"go-task/handler"
	"go-task/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// user /me path to show user profile
	userRoutes := api.Group("/user")
	userRoutes.Post("/register", handler.RegisterUser)
	userRoutes.Post("/login", handler.LoginUser)
	userRoutes.Get("/products", middleware.Protected(), handler.GetUserProducts)

	productRoutes := api.Group("/products")
	productRoutes.Get("/", middleware.Protected(), handler.GetAllProducts)
	productRoutes.Get("/:id", handler.GetProductById)
	productRoutes.Post("/", middleware.Protected(), handler.CreateProduct)
	productRoutes.Patch("/:id", middleware.Protected(), handler.UpdateProduct)
	productRoutes.Delete("/:id", middleware.Protected(), handler.DeleteProductById)

	// admin RBAC Routes
}
