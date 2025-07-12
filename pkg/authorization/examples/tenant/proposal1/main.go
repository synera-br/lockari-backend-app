package main

import (
	"context"
	"fmt"
	"log"

	"github.com/synera-br/lockari-backend-app/pkg/authorization"
)

// This is a simplified example to demonstrate the direct integration approach.
// In a real application, the authorization.Service would be properly initialized.
type MockAuthService struct{}

func (s *MockAuthService) CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error {
	log.Printf("Tenant '%s' created for user '%s'", tenantID, userID)
	log.Printf("User '%s' assigned as admin to tenant '%s'", userID, tenantID)
	return nil
}

func main() {
	// In a real application, you would initialize the authorization service here.
	authSvc := &MockAuthService{}

	// Simulate a signup event
	tenantID := "new-tenant-123"
	userID := "user-abc-789"

	fmt.Println("Simulating signup for a new user...")
	err := authSvc.CreateTenantWithAdmin(context.Background(), tenantID, userID)
	if err != nil {
		log.Fatalf("Failed to create tenant with admin: %v", err)
	}

	fmt.Println("Signup simulation completed successfully!")
}
