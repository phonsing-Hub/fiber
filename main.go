package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/fiber/docs"
	"github.com/fiber/src/database"          // แก้เป็น path ที่ถูกต้องสำหรับ database
	"github.com/fiber/src/router"             // แก้เป็น path ที่ถูกต้องสำหรับ router
)





func main() {
	// Initializing Fiber app
	app := fiber.New()

	// Setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://192.168.1.46:5173, http://localhost:8000", // ลบ '/' ที่ไม่จำเป็นออก
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Swagger documentation route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setting up routes
	router.SetupRoutes(app, db.DB)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
