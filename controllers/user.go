package controllers

import (
	"context"
	"net/http"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Create(c *gin.Context, customers *mongo.Collection, merchants *mongo.Collection, store *sessions.CookieStore, logger zerolog.Logger) {
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
		// AgentName string `json:"agent_name"`
		// Verified  bool   `json:"verified"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&userData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password using the bcrypt package
	hashedPassword, _ := utils.HashPassword(userData.Password)

	// UID of the document to be inserted
	var uid string

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(c.Request, "auth")
	if err != nil {
		logger.Error().Msg("Unable to get user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if userData.UserType == "customer" {
		userModel := models.Customer{
			Email:           userData.Email,
			Password:        hashedPassword,
			CustomerName:    userData.CustomerName,
			MobileNumber:    userData.MobileNumber,
			DeliveryAddress: userData.DeliveryAddress,
		}

		result, err := customers.InsertOne(context.TODO(), userModel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Assuming the inserted document has an "_id" field
		uid = result.InsertedID.(primitive.ObjectID).Hex()

		// Set session values.
		session.Values["user"] = uid
	} else if userData.UserType == "merchant" {
		userModel := models.Merchant{
			Email:        userData.Email,
			Password:     hashedPassword,
			MerchantName: userData.MerchantName,
			MobileNumber: userData.MobileNumber,
			StoreAddress: userData.StoreAddress,
		}

		result, err := merchants.InsertOne(context.TODO(), userModel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Assuming the inserted document has an "_id" field
		uid = result.InsertedID.(primitive.ObjectID).Hex()

		// Set session values.
		session.Values["user"] = uid
	}

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

func Login(c *gin.Context, customers *mongo.Collection, merchants *mongo.Collection, store *sessions.CookieStore, logger zerolog.Logger) {
	// Get JSON body from request
	var userData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		UserType string `json:"user_type" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&userData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var uid string

	// Gorilla sessions
	session, err := store.Get(c.Request, "auth")
	if err != nil {
		logger.Error().Msg("Unable to get user session")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if userData.UserType == "customer" {
		var user models.Customer
		err := customers.FindOne(context.TODO(), bson.M{"email": userData.Email}).Decode(&user)
		if err != nil {
			logger.Error().Msg("User not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Hash the password using the bcrypt package
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
		if err != nil {
			logger.Error().Msg("Email or password is incorrect")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is incorrect"})
			return
		}
		uid = user.UID.Hex()

		// Set session values.
		session.Values["user"] = uid
	} else if userData.UserType == "merchant" {
		var user models.Merchant
		err := merchants.FindOne(context.TODO(), bson.M{"email": userData.Email}).Decode(&user)
		if err != nil {
			logger.Error().Msg("User not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Hash the password using the bcrypt package
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
		if err != nil {
			logger.Error().Msg("Email or password is incorrect")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is incorrect"})
			return
		}
		uid = user.UID.Hex()

		// Set session values.
		session.Values["user"] = uid
	}

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
func Logout(c *gin.Context, logger zerolog.Logger) {
	// Clear the "auth" cookie
	c.SetCookie("auth", "", -1, "/api/v1", "", true, true)

	// Return success
	logger.Info().Msg("User logout successful")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func Retrieve(c *gin.Context, customers *mongo.Collection, merchants *mongo.Collection, logger zerolog.Logger) {
	// Retrieve UID from the cookie
	uid, err := c.Get("uid")
	if !err {
		logger.Error().Msg("Unable to get UID from cookie")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get UID from cookie"})
		return
	}

	var userData struct {
		UserType string `json:"user_type"`
	}

	// Get JSON body from request
	if err := c.ShouldBindJSON(&userData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Choose the appropriate collection based on user type
	if userData.UserType == "customer" {
		var user models.Customer
		exists := customers.FindOne(context.TODO(), bson.M{"_id": uid}).Decode(user)
		if exists != nil {
			logger.Error().Msg("Unable to find user")
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user.Password = "*******"

		// Return user details
		c.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
	} else if userData.UserType == "merchant" {
		var user models.Merchant
		exists := merchants.FindOne(context.TODO(), bson.M{"_id": uid}).Decode(user)
		if exists != nil {
			logger.Error().Msg("Unable to find user")
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user.Password = "*******"

		// Return user details
		c.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
	} else {
		logger.Error().Msg("Invalid user type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		return
	}
}

// Update user details endpoint
func Update(c *gin.Context, customers *mongo.Collection, merchants *mongo.Collection, logger zerolog.Logger) {
	// Retrieve UID from the cookie
	uid, err := c.Get("uid")
	if !err {
		logger.Error().Msg("Unable to get UID from cookie")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get UID from cookie"})
		return
	}

	// Get JSON body from request
	var userData struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
		UserType     string `json:"user_type" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Address      string `json:"address" binding:"required"`

		// Fields for delivery agents
		// AgentName string `json:"agent_name"`
		// Verified  bool   `json:"verified"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user based on user type
	if userData.UserType == "customer" {
		var userModel models.Customer
		err := customers.FindOne(context.TODO(), bson.M{"_id": uid}).Decode(&userModel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		userModel.MobileNumber = userData.MobileNumber
		userModel.CustomerName = userData.Name
		userModel.DeliveryAddress = userData.Address

		_, err = customers.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": userModel})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update details"})
			return
		}
	} else if userData.UserType == "merchant" {
		var userModel models.Merchant
		err := merchants.FindOne(context.TODO(), bson.M{"_id": uid}).Decode(&userModel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		userModel.MobileNumber = userData.MobileNumber
		userModel.MerchantName = userData.Name
		userModel.StoreAddress = userData.Address

		_, err = merchants.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": userModel})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update details"})
			return
		}
	}

	// Return success
	logger.Info().Msg("User successfully updated")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
