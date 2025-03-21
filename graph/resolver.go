package graph

import (
	"github.com/taufiqfebriant/tripatrasvc/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// MongoDB collections
	userCollection    *mongo.Collection
	productCollection *mongo.Collection
}

// NewResolver creates a new resolver with all required dependencies
func NewResolver() *Resolver {
	return &Resolver{
		userCollection:    db.GetCollection("users"),
		productCollection: db.GetCollection("products"),
	}
}
