package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateTenantID generates a cryptographically secure 32-character tenant ID
func GenerateTenantID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate tenant ID: %w", err)
	}
	return fmt.Sprintf("%x", bytes), nil
}
