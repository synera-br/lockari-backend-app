# 🚀 Plano de Implementação - Ordem de Desenvolvimento

## 📅 Cronograma de Desenvolvimento

### **Semana 1: Fundação e Estrutura Base**

#### **Dia 1: Setup e Tipos Base**
- [ ] Criar estrutura de arquivos
- [ ] Implementar `types.go` com todos os enums e structs
- [ ] Implementar `interfaces.go` com contratos públicos
- [ ] Criar testes unitários básicos para tipos
- [ ] Configurar CI/CD para o package

**Arquivos a criar:**
- `types.go` - Todos os tipos e enums
- `interfaces.go` - Interfaces públicas
- `errors.go` - Erros customizados
- `types_test.go` - Testes para tipos

#### **Dia 2: Configuração**
- [ ] Implementar `config.go` com validação
- [ ] Criar sistema de loading de configuração
- [ ] Implementar defaults sensatos
- [ ] Testes de configuração

**Arquivos a criar:**
- `config.go` - Configuração completa
- `config_test.go` - Testes de configuração

#### **Dia 3: Cliente OpenFGA Base**
- [ ] Implementar `client.go` com operações básicas
- [ ] Configurar cliente OpenFGA SDK
- [ ] Implementar método `Check` básico
- [ ] Testes com mock do OpenFGA

**Arquivos a criar:**
- `client.go` - Cliente OpenFGA
- `client_test.go` - Testes do cliente

#### **Dia 4: Serviço Principal**
- [ ] Implementar `service.go` com interface básica
- [ ] Método `Check` com validação
- [ ] Integração cliente + serviço
- [ ] Testes de integração

**Arquivos a criar:**
- `service.go` - Serviço principal
- `service_test.go` - Testes do serviço

#### **Dia 5: Cache System**
- [ ] Implementar `cache.go` com cache em memória
- [ ] Integrar cache no serviço
- [ ] Métricas básicas de cache
- [ ] Testes de cache

**Arquivos a criar:**
- `cache.go` - Sistema de cache
- `cache_test.go` - Testes do cache

### **Semana 2: Recursos Avançados**

#### **Dia 6: Middleware Gin**
- [ ] Implementar `middleware.go` para Gin
- [ ] Middleware para vault permissions
- [ ] Middleware para tenant isolation
- [ ] Testes de middleware

**Arquivos a criar:**
- `middleware.go` - Middleware Gin
- `middleware_test.go` - Testes do middleware

#### **Dia 7: Sistema de Auditoria**
- [ ] Implementar `audit.go` com logging estruturado
- [ ] Eventos de auditoria
- [ ] Integração com observabilidade
- [ ] Testes de auditoria

**Arquivos a criar:**
- `audit.go` - Sistema de auditoria
- `audit_test.go` - Testes de auditoria

#### **Dia 8: Serviço Lockari**
- [ ] Implementar `lockari_service.go` com operações de alto nível
- [ ] Operações de vault e tenant
- [ ] Tokens de API
- [ ] Testes específicos do domínio

**Arquivos a criar:**
- `lockari_service.go` - Serviço específico do Lockari
- `lockari_service_test.go` - Testes específicos

#### **Dia 9: Recursos Enterprise**
- [ ] Compartilhamento externo
- [ ] Aprovação dupla
- [ ] Grupos e roles
- [ ] Testes enterprise

#### **Dia 10: Otimizações e Performance**
- [ ] Batch operations
- [ ] Connection pooling
- [ ] Circuit breaker
- [ ] Benchmarks

### **Semana 3: Integração e Testes**

#### **Dia 11-12: Testes de Integração**
- [ ] Testes com OpenFGA real
- [ ] Testes Docker Compose
- [ ] Testes de performance
- [ ] Testes de failover

#### **Dia 13-14: Documentação e Exemplos**
- [ ] Documentação completa
- [ ] Exemplos de uso
- [ ] Guias de implementação
- [ ] Troubleshooting

#### **Dia 15: Integração com Projeto**
- [ ] Integração com handlers existentes
- [ ] Configuração de produção
- [ ] Deploy e testes finais

## 🔧 Implementação Detalhada

### **1. types.go - Primeira Implementação**

```go
package authorization

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

// ===== VAULT PERMISSIONS =====

// VaultPermission representa permissões granulares de vault
type VaultPermission string

const (
    VaultPermissionView      VaultPermission = "can_view"
    VaultPermissionRead      VaultPermission = "can_read"
    VaultPermissionCopy      VaultPermission = "can_copy"
    VaultPermissionDownload  VaultPermission = "can_download"
    VaultPermissionWrite     VaultPermission = "can_write"
    VaultPermissionDelete    VaultPermission = "can_delete"
    VaultPermissionShare     VaultPermission = "can_share"
    VaultPermissionManage    VaultPermission = "can_manage"
)

// AllVaultPermissions retorna todas as permissões válidas
func AllVaultPermissions() []VaultPermission {
    return []VaultPermission{
        VaultPermissionView, VaultPermissionRead, VaultPermissionCopy,
        VaultPermissionDownload, VaultPermissionWrite, VaultPermissionDelete,
        VaultPermissionShare, VaultPermissionManage,
    }
}

// IsValid verifica se a permissão é válida
func (vp VaultPermission) IsValid() bool {
    for _, valid := range AllVaultPermissions() {
        if vp == valid {
            return true
        }
    }
    return false
}

// String implementa fmt.Stringer
func (vp VaultPermission) String() string {
    return string(vp)
}

// ===== SECRET PERMISSIONS =====

type SecretPermission string

const (
    SecretPermissionView           SecretPermission = "can_view"
    SecretPermissionRead           SecretPermission = "can_read"
    SecretPermissionCopy           SecretPermission = "can_copy"
    SecretPermissionDownload       SecretPermission = "can_download"
    SecretPermissionWrite          SecretPermission = "can_write"
    SecretPermissionDelete         SecretPermission = "can_delete"
    SecretPermissionReadSensitive  SecretPermission = "can_read_sensitive"
    SecretPermissionCopySensitive  SecretPermission = "can_copy_sensitive"
    SecretPermissionCopyProduction SecretPermission = "can_copy_production"
)

// ===== TENANT ROLES =====

type TenantRole string

const (
    TenantRoleOwner  TenantRole = "owner"
    TenantRoleAdmin  TenantRole = "admin"
    TenantRoleMember TenantRole = "member"
    TenantRoleGuest  TenantRole = "guest"
)

// ===== GROUP ROLES =====

type GroupRole string

const (
    GroupRoleOwner  GroupRole = "owner"
    GroupRoleAdmin  GroupRole = "admin"
    GroupRoleMember GroupRole = "member"
)

// ===== PLAN FEATURES =====

type PlanFeature string

const (
    PlanFeatureBasic                PlanFeature = "basic"
    PlanFeatureAdvancedPermissions  PlanFeature = "advanced_permissions"
    PlanFeatureCrossTenantSharing   PlanFeature = "cross_tenant_sharing"
    PlanFeatureAuditLogs           PlanFeature = "audit_logs"
    PlanFeatureBackup              PlanFeature = "backup"
    PlanFeatureExternalSharing     PlanFeature = "external_sharing"
)

// ===== TOKEN PERMISSIONS =====

type TokenPermission string

const (
    TokenPermissionUse             TokenPermission = "can_use"
    TokenPermissionReadSecrets     TokenPermission = "can_read_secrets"
    TokenPermissionWriteSecrets    TokenPermission = "can_write_secrets"
    TokenPermissionManageVault     TokenPermission = "can_manage_vault"
    TokenPermissionRevoke          TokenPermission = "can_revoke"
    TokenPermissionRegenerate      TokenPermission = "can_regenerate"
)

// ===== REQUEST/RESPONSE TYPES =====

// CheckRequest representa uma solicitação de verificação
type CheckRequest struct {
    User     string `json:"user"`
    Relation string `json:"relation"`
    Object   string `json:"object"`
}

// Validate valida a solicitação
func (cr *CheckRequest) Validate() error {
    if cr.User == "" {
        return errors.New("user cannot be empty")
    }
    if cr.Relation == "" {
        return errors.New("relation cannot be empty")
    }
    if cr.Object == "" {
        return errors.New("object cannot be empty")
    }
    
    // Validar formato do usuário
    if !strings.HasPrefix(cr.User, "user:") && !strings.HasPrefix(cr.User, "token:") {
        return errors.New("user must start with 'user:' or 'token:'")
    }
    
    // Validar formato do objeto
    parts := strings.Split(cr.Object, ":")
    if len(parts) != 2 {
        return errors.New("object must be in format 'type:id'")
    }
    
    return nil
}

// CheckResponse representa a resposta de uma verificação
type CheckResponse struct {
    Allowed bool   `json:"allowed"`
    Reason  string `json:"reason,omitempty"`
}

// WriteRequest representa uma solicitação para criar tuplas
type WriteRequest struct {
    Tuples []Tuple `json:"tuples"`
}

// DeleteRequest representa uma solicitação para deletar tuplas
type DeleteRequest struct {
    Tuples []Tuple `json:"tuples"`
}

// Tuple representa uma tupla OpenFGA
type Tuple struct {
    User     string `json:"user"`
    Relation string `json:"relation"`
    Object   string `json:"object"`
}

// String implementa fmt.Stringer para Tuple
func (t Tuple) String() string {
    return fmt.Sprintf("%s#%s@%s", t.User, t.Relation, t.Object)
}

// Validate valida a tupla
func (t *Tuple) Validate() error {
    if t.User == "" {
        return errors.New("user cannot be empty")
    }
    if t.Relation == "" {
        return errors.New("relation cannot be empty")
    }
    if t.Object == "" {
        return errors.New("object cannot be empty")
    }
    return nil
}

// ListObjectsRequest representa solicitação para listar objetos
type ListObjectsRequest struct {
    User     string `json:"user"`
    Relation string `json:"relation"`
    Type     string `json:"type"`
}

// ListObjectsResponse representa a resposta da listagem
type ListObjectsResponse struct {
    Objects []string `json:"objects"`
}

// ===== UTILITY FUNCTIONS =====

// FormatUser formata um ID de usuário para o formato OpenFGA
func FormatUser(userID string) string {
    if strings.HasPrefix(userID, "user:") {
        return userID
    }
    return fmt.Sprintf("user:%s", userID)
}

// FormatObject formata um objeto para o formato OpenFGA
func FormatObject(objectType, objectID string) string {
    return fmt.Sprintf("%s:%s", objectType, objectID)
}

// FormatTenant formata um tenant para o formato OpenFGA
func FormatTenant(tenantID string) string {
    return FormatObject("tenant", tenantID)
}

// FormatVault formata um vault para o formato OpenFGA
func FormatVault(vaultID string) string {
    return FormatObject("vault", vaultID)
}

// FormatSecret formata um secret para o formato OpenFGA
func FormatSecret(secretID string) string {
    return FormatObject("secret", secretID)
}

// FormatToken formata um token para o formato OpenFGA
func FormatToken(tokenID string) string {
    return FormatObject("token", tokenID)
}

// ParseObject extrai tipo e ID de um objeto OpenFGA
func ParseObject(object string) (objectType, objectID string, err error) {
    parts := strings.Split(object, ":")
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid object format: %s", object)
    }
    return parts[0], parts[1], nil
}

// ParseUser extrai o ID do usuário de um formato OpenFGA
func ParseUser(user string) (userID string, err error) {
    if strings.HasPrefix(user, "user:") {
        return strings.TrimPrefix(user, "user:"), nil
    }
    return "", fmt.Errorf("invalid user format: %s", user)
}
```

### **2. interfaces.go - Contratos Públicos**

```go
package authorization

import (
    "context"
)

// AuthorizationService é a interface básica para operações OpenFGA
type AuthorizationService interface {
    // Check verifica uma permissão única
    Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error)
    
    // CheckBatch verifica múltiplas permissões em uma chamada
    CheckBatch(ctx context.Context, reqs []*CheckRequest) ([]*CheckResponse, error)
    
    // ListObjects lista todos os objetos que o usuário pode acessar
    ListObjects(ctx context.Context, req *ListObjectsRequest) (*ListObjectsResponse, error)
    
    // Write cria relacionamentos (tuplas)
    Write(ctx context.Context, req *WriteRequest) error
    
    // Delete remove relacionamentos
    Delete(ctx context.Context, req *DeleteRequest) error
    
    // Health verifica se o OpenFGA está disponível
    Health(ctx context.Context) error
}

// LockariAuthorizationService é a interface específica do domínio Lockari
type LockariAuthorizationService interface {
    AuthorizationService
    
    // === VAULT OPERATIONS ===
    
    // CanAccessVault verifica se o usuário pode acessar um vault
    CanAccessVault(ctx context.Context, userID, vaultID string, permission VaultPermission) (bool, error)
    
    // SetupVault configura um novo vault com permissões iniciais
    SetupVault(ctx context.Context, vaultID, tenantID, ownerID string) error
    
    // ShareVault compartilha um vault com outro usuário
    ShareVault(ctx context.Context, vaultID, ownerID, targetUserID string, permission VaultPermission) error
    
    // ListAccessibleVaults lista todos os vaults acessíveis para o usuário
    ListAccessibleVaults(ctx context.Context, userID string) ([]string, error)
    
    // === TENANT OPERATIONS ===
    
    // SetupTenant configura um novo tenant
    SetupTenant(ctx context.Context, tenantID, ownerID string, features []PlanFeature) error
    
    // AddUserToTenant adiciona um usuário ao tenant
    AddUserToTenant(ctx context.Context, userID, tenantID string, role TenantRole) error
    
    // RemoveUserFromTenant remove um usuário do tenant
    RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error
    
    // === GROUP OPERATIONS ===
    
    // CreateGroup cria um novo grupo
    CreateGroup(ctx context.Context, groupID, tenantID, ownerID string) error
    
    // AddUserToGroup adiciona um usuário ao grupo
    AddUserToGroup(ctx context.Context, userID, groupID string, role GroupRole) error
    
    // === TOKEN OPERATIONS ===
    
    // CreateAPIToken cria um token de API com permissões específicas
    CreateAPIToken(ctx context.Context, userID, vaultID string, permissions []TokenPermission) (string, error)
    
    // CheckTokenPermission verifica se um token tem uma permissão específica
    CheckTokenPermission(ctx context.Context, tokenID string, permission TokenPermission) (bool, error)
    
    // RevokeToken revoga um token de API
    RevokeToken(ctx context.Context, tokenID string) error
    
    // === EXTERNAL SHARING (Enterprise) ===
    
    // InitiateExternalSharing inicia processo de compartilhamento externo
    InitiateExternalSharing(ctx context.Context, userID, vaultID, targetTenantID string) (string, error)
    
    // ApproveExternalSharing aprova compartilhamento externo
    ApproveExternalSharing(ctx context.Context, userID, requestID string) error
    
    // RejectExternalSharing rejeita compartilhamento externo
    RejectExternalSharing(ctx context.Context, userID, requestID string) error
}

// Cache é a interface para cache de permissões
type Cache interface {
    // Get recupera um valor do cache
    Get(key string) (bool, bool)
    
    // Set armazena um valor no cache
    Set(key string, value bool, ttl time.Duration)
    
    // Delete remove um valor do cache
    Delete(key string)
    
    // Clear limpa todo o cache
    Clear()
    
    // Stats retorna estatísticas do cache
    Stats() CacheStats
}

// CacheStats contém estatísticas do cache
type CacheStats struct {
    Hits     int64 `json:"hits"`
    Misses   int64 `json:"misses"`
    Entries  int64 `json:"entries"`
    HitRate  float64 `json:"hit_rate"`
    Size     int64 `json:"size"`
}

// AuditService é a interface para sistema de auditoria
type AuditService interface {
    // LogPermissionCheck registra uma verificação de permissão
    LogPermissionCheck(ctx context.Context, event PermissionCheckEvent)
    
    // LogPermissionGrant registra concessão de permissão
    LogPermissionGrant(ctx context.Context, event PermissionGrantEvent)
    
    // LogPermissionRevoke registra revogação de permissão
    LogPermissionRevoke(ctx context.Context, event PermissionRevokeEvent)
    
    // LogSuspiciousActivity registra atividade suspeita
    LogSuspiciousActivity(ctx context.Context, event SuspiciousActivityEvent)
}

// === AUDIT EVENTS ===

// PermissionCheckEvent representa um evento de verificação de permissão
type PermissionCheckEvent struct {
    User     string `json:"user"`
    Relation string `json:"relation"`
    Object   string `json:"object"`
    Result   string `json:"result"`
    Error    string `json:"error,omitempty"`
}

// PermissionGrantEvent representa um evento de concessão de permissão
type PermissionGrantEvent struct {
    Grantor  string `json:"grantor"`
    Grantee  string `json:"grantee"`
    Relation string `json:"relation"`
    Object   string `json:"object"`
}

// PermissionRevokeEvent representa um evento de revogação de permissão
type PermissionRevokeEvent struct {
    Revoker  string `json:"revoker"`
    Revokee  string `json:"revokee"`
    Relation string `json:"relation"`
    Object   string `json:"object"`
}

// SuspiciousActivityEvent representa atividade suspeita
type SuspiciousActivityEvent struct {
    User        string `json:"user"`
    Activity    string `json:"activity"`
    Details     string `json:"details"`
    Severity    string `json:"severity"`
    ClientIP    string `json:"client_ip"`
    UserAgent   string `json:"user_agent"`
}
```

### **3. Testes Base**

```go
// types_test.go
package authorization

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestVaultPermission_IsValid(t *testing.T) {
    tests := []struct {
        name       string
        permission VaultPermission
        want       bool
    }{
        {"valid permission", VaultPermissionRead, true},
        {"invalid permission", VaultPermission("invalid"), false},
        {"empty permission", VaultPermission(""), false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.permission.IsValid()
            assert.Equal(t, tt.want, got)
        })
    }
}

func TestCheckRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     *CheckRequest
        wantErr bool
    }{
        {
            name: "valid request",
            req: &CheckRequest{
                User:     "user:alice",
                Relation: "can_read",
                Object:   "vault:123",
            },
            wantErr: false,
        },
        {
            name: "empty user",
            req: &CheckRequest{
                User:     "",
                Relation: "can_read",
                Object:   "vault:123",
            },
            wantErr: true,
        },
        {
            name: "invalid user format",
            req: &CheckRequest{
                User:     "alice",
                Relation: "can_read",
                Object:   "vault:123",
            },
            wantErr: true,
        },
        {
            name: "invalid object format",
            req: &CheckRequest{
                User:     "user:alice",
                Relation: "can_read",
                Object:   "vault123",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestFormatUser(t *testing.T) {
    tests := []struct {
        name   string
        userID string
        want   string
    }{
        {"without prefix", "alice", "user:alice"},
        {"with prefix", "user:alice", "user:alice"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FormatUser(tt.userID)
            assert.Equal(t, tt.want, got)
        })
    }
}

func TestParseObject(t *testing.T) {
    tests := []struct {
        name       string
        object     string
        wantType   string
        wantID     string
        wantErr    bool
    }{
        {"valid object", "vault:123", "vault", "123", false},
        {"invalid object", "vault123", "", "", true},
        {"empty object", "", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotType, gotID, err := ParseObject(tt.object)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantType, gotType)
                assert.Equal(t, tt.wantID, gotID)
            }
        })
    }
}
```

## 📋 Checklist de Implementação

### **✅ Fase 1: Fundação (Semana 1)**
- [ ] `types.go` - Tipos e enums completos
- [ ] `interfaces.go` - Contratos públicos
- [ ] `errors.go` - Erros customizados
- [ ] `config.go` - Configuração com validação
- [ ] `client.go` - Cliente OpenFGA básico
- [ ] `service.go` - Serviço principal
- [ ] `cache.go` - Cache em memória
- [ ] Testes unitários para todos os componentes

### **✅ Fase 2: Recursos Avançados (Semana 2)**
- [ ] `middleware.go` - Middleware Gin
- [ ] `audit.go` - Sistema de auditoria
- [ ] `lockari_service.go` - Serviço específico do domínio
- [ ] Recursos enterprise (compartilhamento externo)
- [ ] Otimizações de performance
- [ ] Testes de integração

### **✅ Fase 3: Finalização (Semana 3)**
- [ ] Documentação completa
- [ ] Exemplos de uso
- [ ] Integração com projeto existente
- [ ] Benchmarks e otimizações
- [ ] Deploy e configuração de produção

Este plano de implementação garante que o desenvolvimento seja estruturado, testado e bem documentado, seguindo as melhores práticas de desenvolvimento Go.
