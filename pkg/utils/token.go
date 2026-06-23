package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken hashes a token using SHA-256. Suitable for refresh tokens which are
// already cryptographically random; bcrypt/argon2 are unnecessary here.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
