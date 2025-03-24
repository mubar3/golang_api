package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang_api/config"
	"golang_api/routes"
	"golang_api/utils"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Fungsi untuk menghubungkan ke MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	// Konfigurasi URI MongoDB (ubah sesuai kebutuhan)
	uri := ""
	if utils.Isnotnull(os.Getenv("DB_USER")) {
		uri = fmt.Sprintf(
			"mongodb://%v:%v@%v:%v",
			os.Getenv("DB_USER"),     // DB username
			os.Getenv("DB_PASSWORD"), // DB password
			os.Getenv("DB_HOST"),     // DB host
			os.Getenv("DB_PORT"),     // DB port
		)
	} else {
		uri = fmt.Sprintf(
			"mongodb://%v:%v",
			os.Getenv("DB_HOST"), // DB host
			os.Getenv("DB_PORT"), // DB port
		)
	}

	// Konfigurasi opsi koneksi
	clientOptions := options.Client().ApplyURI(uri)

	// Membuat koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Fatal(err.Error())
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	return client, nil
}

var Timezone *time.Location

func main() {

	// Inisialisasi logger
	err := utils.InitLogger("./logs")
	if err != nil {
		log.Fatalf("Gagal menginisialisasi logger: %v", err)
		return
	}
	defer utils.Logger.Close()

	// Memuat file .env
	err = godotenv.Load()
	if err != nil {
		utils.Logger.LogMessage("ERROR", err.Error())
		// logrus.Fatal(err.Error())
		log.Fatal("Error loading .env file")
		return
	}

	// Set timezone menggunakan package config
	err = config.InitTimezone(os.Getenv("TIMEZONE"))
	if err != nil {
		log.Fatalf("Failed to initialize timezone: %v", err)
		return
	}

	// inisiasi mongodb
	client, err := connectToMongoDB()
	if err != nil {
		utils.Logger.LogMessage("ERROR", err.Error())
		// logrus.Fatal(err.Error())
		log.Fatalf("Error initializing MongoDB connection: %v", err)
		return
	}

	// Jalankan cronjobs menggunakan goroutine
	go utils.StartCronJobs()

	// Setup routing
	router := routes.SetupRoutes(client)

	// Jalankan API di port yang ditentukan
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port jika tidak ada di .env
	}
	if err := router.Run(":" + port); err != nil {
		utils.Logger.LogMessage("ERROR", err.Error())
		// logrus.Fatal(err.Error())
		log.Fatalf("Failed to start the server: %v", err)
		return
	}
}
