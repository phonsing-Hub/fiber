package main

import (
	"log"
	"github.com/fiber/src/database"
	"github.com/fiber/src/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/fiber/docs"
)

func main() {
	// Initializing Fiber app
	app := fiber.New()

	// Connecting to the database
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	app.Get("/swagger/*", swagger.HandlerDefault) 
	// Setting up routes
	router.SetupRoutes(app, db.DB)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
