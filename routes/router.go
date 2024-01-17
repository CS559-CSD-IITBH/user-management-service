package routes

import (
	"github.com/CS559-CSD-IITBH/user-management-service/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, store *sessions.FilesystemStore) *gin.Engine {
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
				controllers.Create(c, db, store)
			})
			user.PATCH("/update/:uid", func(c *gin.Context) {
				controllers.Update(c, db, store)
			})
			user.POST("/login", func(c *gin.Context) {
				controllers.Login(c, db, store)
			})
			user.GET("/logout", func(c *gin.Context) {
				controllers.Logout(c, db, store)
			})
		}

		password := v1.Group("/password")
		{
			// Pass 'db' and 'store' to controllers
			password.POST("/forgot", func(c *gin.Context) {
				controllers.ForgotPassword(c, db, store)
			})
			password.POST("/reset", func(c *gin.Context) {
				controllers.ResetPassword(c, db, store)
			})
		}
	}

	return r
}
