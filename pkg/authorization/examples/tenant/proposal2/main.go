package main

import (
	"context"
	"fmt"
	"log"
)

// TenantManager defines the interface for managing tenants.
type TenantManager interface {
	CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error
}

// This is a simplified example to demonstrate the decoupled integration approach.
// In a real application, the authorization.Service would be properly initialized.
type MockAuthService struct{}

func (s *MockAuthService) CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error {
	log.Printf("Tenant '%s' created for user '%s'", tenantID, userID)
	log.Printf("User '%s' assigned as admin to tenant '%s'", userID, tenantID)
	return nil
}

// SignupService depends on the TenantManager interface, not the concrete implementation.
type SignupService struct {
	tenantManager TenantManager
}

func (s *SignupService) SignUp(ctx context.Context, tenantID, userID string) error {
	fmt.Println("Simulating signup for a new user...")
	err := s.tenantManager.CreateTenantWithAdmin(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to create tenant with admin: %w", err)
	}
	return nil
}

func main() {
	// In a real application, you would initialize the authorization service here.
	authSvc := &MockAuthService{}

	// The SignupService is initialized with the TenantManager implementation.
	signupSvc := &SignupService{
		tenantManager: authSvc,
	}

	// Simulate a signup event
	tenantID := "new-tenant-456"
	userID := "user-def-456"

	err := signupSvc.SignUp(context.Background(), tenantID, userID)
	if err != nil {
		log.Fatalf("Signup failed: %v", err)
	}

	fmt.Println("Signup simulation completed successfully!")
}
