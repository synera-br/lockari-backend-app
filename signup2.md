# Proposal 2: Decoupled Integration using an Interface

This proposal suggests creating a new interface for tenant management in the `pkg/authorization` package and using that interface in the `SignupEventService`.

## Why this solution?

This approach promotes loose coupling and better separation of concerns. By depending on an interface instead of a concrete implementation, the `SignupEventService` becomes more modular and easier to test in isolation. This is a more robust and maintainable solution in the long run.

## Implementation Steps

1.  **Create a new interface in `pkg/authorization/interfaces.go`**:
    *   Create a `TenantManager` interface:
        ```go
        type TenantManager interface {
            CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error
        }
        ```

2.  **Implement the interface in `pkg/authorization/service.go`**:
    *   The `Service` struct will implement the `TenantManager` interface.
    *   The implementation of `CreateTenantWithAdmin` will be the same as in Proposal 1.

3.  **Update `internal/core/service/auth/signup.go`**:
    *   Add a new field `tenantManager authorization.TenantManager` to the `SignupEvent` struct.
    *   Update the `InitializeSignupEventService` function to accept and initialize the new field.
    *   In the `Create` method, after the `s.auth.SetTenantId` call, add a call to `s.tenantManager.CreateTenantWithAdmin(ctx, tenantId, signupData.User.Uid)`.

4.  **Update `cmd/main.go`**:
    *   When initializing the `SignupEventService`, pass the `authorization.Service` instance (which implements the `TenantManager` interface).

## Code Example (`pkg/authorization/interfaces.go`)

```go
package authorization

import "context"

// TenantManager defines the interface for managing tenants.
type TenantManager interface {
    CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error
}
```

## Code Example (`pkg/authorization/service.go`)

```go
// CreateTenantWithAdmin creates a new tenant and assigns an admin to it.
func (s *Service) CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error {
	writeReq := &WriteRequest{
		Tuples: []Tuple{
			{
				User:     FormatUser(userID),
				Relation: "owner",
				Object:   FormatTenant(tenantID),
			},
			{
				User:     FormatUser(userID),
				Relation: "admin",
				Object:   FormatTenant(tenantID),
			},
		},
	}

	err := s.writeRelationship(ctx, writeReq)
	if err != nil {
		return fmt.Errorf("failed to create tenant with admin: %w", err)
	}

	return nil
}
```
