package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/fiber/src/models"
)


func GetUsers(c *fiber.Ctx, db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return c.Status(500).SendString("Error fetching users")
	}
	return c.JSON(users)
}
