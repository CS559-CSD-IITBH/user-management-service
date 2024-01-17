package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/sqlite"
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

	// Add user to database
	db, err := gorm.Open(sqlite.Open(os.Getenv("DSN")), &gorm.Config{})
	db.AutoMigrate(&models.Customer{})
	db.AutoMigrate(&models.Merchant{})
	db.AutoMigrate(&models.DeliveryAgent{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.PasswordResetToken{})

	// Session store in  NewFilesystemStore
	store := sessions.NewFilesystemStore("sessions/", []byte("secret-key"))

	// Set max age for cookie
	store.Options = &sessions.Options{
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	if err != nil {
		log.Fatalln("Internal server error: Unable to connect to the DB")
	}

	r := routes.SetupRouter(db, store)
	r.Run(":" + os.Getenv("PORT"))
}
