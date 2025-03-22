package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ctx = context.Background()

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	fmt.Println("mongoURI", mongoURI)
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	// Create MongoDB client
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Set global client
	Client = client

	// Get database name from environment
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		log.Fatal("MONGO_DB environment variable is not set")
	}

	// Set global database
	DB = client.Database(dbName)

	fmt.Println("Successfully connected to MongoDB")
}

// Disconnect closes the MongoDB connection
func Disconnect() {
	if Client != nil {
		if err := Client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
