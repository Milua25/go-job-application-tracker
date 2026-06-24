package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		slog.Error("failed to generate salt", "error", err)
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	return hashBase64 + "." + saltBase64, nil
}

func VerifyPassword(hashedPassword, password string) bool {
	parts := strings.Split(hashedPassword, ".")
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

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if len(computedHash) != len(hash) {
		return false
	}
	return subtle.ConstantTimeCompare(computedHash, hash) == 1
}
