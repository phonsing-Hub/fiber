package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fiber/src/handlers"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/register", func(c *fiber.Ctx) error {
		return handlers.RegisterHandler(c, db)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return handlers.LoginHandler(c, db)
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		return handlers.GetUsers(c, db)
	})
}