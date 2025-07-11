package authorization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// === EXAMPLE USAGE ===

// ExampleBasicSetup demonstra como configurar o sistema de autorização
func ExampleBasicSetup() {
	// 1. Configurar o OpenFGA
	config := &Config{
		APIURL:               "https://api.fga.example.com",
		StoreID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		AuthorizationModelID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		ClientID:             "your-client-id",
		ClientSecret:         "your-client-secret",
		APITokenIssuer:       "https://auth.example.com",
		APIAudience:          "https://api.fga.example.com",

		// Cache configuration
		CacheEnabled: true,
		CacheTTL:     5 * time.Minute,
		CacheMaxSize: 1000,

		// Performance settings
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		RetryDelay:    1 * time.Second,
		BatchSize:     10,

		// Audit configuration
		AuditEnabled: true,
		AuditLevel:   "info",
	}

	// 2. Criar logger simples
	logger := &SimpleLogger{}

	// 3. Criar cliente OpenFGA
	client, err := NewOpenFGAClient(ClientOptions{
		Config: config,
		Logger: logger,
	})
	if err != nil {
		log.Fatal("Failed to create OpenFGA client:", err)
	}

	// 4. Criar cache
	cache := NewMemoryCache(CacheOptions{
		MaxSize:         1000,
		CleanupInterval: 5 * time.Minute,
	})

	// 5. Criar serviço de auditoria
	audit := NewAuditLogger(logger)

	// 6. Criar serviço principal
	service := NewService(ServiceOptions{
		Client: client,
		Config: config,
		Logger: logger,
		Audit:  audit,
		Cache:  cache,
	})

	// 7. Criar middleware
	middleware := NewAuthorizationMiddleware(MiddlewareOptions{
		Client:       client,
		Logger:       logger,
		AuditService: audit,
	})

	// 8. Configurar rotas Gin
	router := gin.Default()

	// Middleware para extrair user_id do JWT (implementar conforme necessário)
	router.Use(ExtractUserFromJWT())
	router.Use(ExtractTenantFromJWT())

	// Rotas protegidas
	vaultRoutes := router.Group("/api/v1/vaults")
	{
		// Listar vaults (requer permissão de visualização)
		vaultRoutes.GET("", middleware.RequireAnyVaultPermission(VaultPermissionView))

		// Obter vault específico
		vaultRoutes.GET("/:vaultId", middleware.RequireVaultPermission(VaultPermissionView))

		// Criar vault (requer papel de membro no tenant)
		vaultRoutes.POST("", middleware.RequireTenantRole(TenantRoleMember))

		// Atualizar vault
		vaultRoutes.PUT("/:vaultId", middleware.RequireVaultPermission(VaultPermissionWrite))

		// Deletar vault
		vaultRoutes.DELETE("/:vaultId", middleware.RequireVaultPermission(VaultPermissionDelete))

		// Compartilhar vault
		vaultRoutes.POST("/:vaultId/share", middleware.RequireVaultPermission(VaultPermissionShare))

		// Gerenciar permissões
		vaultRoutes.POST("/:vaultId/permissions", middleware.RequireVaultPermission(VaultPermissionManage))
	}

	// Rotas de secrets
	secretRoutes := router.Group("/api/v1/secrets")
	{
		// Listar secrets
		secretRoutes.GET("", middleware.RequireAnyVaultPermission(VaultPermissionView))

		// Obter secret
		secretRoutes.GET("/:secretId", middleware.RequireSecretPermission(SecretPermissionView))

		// Ler conteúdo do secret
		secretRoutes.GET("/:secretId/content", middleware.RequireSecretPermission(SecretPermissionRead))

		// Criar secret
		secretRoutes.POST("", middleware.RequireVaultPermission(VaultPermissionWrite))

		// Atualizar secret
		secretRoutes.PUT("/:secretId", middleware.RequireSecretPermission(SecretPermissionWrite))

		// Deletar secret
		secretRoutes.DELETE("/:secretId", middleware.RequireSecretPermission(SecretPermissionDelete))

		// Copiar secret
		secretRoutes.POST("/:secretId/copy", middleware.RequireSecretPermission(SecretPermissionCopy))

		// Download secret
		secretRoutes.GET("/:secretId/download", middleware.RequireSecretPermission(SecretPermissionDownload))
	}

	// Rotas administrativas
	adminRoutes := router.Group("/api/v1/admin")
	{
		// Apenas admins podem acessar
		adminRoutes.Use(middleware.RequireTenantRole(TenantRoleAdmin))

		adminRoutes.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin users endpoint"})
		})

		adminRoutes.GET("/audit", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin audit endpoint"})
		})
	}

	// Exemplo de uso do serviço diretamente
	ExampleServiceUsage(service)

	fmt.Println("Authorization system setup complete!")
}

// ExampleServiceUsage demonstra como usar o serviço de autorização diretamente
func ExampleServiceUsage(service *Service) {
	ctx := context.Background()

	// Exemplo 1: Verificar permissão de vault
	userID := "user123"
	vaultID := "vault456"

	canView, err := service.CanAccessVault(ctx, userID, vaultID, VaultPermissionView)
	if err != nil {
		log.Printf("Error checking vault permission: %v", err)
		return
	}

	if canView {
		fmt.Printf("User %s can view vault %s\n", userID, vaultID)
	} else {
		fmt.Printf("User %s cannot view vault %s\n", userID, vaultID)
	}

	// Exemplo 2: Verificar permissão de secret
	secretID := "secret789"

	canRead, err := service.CanAccessSecret(ctx, userID, secretID, SecretPermissionRead)
	if err != nil {
		log.Printf("Error checking secret permission: %v", err)
		return
	}

	if canRead {
		fmt.Printf("User %s can read secret %s\n", userID, secretID)
	} else {
		fmt.Printf("User %s cannot read secret %s\n", userID, secretID)
	}

	// Exemplo 3: Verificar papel no tenant
	tenantID := "tenant123"

	isAdmin, err := service.HasTenantRole(ctx, userID, tenantID, TenantRoleAdmin)
	if err != nil {
		log.Printf("Error checking tenant role: %v", err)
		return
	}

	if isAdmin {
		fmt.Printf("User %s is admin of tenant %s\n", userID, tenantID)
	} else {
		fmt.Printf("User %s is not admin of tenant %s\n", userID, tenantID)
	}

	// Exemplo 4: Listar vaults do usuário
	userVaults, err := service.ListUserVaults(ctx, userID)
	if err != nil {
		log.Printf("Error listing user vaults: %v", err)
		return
	}

	fmt.Printf("User %s has access to %d vaults: %v\n", userID, len(userVaults), userVaults)

	// Exemplo 5: Verificação em lote
	batchRequests := []*CheckRequest{
		{
			User:     FormatUser(userID),
			Relation: string(VaultPermissionView),
			Object:   FormatVault("vault1"),
		},
		{
			User:     FormatUser(userID),
			Relation: string(VaultPermissionWrite),
			Object:   FormatVault("vault2"),
		},
		{
			User:     FormatUser(userID),
			Relation: string(SecretPermissionRead),
			Object:   FormatSecret("secret1"),
		},
	}

	batchResults, err := service.CheckBatch(ctx, batchRequests)
	if err != nil {
		log.Printf("Error in batch check: %v", err)
		return
	}

	fmt.Printf("Batch check completed in %v\n", batchResults.Duration)
	for i, result := range batchResults.Results {
		if result != nil {
			fmt.Printf("Request %d: %v\n", i, result.Allowed)
		}
	}
}

// ExampleAdvancedMiddleware demonstra uso avançado do middleware
func ExampleAdvancedMiddleware() {
	router := gin.Default()

	// Criar middleware
	middleware := NewAuthorizationMiddleware(MiddlewareOptions{
		// ... configuração ...
	})

	// Exemplo 1: Middleware customizado
	router.GET("/api/v1/custom/:resourceId",
		middleware.RequireCustomPermission(
			"{user_id}",              // Placeholder será substituído pelo user_id do contexto
			"can_access",             // Relação fixa
			"resource:{resource_id}", // Placeholder será substituído pelo parâmetro da URL
		),
		func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Custom authorization passed"})
		},
	)

	// Exemplo 2: Múltiplas permissões necessárias
	router.POST("/api/v1/vaults/:vaultId/secrets",
		middleware.RequireAllVaultPermissions(
			VaultPermissionView,
			VaultPermissionWrite,
		),
		func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "All permissions verified"})
		},
	)

	// Exemplo 3: Qualquer uma das permissões
	router.GET("/api/v1/vaults/:vaultId/metadata",
		middleware.RequireAnyVaultPermission(
			VaultPermissionView,
			VaultPermissionRead,
			VaultPermissionManage,
		),
		func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Any permission verified"})
		},
	)

	// Exemplo 4: Verificação em lote
	batchPermissions := []PermissionCheck{
		{
			User:     "{user_id}",
			Relation: string(VaultPermissionView),
			Object:   "vault:{vault_id}",
		},
		{
			User:     "{user_id}",
			Relation: string(TenantRoleMember),
			Object:   "tenant:{tenant_id}",
		},
	}

	router.POST("/api/v1/complex-operation",
		middleware.BatchRequirePermissions(batchPermissions),
		func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Batch permissions verified"})
		},
	)

	fmt.Println("Advanced middleware setup complete!")
}

// ExampleCacheUsage demonstra uso do cache
func ExampleCacheUsage() {
	// Criar cache
	cache := NewAuthorizationCache(CacheOptions{
		MaxSize:         1000,
		CleanupInterval: 5 * time.Minute,
	})

	// Exemplo 1: Usar cache diretamente
	user := "user:123"
	relation := "can_view"
	object := "vault:456"

	// Verificar se está no cache
	if allowed, found := cache.GetPermission(user, relation, object); found {
		fmt.Printf("Cache hit: %v\n", allowed)
	} else {
		fmt.Println("Cache miss")

		// Simular verificação no OpenFGA
		allowed := true // Resultado da verificação

		// Armazenar no cache
		cache.SetPermission(user, relation, object, allowed, 5*time.Minute)
	}

	// Exemplo 2: Invalidar cache
	cache.InvalidateUser("user:123")    // Invalida todas as permissões do usuário
	cache.InvalidateObject("vault:456") // Invalida todas as permissões do objeto

	// Exemplo 3: Estatísticas do cache
	stats := cache.Stats()
	fmt.Printf("Cache stats: Hits=%d, Misses=%d, Hit Rate=%.2f%%\n",
		stats.Hits, stats.Misses, stats.HitRate*100)

	// Limpeza
	cache.Close()
}

// ExampleAuditUsage demonstra uso do sistema de auditoria
func ExampleAuditUsage() {
	// Criar logger
	logger := &SimpleLogger{}

	// Criar serviço de auditoria
	audit := NewAuditLogger(logger)

	ctx := context.Background()

	// Exemplo 1: Registrar verificação de permissão
	checkEvent := PermissionCheckEvent{
		User:      "user:123",
		Relation:  "can_view",
		Object:    "vault:456",
		Result:    "allowed",
		Timestamp: time.Now(),
		Duration:  50 * time.Millisecond,
	}
	audit.LogPermissionCheck(ctx, checkEvent)

	// Exemplo 2: Registrar concessão de permissão
	grantEvent := PermissionGrantEvent{
		Grantor:   "user:admin",
		Grantee:   "user:123",
		Relation:  "can_write",
		Object:    "vault:456",
		Timestamp: time.Now(),
	}
	audit.LogPermissionGrant(ctx, grantEvent)

	// Exemplo 3: Registrar atividade suspeita
	suspiciousEvent := SuspiciousActivityEvent{
		User:      "user:123",
		Activity:  "multiple_failed_access_attempts",
		Details:   "5 failed attempts in 1 minute",
		Severity:  "high",
		ClientIP:  "192.168.1.100",
		UserAgent: "curl/7.68.0",
		Timestamp: time.Now(),
	}
	audit.LogSuspiciousActivity(ctx, suspiciousEvent)

	fmt.Println("Audit usage examples complete!")
}

// === HELPER IMPLEMENTATIONS ===

// SimpleLogger é uma implementação simples de logger para exemplos
type SimpleLogger struct{}

func (sl *SimpleLogger) Debug(msg string, fields ...interface{}) {
	fmt.Printf("DEBUG: %s - %+v\n", msg, fields)
}

func (sl *SimpleLogger) Info(msg string, fields ...interface{}) {
	fmt.Printf("INFO: %s - %+v\n", msg, fields)
}

func (sl *SimpleLogger) Warn(msg string, fields ...interface{}) {
	fmt.Printf("WARN: %s - %+v\n", msg, fields)
}

func (sl *SimpleLogger) Error(msg string, fields ...interface{}) {
	fmt.Printf("ERROR: %s - %+v\n", msg, fields)
}

func (sl *SimpleLogger) With(fields ...interface{}) Logger {
	return sl // Implementação simples
}

// === INTEGRATION EXAMPLES ===

// ExampleIntegrationWithExistingHandlers demonstra integração com handlers existentes
func ExampleIntegrationWithExistingHandlers() {
	// Assumindo que você já tem handlers existentes
	router := gin.Default()

	// Middleware de autorização
	middleware := NewAuthorizationMiddleware(MiddlewareOptions{
		// ... configuração ...
	})

	// Exemplo: Integrar com handler existente de vault
	router.GET("/api/v1/vaults/:vaultId",
		middleware.RequireVaultPermission(VaultPermissionView),
		func(c *gin.Context) {
			// Seu handler existente aqui
			vaultID := c.Param("vaultId")

			// Obter informações de autorização do contexto
			authorizedVaultID := c.GetString("authorized_vault_id")
			authorizedPermission := c.GetString("authorized_permission")

			// Verificar se os IDs coincidem (paranoia)
			if authorizedVaultID != vaultID {
				c.JSON(500, gin.H{"error": "authorization_mismatch"})
				return
			}

			// Seu código de negócio aqui
			c.JSON(200, gin.H{
				"vault_id":   vaultID,
				"permission": authorizedPermission,
				"message":    "Access granted",
			})
		},
	)

	fmt.Println("Integration examples complete!")
}

// === TESTING EXAMPLES ===

// ExampleTestingAuthorization demonstra como testar o sistema de autorização
func ExampleTestingAuthorization() {
	// Criar configuração de teste
	config := &Config{
		APIURL:               "https://test.fga.example.com",
		StoreID:              "test-store",
		AuthorizationModelID: "test-model",
		CacheEnabled:         false, // Desabilitar cache para testes
		AuditEnabled:         true,
	}

	// Criar mock logger
	logger := &SimpleLogger{}

	// Criar cliente (em testes reais, você usaria mocks)
	client, err := NewOpenFGAClient(ClientOptions{
		Config: config,
		Logger: logger,
	})
	if err != nil {
		log.Fatal("Failed to create test client:", err)
	}

	// Criar serviço de teste
	service := NewService(ServiceOptions{
		Client: client,
		Config: config,
		Logger: logger,
	})

	// Exemplo de teste
	ctx := context.Background()

	// Verificar permissão que deve ser permitida
	allowed, err := service.CanAccessVault(ctx, "user123", "vault456", VaultPermissionView)
	if err != nil {
		log.Printf("Test error: %v", err)
		return
	}

	if allowed {
		fmt.Println("Test passed: User can access vault")
	} else {
		fmt.Println("Test failed: User should be able to access vault")
	}

	fmt.Println("Testing examples complete!")
}

// === PERFORMANCE EXAMPLES ===

// ExamplePerformanceOptimization demonstra otimizações de performance
func ExamplePerformanceOptimization() {
	// Criar cache com limpeza frequente
	cache := NewAuthorizationCache(CacheOptions{
		MaxSize:         10000,
		CleanupInterval: 1 * time.Minute,
	})

	// Exemplo de uso de cache warming
	warmer := NewCacheWarmer(cache, nil)

	// Pré-carregar permissões para usuários ativos
	activeUsers := []string{"user1", "user2", "user3"}
	commonVaults := []string{"vault1", "vault2", "vault3"}

	for _, userID := range activeUsers {
		err := warmer.WarmUserPermissions(userID, commonVaults)
		if err != nil {
			log.Printf("Cache warming error: %v", err)
		}
	}

	fmt.Println("Performance optimization examples complete!")
}

// === MONITORING EXAMPLES ===

// ExampleMonitoring demonstra monitoramento do sistema
func ExampleMonitoring() {
	// Criar métricas
	auditMetrics := NewAuditMetrics()

	// Simular alguns eventos
	auditMetrics.IncrementPermissionChecks()
	auditMetrics.IncrementPermissionGrants()
	auditMetrics.IncrementSuspiciousEvents()

	// Obter métricas
	metrics := auditMetrics.GetMetrics()
	fmt.Printf("Audit metrics: %+v\n", metrics)

	// Exemplo de cache metrics
	cache := NewAuthorizationCache(CacheOptions{
		MaxSize:         1000,
		CleanupInterval: 5 * time.Minute,
	})

	cacheMetrics := NewCacheMetrics(cache)
	cacheStats := cacheMetrics.CollectMetrics()
	fmt.Printf("Cache metrics: %+v\n", cacheStats)

	fmt.Println("Monitoring examples complete!")
}
