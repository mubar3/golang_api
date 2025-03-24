package controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"golang_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User_mobile struct {
	Id         string    `bson:"_id"`
	Username   string    `bson:"username"`
	Password   string    `bson:"password"` // Ini harus berupa hash
	Keyz       string    `bson:"keyz"`
	Session    *string   `bson:"session_id"` // Jika tidak ada, akan menjadi nil
	Foto       *string   `bson:"foto"`       // Jika tidak ada, akan menjadi nil
	Jabatan    *string   `bson:"jabatan"`    // Jika tidak ada, akan menjadi nil
	Created_at time.Time `bson:"created_time"`
}

// Fungsi Login
func Login(connection *mongo.Database, c *gin.Context) {
	status, eror := utils.NullValidation(map[string]interface{}{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
	if !status {
		utils.Response(c, http.StatusBadRequest, eror, nil)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	// cek db
	collection := connection.Collection("user_mobile")
	var user_mobile User_mobile
	filter := bson.M{"username": username}
	err := collection.FindOne(context.TODO(), filter).Decode(&user_mobile)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "Username Salah", nil)
		return
	}

	// Hash password
	password = utils.HashPassword(password, user_mobile.Keyz)
	if password != user_mobile.Password {
		utils.Response(c, http.StatusBadRequest, "Password Salah", user_mobile.Id)
		return
	}

	// logs
	utils.Logger.WithField("user_mobile", user_mobile.Id).LogMessage("LOG", "Accessing API endpoint /login")

	session_id := uuid.NewString()
	update := bson.M{
		"$set": bson.M{
			"session_id": session_id, // Nilai session ID baru
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Proses unggah data gagal. Silakan coba lagi nanti.", user_mobile.Id)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil!",
		"session": session_id,
	})
}

func ChangePassword(connection *mongo.Database, c *gin.Context) {
	status, eror := utils.NullValidation(map[string]interface{}{
		"session":  c.PostForm("session_id"),
		"password": c.PostForm("password"),
	})
	if !status {
		utils.Response(c, http.StatusBadRequest, eror, nil)
		return
	}

	// cek db
	collection := connection.Collection("user_mobile")
	var user_mobile User_mobile
	filter := bson.M{"session_id": c.PostForm("session_id")}
	err := collection.FindOne(context.TODO(), filter).Decode(&user_mobile)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "Session tidak tersedia", nil)
		return
	}

	// Hash password
	password := utils.HashPassword(c.PostForm("password"), user_mobile.Keyz)
	if password == user_mobile.Password {
		utils.Response(c, http.StatusBadRequest, "Password tidak boleh sama dengan password yang lama", user_mobile.Id)
		return
	}
	// logs
	utils.Logger.WithField("user_mobile", user_mobile.Id).LogMessage("LOG", "Accessing API endpoint /change-password")

	update := bson.M{
		"$set": bson.M{
			"password": password, // Nilai session ID baru
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "Proses unggah data gagal. Silakan coba lagi nanti.", user_mobile.Id)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Ubah password berhasil!",
	})
}

// Fungsi untuk menyisipkan user ke koleksi MongoDB
func InsertUser(connection *mongo.Database, c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	keyz := "265"

	status, eror := utils.NullValidation(map[string]interface{}{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
	if !status {
		utils.Response(c, http.StatusBadRequest, eror, nil)
		return
	}

	// Hash password
	hashedPassword := utils.HashPassword(password, keyz)

	// Dapatkan koleksi user_mobile
	collection := connection.Collection("user_mobile")

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
		utils.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User inserted successfully!",
	})
}

func GetData(connection *mongo.Database, c *gin.Context) {
	status, eror := utils.NullValidation(map[string]interface{}{
		"session": c.PostForm("session_id"),
	})
	if !status {
		utils.Response(c, http.StatusBadRequest, eror, nil)
		return
	}

	// cek db
	collection := connection.Collection("user_mobile")
	var user_mobile User_mobile
	filter := bson.M{"session_id": c.PostForm("session_id")}
	err := collection.FindOne(context.TODO(), filter).Decode(&user_mobile)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "Session tidak tersedia", nil)
		return
	}

	// logs
	utils.Logger.WithField("user_mobile", user_mobile.Id).LogMessage("LOG", "Accessing API endpoint /get-data")

	foto := os.Getenv("URL") + "/static/profil/blank.png"
	if user_mobile.Foto != nil {
		foto = os.Getenv("URL") + "/static/profil/" + *user_mobile.Foto // Dereferensi pointer
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{ // Gunakan gin.H untuk struktur key-value
			"username": user_mobile.Username,
			"jabatan":  user_mobile.Jabatan,
			"foto":     foto,
		},
	})
}
