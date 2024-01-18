package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/routes"
	"github.com/gorilla/sessions"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Internal server error: Unable to load the env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("Internal server error: Unable to connect to Postgres")
	}

	err = db.AutoMigrate(
		&models.Customer{},
		&models.Merchant{},
		&models.DeliveryAgent{},
		&models.User{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		log.Fatal("Internal server error: Unable to migrate models to Postgres")
	}

	// Session store in  NewFilesystemStore
	store := sessions.NewFilesystemStore("sessions/", []byte("secret-key"))

	// Set max age for cookie
	store.Options = &sessions.Options{
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	r := routes.SetupRouter(db, store)
	r.Run(":" + os.Getenv("PORT"))
}
