package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/argon2"
)

const (
	// Argon2 parameters
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4
	argon2KeyLen  = 32
	hashDelimiter = "$"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		slog.Error("failed to generate salt", "error", err)
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	return hashBase64 + hashDelimiter + saltBase64, nil
}

func VerifyPassword(hashedPassword, password string) bool {
	parts := strings.Split(hashedPassword, hashDelimiter)
	if len(parts) != 2 {
		slog.Error("invalid hashed password format")
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		slog.Error("failed to decode salt", "error", err)
		return false
	}

	hash, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		slog.Error("failed to decode hash", "error", err)
		return false
	}

	computedHash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	if len(computedHash) != len(hash) {
		return false
	}
	return subtle.ConstantTimeCompare(computedHash, hash) == 1
}

func StrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
