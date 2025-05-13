# P4L Assessment Project

This repository contains the P4L Assessment application built with Go and [Fiber](https://gofiber.io/). It provides a RESTful API for user and product management, with JWT authentication and validation.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

- Go (v1.16 or later)
- PostgreSQL (v12.x or later)
- Git

---

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/p4l-assessment.git
   cd p4l-assesment
   ```

2. **Install dependencies:**
   ```sh
   go mod download
   ```
---

## Configuration

1. **Environment Variables:**

   Copy `.env.example` to `.env` and update the values as needed:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=p4l_db
   APP_PORT=3000
   JWT_SECRET=your_jwt_secret
   ```

---

## Database Setup

1. **Create the PostgreSQL database:**
   ```sh
   psql -U postgres
   CREATE DATABASE p4l_db;
   \q
   ```

2. **Migrations**

   Migration automatically run when the application starts. you can see the migration files in the `database/conn.go` directory.


---

## Running the Application

1. **Start the development server:**
   ```sh
   go run main.go
   ```

   The server will run on `http://localhost:3000` (or the port specified in `.env`).

3. **Running in Development Mode:**
   ```sh
   air
   ```

   The fiber running with powered by air will automatically reload the server on file changes.

2. **Build and run for production:**
   ```sh
   go build -o main.exe main.go
   ./main.exe
   ```

---

## Running Tests

1. **Set up a test database:**
   ```sh
   psql -U postgres
   CREATE DATABASE p4l_test_db;
   \q
   ```

2. **Update your `.env` for testing:**
   ```
   DB_NAME=p4l_test_db
   ```

3. **Run tests:**
   ```sh
   go test -v ./tests
   ```

---

## Project Structure

```
.
├── config/         # Environment and configuration helpers
├── database/       # Database connection logic
├── handler/        # HTTP handlers for admin, product, user
├── middleware/     # Fiber middleware (e.g., JWT auth)
├── models/         # GORM models for User, Product, etc.
├── routes/         # API route definitions
├── tests/          # Integration and helper tests
├── utils/          # JWT, validation, and utility functions
├── main.go         # Application entry point
└── ...
```

---

## API Documentation

### User Endpoints

- `GET /api/users` — Get all users
- `GET /api/users/:id` — Get user by ID
- `POST /api/users` — Create a new user
- `PUT /api/users/:id` — Update a user
- `DELETE /api/users/:id` — Delete a user

### Product Endpoints

- `GET /api/products` — Get all products
- `GET /api/products/:id` — Get product by ID
- `POST /api/products` — Create a new product
- `PUT /api/products/:id` — Update a product
- `DELETE /api/products/:id` — Delete a product

### Admin Endpoints

- `GET /api/admin/all-user` — Get all users (admin only)
- `GET /api/admin/user/:id` — Get user by ID (admin only)
- `POST /api/admin/user` — Create a new user (admin only)
- `PATCH /api/admin/user/:id` — Update a user (admin only)
- `DELETE /api/admin/user/:id` — Delete a user (admin only)



> **Note:** Most endpoints require JWT authentication. Obtain a token via the login endpoint and include it as a cookie named `_token`.

---

## Troubleshooting

### Common Issues

- **Database Connection:**
  Ensure PostgreSQL is running and credentials in `.env` are correct.

- **Dependency Issues:**
  Run `go mod tidy` to resolve missing dependencies.

- **Port Conflicts:**
  Change `APP_PORT` in `.env` if the default port is in use.

### Getting Help

If you encounter issues not covered here, open an issue on this Repository.

---
