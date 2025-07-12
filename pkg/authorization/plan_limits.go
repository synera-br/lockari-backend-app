package authorization

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// PlanLimits representa os limites de um plano
type PlanLimits struct {
	VaultLimit     int                    `json:"vault_limit"`
	UserLimit      int                    `json:"user_limit"`
	IsUnlimited    bool                   `json:"is_unlimited"`
	Features       []PlanFeature          `json:"features"`
	CustomFeatures map[string]interface{} `json:"custom_features,omitempty"`
	Description    string                 `json:"description,omitempty"`
}

// DefaultPlanLimits define os limites padrão por plano
var DefaultPlanLimits = map[PlanType]PlanLimits{
	PlanFree: {
		VaultLimit:  3,
		UserLimit:   1,
		IsUnlimited: false,
		Features: []PlanFeature{
			PlanFeatureBasic,
			PlanFeatureVaultLimit,
			PlanFeatureUserLimit,
		},
		Description: "Plano gratuito com recursos básicos",
	},
	PlanPro: {
		VaultLimit:  50,
		UserLimit:   10,
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
		Description: "Plano profissional com recursos avançados",
	},
	PlanEnterprise: {
		VaultLimit:  0, // 0 = unlimited
		UserLimit:   0, // 0 = unlimited
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
		Description: "Plano enterprise com recursos ilimitados",
	},
}

// TenantPlanOverride representa customizações específicas por tenant
type TenantPlanOverride struct {
	ID        string     `json:"id" db:"id"`
	TenantID  string     `json:"tenant_id" db:"tenant_id"`
	PlanType  PlanType   `json:"plan_type" db:"plan_type"`
	Limits    PlanLimits `json:"limits" db:"limits"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// PlanLimitService gerencia os limites dos planos
type PlanLimitService struct {
	db           *sql.DB
	cache        map[string]PlanLimits
	cacheTTL     time.Duration
	cacheUpdated map[string]time.Time
	mu           sync.RWMutex
}

// NewPlanLimitService cria uma nova instância do serviço
func NewPlanLimitService(db *sql.DB) *PlanLimitService {
	return &PlanLimitService{
		db:           db,
		cache:        make(map[string]PlanLimits),
		cacheTTL:     15 * time.Minute, // Cache por 15 minutos
		cacheUpdated: make(map[string]time.Time),
	}
}

// GetPlanLimits obtém os limites de um plano, considerando customizações
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

// SetCustomLimits define limites customizados para um tenant
func (s *PlanLimitService) SetCustomLimits(ctx context.Context, tenantID string, planType PlanType, limits PlanLimits) error {
	// 1. Salvar no banco
	override := &TenantPlanOverride{
		TenantID:  tenantID,
		PlanType:  planType,
		Limits:    limits,
		UpdatedAt: time.Now(),
	}

	err := s.saveTenantPlanOverride(ctx, override)
	if err != nil {
		return fmt.Errorf("failed to save tenant plan override: %w", err)
	}

	// 2. Atualizar cache
	cacheKey := fmt.Sprintf("%s:%s", tenantID, planType)
	s.updateCache(cacheKey, limits)

	return nil
}

// IsWithinLimits verifica se um tenant está dentro dos limites
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

// GetRemainingLimits retorna quantos recursos ainda podem ser criados
func (s *PlanLimitService) GetRemainingLimits(ctx context.Context, tenantID string, planType PlanType, resourceType string, currentCount int) (int, error) {
	limits, err := s.GetPlanLimits(ctx, tenantID, planType)
	if err != nil {
		return 0, err
	}

	if limits.IsUnlimited {
		return -1, nil // -1 indica ilimitado
	}

	switch resourceType {
	case "vault":
		remaining := limits.VaultLimit - currentCount
		if remaining < 0 {
			return 0, nil
		}
		return remaining, nil
	case "user":
		remaining := limits.UserLimit - currentCount
		if remaining < 0 {
			return 0, nil
		}
		return remaining, nil
	default:
		return 0, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

// HasFeature verifica se um plano tem uma funcionalidade específica
func (s *PlanLimitService) HasFeature(ctx context.Context, tenantID string, planType PlanType, feature PlanFeature) (bool, error) {
	limits, err := s.GetPlanLimits(ctx, tenantID, planType)
	if err != nil {
		return false, err
	}

	for _, f := range limits.Features {
		if f == feature {
			return true, nil
		}
	}

	return false, nil
}

// updateCache atualiza o cache com timestamp
func (s *PlanLimitService) updateCache(key string, limits PlanLimits) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = limits
	s.cacheUpdated[key] = time.Now()
}

// clearCache limpa o cache para um tenant específico
func (s *PlanLimitService) ClearCache(tenantID string, planType PlanType) {
	cacheKey := fmt.Sprintf("%s:%s", tenantID, planType)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.cache, cacheKey)
	delete(s.cacheUpdated, cacheKey)
}

// getTenantPlanOverride busca customizações no banco
func (s *PlanLimitService) getTenantPlanOverride(ctx context.Context, tenantID string, planType PlanType) (*TenantPlanOverride, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not available")
	}

	query := `
		SELECT id, tenant_id, plan_type, limits, created_at, updated_at
		FROM tenant_plan_overrides
		WHERE tenant_id = $1 AND plan_type = $2
		ORDER BY updated_at DESC
		LIMIT 1
	`

	var override TenantPlanOverride
	var limitsJSON []byte

	err := s.db.QueryRowContext(ctx, query, tenantID, planType).Scan(
		&override.ID,
		&override.TenantID,
		&override.PlanType,
		&limitsJSON,
		&override.CreatedAt,
		&override.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Não é erro, só não tem customização
		}
		return nil, fmt.Errorf("failed to query tenant plan override: %w", err)
	}

	err = json.Unmarshal(limitsJSON, &override.Limits)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal limits: %w", err)
	}

	return &override, nil
}

// saveTenantPlanOverride salva customizações no banco
func (s *PlanLimitService) saveTenantPlanOverride(ctx context.Context, override *TenantPlanOverride) error {
	if s.db == nil {
		return fmt.Errorf("database not available")
	}

	limitsJSON, err := json.Marshal(override.Limits)
	if err != nil {
		return fmt.Errorf("failed to marshal limits: %w", err)
	}

	query := `
		INSERT INTO tenant_plan_overrides (tenant_id, plan_type, limits, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (tenant_id, plan_type)
		DO UPDATE SET
			limits = EXCLUDED.limits,
			updated_at = NOW()
		RETURNING id
	`

	err = s.db.QueryRowContext(ctx, query, override.TenantID, override.PlanType, limitsJSON).Scan(&override.ID)
	if err != nil {
		return fmt.Errorf("failed to save tenant plan override: %w", err)
	}

	return nil
}

// GetAllPlanTypes retorna todos os tipos de plano disponíveis
func GetAllPlanTypes() []PlanType {
	return []PlanType{PlanFree, PlanPro, PlanEnterprise}
}

// GetPlanDescription retorna a descrição de um plano
func GetPlanDescription(planType PlanType) string {
	if limits, exists := DefaultPlanLimits[planType]; exists {
		return limits.Description
	}
	return ""
}

// ValidatePlanType verifica se um tipo de plano é válido
func ValidatePlanType(planType PlanType) bool {
	_, exists := DefaultPlanLimits[planType]
	return exists
}
