package handlers

import (
	"github.com/fiber/src/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	UserID string
	Name   string
	Email  string
}

func GetUser(c *fiber.Ctx, db *gorm.DB) error {
	// ดึง user_id จาก JWT token ที่อยู่ใน cookie
	userID, err := utils.Decoded(c.Cookies("auth"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	// สร้างตัวแปร user สำหรับเก็บข้อมูลผู้ใช้
	var user User
	// ค้นหาข้อมูลผู้ใช้ในฐานข้อมูลโดยใช้ user_id
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		// หากไม่พบข้อมูลผู้ใช้
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	// ส่งข้อมูลผู้ใช้กลับไปเป็น JSON
	return c.Status(200).JSON(user)
}
