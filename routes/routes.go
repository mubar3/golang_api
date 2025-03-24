package routes

import (
	"golang_api/controller"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(client *mongo.Client) *gin.Engine {
	router := gin.Default()
	connection := client.Database(os.Getenv("DB_DATABASE"))

	// Custom handler untuk file tidak ada
	router.GET("/static/*filepath", func(c *gin.Context) {
		assetPath := "./asset"                                    // Path folder asset
		defaultImage := filepath.Join(assetPath, "none.jpg")      // Path gambar default
		filePath := filepath.Join(assetPath, c.Param("filepath")) // Path file yang diminta

		// Periksa keberadaan file
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Jika file tidak ditemukan, arahkan ke gambar default
			c.File(defaultImage)
		} else {
			// Jika file ditemukan, kirim file tersebut
			c.File(filePath)
		}
	})

	router.POST("/login", func(c *gin.Context) {
		controller.Login(connection, c)
	})
	router.POST("/insert-user", func(c *gin.Context) {
		controller.InsertUser(connection, c)
	})
	router.POST("/change-password", func(c *gin.Context) {
		controller.ChangePassword(connection, c)
	})
	router.POST("/get-data", func(c *gin.Context) {
		controller.GetData(connection, c)
	})
	router.POST("/upload-img", func(c *gin.Context) {
		controller.UploadImg(connection, c)
	})

	return router
}
