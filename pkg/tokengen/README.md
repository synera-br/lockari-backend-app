# TokenGen Package

Package para geração e validação de tokens JWT seguros em Go.

## Características

- ✅ Interface limpa e fácil de usar
- ✅ Secret configurável por ambiente
- ✅ Tokens JWT com assinatura HMAC-SHA256
- ✅ Suporte a claims customizadas
- ✅ Validação de expiração automática
- ✅ Type-safe com structs bem definidas
- ✅ Tratamento de erros específicos
- ✅ **Tokens sem expiração (non-expiring)**
- ✅ **Tokens com expiração muito longa (long-lived)**

## Instalação

Este package faz parte do projeto Lockari e usa as dependências já presentes no `go.mod`.

## Uso Básico

### 1. Inicialização

```go
import (
    "os"
    "time"
    "github.com/synera-br/lockari-backend-app/pkg/tokengen"
)

// Secret vem de variável de ambiente
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    panic("JWT_SECRET environment variable is required")
}

// Criar o gerador de tokens
tokenGen := tokengen.NewTokenGenerator(
    secret,           // Secret do ambiente
    "lockari-app",    // Issuer da aplicação
    time.Hour,        // Duração padrão: 1 hora
)
```

### 2. Gerar Token

```go
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
    // Tratar erro
}
```

### 3. Validar Token

```go
claims, err := tokenGen.Validate(token)
if err != nil {
    // Token inválido ou expirado
}

// Usar claims validadas
userID := claims.UserID
tenantID := claims.TenantID
```

## Tokens Especiais

### 1. Token SEM Expiração

```go
// Para service accounts, APIs internas, etc.
claims := tokengen.TokenClaims{
    UserID:   "service-account-123",
    TenantID: "internal",
    Scope:    []string{"admin", "service"},
    Metadata: map[string]interface{}{
        "type": "service_account",
        "role": "admin",
    },
}

// Gera token que NUNCA expira
nonExpiringToken, err := tokenGen.GenerateNonExpiring(claims)
if err != nil {
    // Tratar erro
}

// Validar token não-expirável
validatedClaims, err := tokenGen.Validate(nonExpiringToken)
if err != nil {
    // Token inválido
}

// Verificar se é não-expirável
if validatedClaims.IsNonExpiring() {
    log.Println("Token não expira!")
}

// Verificar expiração (sempre retorna false para tokens não-expiráveis)
expired := tokenGen.IsExpired(nonExpiringToken) // false
```

### 2. Token com Expiração Muito Longa

```go
// Para tokens que precisam durar muito tempo (ex: 100 anos)
longLivedToken, err := tokenGen.GenerateLongLived(claims, time.Hour*24*365*100)
if err != nil {
    // Tratar erro
}

// Validar normalmente
validatedClaims, err := tokenGen.Validate(longLivedToken)
```

## Casos de Uso

### Tokens Não-Expiráveis
- 🔧 **Service Accounts**: Autenticação entre serviços
- 🤖 **APIs Internas**: Comunicação entre microserviços
- 📊 **Monitoramento**: Tokens para sistemas de observabilidade
- 🔄 **Integração**: Webhooks e sistemas externos

### Tokens de Longa Duração
- 📱 **Apps Mobile**: Tokens que duram meses/anos
- 🖥️ **Desktop Apps**: Aplicações que ficam logadas por muito tempo
- 🔗 **Refresh Tokens**: Tokens de renovação com TTL longo

## Middleware para Gin

```go
func TokenValidationMiddleware(tokenGen tokengen.TokenGenerator) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("X-Token")
        if token == "" {
            c.JSON(401, gin.H{"error": "Missing X-Token header"})
            c.Abort()
            return
        }
        
        claims, err := tokenGen.Validate(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Adicionar ao contexto
        c.Set("user_id", claims.UserID)
        c.Set("tenant_id", claims.TenantID)
        c.Set("scope", claims.Scope)
        c.Next()
    }
}
```

## Uso no Handler

```go
func (h *signupHandler) Create(c *gin.Context) {
    // Dados do middleware
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")
    
    // Verificar permissões
    scopes := c.GetStringSlice("scope")
    // ...
}
```

## Configuração por Ambiente

### Development
```bash
JWT_SECRET=dev-secret-key-very-long-and-secure-string
```

### Production  
```bash
JWT_SECRET=prod-ultra-secure-256-bit-secret-key-from-vault
```

### Testing
```bash
JWT_SECRET=test-secret-for-unit-tests-only
```

## Interface Completa

```go
type TokenGenerator interface {
    Generate(claims TokenClaims) (string, error)
    Validate(token string) (*TokenClaims, error)  
    Refresh(token string, newDuration time.Duration) (string, error)
    IsExpired(token string) bool
}
```

## Estrutura de Claims

```go
type TokenClaims struct {
    UserID    string                 `json:"user_id"`
    TenantID  string                 `json:"tenant_id,omitempty"`
    Scope     []string               `json:"scope,omitempty"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
    IssuedAt  time.Time              `json:"iat"`
    ExpiresAt time.Time              `json:"exp"`
}
```

## Tratamento de Erros

O package fornece erros específicos para diferentes cenários:

- `ErrInvalidToken`: Token malformado ou inválido
- `ErrExpiredToken`: Token expirado  
- `ErrMalformedToken`: Token com formato incorreto
- `ErrInvalidClaims`: Claims inválidas

Além disso, `TokenError` fornece códigos e mensagens específicas para debugging.

## Segurança

- ✅ Tokens assinados com HMAC-SHA256
- ✅ Secret configurável por ambiente
- ✅ Validação de expiração automática
- ✅ Verificação de assinatura em toda validação
- ✅ Claims estruturadas e type-safe

## Arquivos do Package

- `interface.go`: Interface principal e função de criação
- `types.go`: Structs de claims e configuração
- `errors.go`: Definições de erros customizados  
- `jwt_impl.go`: Implementação JWT
- `examples.go`: Exemplos de uso (pode ser removido)
- `README.md`: Esta documentação
