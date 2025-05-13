package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-task/database"
	"go-task/models"
	"go-task/routes"
	"go-task/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App
var authToken string
var adminAuthToken string

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Failed to load .env file")
	}

	// Set up test database connection
	database.ConnectDB()

	// Initialize app for testing
	app = fiber.New()
	routes.SetupRoutes(app)

	// Run tests
	m.Run()
}

// Helper function to make authenticated requests
func makeAuthenticatedRequest(method, url string, body io.Reader) (*http.Response, error) {
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")

	cookie := &http.Cookie{
		Name:     "_token",
		Value:    authToken,
		Path:     "/",
		MaxAge:   int(utils.JwtExpire),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	req.AddCookie(cookie)

	return app.Test(req)
}

// Test user registration
func TestRegisterUser(t *testing.T) {
	CleanupDatabase(database.DB)
	fmt.Println("Database cleanup successfully.")

	var admin models.Role = "admin"
	var lastName = "dummy"

	userData := models.Users{
		Username:  "john",
		Email:     "janeDoe@example.com",
		Password:  "password12345678",
		FirstName: "new",
		LastName:  &lastName,
		Role:      (*models.Role)(&admin),
	}

	jsonData, _ := json.Marshal(userData)

	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test user login
func TestLoginUser(t *testing.T) {
	loginData := map[string]string{
		"email":    "janeDoe@example.com",
		"password": "password12345678",
	}

	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if data := result["data"].(string); data != "" {
		authToken = data
	}

	assert.Contains(t, result, "data")
}

// Test getting user products
func TestGetUserProducts(t *testing.T) {
	resp, err := makeAuthenticatedRequest(http.MethodGet, "/api/user/products", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test getting all products
func TestGetAllProducts(t *testing.T) {
	resp, err := makeAuthenticatedRequest(http.MethodGet, "/api/products", nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test getting product by ID
func TestGetProductById(t *testing.T) {
	// First create a product
	productData := models.Products{
		Name:     "Dummy Get By Id",
		Quantity: 99,
		Price:    10.2,
	}

	jsonData, _ := json.Marshal(productData)

	// Create a product first
	createResp, err := makeAuthenticatedRequest(http.MethodPost, "/api/products", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	var createResult map[string]interface{}
	createBody, _ := io.ReadAll(createResp.Body)

	json.Unmarshal(createBody, &createResult)

	data, ok := createResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Failed to parse data field from response")
	}

	productID := fmt.Sprintf("%v", data["Id"])

	// Now test getting that product
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/products/"+productID, nil))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test creating a product
func TestCreateProduct(t *testing.T) {
	productData := models.Products{
		Name:     "Test Product",
		Quantity: 10,
		Price:    99.99,
	}

	jsonData, _ := json.Marshal(productData)

	resp, err := makeAuthenticatedRequest(http.MethodPost, "/api/products", bytes.NewReader(jsonData))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// Test updating a product
func TestUpdateProduct(t *testing.T) {
	// First create a product
	productData := models.Products{
		Name:     "Test Product Dummy",
		Quantity: 10,
		Price:    99.99,
	}

	jsonData, _ := json.Marshal(productData)

	createResp, err := makeAuthenticatedRequest(http.MethodPost, "/api/products", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	// fmt.Println("Create response:", createResp.StatusCode)

	var createResult map[string]interface{}
	createBody, _ := io.ReadAll(createResp.Body)

	json.Unmarshal(createBody, &createResult)

	data, ok := createResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Failed to parse data field from response")
	}

	productID := fmt.Sprintf("%v", data["Id"])

	// Now update the product
	updateData := models.Products{
		Name:     "Test Product Dummy Edit",
		Quantity: 100,
		Price:    100,
	}

	jsonUpdateData, _ := json.Marshal(updateData)

	resp, err := makeAuthenticatedRequest(http.MethodPatch, "/api/products/"+productID, bytes.NewReader(jsonUpdateData))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test deleting a product
func TestDeleteProductById(t *testing.T) {
	// First create a product
	productData := models.Products{
		Name:     "Dummy Delete",
		Quantity: 99,
		Price:    10.2,
	}

	jsonData, _ := json.Marshal(productData)

	createResp, err := makeAuthenticatedRequest(http.MethodPost, "/api/products", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	var createResult map[string]interface{}
	createBody, _ := io.ReadAll(createResp.Body)

	json.Unmarshal(createBody, &createResult)

	data, ok := createResult["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Failed to parse data field from response")
	}

	productID := fmt.Sprintf("%v", data["Id"])

	// Now delete the product
	resp, err := makeAuthenticatedRequest(http.MethodDelete, "/api/products/"+productID, nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify it's deleted by trying to get it
	getResp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/api/products/"+productID, nil))
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestSetupTestUsers(t *testing.T) {
	// Create admin user
	adminRole := models.Role("admin")
	adminUser := models.Users{
		Username:  "adminuser",
		Email:     "admin@example.com",
		Password:  "password12345678",
		FirstName: "Admin",
		Role:      &adminRole,
	}

	jsonData, _ := json.Marshal(adminUser)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var lastName = "dummy"

	// Create regular user
	userRole := models.Role("user")
	regularUser := models.Users{
		Username:  "regularuser",
		Email:     "user@example.com",
		Password:  "password12345678",
		FirstName: "Regular",
		LastName:  &lastName,
		Role:      &userRole,
	}

	jsonData, _ = json.Marshal(regularUser)
	req = httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Login as admin
	loginData := map[string]string{
		"email":    "admin@example.com",
		"password": "password12345678",
	}
	jsonData, _ = json.Marshal(loginData)
	req = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	data := result["data"].(string)
	if data != "" {
		adminAuthToken = data
	}

	// Login as regular user
	loginData = map[string]string{
		"email":    "user@example.com",
		"password": "password12345678",
	}
	jsonData, _ = json.Marshal(loginData)
	req = httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	if data := result["data"].(string); data != "" {
		authToken = data
	}
}

// Test access to admin routes without authentication
func TestAdminRoutesNoAuth(t *testing.T) {
	// Try each route without auth
	routes := []string{
		"/api/admin/all-user",
		"/api/admin/user/1",
		"/api/admin/user",
	}
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, route := range routes {
		for _, method := range methods {
			if (route == "/api/admin/all-user" || route == "/api/admin/user/1") && method != http.MethodGet && method != http.MethodDelete && method != http.MethodPatch {
				continue
			}
			if route == "/api/admin/user" && method != http.MethodPost {
				continue
			}

			req := httptest.NewRequest(method, route, nil)
			resp, _ := app.Test(req)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, fmt.Sprintf("Route %s with method %s should require auth", route, method))
		}
	}
}
