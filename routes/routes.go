package routes

import (
	"golang_api/controller"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
)

// func SetupRoutes(client *mongo.Client) *gin.Engine {
// 	router := gin.Default()
// 	connection := client.Database(os.Getenv("DB_DATABASE"))

// 	// Custom handler untuk file tidak ada
// 	router.GET("/static/*filepath", func(c *gin.Context) {
// 		assetPath := "./asset"                                    // Path folder asset
// 		defaultImage := filepath.Join(assetPath, "none.jpg")      // Path gambar default
// 		filePath := filepath.Join(assetPath, c.Param("filepath")) // Path file yang diminta

// 		// Periksa keberadaan file
// 		if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 			// Jika file tidak ditemukan, arahkan ke gambar default
// 			c.File(defaultImage)
// 		} else {
// 			// Jika file ditemukan, kirim file tersebut
// 			c.File(filePath)
// 		}
// 	})

// 	router.POST("/login", func(c *gin.Context) {
// 		controller.Login(connection, c)
// 	})
// 	router.POST("/insert-user", func(c *gin.Context) {
// 		controller.InsertUser(connection, c)
// 	})
// 	router.POST("/change-password", func(c *gin.Context) {
// 		controller.ChangePassword(connection, c)
// 	})
// 	router.POST("/get-data", func(c *gin.Context) {
// 		controller.GetData(connection, c)
// 	})
// 	router.POST("/upload-img", func(c *gin.Context) {
// 		controller.UploadImg(connection, c)
// 	})

// 	return router
// }

func SetupRoutes(client *mongo.Client) http.Handler {
	router := http.NewServeMux()
	connection := client.Database(os.Getenv("DB_DATABASE"))

	// Rute untuk file statis
	router.HandleFunc("/static/", StaticFileHandler)

	AddRoute(router, http.MethodPost, "/login", func(w http.ResponseWriter, r *http.Request) {
		controller.Login(connection, w, r)
	})
	AddRoute(router, http.MethodPost, "/insert-user", func(w http.ResponseWriter, r *http.Request) {
		controller.InsertUser(connection, w, r)
	})
	AddRoute(router, http.MethodPost, "/change-password", func(w http.ResponseWriter, r *http.Request) {
		controller.ChangePassword(connection, w, r)
	})
	AddRoute(router, http.MethodPost, "/get-data", func(w http.ResponseWriter, r *http.Request) {
		controller.GetData(connection, w, r)
	})
	AddRoute(router, http.MethodPost, "/upload-img", func(w http.ResponseWriter, r *http.Request) {
		controller.UploadImg(connection, w, r)
	})

	return router
}

// Custom handler untuk file tidak ada
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	assetPath := "./asset"                                             // Path folder asset
	defaultImage := filepath.Join(assetPath, "none.jpg")               // Path gambar default
	filePath := filepath.Join(assetPath, r.URL.Path[len("/static/"):]) // Path file yang diminta

	// Periksa keberadaan file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Jika file tidak ditemukan, arahkan ke gambar default
		http.ServeFile(w, r, defaultImage)
		return
	}

	// Jika file ditemukan, kirim file tersebut
	http.ServeFile(w, r, filePath)
}

func AddRoute(mux *http.ServeMux, method, path string, handlerFunc http.HandlerFunc) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handlerFunc(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func logRequest(r *http.Request) {
	log.Printf("%s %s %s from %s\n", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)
}
