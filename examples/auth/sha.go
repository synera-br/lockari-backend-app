package main

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

func main() {
	tenantID, err := GenerateTenantID()
	if err != nil {
		fmt.Println("Error generating tenant ID:", err)
		return
	}
	fmt.Println("Generated Tenant ID:", tenantID)
}
