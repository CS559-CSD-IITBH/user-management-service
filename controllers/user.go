package controllers

import (
	"net/http"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Create(c *gin.Context, db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) {

	// Get JSON body from request
	var userData struct {
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		MobileNumber string `json:"mobile_number" binding:"required"`
		UserType     string `json:"user_type" binding:"required"`

		// Fields for customers
		CustomerName    string `json:"customer_name"`
		DeliveryAddress string `json:"delivery_address"`

		// Fields for merchants
		MerchantName string `json:"merchant_name"`
		StoreAddress string `json:"store_address"`

		// Fields for delivery agents
		LicenseNumber string `json:"license_number"`
		VehicleType   string `json:"vehicle_type"`
		VehicleNumber string `json:"vehicle_number"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&userData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password using the bcrypt package
	hashedPassword, _ := utils.HashPassword(userData.Password)

	// Generate a random UID using the crypto/rand package
	uid, _ := utils.GenerateUID()

	// Store the user in the database
	result := db.Create(&models.User{UID: uid, Email: userData.Email, UserType: userData.UserType, Password: hashedPassword})

	// Check for errors
	if result.Error != nil {
		logger.Error().Msg("Unable to create user in DB")
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	// Now we can use the UID to store additional information about the user in a separate table
	if userData.UserType == "customer" {
		result = db.Create(&models.Customer{UID: uid, CustomerName: userData.CustomerName, MobileNumber: userData.MobileNumber, DeliveryAddress: userData.DeliveryAddress})
	} else if userData.UserType == "merchant" {
		result = db.Create(&models.Merchant{UID: uid, MerchantName: userData.MerchantName, MobileNumber: userData.MobileNumber, StoreAddress: userData.StoreAddress})
	} else if userData.UserType == "delivery_agent" {
		result = db.Create(&models.DeliveryAgent{UID: uid, LicenseNumber: userData.LicenseNumber, VehicleType: userData.VehicleType, VehicleNumber: userData.VehicleNumber})
	}

	// Check for errors
	if result.Error != nil {
		logger.Error().Msg("Unable to add to user-specific DB")
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(c.Request, "auth")
	if err != nil {
		logger.Error().Msg("Unable to get user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Set session values.
	session.Values["user"] = uid

	// Save it
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logger.Error().Msg("Unable to save user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Return success
	logger.Info().Msg("User successfully created")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func Update(c *gin.Context, db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) {
	uid := c.Param("uid")

	// Get JSON body from request
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists in any of the tables
	var user models.User
	db.Where("uid = ?", uid).First(&user)

	if user.UID == "" {
		logger.Error().Msg("Unable to find user")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Update user based on user type
	if user.UserType == "customer" {
		var customer models.Customer
		db.Where("uid = ?", uid).First(&customer)
		utils.UpdateFields(&customer, updateData)
		db.Save(&customer)
	} else if user.UserType == "merchant" {
		var merchant models.Merchant
		db.Where("uid = ?", uid).First(&merchant)
		utils.UpdateFields(&merchant, updateData)
		db.Save(&merchant)
	} else if user.UserType == "delivery_agent" {
		var deliveryAgent models.DeliveryAgent
		db.Where("uid = ?", uid).First(&deliveryAgent)
		utils.UpdateFields(&deliveryAgent, updateData)
		db.Save(&deliveryAgent)
	}

	// Return success
	logger.Info().Msg("User successfully updated")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func Login(c *gin.Context, db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) {

	// Get JSON body from request
	var userData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&userData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	var user models.User
	db.Where("email = ?", userData.Email).First(&user)

	// Hash the password using the bcrypt package
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		logger.Error().Msg("Email or password is incorrect")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Gorilla sessions
	session, err := store.Get(c.Request, "auth")
	if err != nil {
		logger.Error().Msg("Unable to get user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Set session values.
	session.Values["user"] = user.UID

	// Log session information
	logger.Info().Msgf("Session user: %v", session.Values["user"])

	// Save it
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logger.Error().Msg("Unable to save user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Return success
	logger.Info().Msg("User login successful")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Logout endpoint
func Logout(c *gin.Context, db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) {
	session, err := store.Get(c.Request, "auth")
	if err != nil {
		logger.Error().Msg("Unable to get user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	session.Options.MaxAge = -1

	// Save it
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logger.Error().Msg("Unable to save user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Return success
	logger.Info().Msg("User logout successful")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
