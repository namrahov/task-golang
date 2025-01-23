package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

// GenerateToken generates a random token by hashing a UUID.
func GenerateToken() string {
	randomUUID := uuid.New()
	return GenerateSHA(randomUUID.String())
}

// GenerateSHA generates a SHA-256 hash for the given text.
func GenerateSHA(text string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
