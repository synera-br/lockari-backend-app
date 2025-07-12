# 🏗️ Análise: Armazenamento de Limites dos Planos

## 📊 Comparação: Arquivo vs Banco de Dados

### **🗂️ Armazenamento em Arquivo (Atual)**

#### ✅ **Vantagens:**
- **Performance**: Acesso instantâneo, sem queries ao banco
- **Simplicidade**: Fácil de gerenciar e versionar
- **Consistência**: Valores sempre disponíveis, mesmo sem conexão DB
- **Deployment**: Valores viajam com o código
- **Type Safety**: Validação em tempo de compilação

#### ❌ **Desvantagens:**
- **Flexibilidade**: Precisa redeploy para alterar limites
- **Personalização**: Difícil ter limites específicos por tenant
- **Configuração**: Valores fixos para todos os ambientes

### **🗄️ Armazenamento em Banco de Dados**

#### ✅ **Vantagens:**
- **Flexibilidade**: Alterar limites sem redeploy
- **Customização**: Limites específicos por tenant/cliente
- **Configuração**: Valores diferentes por ambiente
- **Auditoria**: Histórico de mudanças nos limites
- **Admin Panel**: Interface para ajustar limites

#### ❌ **Desvantagens:**
- **Performance**: Query adicional para verificar limites
- **Complexidade**: Mais código, cache, fallbacks
- **Dependência**: Precisa do banco sempre disponível
- **Consistência**: Risco de valores inconsistentes

## 🎯 **Recomendação: Abordagem Híbrida**

### **Estratégia Recomendada:**

```go
// 1. Valores padrão em arquivo (fallback)
type PlanLimits struct {
    VaultLimit int `json:"vault_limit"`
    UserLimit  int `json:"user_limit"`
    IsUnlimited bool `json:"is_unlimited"`
}

var DefaultPlanLimits = map[PlanType]PlanLimits{
    PlanFree: {
        VaultLimit: 3,
        UserLimit: 1,
        IsUnlimited: false,
    },
    PlanPro: {
        VaultLimit: 50,
        UserLimit: 10,
        IsUnlimited: false,
    },
    PlanEnterprise: {
        VaultLimit: 0, // 0 = unlimited
        UserLimit: 0,  // 0 = unlimited
        IsUnlimited: true,
    },
}

// 2. Override por banco (customização)
type TenantPlanOverride struct {
    TenantID    string     `json:"tenant_id"`
    PlanType    PlanType   `json:"plan_type"`
    Limits      PlanLimits `json:"limits"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

### **Implementação Prática:**

```go
type PlanLimitService struct {
    db    *sql.DB
    cache map[string]PlanLimits // Cache em memória
    mu    sync.RWMutex
}

func (s *PlanLimitService) GetPlanLimits(ctx context.Context, tenantID string, planType PlanType) (PlanLimits, error) {
    // 1. Verificar cache primeiro
    s.mu.RLock()
    if cached, exists := s.cache[tenantID]; exists {
        s.mu.RUnlock()
        return cached, nil
    }
    s.mu.RUnlock()

    // 2. Buscar override no banco
    override, err := s.getTenantPlanOverride(ctx, tenantID)
    if err == nil && override != nil {
        s.mu.Lock()
        s.cache[tenantID] = override.Limits
        s.mu.Unlock()
        return override.Limits, nil
    }

    // 3. Usar valores padrão do arquivo
    defaultLimits, exists := DefaultPlanLimits[planType]
    if !exists {
        return PlanLimits{}, fmt.Errorf("plan type not found: %s", planType)
    }

    // 4. Cache do valor padrão
    s.mu.Lock()
    s.cache[tenantID] = defaultLimits
    s.mu.Unlock()

    return defaultLimits, nil
}

func (s *PlanLimitService) SetCustomLimits(ctx context.Context, tenantID string, planType PlanType, limits PlanLimits) error {
    // 1. Salvar no banco
    override := &TenantPlanOverride{
        TenantID: tenantID,
        PlanType: planType,
        Limits:   limits,
        UpdatedAt: time.Now(),
    }
    
    err := s.saveTenantPlanOverride(ctx, override)
    if err != nil {
        return err
    }

    // 2. Atualizar cache
    s.mu.Lock()
    s.cache[tenantID] = limits
    s.mu.Unlock()

    return nil
}
```

## 📋 **Estrutura de Tabela Sugerida:**

```sql
-- Tabela para overrides personalizados
CREATE TABLE tenant_plan_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    plan_type VARCHAR(50) NOT NULL,
    vault_limit INTEGER,
    user_limit INTEGER,
    is_unlimited BOOLEAN DEFAULT FALSE,
    features JSONB, -- Features específicas
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, plan_type)
);

-- Índices para performance
CREATE INDEX idx_tenant_plan_overrides_tenant_id ON tenant_plan_overrides(tenant_id);
CREATE INDEX idx_tenant_plan_overrides_plan_type ON tenant_plan_overrides(plan_type);
```

## 🚀 **Implementação Completa:**

```go
// config/plan_limits.go
package config

type PlanLimits struct {
    VaultLimit      int            `json:"vault_limit"`
    UserLimit       int            `json:"user_limit"`
    IsUnlimited     bool           `json:"is_unlimited"`
    Features        []PlanFeature  `json:"features"`
    CustomFeatures  map[string]any `json:"custom_features,omitempty"`
}

var DefaultPlanLimits = map[PlanType]PlanLimits{
    PlanFree: {
        VaultLimit: 3,
        UserLimit: 1,
        IsUnlimited: false,
        Features: []PlanFeature{
            PlanFeatureBasic,
            PlanFeatureVaultLimit,
            PlanFeatureUserLimit,
        },
    },
    PlanPro: {
        VaultLimit: 50,
        UserLimit: 10,
        IsUnlimited: false,
        Features: []PlanFeature{
            PlanFeatureBasic,
            PlanFeatureVaultLimit,
            PlanFeatureUserLimit,
            PlanFeatureAdvancedPermissions,
            PlanFeatureAuditLogs,
            PlanFeatureAPIAccess,
            PlanFeatureGroupManagement,
        },
    },
    PlanEnterprise: {
        VaultLimit: 0, // 0 = unlimited
        UserLimit: 0,  // 0 = unlimited
        IsUnlimited: true,
        Features: []PlanFeature{
            PlanFeatureBasic,
            PlanFeatureUnlimitedVaults,
            PlanFeatureUnlimitedUsers,
            PlanFeatureAdvancedPermissions,
            PlanFeatureAuditLogs,
            PlanFeatureAPIAccess,
            PlanFeatureGroupManagement,
            PlanFeatureCrossTenantSharing,
            PlanFeatureExternalSharing,
            PlanFeatureSSO,
            PlanFeatureAdvancedSecurity,
        },
    },
}

// service/plan_limit_service.go
package service

type PlanLimitService struct {
    db           *sql.DB
    cache        map[string]PlanLimits
    cacheTTL     time.Duration
    cacheUpdated map[string]time.Time
    mu           sync.RWMutex
}

func NewPlanLimitService(db *sql.DB) *PlanLimitService {
    return &PlanLimitService{
        db:           db,
        cache:        make(map[string]PlanLimits),
        cacheTTL:     15 * time.Minute, // Cache por 15 minutos
        cacheUpdated: make(map[string]time.Time),
    }
}

func (s *PlanLimitService) GetPlanLimits(ctx context.Context, tenantID string, planType PlanType) (PlanLimits, error) {
    cacheKey := fmt.Sprintf("%s:%s", tenantID, planType)
    
    // 1. Verificar cache com TTL
    s.mu.RLock()
    if cached, exists := s.cache[cacheKey]; exists {
        if updatedAt, hasTime := s.cacheUpdated[cacheKey]; hasTime {
            if time.Since(updatedAt) < s.cacheTTL {
                s.mu.RUnlock()
                return cached, nil
            }
        }
    }
    s.mu.RUnlock()

    // 2. Buscar override no banco
    override, err := s.getTenantPlanOverride(ctx, tenantID, planType)
    if err == nil && override != nil {
        s.updateCache(cacheKey, override.Limits)
        return override.Limits, nil
    }

    // 3. Usar valores padrão
    defaultLimits, exists := DefaultPlanLimits[planType]
    if !exists {
        return PlanLimits{}, fmt.Errorf("plan type not found: %s", planType)
    }

    s.updateCache(cacheKey, defaultLimits)
    return defaultLimits, nil
}

func (s *PlanLimitService) updateCache(key string, limits PlanLimits) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.cache[key] = limits
    s.cacheUpdated[key] = time.Now()
}

func (s *PlanLimitService) IsWithinLimits(ctx context.Context, tenantID string, planType PlanType, resourceType string, currentCount int) (bool, error) {
    limits, err := s.GetPlanLimits(ctx, tenantID, planType)
    if err != nil {
        return false, err
    }

    if limits.IsUnlimited {
        return true, nil
    }

    switch resourceType {
    case "vault":
        return currentCount < limits.VaultLimit, nil
    case "user":
        return currentCount < limits.UserLimit, nil
    default:
        return false, fmt.Errorf("unknown resource type: %s", resourceType)
    }
}
```

## 🎯 **Vantagens da Abordagem Híbrida:**

1. **Performance**: Cache em memória + valores padrão em arquivo
2. **Flexibilidade**: Customização por tenant quando necessário
3. **Confiabilidade**: Fallback para valores padrão se banco falhar
4. **Simplicidade**: Maioria dos casos usa valores padrão
5. **Escalabilidade**: Apenas tenants com limites customizados vão ao banco

## 💡 **Quando Usar Cada Abordagem:**

### **Arquivo (Recomendado para início):**
- Produto novo/MVP
- Limites padronizados
- Equipe pequena
- Poucas mudanças nos planos

### **Banco (Recomendado para escala):**
- Produto maduro
- Clientes enterprise com necessidades específicas
- Equipe de suporte que ajusta limites
- Múltiplos ambientes com configurações diferentes

### **Híbrido (Melhor dos dois mundos):**
- Padrões em arquivo para performance
- Customizações no banco para flexibilidade
- Cache para otimização
- Fallback para confiabilidade

## 📝 **Conclusão:**

**Recomendo começar com a abordagem híbrida**, pois oferece:
- Simplicidade inicial (valores em arquivo)
- Flexibilidade futura (customização no banco)
- Performance otimizada (cache)
- Confiabilidade (fallback)

Isso permite evoluir o sistema conforme a necessidade sem refatoração drástica.
