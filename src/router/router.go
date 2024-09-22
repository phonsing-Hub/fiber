package router

import (
	"github.com/fiber/src/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/api/v1/register", func(c *fiber.Ctx) error {
		return handlers.RegisterHandler(c, db)
	})
	app.Post("/api/v1/login", func(c *fiber.Ctx) error {
		return handlers.LoginHandler(c, db)
	})

	app.Get("/api/v1/check-auth", func(c *fiber.Ctx) error {
		return handlers.CheckAuth(c)
	})

	app.Get("/api/v1/user", func(c *fiber.Ctx) error {
		return handlers.GetUser(c, db)
	})

	app.Get("/api/v1/checkout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "auth",
			MaxAge:   -1, // Token expires in 24 hours
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Strict",
		})
		return c.Status(200).SendString("logout successful!")
	})
}
