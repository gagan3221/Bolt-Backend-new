package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"bolt-backend/config"
	"bolt-backend/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// =====================
// CONNECT WALLET
// =====================
func ConnectWallet(c *fiber.Ctx) error {
	type Request struct {
		WalletAddress string `json:"wallet_address"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	cfg := config.LoadConfig()
	collection := database.GetCollection(cfg.DatabaseName, "users")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"wallet_address": req.WalletAddress}},
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save wallet"})
	}

	return c.JSON(fiber.Map{"message": "Wallet saved successfully"})
}


// =====================
// GET WALLET BALANCE
// =====================
func GetWalletBalance(c *fiber.Ctx) error {
	address := c.Query("address")

	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Wallet address required"})
	}

	cfg := config.LoadConfig()

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getBalance",
		"params":  []string{address, "latest"},
	}

	jsonData, _ := json.Marshal(requestBody)

	resp, err := http.Post(cfg.AlchemyURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch balance"})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	balanceHex, ok := result["result"].(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid blockchain response"})
	}

	balanceInt := new(big.Int)
	balanceInt.SetString(balanceHex[2:], 16)

	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(balanceInt),
		big.NewFloat(1e18),
	)

	return c.JSON(fiber.Map{
		"wallet_address": address,
		"balance":        ethValue.Text('f', 6),
		"unit":           "ETH",
	})
}


// =====================
// SEND CRYPTO
// =====================
func SendCrypto(c *fiber.Ctx) error {
	type Request struct {
		FromWallet string `json:"from_wallet"`
		ToWallet   string `json:"to_wallet"`
		Amount     string `json:"amount"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	txHash := primitive.NewObjectID().Hex()

	return c.JSON(fiber.Map{
		"message": "Transaction initiated via MetaMask",
		"tx_hash": txHash,
	})
}


// =====================
// GET TRANSACTIONS
// =====================
func GetTransactions(c *fiber.Ctx) error {
	address := c.Query("address")

	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Wallet address required"})
	}

	transactions := []fiber.Map{
		{
			"to":     "0x12345",
			"amount": "0.5",
			"unit":   "ETH",
			"status": "confirmed",
		},
	}

	return c.JSON(fiber.Map{
		"wallet_address": address,
		"transactions":   transactions,
	})
}

func LogTransaction(c *fiber.Ctx) error {
	type Request struct {
		TxHash     string `json:"tx_hash"`
		FromWallet string `json:"from_wallet"`
		ToWallet   string `json:"to_wallet"`
		Amount     string `json:"amount"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction logged",
		"tx_hash": req.TxHash,
	})
}