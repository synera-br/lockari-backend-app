# 📋 Plano de Desenvolvimento - Package OpenFGA Authorization

## 🎯 Objetivo Geral

Desenvolver um package completo de autorização para o Lockari usando OpenFGA, que implemente o modelo refinado de permissões granulares, suporte multi-tenancy e integração com o sistema existente.

## 🏗️ Arquitetura Proposta

### Visão Geral
```
┌─────────────────────────────────────────────────────────────────┐
│                    Lockari Backend Go                          │
├─────────────────────────────────────────────────────────────────┤
│  pkg/authorization/                                             │
│  ├── interfaces.go          # Interfaces principais             │
│  ├── types.go               # Tipos e enums                     │
│  ├── config.go              # Configuração                      │
│  ├── client.go              # Cliente OpenFGA                   │
│  ├── service.go             # Serviço de autorização            │
│  ├── middleware.go          # Middleware Gin                    │
│  ├── cache.go               # Cache de permissões               │
│  ├── audit.go               # Auditoria e logging               │
│  └── examples.go            # Exemplos de uso                   │
├─────────────────────────────────────────────────────────────────┤
│  OpenFGA SDK + HTTP Client                                     │
├─────────────────────────────────────────────────────────────────┤
│  OpenFGA Server (Docker/Cloud)                                 │
│  ├── PostgreSQL Store                                          │
│  └── Authorization Model                                       │
└─────────────────────────────────────────────────────────────────┘
```

## 🗂️ Estrutura de Arquivos

### 1. Arquivos Principais
```
pkg/authorization/
├── interfaces.go              # Interfaces públicas
├── types.go                   # Tipos, enums e constantes
├── config.go                  # Configuração OpenFGA
├── client.go                  # Cliente OpenFGA com retry/cache
├── service.go                 # Serviço principal de autorização
├── lockari_service.go         # Serviço específico do domínio
├── middleware.go              # Middleware Gin para autorização
├── cache.go                   # Cache em memória para performance
├── audit.go                   # Sistema de auditoria
├── errors.go                  # Erros customizados
├── examples.go                # Exemplos de uso completos
├── README.md                  # Documentação existente
├── authorization_test.go      # Testes unitários
└── integration_test.go        # Testes de integração
```

### 2. Subpackages (Futuro)
```
pkg/authorization/
├── openfga/                   # Específico do OpenFGA
│   ├── client.go              # Cliente OpenFGA otimizado
│   ├── models.go              # Modelos e tuplas
│   └── utils.go               # Utilitários OpenFGA
├── cache/                     # Cache especializado
│   ├── memory.go              # Cache em memória
│   └── redis.go               # Cache Redis (futuro)
└── audit/                     # Auditoria especializada
    ├── logger.go              # Logger de auditoria
    └── events.go              # Eventos de auditoria
```

## 🔧 Componentes Principais

### 1. **Interfaces Públicas** (`interfaces.go`)

#### Por que desenvolver:
- **Abstração**: Permite trocar implementações sem quebrar código existente
- **Testabilidade**: Facilita criação de mocks para testes
- **Flexibilidade**: Suporta diferentes provedores de autorização no futuro

#### O que será implementado:
```go
// Interface básica para operações OpenFGA
type AuthorizationService interface {
    Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error)
    CheckBatch(ctx context.Context, reqs []*CheckRequest) ([]*CheckResponse, error)
    ListObjects(ctx context.Context, req *ListObjectsRequest) (*ListObjectsResponse, error)
    Write(ctx context.Context, req *WriteRequest) error
    Delete(ctx context.Context, req *DeleteRequest) error
}

// Interface específica do domínio Lockari
type LockariAuthorizationService interface {
    AuthorizationService
    
    // Vault Operations
    CanAccessVault(ctx context.Context, userID, vaultID string, permission VaultPermission) (bool, error)
    SetupVault(ctx context.Context, vaultID, tenantID, ownerID string) error
    ShareVault(ctx context.Context, vaultID, ownerID, targetUserID string, permission VaultPermission) error
    
    // Tenant Operations
    SetupTenant(ctx context.Context, tenantID, ownerID string, features []PlanFeature) error
    AddUserToTenant(ctx context.Context, userID, tenantID string, role TenantRole) error
    
    // Token Operations
    CreateAPIToken(ctx context.Context, userID, vaultID string, permissions []TokenPermission) (string, error)
    CheckTokenPermission(ctx context.Context, tokenID string, permission TokenPermission) (bool, error)
    
    // External Sharing (Enterprise)
    InitiateExternalSharing(ctx context.Context, userID, vaultID, targetTenantID string) (string, error)
    ApproveExternalSharing(ctx context.Context, userID, requestID string) error
}
```

### 2. **Tipos e Enums** (`types.go`)

#### Por que desenvolver:
- **Type Safety**: Evita erros de string incorretas
- **Autocompletar**: Melhor experiência de desenvolvimento
- **Validação**: Garante que apenas valores válidos sejam usados

#### O que será implementado:
```go
// Vault Permissions
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

// Secret Permissions
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

// Tenant Roles
type TenantRole string

const (
    TenantRoleOwner  TenantRole = "owner"
    TenantRoleAdmin  TenantRole = "admin"
    TenantRoleMember TenantRole = "member"
    TenantRoleGuest  TenantRole = "guest"
)

// Plan Features
type PlanFeature string

const (
    PlanFeatureBasic                PlanFeature = "basic"
    PlanFeatureAdvancedPermissions  PlanFeature = "advanced_permissions"
    PlanFeatureCrossTenantSharing   PlanFeature = "cross_tenant_sharing"
    PlanFeatureAuditLogs           PlanFeature = "audit_logs"
    PlanFeatureBackup              PlanFeature = "backup"
    PlanFeatureExternalSharing     PlanFeature = "external_sharing"
)
```

### 3. **Configuração** (`config.go`)

#### Por que desenvolver:
- **Flexibilidade**: Suporta diferentes ambientes (dev, staging, prod)
- **Segurança**: Configuração centralizada de credenciais
- **Manutenibilidade**: Fácil alteração de parâmetros

#### O que será implementado:
```go
type Config struct {
    // OpenFGA Server
    APIURL               string        `mapstructure:"api_url"`
    StoreID              string        `mapstructure:"store_id"`
    AuthorizationModelID string        `mapstructure:"authorization_model_id"`
    
    // Authentication
    APITokenIssuer   string `mapstructure:"api_token_issuer"`
    APIAudience      string `mapstructure:"api_audience"`
    ClientID         string `mapstructure:"client_id"`
    ClientSecret     string `mapstructure:"client_secret"`
    Scopes           string `mapstructure:"scopes"`
    
    // Performance
    Timeout         time.Duration `mapstructure:"timeout"`
    RetryAttempts   int           `mapstructure:"retry_attempts"`
    RetryDelay      time.Duration `mapstructure:"retry_delay"`
    
    // Cache
    CacheEnabled    bool          `mapstructure:"cache_enabled"`
    CacheTTL        time.Duration `mapstructure:"cache_ttl"`
    CacheMaxSize    int           `mapstructure:"cache_max_size"`
    
    // Audit
    AuditEnabled    bool   `mapstructure:"audit_enabled"`
    AuditLevel      string `mapstructure:"audit_level"`
}
```

### 4. **Cliente OpenFGA** (`client.go`)

#### Por que desenvolver:
- **Abstração**: Encapsula complexidade do SDK OpenFGA
- **Retry Logic**: Implementa lógica de retry para falhas de rede
- **Cache**: Otimiza performance com cache inteligente
- **Logging**: Facilita debug e monitoramento

#### O que será implementado:
```go
type Client struct {
    fgaClient *client.OpenFgaClient
    config    *Config
    cache     Cache
    logger    *slog.Logger
}

func NewClient(config *Config) (*Client, error) {
    // Configuração do cliente OpenFGA
    // Inicialização do cache
    // Setup de logging
    // Validação de conectividade
}

func (c *Client) Check(ctx context.Context, user, relation, object string) (bool, error) {
    // Cache check primeiro
    // Chamada para OpenFGA com retry
    // Cache do resultado
    // Auditoria da operação
}
```

### 5. **Serviço Principal** (`service.go`)

#### Por que desenvolver:
- **Orquestração**: Combina cliente, cache e auditoria
- **Business Logic**: Implementa regras de negócio específicas
- **Validation**: Valida inputs e formatos
- **Error Handling**: Tratamento consistente de erros

#### O que será implementado:
```go
type Service struct {
    client *Client
    cache  Cache
    audit  AuditService
    config *Config
}

func (s *Service) CheckVaultPermission(ctx context.Context, userID, vaultID string, permission VaultPermission) (bool, error) {
    // Validação de inputs
    // Formatação de usuário/objeto
    // Verificação via cliente
    // Auditoria do resultado
    // Cache do resultado
}
```

### 6. **Middleware Gin** (`middleware.go`)

#### Por que desenvolver:
- **Integração**: Integra diretamente com framework Gin
- **Automatização**: Autorização automática em rotas
- **Flexibilidade**: Suporta diferentes tipos de verificação
- **Performance**: Evita código duplicado nos handlers

#### O que será implementado:
```go
type AuthorizationMiddleware struct {
    service LockariAuthorizationService
    logger  *slog.Logger
}

func (m *AuthorizationMiddleware) RequireVaultPermission(permission VaultPermission) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        vaultID := c.Param("vaultId")
        
        allowed, err := m.service.CanAccessVault(c.Request.Context(), userID, vaultID, permission)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "authorization_error"})
            return
        }
        
        if !allowed {
            c.AbortWithStatusJSON(403, gin.H{"error": "insufficient_permissions"})
            return
        }
        
        c.Next()
    }
}
```

### 7. **Cache System** (`cache.go`)

#### Por que desenvolver:
- **Performance**: Reduz latência em verificações frequentes
- **Scalability**: Reduz carga no OpenFGA
- **Reliability**: Fallback para operações críticas
- **Cost Optimization**: Reduz número de chamadas pagas

#### O que será implementado:
```go
type Cache interface {
    Get(key string) (bool, bool)
    Set(key string, value bool, ttl time.Duration)
    Delete(key string)
    Clear()
    Stats() CacheStats
}

type MemoryCache struct {
    mu    sync.RWMutex
    data  map[string]cacheEntry
    stats CacheStats
}

type cacheEntry struct {
    value     bool
    timestamp time.Time
    ttl       time.Duration
}
```

### 8. **Sistema de Auditoria** (`audit.go`)

#### Por que desenvolver:
- **Compliance**: Atende requisitos de auditoria
- **Security**: Rastreia tentativas de acesso
- **Debug**: Facilita troubleshooting
- **Analytics**: Permite análise de uso

#### O que será implementado:
```go
type AuditService interface {
    LogPermissionCheck(ctx context.Context, event PermissionCheckEvent)
    LogPermissionGrant(ctx context.Context, event PermissionGrantEvent)
    LogPermissionRevoke(ctx context.Context, event PermissionRevokeEvent)
    LogSuspiciousActivity(ctx context.Context, event SuspiciousActivityEvent)
}

type AuditEvent struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    UserID      string                 `json:"user_id"`
    Action      string                 `json:"action"`
    Resource    string                 `json:"resource"`
    Result      string                 `json:"result"`
    Metadata    map[string]interface{} `json:"metadata"`
    RequestID   string                 `json:"request_id"`
    ClientIP    string                 `json:"client_ip"`
    UserAgent   string                 `json:"user_agent"`
}
```

## 🔄 Fluxo de Desenvolvimento

### Fase 1: Fundação (Semana 1)
1. **Setup inicial**
   - Estrutura de arquivos
   - Configuração básica
   - Testes unitários iniciais

2. **Tipos e interfaces**
   - Definição de todas as interfaces
   - Tipos e enums
   - Validação de tipos

3. **Cliente básico**
   - Conexão com OpenFGA
   - Operações básicas (Check, Write, Delete)
   - Tratamento de erros

### Fase 2: Serviços Core (Semana 2)
1. **Serviço principal**
   - Implementação das interfaces
   - Lógica de negócio
   - Validações

2. **Cache system**
   - Cache em memória
   - Estratégias de invalidação
   - Métricas de performance

3. **Middleware Gin**
   - Integração com framework
   - Diferentes tipos de verificação
   - Error handling

### Fase 3: Recursos Avançados (Semana 3)
1. **Sistema de auditoria**
   - Logging estruturado
   - Eventos de auditoria
   - Integração com o8y

2. **Serviço específico Lockari**
   - Operações de alto nível
   - Compartilhamento externo
   - Tokens de API

3. **Otimizações**
   - Batch operations
   - Connection pooling
   - Retry logic avançado

### Fase 4: Testes e Documentação (Semana 4)
1. **Testes completos**
   - Testes unitários
   - Testes de integração
   - Benchmarks

2. **Documentação**
   - Exemplos de uso
   - Guides de implementação
   - Troubleshooting

3. **Integração com projeto**
   - Handlers existentes
   - Configuração
   - Deploy

## 🎯 Benefícios da Implementação

### 1. **Segurança**
- Autorização granular baseada em relacionamentos
- Isolamento completo entre tenants
- Auditoria completa de todas as operações
- Proteção contra privilege escalation

### 2. **Performance**
- Cache inteligente reduz latência
- Batch operations otimizam throughput
- Connection pooling melhora eficiência
- Retry logic garante confiabilidade

### 3. **Manutenibilidade**
- Código bem estruturado e documentado
- Interfaces claras permitem evolução
- Testes abrangentes garantem qualidade
- Logging facilita debug

### 4. **Flexibilidade**
- Suporta diferentes modelos de autorização
- Fácil extensão para novos recursos
- Configuração flexível por ambiente
- Integração simples com handlers existentes

## 📊 Métricas de Sucesso

### 1. **Funcionalidade**
- [ ] 100% das interfaces implementadas
- [ ] Suporte a todos os tipos de permissão
- [ ] Middleware funcionando em todas as rotas
- [ ] Auditoria completa implementada

### 2. **Performance**
- [ ] Cache hit rate > 80%
- [ ] Latência média < 50ms
- [ ] Throughput > 1000 req/s
- [ ] Zero downtime durante deploys

### 3. **Qualidade**
- [ ] Code coverage > 90%
- [ ] Zero bugs críticos
- [ ] Documentação completa
- [ ] Exemplos funcionais

### 4. **Integração**
- [ ] Todos os handlers usando middleware
- [ ] Configuração via environment
- [ ] Logging estruturado
- [ ] Métricas expostas

## 🚀 Próximos Passos Imediatos

1. **Criar estrutura básica de arquivos**
2. **Implementar interfaces e tipos**
3. **Desenvolver cliente OpenFGA básico**
4. **Criar testes unitários iniciais**
5. **Implementar configuração**
6. **Integrar com handlers existentes**

## 📋 Considerações Técnicas

### 1. **Dependências**
- OpenFGA SDK Go: já instalado (`v0.7.1`)
- Gin framework: já em uso
- Context: Para timeout e cancelamento
- Slog: Para logging estruturado

### 2. **Configuração**
- Usar Viper para carregar configurações
- Suportar variáveis de ambiente
- Validação de configuração obrigatória
- Defaults sensatos para desenvolvimento

### 3. **Testes**
- Usar testify para assertions
- Mocks para OpenFGA client
- Testes de integração com Docker
- Benchmarks para performance

### 4. **Documentação**
- Godoc para todas as funções públicas
- README com exemplos completos
- Guias de implementação
- Troubleshooting guide

Este plano fornece uma base sólida para implementar um sistema de autorização robusto, escalável e maintível para o projeto Lockari, aproveitando ao máximo as capacidades do OpenFGA enquanto mantém a simplicidade de uso para os desenvolvedores.
