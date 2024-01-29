package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/routes"
	"github.com/gorilla/sessions"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Msg("Unable to load the env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		logger.Fatal().Msg("Unable to connect to Postgres")
	}

	err = db.AutoMigrate(
		&models.Customer{},
		&models.Merchant{},
		&models.DeliveryAgent{},
		&models.User{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		logger.Fatal().Msg("Unable to migrate models to Postgres")
	}

	// Session store
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Set max age for cookie
	store.Options = &sessions.Options{
		Path:     "/api/v1",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	r := routes.SetupRouter(db, store, logger)
	logger.Info().Msg("Setup Complete. Starting user-service...")
	r.Run(":" + os.Getenv("PORT"))
}
