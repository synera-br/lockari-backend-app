# How to run the Decoupled Integration example

This example demonstrates how to use an interface to decouple the tenant creation logic from a service.

## Running the example

To run the example, execute the following command from the root of the project:

```bash
go run pkg/authorization/examples/tenant/proposal2/main.go
```

## Expected Output

```
Simulating signup for a new user...
Tenant 'new-tenant-456' created for user 'user-def-456'
User 'user-def-456' assigned as admin to tenant 'new-tenant-456'
Signup simulation completed successfully!
```
