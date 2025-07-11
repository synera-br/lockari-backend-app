package authorization

import (
	"fmt"
)

// ErrInvalidUser é retornado quando o formato do usuário é inválido
var ErrInvalidUser = fmt.Errorf("invalid user format")

// ErrInvalidObject é retornado quando o formato do objeto é inválido
var ErrInvalidObject = fmt.Errorf("invalid object format")

// ErrInvalidRelation é retornado quando a relação é inválida
var ErrInvalidRelation = fmt.Errorf("invalid relation")

// ErrInvalidPermission é retornado quando a permissão é inválida
var ErrInvalidPermission = fmt.Errorf("invalid permission")

// ErrInvalidTuple é retornado quando a tupla é inválida
var ErrInvalidTuple = fmt.Errorf("invalid tuple")

// ErrInvalidRequest é retornado quando a solicitação é inválida
var ErrInvalidRequest = fmt.Errorf("invalid request")

// ErrInvalidConfig é retornado quando a configuração é inválida
var ErrInvalidConfig = fmt.Errorf("invalid configuration")

// ErrInvalidToken é retornado quando o token é inválido
var ErrInvalidToken = fmt.Errorf("invalid token")

// ErrInvalidTenant é retornado quando o tenant é inválido
var ErrInvalidTenant = fmt.Errorf("invalid tenant")

// ErrInvalidVault é retornado quando o vault é inválido
var ErrInvalidVault = fmt.Errorf("invalid vault")

// ErrInvalidSecret é retornado quando o secret é inválido
var ErrInvalidSecret = fmt.Errorf("invalid secret")

// ErrInvalidGroup é retornado quando o grupo é inválido
var ErrInvalidGroup = fmt.Errorf("invalid group")

// ErrPermissionDenied é retornado quando a permissão é negada
var ErrPermissionDenied = fmt.Errorf("permission denied")

// ErrNotFound é retornado quando o recurso não é encontrado
var ErrNotFound = fmt.Errorf("not found")

// ErrAlreadyExists é retornado quando o recurso já existe
var ErrAlreadyExists = fmt.Errorf("already exists")

// ErrUnauthorized é retornado quando o usuário não está autorizado
var ErrUnauthorized = fmt.Errorf("unauthorized")

// ErrForbidden é retornado quando o acesso é proibido
var ErrForbidden = fmt.Errorf("forbidden")

// ErrTimeout é retornado quando a operação timeout
var ErrTimeout = fmt.Errorf("operation timeout")

// ErrRateLimited é retornado quando o rate limit é excedido
var ErrRateLimited = fmt.Errorf("rate limited")

// ErrCircuitBreakerOpen é retornado quando o circuit breaker está aberto
var ErrCircuitBreakerOpen = fmt.Errorf("circuit breaker open")

// ErrCacheNotAvailable é retornado quando o cache não está disponível
var ErrCacheNotAvailable = fmt.Errorf("cache not available")

// ErrAuditNotAvailable é retornado quando o sistema de auditoria não está disponível
var ErrAuditNotAvailable = fmt.Errorf("audit system not available")

// ErrClientNotInitialized é retornado quando o cliente não foi inicializado
var ErrClientNotInitialized = fmt.Errorf("client not initialized")

// ErrServiceUnavailable é retornado quando o serviço não está disponível
var ErrServiceUnavailable = fmt.Errorf("service unavailable")

// ErrInternalError é retornado para erros internos
var ErrInternalError = fmt.Errorf("internal error")

// AuthorizationError representa um erro de autorização
type AuthorizationError struct {
	Code    string
	Message string
	Details map[string]interface{}
	Cause   error
}

// Error implementa a interface error
func (e *AuthorizationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implementa a interface Unwrap
func (e *AuthorizationError) Unwrap() error {
	return e.Cause
}

// Is implementa a interface Is
func (e *AuthorizationError) Is(target error) bool {
	t, ok := target.(*AuthorizationError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewAuthorizationError cria um novo erro de autorização
func NewAuthorizationError(code, message string, cause error) *AuthorizationError {
	return &AuthorizationError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adiciona detalhes ao erro
func (e *AuthorizationError) WithDetails(key string, value interface{}) *AuthorizationError {
	e.Details[key] = value
	return e
}

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
	Cause   error
}

// Error implementa a interface error
func (e *ValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("validation error for field '%s': %s (value: %v): %v", e.Field, e.Message, e.Value, e.Cause)
	}
	return fmt.Sprintf("validation error for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// Unwrap implementa a interface Unwrap
func (e *ValidationError) Unwrap() error {
	return e.Cause
}

// NewValidationError cria um novo erro de validação
func NewValidationError(field, message string, value interface{}, cause error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
		Cause:   cause,
	}
}

// ConfigError representa um erro de configuração
type ConfigError struct {
	Key     string
	Message string
	Cause   error
}

// Error implementa a interface error
func (e *ConfigError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("config error for key '%s': %s: %v", e.Key, e.Message, e.Cause)
	}
	return fmt.Sprintf("config error for key '%s': %s", e.Key, e.Message)
}

// Unwrap implementa a interface Unwrap
func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// NewConfigError cria um novo erro de configuração
func NewConfigError(key, message string, cause error) *ConfigError {
	return &ConfigError{
		Key:     key,
		Message: message,
		Cause:   cause,
	}
}

// NetworkError representa um erro de rede
type NetworkError struct {
	Operation string
	Address   string
	Message   string
	Cause     error
}

// Error implementa a interface error
func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("network error during %s to %s: %s: %v", e.Operation, e.Address, e.Message, e.Cause)
	}
	return fmt.Sprintf("network error during %s to %s: %s", e.Operation, e.Address, e.Message)
}

// Unwrap implementa a interface Unwrap
func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// Temporary implementa a interface Temporary
func (e *NetworkError) Temporary() bool {
	return true
}

// NewNetworkError cria um novo erro de rede
func NewNetworkError(operation, address, message string, cause error) *NetworkError {
	return &NetworkError{
		Operation: operation,
		Address:   address,
		Message:   message,
		Cause:     cause,
	}
}

// CacheError representa um erro de cache
type CacheError struct {
	Operation string
	Key       string
	Message   string
	Cause     error
}

// Error implementa a interface error
func (e *CacheError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("cache error during %s for key '%s': %s: %v", e.Operation, e.Key, e.Message, e.Cause)
	}
	return fmt.Sprintf("cache error during %s for key '%s': %s", e.Operation, e.Key, e.Message)
}

// Unwrap implementa a interface Unwrap
func (e *CacheError) Unwrap() error {
	return e.Cause
}

// NewCacheError cria um novo erro de cache
func NewCacheError(operation, key, message string, cause error) *CacheError {
	return &CacheError{
		Operation: operation,
		Key:       key,
		Message:   message,
		Cause:     cause,
	}
}

// AuditError representa um erro de auditoria
type AuditError struct {
	Operation string
	Message   string
	Cause     error
}

// Error implementa a interface error
func (e *AuditError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("audit error during %s: %s: %v", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("audit error during %s: %s", e.Operation, e.Message)
}

// Unwrap implementa a interface Unwrap
func (e *AuditError) Unwrap() error {
	return e.Cause
}

// NewAuditError cria um novo erro de auditoria
func NewAuditError(operation, message string, cause error) *AuditError {
	return &AuditError{
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// Códigos de erro comuns
const (
	ErrCodeInvalidUser          = "INVALID_USER"
	ErrCodeInvalidObject        = "INVALID_OBJECT"
	ErrCodeInvalidRelation      = "INVALID_RELATION"
	ErrCodeInvalidPermission    = "INVALID_PERMISSION"
	ErrCodeInvalidTuple         = "INVALID_TUPLE"
	ErrCodeInvalidRequest       = "INVALID_REQUEST"
	ErrCodeInvalidConfig        = "INVALID_CONFIG"
	ErrCodeInvalidToken         = "INVALID_TOKEN"
	ErrCodeInvalidTenant        = "INVALID_TENANT"
	ErrCodeInvalidVault         = "INVALID_VAULT"
	ErrCodeInvalidSecret        = "INVALID_SECRET"
	ErrCodeInvalidGroup         = "INVALID_GROUP"
	ErrCodePermissionDenied     = "PERMISSION_DENIED"
	ErrCodeNotFound             = "NOT_FOUND"
	ErrCodeAlreadyExists        = "ALREADY_EXISTS"
	ErrCodeUnauthorized         = "UNAUTHORIZED"
	ErrCodeForbidden            = "FORBIDDEN"
	ErrCodeTimeout              = "TIMEOUT"
	ErrCodeRateLimited          = "RATE_LIMITED"
	ErrCodeCircuitBreakerOpen   = "CIRCUIT_BREAKER_OPEN"
	ErrCodeCacheNotAvailable    = "CACHE_NOT_AVAILABLE"
	ErrCodeAuditNotAvailable    = "AUDIT_NOT_AVAILABLE"
	ErrCodeClientNotInitialized = "CLIENT_NOT_INITIALIZED"
	ErrCodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	ErrCodeInternalError        = "INTERNAL_ERROR"
)

// Funções auxiliares para criar erros comuns
func NewInvalidUserError(userID string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeInvalidUser, "invalid user format", cause).
		WithDetails("user_id", userID)
}

func NewInvalidObjectError(objectType, objectID string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeInvalidObject, "invalid object format", cause).
		WithDetails("object_type", objectType).
		WithDetails("object_id", objectID)
}

func NewInvalidRelationError(relation string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeInvalidRelation, "invalid relation", cause).
		WithDetails("relation", relation)
}

func NewInvalidPermissionError(permission string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeInvalidPermission, "invalid permission", cause).
		WithDetails("permission", permission)
}

func NewPermissionDeniedError(userID, relation, object string) *AuthorizationError {
	return NewAuthorizationError(ErrCodePermissionDenied, "permission denied", nil).
		WithDetails("user_id", userID).
		WithDetails("relation", relation).
		WithDetails("object", object)
}

func NewNotFoundError(resourceType, resourceID string) *AuthorizationError {
	return NewAuthorizationError(ErrCodeNotFound, "resource not found", nil).
		WithDetails("resource_type", resourceType).
		WithDetails("resource_id", resourceID)
}

func NewUnauthorizedError(userID string) *AuthorizationError {
	return NewAuthorizationError(ErrCodeUnauthorized, "user not authorized", nil).
		WithDetails("user_id", userID)
}

func NewForbiddenError(userID, action string) *AuthorizationError {
	return NewAuthorizationError(ErrCodeForbidden, "action forbidden", nil).
		WithDetails("user_id", userID).
		WithDetails("action", action)
}

func NewTimeoutError(operation string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeTimeout, "operation timeout", cause).
		WithDetails("operation", operation)
}

func NewRateLimitedError(userID string) *AuthorizationError {
	return NewAuthorizationError(ErrCodeRateLimited, "rate limit exceeded", nil).
		WithDetails("user_id", userID)
}

func NewCircuitBreakerOpenError(service string) *AuthorizationError {
	return NewAuthorizationError(ErrCodeCircuitBreakerOpen, "circuit breaker is open", nil).
		WithDetails("service", service)
}

func NewServiceUnavailableError(service string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeServiceUnavailable, "service unavailable", cause).
		WithDetails("service", service)
}

func NewInternalError(operation string, cause error) *AuthorizationError {
	return NewAuthorizationError(ErrCodeInternalError, "internal error", cause).
		WithDetails("operation", operation)
}
