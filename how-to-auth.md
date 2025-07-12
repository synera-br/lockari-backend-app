# How to Use the Auth Service and Generic Handler

This document explains how to use the new `auth` service and the generic handler.

## Auth Service

The `auth` service is responsible for handling the business logic for login and signup. It is independent of the `gin` framework, so you can use it with any framework or in any context.

### Creating a New Auth Service

To create a new `auth` service, you need to create a new `Service` struct and then call the `New` function:

```go
package main

import (
	"log"

	"github.com/synera-br/lockari-backend-app/internal/core/service/auth"
	"github.com/synera-br/lockari-backend-app/pkg/logger"
)

func main() {
	log, err := logger.New(&logger.Config{})
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	authService := auth.New(
		// loginService
		// signupService
		// encryptor
		// authClient
		// token
		log,
	)
}
```

### Using the Auth Service

The `auth` service has two methods: `Login` and `Signup`.

The `Login` method performs a login. It takes a context and a payload as input. The payload is a base64-encoded string that contains the encrypted login event.

The `Signup` method performs a signup. It takes a context and a payload as input. The payload is a base64-encoded string that contains the encrypted signup event.

## Generic Handler

The generic handler is a wrapper around the `gin` framework. It provides a modular and configurable way to handle HTTP requests and responses.

### Creating a New Generic Handler

To create a new generic handler, you need to create a new `Handler` struct and then call the `New` function:

```go
package main

import (
	"log"

	"github.com/synera-br/lockari-backend-app/internal/handler/web"
	"github.com/synera-br/lockari-backend-app/pkg/logger"
)

func main() {
	log, err := logger.New(&logger.Config{})
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	handler := web.New(log)
}
```

### Using the Generic Handler

The generic handler has a `Handle` method that takes a handler function as input. The handler function takes a payload as input and returns a data object and an error.

The `Handle` method returns a `gin.HandlerFunc` that can be used to register a route.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/synera-br/lockari-backend-app/internal/core/service/auth"
	webhandler "github.com/synera-br/lockari-backend-app/internal/handler/web/auth"
)

func main() {
	// ...

	loginHandler := webhandler.NewLoginHandler(authService, handler)
	signupHandler := webhandler.NewSignupHandler(authService, handler)

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		loginHandler.RegisterRoutes(v1)
		signupHandler.RegisterRoutes(v1)
	}

	router.Run(":8080")
}
```
