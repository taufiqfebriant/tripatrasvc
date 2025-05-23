package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.68

import (
	"context"
	"fmt"
	"time"

	"github.com/taufiqfebriant/tripatrasvc/graph/model"
	"github.com/taufiqfebriant/tripatrasvc/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {
	var user model.User
	err := r.userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %v", err)
	}

	return &model.AuthResponse{
		User:        &user,
		AccessToken: accessToken,
	}, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	// Check if email already exists
	exists, err := r.userCollection.CountDocuments(ctx, bson.M{"email": input.Email})
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	if exists > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user := &model.User{
		ID:       primitive.NewObjectID().Hex(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	_, err = r.userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
	update := bson.M{}
	if input.Name != nil {
		update["name"] = *input.Name
	}
	if input.Email != nil {
		update["email"] = *input.Email
	}
	if input.Password != nil {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %v", err)
		}
		update["password"] = string(hashedPassword)
	}

	var user model.User
	err := r.userCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return &user, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	result, err := r.userCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, fmt.Errorf("error deleting user: %v", err)
	}

	return result.DeletedCount > 0, nil
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	now := time.Now()
	product := &model.Product{
		ID:        primitive.NewObjectID().Hex(),
		Name:      input.Name,
		Price:     input.Price,
		Stock:     int(input.Stock),
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.productCollection.InsertOne(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("error creating product: %v", err)
	}

	return product, nil
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, input model.UpdateProductInput) (*model.Product, error) {
	update := bson.M{"updatedAt": time.Now()}
	if input.Name != nil {
		update["name"] = *input.Name
	}
	if input.Price != nil {
		update["price"] = *input.Price
	}
	if input.Stock != nil {
		update["stock"] = *input.Stock
	}

	var product model.Product
	err := r.productCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&product)

	if err != nil {
		return nil, fmt.Errorf("error updating product: %v", err)
	}

	return &product, nil
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	result, err := r.productCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, fmt.Errorf("error deleting product: %v", err)
	}

	return result.DeletedCount > 0, nil
}

// Stock is the resolver for the stock field.
func (r *productResolver) Stock(ctx context.Context, obj *model.Product) (int32, error) {
	return int32(obj.Stock), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *productResolver) CreatedAt(ctx context.Context, obj *model.Product) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *productResolver) UpdatedAt(ctx context.Context, obj *model.Product) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID := ctx.Value(utils.UserIDKey).(string)
	if userID == "" {
		return nil, fmt.Errorf("user not found")
	}

	var user model.User
	err := r.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	cursor, err := r.userCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding users: %v", err)
	}

	return users, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	return &user, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	cursor, err := r.productCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %v", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("error decoding products: %v", err)
	}

	return products, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	err := r.productCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching product: %v", err)
	}

	return &product, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Product returns ProductResolver implementation.
func (r *Resolver) Product() ProductResolver { return &productResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
