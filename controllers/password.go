package controllers

import (
	"fmt"
	"net/http"

	"github.com/CS559-CSD-IITBH/user-management-service/models"
	"github.com/CS559-CSD-IITBH/user-management-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func ForgotPassword(c *gin.Context, db *gorm.DB, store *sessions.FilesystemStore) {
	// Get JSON body from request
	var requestData struct {
		Email string `json:"email" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email exists in the database
	var user models.User
	db.Where("email = ?", requestData.Email).First(&user)
	if user.UID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Generate a unique token for password reset
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Store the reset token in the database
	db.Create(&models.PasswordResetToken{UID: user.UID, Token: resetToken})

	// TODO: Send reset token to the user (via email, SMS, etc.)
	utils.SendEmail(user.Email, resetToken)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func ResetPassword(c *gin.Context, db *gorm.DB, store *sessions.FilesystemStore) {
	// Get JSON body from request
	var resetData struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Read JSON body
	if err := c.ShouldBindJSON(&resetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user associated with the reset token
	var resetToken models.PasswordResetToken
	db.Where("token = ?", resetData.Token).First(&resetToken)
	if resetToken.UID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid reset token"})
		return
	}

	// Update the user's password
	hashedPassword, _ := utils.HashPassword(resetData.Password)
	result := db.Model(&models.User{}).Where("uid = ?", resetToken.UID).Update("password", hashedPassword)
	if result.Error != nil {
		// Handle the error
		fmt.Println("Error updating password:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Delete the used reset token
	db.Delete(&resetToken)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
