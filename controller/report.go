package controller

import (
	"context"
	"golang_api/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadImg(connection *mongo.Database, w http.ResponseWriter, c *http.Request) {
	status, eror := utils.NullValidation(map[string]interface{}{
		"session": c.PostFormValue("session_id"),
		"foto":    c.PostFormValue("foto"),
	})
	if !status {
		utils.Response(w, http.StatusBadRequest, eror, nil, nil)
		return
	}

	// cek db
	collection := connection.Collection("user_mobile")
	var user_mobile User_mobile
	filter := bson.M{"session_id": c.PostFormValue("session_id")}
	err := collection.FindOne(context.TODO(), filter).Decode(&user_mobile)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, "Session tidak tersedia", nil, nil)
		return
	}

	// cek foto base64
	if !utils.IsBase64ImageValid(c.PostFormValue("foto")) {
		utils.Response(w, http.StatusBadRequest, "Foto tidak valid", user_mobile.Id, nil)
		return
	}

	// logs
	utils.Logger.WithField("user_mobile", user_mobile.Id).LogMessage("LOG", "Accessing API endpoint /upload-img")

	// Dekode, kompres, dan simpan gambar
	path, err := utils.DecodeAndCompressBase64Image(c.PostFormValue("foto"), "./asset")
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err.Error(), user_mobile.Id, nil)
		return
	}

	utils.Response(w, http.StatusOK, "Upload berhasil!", user_mobile.Id, map[string]interface{}{
		"filename": path,
	})
}
