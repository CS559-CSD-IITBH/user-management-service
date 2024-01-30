package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/rs/zerolog"

	"github.com/CS559-CSD-IITBH/user-management-service/routes"
	"github.com/gorilla/sessions"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Msg("Unable to load the env file")
	}

	// Replace this with your MongoDB Atlas connection string
	connectionString := os.Getenv("MONGO_URL")

	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Internal server error: Unable to connect to Mongo")
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Internal server error: Unable to talk to Mongo")
	}

	log.Println("Connected to Mongo!")

	// You can now use the "client" variable to interact with your MongoDB database.
	// For example, you can access a collection:
	customers := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_CUSTOMERS"))
	merchants := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_MERCHANTS"))

	// Session store
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Set max age for cookie
	store.Options = &sessions.Options{
		Path:     "/api/v1",
		MaxAge:   86400 * 7,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	r := routes.SetupRouter(customers, merchants, store, logger)
	logger.Info().Msg("Setup Complete. Starting user-service...")
	r.Run(":" + os.Getenv("PORT"))
}
