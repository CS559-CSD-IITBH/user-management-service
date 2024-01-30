package controllers

// import (
// 	"context"
// 	"net/http"

// 	"github.com/CS559-CSD-IITBH/user-management-service/models"
// 	"github.com/CS559-CSD-IITBH/user-management-service/utils"
// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/zerolog"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func ForgotPassword(c *gin.Context, customers *mongo.Collection, merchants *mongo.Collection, logger zerolog.Logger) {
// 	var requestData struct {
// 		Email string `json:"email" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&requestData); err != nil {
// 		logger.Error().Msg("Unable to bind JSON")
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user models.Customer
// 	err := customers.FindOne(context.TODO(), bson.M{"email": requestData.Email}).Decode(&user)
// 	if err != nil {
// 		err = merchants.FindOne(context.TODO(), bson.M{"email": requestData.Email}).Decode(&user)
// 	}
// 	if user.UID == "" {
// 		logger.Error().Msg("Unable to find user")
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}

// 	resetToken, err := utils.GenerateResetToken()
// 	if err != nil {
// 		logger.Error().Msg("Unable to generate unique token")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
// 		return
// 	}

// 	resetTokenModel := models.PasswordResetToken{UID: user.UID, Token: resetToken}
// 	_, err = db.InsertOne(context.TODO(), resetTokenModel)
// 	if err != nil {
// 		logger.Error().Msg("Unable to store reset token in DB")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
// 		return
// 	}

// 	utils.SendEmail(user.Email, resetToken)

// 	logger.Info().Msg("Password reset token sent successfully")
// 	c.JSON(http.StatusOK, gin.H{"status": "success"})
// }

// func ResetPassword(c *gin.Context, customers *mongo.Collection, merchants *mongo.Collection, logger zerolog.Logger) {
// 	var resetData struct {
// 		Token    string `json:"token" binding:"required"`
// 		Password string `json:"password" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&resetData); err != nil {
// 		logger.Error().Msg("Unable to bind JSON")
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var resetToken models.PasswordResetToken
// 	err := db.FindOne(context.TODO(), bson.M{"token": resetData.Token}).Decode(&resetToken)
// 	if err != nil {
// 		logger.Error().Msg("Invalid reset token")
// 		c.JSON(http.StatusNotFound, gin.H{"error": "invalid reset token"})
// 		return
// 	}

// 	hashedPassword, _ := utils.HashPassword(resetData.Password)
// 	result, err := db.UpdateOne(context.TODO(),
// 		bson.M{"uid": resetToken.UID},
// 		bson.M{"$set": bson.M{"password": hashedPassword}})
// 	if err != nil {
// 		logger.Error().Msg("Unable to update password in DB")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
// 		return
// 	}

// 	if result.ModifiedCount == 0 {
// 		logger.Error().Msg("User not found for the given reset token")
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}

// 	_, err = db.DeleteOne(context.TODO(), bson.M{"token": resetData.Token})
// 	if err != nil {
// 		logger.Error().Msg("Unable to delete reset token from DB")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
// 		return
// 	}

// 	logger.Info().Msg("Password reset successfully")
// 	c.JSON(http.StatusOK, gin.H{"status": "success"})
// }
