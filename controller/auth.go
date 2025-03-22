package controller

import (
	"context"
	"net/http"
	"time"

	"golang_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User_mobile struct {
	Username   string    `bson:"username"`
	Password   string    `bson:"password"`     // Ini harus berupa hash
	Keyz       string    `bson:"keyz"`         // Ini harus berupa hash
	Created_at time.Time `bson:"created_time"` // Ini harus berupa hash
}

// Fungsi Login
func Login(client *mongo.Client, c *gin.Context) {
	// Ambil input username dan password dari form-data
	username := c.PostForm("username")
	password := c.PostForm("password")

	// cek inputan
	if !utils.Isnotnull(username) || !utils.Isnotnull(password) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Username/password kosong",
		})
		return
	}

	// cek db
	collection := client.Database("local").Collection("user_mobile")
	var user_mobile User_mobile
	filter := bson.M{"username": username}
	err := collection.FindOne(context.TODO(), filter).Decode(&user_mobile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Username Salah",
		})
		return
	}

	// Hash password
	password = utils.HashPassword(password, user_mobile.Keyz)
	if password != user_mobile.Password {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Password Salah",
		})
		return
	}

	session_id := uuid.NewString()
	update := bson.M{
		"$set": bson.M{
			"session_id": session_id, // Nilai session ID baru
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logrus.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successfully!",
		"session": session_id,
	})
}

// Fungsi untuk menyisipkan user ke koleksi MongoDB
func InsertUser(client *mongo.Client, c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	keyz := "265"

	// cek inputan
	if !utils.Isnotnull(username) || !utils.Isnotnull(password) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Username/password kosong",
		})
		return
	}

	// Hash password
	hashedPassword := utils.HashPassword(password, keyz)

	// Dapatkan koleksi user_mobile
	collection := client.Database("local").Collection("user_mobile")

	// Data user yang akan disimpan
	user := bson.M{
		"username":   username,
		"password":   hashedPassword,
		"keyz":       keyz,
		"created_at": time.Now(),
	}

	// Insert data ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		logrus.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to insert user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User inserted successfully!",
	})
}
