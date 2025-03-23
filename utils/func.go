package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HashPassword(pass string, key string) string {
	// Hash the password using SHA-256
	hashPass := sha256.Sum256([]byte(pass))

	// Convert hashPass to a hex string
	hashPassHex := hex.EncodeToString(hashPass[:])

	// Concatenate the key and the hashed password, then hash again
	encryptPass := sha256.Sum256([]byte(key + hashPassHex))

	// Convert encryptPass to a hex string and return
	return hex.EncodeToString(encryptPass[:])
}

func Isnotnull(pass any) bool {
	if pass == nil || pass == "null" || pass == "" {
		return false
	}
	return true
}

func NullValidation(data interface{}) (status bool, message string) {
	// Periksa apakah data adalah map
	if items, ok := data.(map[string]interface{}); ok {
		for key, value := range items {
			if !Isnotnull(value) {
				return false, key + " kosong"
			}
		}
		return true, ""
	}
	return false, "data tidak valid"
}

func Response(c *gin.Context, httpStatus int, message string, user any) {
	if Isnotnull(user) {
		logrus.WithField("user_mobile", user).Warn(message)
	}
	if httpStatus == 200 {
		c.JSON(httpStatus, gin.H{
			"status":  "succes",
			"message": message,
		})
	} else {
		c.JSON(httpStatus, gin.H{
			"status":  "error",
			"message": message,
		})
	}
}
