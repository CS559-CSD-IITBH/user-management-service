package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Middleware checks if the user is authenticated as a merchant
func Validate(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Cookie("auth")
		if ok != nil {
			logger.Error().Msg("Unable to validate user session")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Pass the UID to the next handler
		c.Set("uid", uid)
		c.Next()
	}
}
