package handlers

import (
	"context"
	"time"

	"bolt-backend/config"
	"bolt-backend/database"
	"bolt-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

// =====================
// CREATE USER (SIGNUP)
// =====================
func CreateUser(c *fiber.Ctx) error {
	cfg := config.LoadConfig()
	collection := database.GetCollection(cfg.DatabaseName, "users")

	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if user.FirstName == "" || user.LastName == "" || user.EmailID == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields required"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// =====================
// GET USERS
// =====================
func GetUsers(c *fiber.Ctx) error {
	cfg := config.LoadConfig()
	collection := database.GetCollection(cfg.DatabaseName, "users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode users"})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"count": len(users),
	})
}

// =====================
// LOGIN USER (FIXED JWT)
// =====================
func LoginUser(c *fiber.Ctx) error {
	cfg := config.LoadConfig()
	collection := database.GetCollection(cfg.DatabaseName, "users")

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email_id": req.EmailID}).Decode(&user)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User not found"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.EmailID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// =====================
// REFRESH TOKEN
// =====================
func RefreshToken(c *fiber.Ctx) error {
	cfg := config.LoadConfig()

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString = tokenString[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	newTokenString, _ := newToken.SignedString([]byte(cfg.JWTSecret))

	return c.JSON(fiber.Map{
		"token": newTokenString,
	})
}