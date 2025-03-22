package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func DecodeBase64Image(base64Str string) image.Image {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		log.Fatalf("Failed to decode base64: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	return img
}

func CompressImage(img image.Image, quality int) []byte {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		log.Fatalf("Failed to compress image: %v", err)
	}
	return buf.Bytes()
}

func SaveToFile(data []byte, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
}
