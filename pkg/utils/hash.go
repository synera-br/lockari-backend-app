package utils

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

// GenerateTenantID generates a cryptographically secure 32-character tenant ID
func GenerateTenantID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate tenant ID: %w", err)
	}
	return fmt.Sprintf("%x", bytes), nil
}

func GenerateTenant() string {
	uid, _ := uuid.NewV7()

	return uid.String()
}

func ValidateUUID(id string) (bool, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("invalid UUID: %w", err)
	}
	return true, nil
}

func ValidateUUIDOrEmpty(id string) (bool, error) {
	if id == "" {
		return true, nil
	}
	return ValidateUUID(id)
}
