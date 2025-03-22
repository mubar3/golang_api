package utils

import (
	"crypto/sha256"
	"encoding/hex"
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
