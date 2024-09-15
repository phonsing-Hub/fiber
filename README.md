# fiber

MyApp is a web application built using Go, Fiber, and GORM for the backend, with Swagger UI for API documentation. It provides user registration and login functionality with JWT authentication.

## Features

- **User Registration**: Create a new user with a unique user ID, hashed password, and other details.
- **User Login**: Authenticate users and provide a JWT token for secure access.
- **API Documentation**: Interactive API documentation using Swagger UI.

## Prerequisites

- [Go](https://golang.org/dl/) 1.18 or later
- [PostgreSQL](https://www.postgresql.org/download/) for the database

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/phonsing-Hub/fiber.git
    cd fiber
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Install Swagger CLI:**

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

4. **Set up the environment:**

    Create a `.env` file in the root directory and add your database credentials:

    ```dotenv
    DATABASE_URL=postgres://username:password@localhost:5432/myDB?sslmode=disable
    SECRET_KEY=your_secret_key
    ```

5. **Generate Swagger documentation:**

    ```bash
    swag init
    ```

6. **Run the application:**

    ```bash
    go run main.go
    ```

    The server will start on `http://localhost:3000`.

## Endpoints

- **POST /register**: Register a new user. Requires `name`, `email`, and `password` in the request body.
  
- **POST /login**: Log in an existing user. Requires `email` and `password` in the request body.

- **GET /swagger/***: Access the Swagger UI for API documentation.

## Example

### Register a new user

**Request:**

```bash
curl -X POST http://localhost:3000/register \
-H "Content-Type: application/json" \
-d '{"name": "John Doe", "email": "john@example.com", "password": "yourpassword"}'
 ```
## 📜 License

This software is licensed under the [![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
 © [NHN Cloud](https://github.com/nhn).

