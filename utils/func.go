package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		if httpStatus == 400 {
			Logger.WithField("user_mobile", user).LogMessage("WARNING", message)
		} else {
			Logger.WithField("user_mobile", user).LogMessage("ERROR", message)
		}
		// logrus.WithField("user_mobile", user).Warn(message)
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

func RemoveBase64Prefix(base64Image string) string {
	// Cek apakah string mengandung prefix "data:image/"
	if idx := strings.Index(base64Image, ","); idx != -1 {
		return base64Image[idx+1:] // Ambil hanya bagian setelah koma
	}
	return base64Image // Jika tidak ada prefix, kembalikan string asli
}

func IsBase64ImageValid(base64Image string) bool {
	base64Image = RemoveBase64Prefix(base64Image)

	// Coba decode string Base64
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		fmt.Println("Base64 decoding failed:", err)
		return false
	}

	// Coba ubah data menjadi gambar
	_, _, err = image.Decode(bytes.NewReader(imageData))
	if err != nil {
		fmt.Println("Image decoding failed:", err)
		return false
	}

	// Jika tidak ada error, berarti valid
	return true
}

func DecodeAndCompressBase64Image(base64Image, folderPath string) (string, error) {
	// Hilangkan prefix data:image/...
	base64Image = RemoveBase64Prefix(base64Image)

	// Decode Base64 ke byte array
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// Decode byte array ke image.Image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image data: %w", err)
	}

	// Tentukan kualitas JPEG
	quality := 50

	// Buat folder dengan struktur tahun/bulan
	now := time.Now()
	year, month := now.Format("2006"), now.Format("01")
	finalFolderPath := filepath.Join(folderPath, year, month)

	// Pastikan folder target ada
	if err := os.MkdirAll(finalFolderPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}

	// Buat nama file berdasarkan timestamp
	fileName := fmt.Sprintf("%d.jpg", now.Unix())
	fullPath := filepath.Join(finalFolderPath, fileName)

	// Buat file output
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() // Pastikan file tertutup setelah selesai

	// Encode gambar ke file dengan kompresi JPEG
	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: quality}); err != nil {
		return "", fmt.Errorf("failed to save compressed image: %w", err)
	}

	// Kembalikan path relatif (tahun/bulan/nama_file.jpg)
	// return filepath.Join(year, month, fileName), nil
	// Kembalikan jalur relatif tahun/bulan + nama file
	relativePath := filepath.Join(year, month, fileName)
	// Normalisasi separator ke `/`
	normalizedPath := strings.ReplaceAll(relativePath, "\\", "/")
	return normalizedPath, nil

}
