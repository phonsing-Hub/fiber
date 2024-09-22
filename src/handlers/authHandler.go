package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/fiber/src/models"
	"github.com/fiber/src/utils"
	"fmt"
	"time"
	"strings"
)

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with name, email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  name body string true "User Name"
// @Param  email body string true "User Email"
// @Param  password body string true "User Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func RegisterHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find the number of users to generate UserID
	var count int64
	db.Model(&models.User{}).Count(&count)

	// Generate UserID (SP followed by 4 digits)
	userID := fmt.Sprintf("SP%04d", count+1)

	// Hash password
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create user
	user := models.User{
		UserID:   userID,
		Name:     data.Name,
		Email:    data.Email,
		Password: hashedPassword,
	}

	if err := db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
           strings.Contains(err.Error(), "uni_users_email") {
            return c.Status(400).JSON(fiber.Map{"error_mail": "มีอีเมลนี้อยู่แล้ว"})
        }

        // Handle other errors
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully", "user_id": user.UserID})
}

// LoginHandler godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  email body string true "User Email"
// @Param  password body string true "User Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func LoginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find user by email
	var user models.User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"email": "User not found"})
	}

	// Check password
	if !utils.CheckPasswordHash(data.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"password": "Invalid password"})
	}

	// Generate JWT
	token, err := utils.CreateJWT(user.ID, user.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Set JWT in cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24), // Token expires in 24 hours
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

func CheckAuth(c *fiber.Ctx) error {
    // รับ JWT token จาก cookie
    tokenString := c.Cookies("auth")

    // ตรวจสอบว่า token มีอยู่หรือไม่
    if tokenString == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
    }

    // ตรวจสอบความถูกต้องของ JWT token
    token, err := utils.ValidateToken(tokenString)
    
    // ตรวจสอบว่า token ถูกต้องหรือไม่
    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
    }

    // หาก token ถูกต้อง ส่งกลับข้อความ "OK"
    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}
