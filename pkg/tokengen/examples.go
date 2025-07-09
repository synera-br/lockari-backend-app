package tokengen

// Este arquivo contém exemplos de como usar o package tokengen
// Para uso real, remova este arquivo ou renomeie para example_test.go

/*
## Exemplo de Uso Básico:

```go
package main

import (
    "os"
    "time"
    "your-app/pkg/tokengen"
)

func main() {
    // 1. Inicialização com secret do ambiente
    secret := os.Getenv("JWT_SECRET") // Diferente por ambiente
    if secret == "" {
        panic("JWT_SECRET environment variable is required")
    }

    // 2. Criar o gerador de tokens
    tokenGen := tokengen.NewTokenGenerator(
        secret,                    // Secret do ambiente
        "lockari-app",            // Issuer da aplicação
        time.Hour,                // Duração padrão: 1 hora
    )

    // 3. Gerar token para um usuário
    claims := tokengen.TokenClaims{
        UserID:   "user123",
        TenantID: "tenant456",
        Scope:    []string{"read", "write", "signup"},
        Metadata: map[string]interface{}{
            "role": "user",
            "plan": "free",
        },
    }

    token, err := tokenGen.Generate(claims)
    if err != nil {
        panic(err)
    }

    // 4. Validar token
    validatedClaims, err := tokenGen.Validate(token)
    if err != nil {
        panic(err)
    }

    println("Token válido para usuário:", validatedClaims.UserID)
}
```

## Middleware para Gin:

```go
func TokenValidationMiddleware(tokenGen tokengen.TokenGenerator) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Buscar token no header X-Token
        token := c.GetHeader("X-Token")
        if token == "" {
            c.JSON(401, gin.H{"error": "Missing X-Token header"})
            c.Abort()
            return
        }

        // Validar token
        claims, err := tokenGen.Validate(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
            c.Abort()
            return
        }

        // Adicionar claims ao contexto
        c.Set("user_id", claims.UserID)
        c.Set("tenant_id", claims.TenantID)
        c.Set("scope", claims.Scope)
        c.Set("token_claims", claims)

        c.Next()
    }
}
```

## Configuração por Ambiente:

### Development (.env.dev):
```
JWT_SECRET=dev-secret-key-very-long-and-secure-string
```

### Production (.env.prod):
```
JWT_SECRET=prod-ultra-secure-256-bit-secret-key-from-vault
```

### Testing (.env.test):
```
JWT_SECRET=test-secret-for-unit-tests-only
```

## Uso no Handler:

```go
func (h *signupHandler) Create(c *gin.Context) {
    // Dados já validados pelo middleware
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")
    claims := c.MustGet("token_claims").(*tokengen.TokenClaims)

    // Verificar se tem permissão de signup
    hasScope := false
    for _, scope := range claims.Scope {
        if scope == "signup" {
            hasScope = true
            break
        }
    }

    if !hasScope {
        c.JSON(403, gin.H{"error": "Insufficient permissions"})
        return
    }

    // Processar signup...
    c.JSON(200, gin.H{"message": "Signup created successfully"})
}
```

## Dependências Necessárias:

Para usar este package, adicione ao go.mod:

```
go get github.com/golang-jwt/jwt/v5
```

*/
