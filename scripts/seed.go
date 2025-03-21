package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/taufiqfebriant/tripatrasvc/db"
	"github.com/taufiqfebriant/tripatrasvc/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}

	// Load .env file from project root
	envPath := filepath.Join(currentDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func main() {
	// Connect to database
	db.Connect()
	defer db.Disconnect()

	// Get users collection
	users := db.GetCollection("users")

	// Clear existing users
	_, err := users.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Error clearing users:", err)
	}

	// Create dummy users
	dummyUsers := []model.User{
		{
			ID:       primitive.NewObjectID().Hex(),
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: hashPassword("password123"),
		},
		{
			ID:       primitive.NewObjectID().Hex(),
			Name:     "Jane Smith",
			Email:    "jane@example.com",
			Password: hashPassword("password123"),
		},
	}

	// Insert users
	for _, user := range dummyUsers {
		_, err := users.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal("Error inserting user:", err)
		}
		fmt.Printf("Created user: %s (%s)\n", user.Name, user.Email)
	}

	fmt.Println("\nDummy users created successfully!")
	fmt.Println("\nYou can now login with:")
	fmt.Println("Email: john@example.com")
	fmt.Println("Password: password123")
	fmt.Println("\nor")
	fmt.Println("\nEmail: jane@example.com")
	fmt.Println("Password: password123")
}

func hashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}
	return string(hashedBytes)
}
