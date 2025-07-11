package authorization

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthorizationMiddleware é o middleware para autorização usando OpenFGA
type AuthorizationMiddleware struct {
	client       *OpenFGAClient
	logger       Logger
	auditService AuditService
}

// MiddlewareOptions define as opções para o middleware
type MiddlewareOptions struct {
	Client       *OpenFGAClient
	Logger       Logger
	AuditService AuditService
}

// NewAuthorizationMiddleware cria um novo middleware de autorização
func NewAuthorizationMiddleware(opts MiddlewareOptions) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		client:       opts.Client,
		logger:       opts.Logger,
		auditService: opts.AuditService,
	}
}

// RequireVaultPermission verifica se o usuário tem permissão específica no vault
func (m *AuthorizationMiddleware) RequireVaultPermission(permission VaultPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extrair dados do contexto
		userID := c.GetString("user_id")
		vaultID := c.Param("vaultId")

		// Validar dados necessários
		if userID == "" {
			m.logError(c, "user_id not found in context", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_authenticated"})
			c.Abort()
			return
		}

		if vaultID == "" {
			m.logError(c, "vaultId not found in URL parameters", nil)
			c.JSON(http.StatusBadRequest, gin.H{"error": "vault_id_required"})
			c.Abort()
			return
		}

		// Formatar para OpenFGA
		user := FormatUser(userID)
		object := FormatVault(vaultID)
		relation := string(permission)

		// Realizar verificação de autorização
		checkReq := &CheckRequest{
			User:     user,
			Relation: relation,
			Object:   object,
		}

		allowed, err := m.checkPermission(c.Request.Context(), checkReq)
		if err != nil {
			m.logError(c, "Authorization check failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization_check_failed"})
			c.Abort()
			return
		}

		// Auditoria
		m.auditPermissionCheck(c, checkReq, allowed, time.Since(start))

		if !allowed {
			m.logUnauthorized(c, user, relation, object)
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
			c.Abort()
			return
		}

		// Adicionar informações de autorização ao contexto
		c.Set("authorized_vault_id", vaultID)
		c.Set("authorized_permission", string(permission))
		c.Set("authorization_duration", time.Since(start))

		c.Next()
	}
}

// RequireSecretPermission verifica se o usuário tem permissão específica no secret
func (m *AuthorizationMiddleware) RequireSecretPermission(permission SecretPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extrair dados do contexto
		userID := c.GetString("user_id")
		secretID := c.Param("secretId")

		// Validar dados necessários
		if userID == "" {
			m.logError(c, "user_id not found in context", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_authenticated"})
			c.Abort()
			return
		}

		if secretID == "" {
			m.logError(c, "secretId not found in URL parameters", nil)
			c.JSON(http.StatusBadRequest, gin.H{"error": "secret_id_required"})
			c.Abort()
			return
		}

		// Formatar para OpenFGA
		user := FormatUser(userID)
		object := FormatSecret(secretID)
		relation := string(permission)

		// Realizar verificação de autorização
		checkReq := &CheckRequest{
			User:     user,
			Relation: relation,
			Object:   object,
		}

		allowed, err := m.checkPermission(c.Request.Context(), checkReq)
		if err != nil {
			m.logError(c, "Authorization check failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization_check_failed"})
			c.Abort()
			return
		}

		// Auditoria
		m.auditPermissionCheck(c, checkReq, allowed, time.Since(start))

		if !allowed {
			m.logUnauthorized(c, user, relation, object)
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
			c.Abort()
			return
		}

		// Adicionar informações de autorização ao contexto
		c.Set("authorized_secret_id", secretID)
		c.Set("authorized_permission", string(permission))
		c.Set("authorization_duration", time.Since(start))

		c.Next()
	}
}

// RequireTenantRole verifica se o usuário tem o papel específico no tenant
func (m *AuthorizationMiddleware) RequireTenantRole(role TenantRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extrair dados do contexto
		userID := c.GetString("user_id")
		tenantID := c.GetString("tenant_id")

		// Validar dados necessários
		if userID == "" {
			m.logError(c, "user_id not found in context", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_authenticated"})
			c.Abort()
			return
		}

		if tenantID == "" {
			m.logError(c, "tenant_id not found in context", nil)
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id_required"})
			c.Abort()
			return
		}

		// Formatar para OpenFGA
		user := FormatUser(userID)
		object := FormatTenant(tenantID)
		relation := string(role)

		// Realizar verificação de autorização
		checkReq := &CheckRequest{
			User:     user,
			Relation: relation,
			Object:   object,
		}

		allowed, err := m.checkPermission(c.Request.Context(), checkReq)
		if err != nil {
			m.logError(c, "Authorization check failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization_check_failed"})
			c.Abort()
			return
		}

		// Auditoria
		m.auditPermissionCheck(c, checkReq, allowed, time.Since(start))

		if !allowed {
			m.logUnauthorized(c, user, relation, object)
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
			c.Abort()
			return
		}

		// Adicionar informações de autorização ao contexto
		c.Set("authorized_tenant_id", tenantID)
		c.Set("authorized_role", string(role))
		c.Set("authorization_duration", time.Since(start))

		c.Next()
	}
}

// RequireCustomPermission verifica uma permissão customizada
func (m *AuthorizationMiddleware) RequireCustomPermission(user, relation, object string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Substituir placeholders pelos valores reais
		actualUser := m.resolvePlaceholder(c, user)
		actualObject := m.resolvePlaceholder(c, object)

		// Validar dados necessários
		if actualUser == "" {
			m.logError(c, "user could not be resolved", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_authenticated"})
			c.Abort()
			return
		}

		if actualObject == "" {
			m.logError(c, "object could not be resolved", nil)
			c.JSON(http.StatusBadRequest, gin.H{"error": "object_not_found"})
			c.Abort()
			return
		}

		// Realizar verificação de autorização
		checkReq := &CheckRequest{
			User:     actualUser,
			Relation: relation,
			Object:   actualObject,
		}

		allowed, err := m.checkPermission(c.Request.Context(), checkReq)
		if err != nil {
			m.logError(c, "Authorization check failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization_check_failed"})
			c.Abort()
			return
		}

		// Auditoria
		m.auditPermissionCheck(c, checkReq, allowed, time.Since(start))

		if !allowed {
			m.logUnauthorized(c, actualUser, relation, actualObject)
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
			c.Abort()
			return
		}

		// Adicionar informações de autorização ao contexto
		c.Set("authorized_user", actualUser)
		c.Set("authorized_relation", relation)
		c.Set("authorized_object", actualObject)
		c.Set("authorization_duration", time.Since(start))

		c.Next()
	}
}

// BatchRequirePermissions verifica múltiplas permissões em uma só operação
func (m *AuthorizationMiddleware) BatchRequirePermissions(permissions []PermissionCheck) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var requests []*CheckRequest

		// Preparar todas as verificações
		for _, perm := range permissions {
			actualUser := m.resolvePlaceholder(c, perm.User)
			actualObject := m.resolvePlaceholder(c, perm.Object)

			if actualUser == "" || actualObject == "" {
				m.logError(c, "Required parameters not found", nil)
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing_required_parameters"})
				c.Abort()
				return
			}

			requests = append(requests, &CheckRequest{
				User:     actualUser,
				Relation: perm.Relation,
				Object:   actualObject,
			})
		}

		// Realizar verificação em lote
		results, err := m.client.BatchCheck(c.Request.Context(), requests)
		if err != nil {
			m.logError(c, "Batch authorization check failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization_check_failed"})
			c.Abort()
			return
		}

		// Verificar se todas as permissões foram concedidas
		for i, result := range results.Results {
			if result == nil || !result.Allowed {
				req := requests[i]
				m.logUnauthorized(c, req.User, req.Relation, req.Object)
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
				c.Abort()
				return
			}
		}

		// Auditoria para todas as verificações
		for i, req := range requests {
			allowed := results.Results[i] != nil && results.Results[i].Allowed
			m.auditPermissionCheck(c, req, allowed, time.Since(start))
		}

		// Adicionar informações de autorização ao contexto
		c.Set("batch_permissions_checked", len(requests))
		c.Set("authorization_duration", time.Since(start))

		c.Next()
	}
}

// checkPermission realiza a verificação de permissão
func (m *AuthorizationMiddleware) checkPermission(ctx context.Context, req *CheckRequest) (bool, error) {
	response, err := m.client.Check(ctx, req)
	if err != nil {
		return false, err
	}

	return response.Allowed, nil
}

// resolvePlaceholder resolve placeholders no formato {field} para valores reais
func (m *AuthorizationMiddleware) resolvePlaceholder(c *gin.Context, placeholder string) string {
	if !strings.Contains(placeholder, "{") {
		return placeholder
	}

	// Substituir placeholders comuns
	placeholder = strings.ReplaceAll(placeholder, "{user_id}", c.GetString("user_id"))
	placeholder = strings.ReplaceAll(placeholder, "{tenant_id}", c.GetString("tenant_id"))
	placeholder = strings.ReplaceAll(placeholder, "{vault_id}", c.Param("vaultId"))
	placeholder = strings.ReplaceAll(placeholder, "{secret_id}", c.Param("secretId"))
	placeholder = strings.ReplaceAll(placeholder, "{group_id}", c.Param("groupId"))
	placeholder = strings.ReplaceAll(placeholder, "{token_id}", c.Param("tokenId"))

	return placeholder
}

// auditPermissionCheck registra a verificação de permissão na auditoria
func (m *AuthorizationMiddleware) auditPermissionCheck(c *gin.Context, req *CheckRequest, allowed bool, duration time.Duration) {
	if m.auditService == nil {
		return
	}

	event := PermissionCheckEvent{
		User:      req.User,
		Relation:  req.Relation,
		Object:    req.Object,
		Result:    formatResult(allowed),
		Timestamp: time.Now(),
		Duration:  duration,
	}

	m.auditService.LogPermissionCheck(c.Request.Context(), event)
}

// logError registra um erro
func (m *AuthorizationMiddleware) logError(c *gin.Context, message string, err error) {
	if m.logger == nil {
		return
	}

	logData := map[string]interface{}{
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"user_agent": c.Request.UserAgent(),
		"ip":         c.ClientIP(),
		"message":    message,
	}

	if err != nil {
		logData["error"] = err.Error()
	}

	m.logger.Error("Authorization middleware error", logData)
}

// logUnauthorized registra tentativa de acesso não autorizado
func (m *AuthorizationMiddleware) logUnauthorized(c *gin.Context, user, relation, object string) {
	if m.logger == nil {
		return
	}

	m.logger.Warn("Unauthorized access attempt", map[string]interface{}{
		"user":       user,
		"relation":   relation,
		"object":     object,
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"user_agent": c.Request.UserAgent(),
		"ip":         c.ClientIP(),
	})
}

// formatResult formata o resultado da verificação
func formatResult(allowed bool) string {
	if allowed {
		return "allowed"
	}
	return "denied"
}

// PermissionCheck representa uma verificação de permissão
type PermissionCheck struct {
	User     string `json:"user"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}

// === MIDDLEWARE HELPERS ===

// RequireAnyVaultPermission verifica se o usuário tem qualquer uma das permissões listadas
func (m *AuthorizationMiddleware) RequireAnyVaultPermission(permissions ...VaultPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		userID := c.GetString("user_id")
		vaultID := c.Param("vaultId")

		if userID == "" || vaultID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing_required_parameters"})
			c.Abort()
			return
		}

		user := FormatUser(userID)
		object := FormatVault(vaultID)

		// Verificar cada permissão até encontrar uma que seja permitida
		for _, permission := range permissions {
			checkReq := &CheckRequest{
				User:     user,
				Relation: string(permission),
				Object:   object,
			}

			allowed, err := m.checkPermission(c.Request.Context(), checkReq)
			if err != nil {
				m.logError(c, "Authorization check failed", err)
				continue
			}

			if allowed {
				m.auditPermissionCheck(c, checkReq, true, time.Since(start))
				c.Set("authorized_vault_id", vaultID)
				c.Set("authorized_permission", string(permission))
				c.Set("authorization_duration", time.Since(start))
				c.Next()
				return
			}
		}

		// Nenhuma permissão foi encontrada
		m.logUnauthorized(c, user, "any_of", object)
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
		c.Abort()
	}
}

// RequireAllVaultPermissions verifica se o usuário tem todas as permissões listadas
func (m *AuthorizationMiddleware) RequireAllVaultPermissions(permissions ...VaultPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		userID := c.GetString("user_id")
		vaultID := c.Param("vaultId")

		if userID == "" || vaultID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing_required_parameters"})
			c.Abort()
			return
		}

		user := FormatUser(userID)
		object := FormatVault(vaultID)

		var requests []*CheckRequest

		// Preparar todas as verificações
		for _, permission := range permissions {
			requests = append(requests, &CheckRequest{
				User:     user,
				Relation: string(permission),
				Object:   object,
			})
		}

		// Realizar verificação em lote
		results, err := m.client.BatchCheck(c.Request.Context(), requests)
		if err != nil {
			m.logError(c, "Batch authorization check failed", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization_check_failed"})
			c.Abort()
			return
		}

		// Verificar se todas as permissões foram concedidas
		for i, result := range results.Results {
			if result == nil || !result.Allowed {
				req := requests[i]
				m.logUnauthorized(c, req.User, req.Relation, req.Object)
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
				c.Abort()
				return
			}
		}

		// Auditoria para todas as verificações
		for _, req := range requests {
			m.auditPermissionCheck(c, req, true, time.Since(start))
		}

		c.Set("authorized_vault_id", vaultID)
		c.Set("authorized_permissions", permissions)
		c.Set("authorization_duration", time.Since(start))

		c.Next()
	}
}

// ExtractUserFromJWT extrai o user_id do token JWT e adiciona ao contexto
func ExtractUserFromJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implementar extração do JWT
		// Por enquanto, assume que o user_id já está no contexto
		// Esta função deve ser implementada quando a autenticação JWT estiver pronta

		c.Next()
	}
}

// ExtractTenantFromJWT extrai o tenant_id do token JWT e adiciona ao contexto
func ExtractTenantFromJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implementar extração do JWT
		// Por enquanto, assume que o tenant_id já está no contexto
		// Esta função deve ser implementada quando a autenticação JWT estiver pronta

		c.Next()
	}
}
