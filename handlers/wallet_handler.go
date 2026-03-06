package handlers

import (
	"github.com/gofiber/fiber/v2"
)


// ConnectWallet godoc
// @Summary Connect wallet
// @Description Save user's wallet address
// @Tags Wallet
// @Accept json
// @Produce json
// @Param wallet body map[string]string true "Wallet Address"
// @Success 200 {object} map[string]string
// @Router /api/wallet/connect [post]
func ConnectWallet(c *fiber.Ctx) error {

	type Request struct {
		WalletAddress string `json:"wallet_address"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	return c.JSON(fiber.Map{
		"message":        "Wallet connected successfully",
		"wallet_address": req.WalletAddress,
	})
}



// GetWalletBalance godoc
// @Summary Get wallet balance
// @Description Fetch wallet balance
// @Tags Wallet
// @Produce json
// @Param address query string true "Wallet Address"
// @Success 200 {object} map[string]string
// @Router /api/wallet/balance [get]
func GetWalletBalance(c *fiber.Ctx) error {

	address := c.Query("address")

	if address == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Wallet address required",
		})
	}

	// dummy balance for demo
	return c.JSON(fiber.Map{
		"wallet_address": address,
		"balance":        "2.35 ETH",
	})
}



// SendCrypto godoc
// @Summary Send crypto
// @Description Send crypto to another wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Param transaction body map[string]string true "Transaction Data"
// @Success 200 {object} map[string]string
// @Router /api/wallet/send [post]
func SendCrypto(c *fiber.Ctx) error {

	type Request struct {
		FromWallet string `json:"from_wallet"`
		ToWallet   string `json:"to_wallet"`
		Amount     string `json:"amount"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction initiated",
		"tx_hash": "0xabc123456789",
	})
}



// GetTransactions godoc
// @Summary Get wallet transactions
// @Description Fetch wallet transaction history
// @Tags Wallet
// @Produce json
// @Param address query string true "Wallet Address"
// @Success 200 {object} map[string]interface{}
// @Router /api/wallet/transaction [get]
func GetTransactions(c *fiber.Ctx) error {

	address := c.Query("address")

	if address == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Wallet address required",
		})
	}

	transactions := []fiber.Map{
		{
			"to":     "0x12345",
			"amount": "0.5 ETH",
			"status": "confirmed",
		},
		{
			"to":     "0x98765",
			"amount": "1 ETH",
			"status": "pending",
		},
	}

	return c.JSON(fiber.Map{
		"wallet_address": address,
		"transactions":   transactions,
	})
}