package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
)

// AuthMiddleware checks if the user is authenticated
func SessionAuth(store *sessions.CookieStore, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "auth")
		if err != nil {
			logger.Error().Msg("Unable to get user session")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Abort()
			return
		}

		uid, ok := session.Values["user"].(string)
		if !ok {
			logger.Error().Msg("Unable to validate user session")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("uid", uid)
		c.Next()
	}
}
