package controllers

import (
	"net/http"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Create(c *gin.Context, db *gorm.DB, store *sessions.FilesystemStore) {

	// Get JSON body from request
	var userData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Name     string `json:"name" binding:"required"`
		UserType string `json:"user_type" binding:"required"`

		// Fields for customers
		Address string `json:"address"`

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
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	// Now we can use the UID to store additional information about the user in a separate table
	if userData.UserType == "customer" {
		result = db.Create(&models.Customer{UID: uid, Address: userData.Address})
	} else if userData.UserType == "merchant" {
		result = db.Create(&models.Merchant{UID: uid, StoreAddress: userData.StoreAddress})
	} else if userData.UserType == "delivery_agent" {
		result = db.Create(&models.DeliveryAgent{UID: uid, LicenseNumber: userData.LicenseNumber, VehicleType: userData.VehicleType, VehicleNumber: userData.VehicleNumber})
	}

	// Check for errors
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"status": "success"})

}

func Update(c *gin.Context, db *gorm.DB, store *sessions.FilesystemStore) {
	uid := c.Param("uid")

	// Get JSON body from request
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists in any of the tables
	var user models.User
	db.Where("uid = ?", uid).First(&user)

	if user.UID == "" {
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

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Create /login endpoint and return JWT token
func Login(c *gin.Context, db *gorm.DB, store *sessions.FilesystemStore) {

	// Get JSON body from request
	var userData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	var user models.User
	db.Where("email = ?", userData.Email).First(&user)

	// Hash the password using the bcrypt package
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Gorilla sessions
	session, _ := store.Get(c.Request, "session-name")
	session.Values["user"] = user.UID
	session.Save(c.Request, c.Writer)

}

// Logout endpoint
func Logout(c *gin.Context, db *gorm.DB, store *sessions.FilesystemStore) {
	session, _ := store.Get(c.Request, "session-name")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
