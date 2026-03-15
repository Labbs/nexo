package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// generateApiKey generates a secure random API key
func generateApiKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "zk_" + hex.EncodeToString(bytes), nil
}

// hashApiKey creates a SHA-256 hash of the API key
func hashApiKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
