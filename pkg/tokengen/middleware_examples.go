package tokengen

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenValidationMiddleware middleware que valida tokens incluindo non-expiring
func TokenValidationMiddleware(tokenGen TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Buscar token no header Authorization ou X-Token
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.GetHeader("X-Token")
		}

		// Remover "Bearer " se presente
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			c.JSON(401, gin.H{
				"error":   "Missing token",
				"message": "Authorization header or X-Token header is required",
			})
			c.Abort()
			return
		}

		// Validar token
		claims, err := tokenGen.Validate(token)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.JSON(401, gin.H{
				"error":   "Invalid token",
				"message": "Token validation failed",
			})
			c.Abort()
			return
		}

		// Log para tokens não-expiráveis (útil para auditoria)
		if claims.IsNonExpiring() {
			log.Printf("Non-expiring token used by user: %s, tenant: %s",
				claims.UserID, claims.TenantID)
		}

		// Adicionar claims ao contexto
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("token_scope", claims.Scope)
		c.Set("token_metadata", claims.Metadata)
		c.Set("token_claims", claims)
		c.Set("token_non_expiring", claims.IsNonExpiring())

		c.Next()
	}
}

// ServiceAccountMiddleware middleware específico para service accounts
func ServiceAccountMiddleware(tokenGen TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Usar o middleware de validação normal primeiro
		TokenValidationMiddleware(tokenGen)(c)

		// Se já foi abortado, não continuar
		if c.IsAborted() {
			return
		}

		// Verificar se é um service account
		claims := c.MustGet("token_claims").(*TokenClaims)

		serviceType, exists := claims.Metadata["service_type"]
		if !exists {
			c.JSON(403, gin.H{
				"error":   "Access denied",
				"message": "Service account required",
			})
			c.Abort()
			return
		}

		// Verificar se o token é não-expirável (recomendado para service accounts)
		if !claims.IsNonExpiring() {
			log.Printf("Warning: Service account '%s' using expiring token", claims.UserID)
		}

		// Adicionar tipo de serviço ao contexto
		c.Set("service_type", serviceType)

		c.Next()
	}
}

// RequireScopeMiddleware middleware que verifica se o token tem um scope específico
func RequireScopeMiddleware(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("token_claims").(*TokenClaims)

		hasScope := false
		for _, scope := range claims.Scope {
			if scope == requiredScope {
				hasScope = true
				break
			}
		}

		if !hasScope {
			c.JSON(403, gin.H{
				"error":   "Insufficient permissions",
				"message": "Required scope: " + requiredScope,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ExampleRouterSetup demonstra como configurar rotas com os middlewares
func ExampleRouterSetup() {
	/*
		// Exemplo de uso com Gin
		r := gin.Default()

		// Configurar tokengen
		tokenGen := NewTokenGenerator("your-secret", "your-app", time.Hour)

		// Rotas públicas
		r.POST("/auth/login", loginHandler)

		// Rotas que requerem autenticação
		authenticated := r.Group("/api/v1")
		authenticated.Use(TokenValidationMiddleware(tokenGen))
		{
			authenticated.GET("/profile", getProfileHandler)
			authenticated.POST("/data", createDataHandler)
		}

		// Rotas para service accounts
		serviceRoutes := r.Group("/api/v1/service")
		serviceRoutes.Use(ServiceAccountMiddleware(tokenGen))
		{
			serviceRoutes.GET("/metrics", getMetricsHandler)
			serviceRoutes.POST("/alerts", createAlertHandler)
		}

		// Rotas que requerem scope específico
		adminRoutes := r.Group("/api/v1/admin")
		adminRoutes.Use(TokenValidationMiddleware(tokenGen))
		adminRoutes.Use(RequireScopeMiddleware("admin"))
		{
			adminRoutes.GET("/users", listUsersHandler)
			adminRoutes.DELETE("/users/:id", deleteUserHandler)
		}

		r.Run(":8080")
	*/
}

// ExampleHandlerWithTokenInfo demonstra como usar as informações do token no handler
func ExampleHandlerWithTokenInfo() {
	/*
		func MyHandler(c *gin.Context) {
			// Obter informações do token
			userID := c.GetString("user_id")
			tenantID := c.GetString("tenant_id")
			scope := c.GetStringSlice("token_scope")
			metadata := c.GetStringMap("token_metadata")
			isNonExpiring := c.GetBool("token_non_expiring")

			// Usar as informações
			if isNonExpiring {
				log.Printf("Service account %s from tenant %s accessed endpoint", userID, tenantID)
			}

			// Verificar scope específico
			hasAdminScope := false
			for _, s := range scope {
				if s == "admin" {
					hasAdminScope = true
					break
				}
			}

			if hasAdminScope {
				// Processar como admin
			}

			// Usar metadata
			if deviceType, exists := metadata["device_type"]; exists {
				log.Printf("Request from device type: %s", deviceType)
			}

			c.JSON(200, gin.H{
				"message": "Success",
				"user_id": userID,
				"tenant_id": tenantID,
			})
		}
	*/
}
