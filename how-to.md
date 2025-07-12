# How to Use the HTTP Server and Middleware

This document explains how to use the new HTTP server and middleware.

## HTTP Server

The HTTP server is a wrapper around the `gin` framework. It provides a modular and configurable way to create a new HTTP server.

### Creating a New Server

To create a new server, you need to create a new `Config` struct and then call the `New` function:

```go
package main

import (
	"context"
	"log"

	"github.com/synera-br/lockari-backend-app/pkg/http_server"
)

func main() {
	config := &httpserver.Config{
		Port:            8080,
		Mode:            "debug",
		OtelServiceName: "my-app",
		Cors: httpserver.CorsConfig{
			AllowedOrigins: "*",
		},
	}

	server, err := httpserver.New(config)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := server.Run(context.Background()); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
```

### Adding Routes

You can add routes to the server using the `AddRoute` and `AddGroup` methods:

```go
// Add a single route
server.AddRoute("GET", "/ping", func(c *gin.Context) {
	c.String(200, "pong")
})

// Add a group of routes
v1 := server.AddGroup("/v1")
{
	v1.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
```

### Swagger

The server automatically adds a `/swagger/*any` route that serves the Swagger UI. You can customize the Swagger information by modifying the `SwaggerInfo` variable in the `pkg/http_server/types.go` file.

## Middleware

The middleware package provides a way to add middleware to your HTTP handlers. The middleware is independent of the `gin` framework, so you can use it with any framework that supports `http.Handler`.

### Creating a New Middleware Chain

To create a new middleware chain, you can use the `Chain` function:

```go
package main

import (
	"net/http"

	"github.com/synera-br/lockari-backend-app/internal/handler/middleware"
)

func main() {
	chain := middleware.Chain(
		middleware.Auth(auth, log),
		middleware.AuthJWT(token, log),
	)

	http.Handle("/", chain(myHandler))
}
```

### Authentication Middleware

The `Auth` and `AuthJWT` middleware can be used to validate tokens and add user information to the context.

The `Auth` middleware validates a Firebase token and adds the user ID, tenant ID, and claims to the context.

The `AuthJWT` middleware validates a JWT and adds the claims to the context.

### Accessing User Information

You can access the user information from the context using the following keys:

*   `middleware.UserIDKey`: The user ID.
*   `middleware.TenantIDKey`: The tenant ID.
*   `middleware.ClaimsKey`: The claims.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/synera-br/lockari-backend-app/internal/handler/middleware"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	claims := r.Context().Value(middleware.ClaimsKey).(map[string]interface{})

	fmt.Fprintf(w, "User ID: %s\n", userID)
	fmt.Fprintf(w, "Tenant ID: %s\n", tenantID)
	fmt.Fprintf(w, "Claims: %v\n", claims)
}
```
