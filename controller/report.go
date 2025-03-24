package controller

import (
	"context"
	"golang_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadImg(connection *mongo.Database, c *gin.Context) {
	status, eror := utils.NullValidation(map[string]interface{}{
		"session": c.PostForm("session_id"),
		"foto":    c.PostForm("foto"),
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

	// cek foto base64
	if !utils.IsBase64ImageValid(c.PostForm("foto")) {
		utils.Response(c, http.StatusBadRequest, "Foto tidak valid", user_mobile.Id)
		return
	}

	// logs
	utils.Logger.WithField("user_mobile", user_mobile.Id).LogMessage("LOG", "Accessing API endpoint /upload-img")

	// Dekode, kompres, dan simpan gambar
	path, err := utils.DecodeAndCompressBase64Image(c.PostForm("foto"), "./asset")
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err.Error(), user_mobile.Id)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Upload berhasil!",
		"filename": path,
	})
}
