package controller

import (
	"context"
	"golang_api/config"
	"golang_api/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Visit struct {
	Id            string     `bson:"_id"`
	User          string     `bson:"user"`
	Foto          string     `bson:"foto"` // Ini harus berupa hash
	Time_checkin  time.Time  `bson:"time_checkin"`
	Time_checkout *time.Time `bson:"time_checout"`
	Created_at    time.Time  `bson:"created_time"`
}

func VisitIn(connection *mongo.Database, w http.ResponseWriter, c *http.Request) {
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

	// logs
	utils.Logger.WithField("user_mobile", user_mobile.Id).LogMessage("LOG", "Accessing API endpoint /visit-in")

	// Dapatkan koleksi user_mobile
	collection = connection.Collection("visit")

	// Data user yang akan disimpan
	visit := bson.M{
		"user":         user_mobile.Id,
		"foto":         c.PostFormValue("foto"),
		"time_checkin": time.Now().In(config.Timezone),
		"created_at":   time.Now().In(config.Timezone),
	}

	// memberi timeout 5 detik
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, visit)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.Response(w, http.StatusOK, "Berhasil visit in!", nil, nil)
}
