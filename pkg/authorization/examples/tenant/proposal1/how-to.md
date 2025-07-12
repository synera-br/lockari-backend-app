# How to run the Direct Integration example

This example demonstrates how to directly integrate the tenant creation logic into a service.

## Running the example

To run the example, execute the following command from the root of the project:

```bash
go run pkg/authorization/examples/tenant/proposal1/main.go
```

## Expected Output

```
Simulating signup for a new user...
Tenant 'new-tenant-123' created for user 'user-abc-789'
User 'user-abc-789' assigned as admin to tenant 'new-tenant-123'
Signup simulation completed successfully!
```
