package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/fiber/docs"
	"github.com/fiber/src/database"          
	"github.com/fiber/src/router"            
)


func main() {
	app := fiber.New();
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://192.168.1.46:5173, http://localhost:8000", 
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(logger.New());
	db, err := database.NewDatabase();
	if err != nil {
		log.Fatal("Failed to connect to database:", err);
	}
	api := app.Group("/api");
    v1 := api.Group("/v1");
	router.SetupRoutes(app, db.DB);
	router.GroupRouterUser(v1, db.DB);
	router.GroupRouterPosts(v1, db.DB);
	router.GroupRouterComments(v1, db.DB);

	app.Get("/swagger/*", swagger.HandlerDefault);
	log.Fatal(app.Listen(":3000"));
}
