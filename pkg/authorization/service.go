package authorization

import (
	"context"
	"fmt"
	"time"
)

// Service é o serviço principal de autorização
type Service struct {
	client *OpenFGAClient
	config *Config
	logger Logger
	audit  AuditService
	cache  interface{} // Simplified for now
}

// ServiceOptions define as opções para o serviço
type ServiceOptions struct {
	Client *OpenFGAClient
	Config *Config
	Logger Logger
	Audit  AuditService
	Cache  interface{} // Simplified for now
}

// NewService cria um novo serviço de autorização
func NewService(opts ServiceOptions) *Service {
	return &Service{
		client: opts.Client,
		config: opts.Config,
		logger: opts.Logger,
		audit:  opts.Audit,
		cache:  opts.Cache,
	}
}

// Check verifica se um usuário tem uma permissão específica
func (s *Service) Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Realizar verificação
	response, err := s.client.Check(ctx, req)
	if err != nil {
		s.logError("Authorization check failed", err)
		return nil, fmt.Errorf("authorization check failed: %w", err)
	}

	// Auditoria
	s.auditPermissionCheck(ctx, req, response.Allowed, time.Now())

	return response, nil
}

// CheckBatch verifica múltiplas permissões em lote
func (s *Service) CheckBatch(ctx context.Context, requests []*CheckRequest) (*BatchCheckResponse, error) {
	if len(requests) == 0 {
		return &BatchCheckResponse{
			Results: make([]*CheckResponse, 0),
		}, nil
	}

	// Validar todas as requests
	for i, req := range requests {
		if err := req.Validate(); err != nil {
			return nil, fmt.Errorf("invalid request at index %d: %w", i, err)
		}
	}

	// Realizar verificação em lote
	response, err := s.client.BatchCheck(ctx, requests)
	if err != nil {
		s.logError("Batch authorization check failed", err)
		return nil, fmt.Errorf("batch authorization check failed: %w", err)
	}

	// Auditoria para todas as verificações
	for i, req := range requests {
		if i < len(response.Results) && response.Results[i] != nil {
			s.auditPermissionCheck(ctx, req, response.Results[i].Allowed, time.Now())
		}
	}

	return response, nil
}

// ListObjects lista objetos que o usuário tem acesso
func (s *Service) ListObjects(ctx context.Context, req *ListObjectsRequest) (*ListObjectsResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Realizar listagem
	response, err := s.client.ListObjects(ctx, req)
	if err != nil {
		s.logError("List objects failed", err)
		return nil, fmt.Errorf("list objects failed: %w", err)
	}

	return response, nil
}

// CanAccessVault verifica se o usuário pode acessar um vault com uma permissão específica
func (s *Service) CanAccessVault(ctx context.Context, userID, vaultID string, permission VaultPermission) (bool, error) {
	if userID == "" || vaultID == "" {
		return false, fmt.Errorf("userID and vaultID are required")
	}

	if !permission.IsValid() {
		return false, fmt.Errorf("invalid vault permission: %s", permission)
	}

	req := &CheckRequest{
		User:     FormatUser(userID),
		Relation: string(permission),
		Object:   FormatVault(vaultID),
	}

	response, err := s.Check(ctx, req)
	if err != nil {
		return false, err
	}

	return response.Allowed, nil
}

// CanAccessSecret verifica se o usuário pode acessar um secret com uma permissão específica
func (s *Service) CanAccessSecret(ctx context.Context, userID, secretID string, permission SecretPermission) (bool, error) {
	if userID == "" || secretID == "" {
		return false, fmt.Errorf("userID and secretID are required")
	}

	if !permission.IsValid() {
		return false, fmt.Errorf("invalid secret permission: %s", permission)
	}

	req := &CheckRequest{
		User:     FormatUser(userID),
		Relation: string(permission),
		Object:   FormatSecret(secretID),
	}

	response, err := s.Check(ctx, req)
	if err != nil {
		return false, err
	}

	return response.Allowed, nil
}

// HasTenantRole verifica se o usuário tem um papel específico no tenant
func (s *Service) HasTenantRole(ctx context.Context, userID, tenantID string, role TenantRole) (bool, error) {
	if userID == "" || tenantID == "" {
		return false, fmt.Errorf("userID and tenantID are required")
	}

	if !role.IsValid() {
		return false, fmt.Errorf("invalid tenant role: %s", role)
	}

	req := &CheckRequest{
		User:     FormatUser(userID),
		Relation: string(role),
		Object:   FormatTenant(tenantID),
	}

	response, err := s.Check(ctx, req)
	if err != nil {
		return false, err
	}

	return response.Allowed, nil
}

// GrantVaultPermission concede uma permissão de vault a um usuário
func (s *Service) GrantVaultPermission(ctx context.Context, grantor, grantee, vaultID string, permission VaultPermission) error {
	if grantor == "" || grantee == "" || vaultID == "" {
		return fmt.Errorf("grantor, grantee, and vaultID are required")
	}

	if !permission.IsValid() {
		return fmt.Errorf("invalid vault permission: %s", permission)
	}

	// Verificar se o grantor tem permissão para gerenciar o vault
	canManage, err := s.CanAccessVault(ctx, grantor, vaultID, VaultPermissionManage)
	if err != nil {
		return fmt.Errorf("failed to check grantor permissions: %w", err)
	}

	if !canManage {
		return fmt.Errorf("grantor does not have manage permission on vault")
	}

	// Criar a tupla de permissão
	tuple := Tuple{
		User:     FormatUser(grantee),
		Relation: string(permission),
		Object:   FormatVault(vaultID),
	}

	writeReq := &WriteRequest{
		Tuples: []Tuple{tuple},
	}

	// Executar a operação (Note: WriteRequest precisa ser implementado no client)
	err = s.writeRelationship(ctx, writeReq)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	// Auditoria
	s.auditPermissionGrant(ctx, grantor, grantee, string(permission), FormatVault(vaultID))

	return nil
}

// RevokeVaultPermission revoga uma permissão de vault de um usuário
func (s *Service) RevokeVaultPermission(ctx context.Context, revoker, revokee, vaultID string, permission VaultPermission) error {
	if revoker == "" || revokee == "" || vaultID == "" {
		return fmt.Errorf("revoker, revokee, and vaultID are required")
	}

	if !permission.IsValid() {
		return fmt.Errorf("invalid vault permission: %s", permission)
	}

	// Verificar se o revoker tem permissão para gerenciar o vault
	canManage, err := s.CanAccessVault(ctx, revoker, vaultID, VaultPermissionManage)
	if err != nil {
		return fmt.Errorf("failed to check revoker permissions: %w", err)
	}

	if !canManage {
		return fmt.Errorf("revoker does not have manage permission on vault")
	}

	// Criar a tupla de permissão para revogar
	tuple := Tuple{
		User:     FormatUser(revokee),
		Relation: string(permission),
		Object:   FormatVault(vaultID),
	}

	deleteReq := &DeleteRequest{
		Tuples: []Tuple{tuple},
	}

	// Executar a operação (Note: DeleteRequest precisa ser implementado no client)
	err = s.deleteRelationship(ctx, deleteReq)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	// Auditoria
	s.auditPermissionRevoke(ctx, revoker, revokee, string(permission), FormatVault(vaultID))

	return nil
}

// ListUserVaults lista todos os vaults que o usuário tem acesso
func (s *Service) ListUserVaults(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	req := &ListObjectsRequest{
		User:     FormatUser(userID),
		Relation: string(VaultPermissionView),
		Type:     "vault",
	}

	response, err := s.ListObjects(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list user vaults: %w", err)
	}

	return response.Objects, nil
}

// ListUserSecrets lista todos os secrets que o usuário tem acesso
func (s *Service) ListUserSecrets(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	req := &ListObjectsRequest{
		User:     FormatUser(userID),
		Relation: string(SecretPermissionView),
		Type:     "secret",
	}

	response, err := s.ListObjects(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list user secrets: %w", err)
	}

	return response.Objects, nil
}

// GetVaultUsers lista todos os usuários que têm acesso a um vault
func (s *Service) GetVaultUsers(ctx context.Context, vaultID string) ([]string, error) {
	if vaultID == "" {
		return nil, fmt.Errorf("vaultID is required")
	}

	// Esta operação requer funcionalidade específica do OpenFGA
	// Por enquanto, retornaremos um placeholder
	// TODO: Implementar quando o OpenFGA SDK suportar reverse queries
	return []string{}, fmt.Errorf("not implemented: get vault users")
}

// HealthCheck verifica a saúde do serviço
func (s *Service) HealthCheck(ctx context.Context) *HealthCheckResponse {
	return s.client.HealthCheck(ctx)
}

// writeRelationship escreve uma relação no OpenFGA
func (s *Service) writeRelationship(ctx context.Context, req *WriteRequest) error {
	// TODO: Implementar quando o client suportar write operations
	return fmt.Errorf("not implemented: write relationship")
}

// deleteRelationship remove uma relação do OpenFGA
func (s *Service) deleteRelationship(ctx context.Context, req *DeleteRequest) error {
	// TODO: Implementar quando o client suportar delete operations
	return fmt.Errorf("not implemented: delete relationship")
}

// auditPermissionCheck registra uma verificação de permissão
func (s *Service) auditPermissionCheck(ctx context.Context, req *CheckRequest, allowed bool, timestamp time.Time) {
	if s.audit == nil {
		return
	}

	event := PermissionCheckEvent{
		User:      req.User,
		Relation:  req.Relation,
		Object:    req.Object,
		Result:    formatBoolResult(allowed),
		Timestamp: timestamp,
	}

	s.audit.LogPermissionCheck(ctx, event)
}

// auditPermissionGrant registra uma concessão de permissão
func (s *Service) auditPermissionGrant(ctx context.Context, grantor, grantee, relation, object string) {
	if s.audit == nil {
		return
	}

	event := PermissionGrantEvent{
		Grantor:   grantor,
		Grantee:   grantee,
		Relation:  relation,
		Object:    object,
		Timestamp: time.Now(),
	}

	s.audit.LogPermissionGrant(ctx, event)
}

// auditPermissionRevoke registra uma revogação de permissão
func (s *Service) auditPermissionRevoke(ctx context.Context, revoker, revokee, relation, object string) {
	if s.audit == nil {
		return
	}

	event := PermissionRevokeEvent{
		Revoker:   revoker,
		Revokee:   revokee,
		Relation:  relation,
		Object:    object,
		Timestamp: time.Now(),
	}

	s.audit.LogPermissionRevoke(ctx, event)
}

// logError registra um erro
func (s *Service) logError(message string, err error) {
	if s.logger == nil {
		return
	}

	s.logger.Error(message, map[string]interface{}{
		"error": err.Error(),
	})
}

// formatBoolResult formata um resultado booleano como string
func formatBoolResult(allowed bool) string {
	if allowed {
		return "allowed"
	}
	return "denied"
}
