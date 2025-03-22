package routes

import (
	"golang_api/controller"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(client *mongo.Client) *gin.Engine {
	router := gin.Default()
	connection := client.Database(os.Getenv("DB_DATABASE"))

	router.POST("/login", func(c *gin.Context) {
		controller.Login(connection, c)
	})
	router.POST("/insert-user", func(c *gin.Context) {
		controller.InsertUser(connection, c)
	})
	router.POST("/change-password", func(c *gin.Context) {
		controller.ChangePassword(connection, c)
	})

	return router
}
