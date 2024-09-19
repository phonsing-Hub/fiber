package router

import (
	"github.com/fiber/src/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/api/register", func(c *fiber.Ctx) error {
		return handlers.RegisterHandler(c, db)
	})
	app.Post("/api/login", func(c *fiber.Ctx) error {
		return handlers.LoginHandler(c, db)
	})

	app.Get("/api/check-auth", func(c *fiber.Ctx) error {
		return handlers.CheckAuth(c)
	})

	app.Get("/api/user", func(c *fiber.Ctx) error {
		return handlers.GetUser(c, db)
	})

	app.Get("/api/checkout", func(c *fiber.Ctx) error {
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
