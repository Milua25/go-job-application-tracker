package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Println("Failed to generate salt", "error", err)
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	return hashBase64 + "." + saltBase64, nil
}

func VerifyPassword(hashedPassword, password string) bool {

	// split the stored hash into its components (hash and salt)
	parts := strings.Split(hashedPassword, ".")
	if len(parts) != 2 {
		log.Println("Invalid hashed password format")
		return false
	}
	hashBase64 := parts[0]
	saltBase64 := parts[1]

	// decode the salt from base64
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		log.Println("Failed to decode salt", "error", err)
		return false
	}
	// decode the hash from base64
	hash, err := base64.StdEncoding.DecodeString(hashBase64)
	if err != nil {
		log.Println("Failed to decode hash", "error", err)
		return false
	}
	// generate a hash of the provided password using the same salt and parameters
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	// compare the computed hash with the stored hash

	// check if the lengths of the hashes are the same to prevent timing attacks
	if len(computedHash) != len(hash) {
		return false
	}

	// use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(computedHash, hash) == 1 {
		// do nothing, the hashes match
	} else {
		return false
	}

	return true
}
