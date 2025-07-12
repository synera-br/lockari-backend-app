# Proposal 1: Direct Integration

This proposal suggests integrating the OpenFGA tenant creation logic directly into the `SignupEventService`.

## Why this solution?

This approach is the simplest and most direct way to achieve the goal. It requires minimal changes to the existing code structure and is easy to understand. However, it does create a tighter coupling between the signup service and the authorization service.

## Implementation Steps

1.  **Extend `pkg/authorization/service.go`**:
    *   Add a new method `CreateTenantWithAdmin(ctx context.Context, tenantID, userID string) error`.
    *   This method will write two tuples to OpenFGA:
        1.  A tuple that defines the tenant: `(user:tenant:<tenantID>, relation:owner, user:user:<userID>)`
        2.  A tuple that assigns the user as an admin of the tenant: `(user:user:<userID>, relation:admin, object:tenant:<tenantID>)`

2.  **Update `internal/core/service/auth/signup.go`**:
    *   Add a new field `authzSvc *authorization.Service` to the `SignupEvent` struct.
    *   Update the `InitializeSignupEventService` function to accept and initialize the new field.
    *   In the `Create` method, after the `s.auth.SetTenantId` call, add a call to `s.authzSvc.CreateTenantWithAdmin(ctx, tenantId, signupData.User.Uid)`.

3.  **Update `cmd/main.go`**:
    *   When initializing the `SignupEventService`, pass the `authorization.Service` instance.

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
