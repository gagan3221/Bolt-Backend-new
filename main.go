package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"bolt-backend/config"
	"bolt-backend/database"
	"bolt-backend/routes"



	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title Bolt Backend API
// @version 1.0
// @description REST API for Bolt Backend with MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@boltbackend.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @schemes http https

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	if err := database.ConnectDB(cfg.MongoURI); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ensure graceful shutdown
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Bolt Backend API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Bolt Backend API",
			"status":  "running",
			"swagger": "http://localhost:" + cfg.Port + "/swagger/index.html",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Setup API routes
	routes.SetupRoutes(app)

	// Swagger documentation route (must be after API routes)
	app.Get("/swagger/*", swagger.WrapHandler)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Printf("📚 Swagger documentation available at http://localhost:%s/swagger/index.html", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
