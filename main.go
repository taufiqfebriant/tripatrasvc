package main

import (
	"context"
	"fmt"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/taufiqfebriant/tripatrasvc/db"
	"github.com/taufiqfebriant/tripatrasvc/graph"
	"github.com/taufiqfebriant/tripatrasvc/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

// Defining the Graphql handler
func graphqlHandler() echo.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := graph.Config{Resolvers: graph.NewResolver()}
	c.Directives.Auth = func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
		authHeader := ctx.Value(utils.AuthHeaderKey)
		if authHeader == nil {
			return nil, fmt.Errorf("access denied: no authorization header")
		}

		tokenString := authHeader.(string)
		if tokenString == "" {
			return nil, fmt.Errorf("access denied: no authorization header")
		}

		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			fmt.Println("err", err)
			return nil, fmt.Errorf("access denied: invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["sub"].(string)
			ctx = context.WithValue(ctx, utils.UserIDKey, userID)
			return next(ctx)
		}

		return nil, fmt.Errorf("access denied: invalid token")
	}

	h := handler.New(graph.NewExecutableSchema(c))

	// Server setup:
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c echo.Context) error {
		// Add Authorization header to context
		ctx := context.WithValue(c.Request().Context(), utils.AuthHeaderKey, c.Request().Header.Get("Authorization"))
		c.SetRequest(c.Request().WithContext(ctx))
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// Defining the Playground handler
func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func main() {
	db.Connect()
	defer db.Disconnect()

	e := echo.New()

	e.POST("/graphql", graphqlHandler())
	e.GET("/", playgroundHandler())

	port := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
