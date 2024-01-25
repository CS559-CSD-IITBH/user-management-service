package routes

import (
	"github.com/CS559-CSD-IITBH/user-management-service/controllers"
	"github.com/CS559-CSD-IITBH/user-management-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, store *sessions.CookieStore, logger zerolog.Logger) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			// Pass 'db' and 'store' to controllers
			user.POST("/create", func(c *gin.Context) {
				controllers.Create(c, db, store, logger)
			})
			user.PATCH("/update", middlewares.SessionAuth(store, logger), func(c *gin.Context) {
				controllers.Update(c, db, store, logger)
			})
			user.POST("/login", func(c *gin.Context) {
				controllers.Login(c, db, store, logger)
			})
			user.GET("/logout", func(c *gin.Context) {
				controllers.Logout(c, db, store, logger)
			})
		}

		password := v1.Group("/password")
		{
			// Pass 'db' and 'store' to controllers
			password.POST("/forgot", func(c *gin.Context) {
				controllers.ForgotPassword(c, db, store, logger)
			})
			password.POST("/reset", func(c *gin.Context) {
				controllers.ResetPassword(c, db, store, logger)
			})
		}
	}

	return r
}
