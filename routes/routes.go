package routes

import (
	"bolt-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// API routes
	api := app.Group("/api")

	// User routes
	users := api.Group("/users")
	users.Post("/", handlers.CreateUser)
	users.Get("/", handlers.GetUsers)
	users.Post("/login", handlers.LoginUser)
	users.Post("/refresh", handlers.RefreshToken)
	// WALLET ROUTES
	wallet := api.Group("/wallet")
	wallet.Post("/connect", handlers.ConnectWallet)
	wallet.Get("/balance", handlers.GetWalletBalance)
	wallet.Post("/send", handlers.SendCrypto)
	wallet.Get("/transaction", handlers.GetTransactions)
}

