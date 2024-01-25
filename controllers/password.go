package controllers

import (
	"net/http"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func ForgotPassword(c *gin.Context, db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) {
	// Get JSON body from request
	var requestData struct {
		Email string `json:"email" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&requestData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email exists in the database
	var user models.User
	db.Where("email = ?", requestData.Email).First(&user)
	if user.UID == "" {
		logger.Error().Msg("Unable to find user")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Generate a unique token for password reset
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		logger.Error().Msg("Unable to generate unique token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Store the reset token in the database
	db.Create(&models.PasswordResetToken{UID: user.UID, Token: resetToken})

	// TODO: Send reset token to the user (via email, SMS, etc.)
	utils.SendEmail(user.Email, resetToken)

	logger.Info().Msg("Password reset token sent successfully")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func ResetPassword(c *gin.Context, db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) {
	// Get JSON body from request
	var resetData struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&resetData); err != nil {
		logger.Error().Msg("Unable to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user associated with the reset token
	var resetToken models.PasswordResetToken
	db.Where("token = ?", resetData.Token).First(&resetToken)
	if resetToken.UID == "" {
		logger.Error().Msg("Invalid reset token")
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid reset token"})
		return
	}

	// Update the user's password
	hashedPassword, _ := utils.HashPassword(resetData.Password)
	result := db.Model(&models.User{}).Where("uid = ?", resetToken.UID).Update("password", hashedPassword)
	if result.Error != nil {
		// Handle the error
		logger.Error().Msg("Unable to update password in DB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Delete the used reset token
	db.Delete(&resetToken)

	logger.Info().Msg("Password reset successfully")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
