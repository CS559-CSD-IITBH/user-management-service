package routes

import (
	"os"

	"github.com/CS559-CSD-IITBH/user-management-service/controllers"
	"github.com/CS559-CSD-IITBH/user-management-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(customers *mongo.Collection, merchants *mongo.Collection, store *sessions.CookieStore, logger zerolog.Logger) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Content-Type", "*"}
	r.Use(cors.New(config))

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			// Pass 'db' and 'store' to controllers
			user.POST("/create", func(c *gin.Context) {
				controllers.Create(c, customers, merchants, store, logger)
			})
			user.POST("/update", middlewares.Validate(logger), func(c *gin.Context) {
				controllers.Update(c, customers, merchants, logger)
			})
			user.POST("/login", func(c *gin.Context) {
				controllers.Login(c, customers, merchants, store, logger)
			})
			user.GET("/logout", middlewares.Validate(logger), func(c *gin.Context) {
				controllers.Logout(c, logger)
			})
			user.POST("/retrieve", middlewares.Validate(logger), func(c *gin.Context) {
				controllers.Retrieve(c, customers, merchants, logger)
			})
		}

		// password := v1.Group("/password")
		// {
		// 	// Pass 'db' and 'store' to controllers
		// 	password.POST("/forgot", func(c *gin.Context) {
		// 		controllers.ForgotPassword(c, customers, merchants, logger)
		// 	})
		// 	password.POST("/reset", func(c *gin.Context) {
		// 		controllers.ResetPassword(c, customers, merchants, logger)
		// 	})
		// }
	}

	return r
}
