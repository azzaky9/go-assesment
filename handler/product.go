package handler

import (
	"errors"
	"go-task/database"
	"go-task/models"
	"go-task/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllProducts(c *fiber.Ctx) error {

	isUserAdmin := utils.IsAdmin(c)
	if !isUserAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized access this resource.",
		})
	}

	products := []models.Products{}
	if err := database.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve product data.",
			"data":    make([]interface{}, 0),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Products retrieved.",
		"data":    products,
	})
}

func GetUserProducts(c *fiber.Ctx) error {

	userId, err := utils.GetUserIDFromToken(c)
	log.Printf("Retrieved userId: %v", userId)
	if err != nil {
		log.Printf("Failed to format userId: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	products := []models.Products{}
	if err := database.DB.Where("user_id = ?", userId).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve product data.",
			"data":    make([]interface{}, 0),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Products retrieved.",
		"data":    products,
	})
}

func GetProductById(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid params :id is required",
		})
	}

	result := models.Products{}
	err := database.DB.First(&result, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data you try to search not found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    result,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	type CreateProductInput struct {
		Name     string  `json:"name" validate:"required,min=3,max=25"`
		Quantity int     `json:"quantity" validate:"required,min=1"`
		Price    float64 `json:"price" validate:"required,min=1"`
	}

	userId, err := utils.GetUserIDFromToken(c)
	if err != nil {
		log.Printf("Failed to format userId: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	var input CreateProductInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid body request",
		})
	}

	if errs := utils.ValidationHandler(input); errs != nil {
		return errs
	}

	product := models.Products{
		Name:     input.Name,
		Quantity: uint(input.Quantity),
		Price:    input.Price,
		UserID:   uint(userId),
	}

	if err := database.DB.Create(&product).Error; err != nil {
		log.Printf("Failed to format userId: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": true,
			"message": "Failed to create product.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Product successfully created",
	})

}

func UpdateProduct(c *fiber.Ctx) error {
	type UpdateProductInput struct {
		Name     *string  `json:"name"`
		Quantity *int     `json:"quantity"`
		Price    *float64 `json:"price"`
	}

	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"succcess": false,
			"message":  "Product ID must be required to update",
		})
	}

	userId, err := utils.GetUserIDFromToken(c)
	if err != nil {
		log.Printf("Failed to format userId: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	var input UpdateProductInput
	log.Printf("Incoming Update product input: %v", input)

	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse body request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid body request.",
		})
	}

	productId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Printf("Failed to format userId: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":  false,
			"messaage": "Failed to update product",
		})
	}

	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Quantity != nil {
		updates["quantity"] = *input.Quantity
	}
	if input.Price != nil {
		updates["price"] = *input.Price
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&models.Products{Id: uint(productId), UserID: userId}).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "success",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
	})
}

func DeleteProductById(c *fiber.Ctx) error {
	id := c.Params("id")
	productId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": true,
			"message": "Product ID must be required",
		})
	}

	_, err = utils.GetUserIDFromToken(c)
	if err != nil {
		log.Printf("Failed to format userId: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}

	product := models.Products{
		Id: uint(productId),
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		log.Printf("Failed to delete product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete the product.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Product successfully deleted.",
	})
}
