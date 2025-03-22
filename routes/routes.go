package routes

import (
	"golang_api/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(client *mongo.Client) *gin.Engine {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		controller.Login(client, c)
	})
	router.POST("/insert-user", func(c *gin.Context) {
		controller.InsertUser(client, c)
	})

	return router
}
