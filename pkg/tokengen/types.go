package tokengen

import (
	"time"
)

// TokenClaims representa as informações contidas no token
type TokenClaims struct {
	UserID      string                 `json:"user_id"`
	TenantID    string                 `json:"tenant_id,omitempty"`
	Scope       []string               `json:"scope,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	IssuedAt    time.Time              `json:"iat"`
	ExpiresAt   time.Time              `json:"exp"`
	NonExpiring bool                   `json:"non_expiring,omitempty"`
	AppID       string                 `json:"app_id,omitempty"` // ID do aplicativo, se aplicável
}

// IsNonExpiring retorna true se o token não deve expirar
func (tc *TokenClaims) IsNonExpiring() bool {
	return tc.NonExpiring
}

// SetNonExpiring define se o token deve ser não-expirante
func (tc *TokenClaims) SetNonExpiring(nonExpiring bool) {
	tc.NonExpiring = nonExpiring
	if nonExpiring {
		// Remove a expiração
		tc.ExpiresAt = time.Time{}
	}
}

// TokenConfig configurações para geração de tokens
type TokenConfig struct {
	Secret          string
	Issuer          string
	DefaultDuration time.Duration
	Algorithm       string // "HS256", "HS384", "HS512"
}

// DefaultConfig retorna configuração padrão
func DefaultConfig(secret, issuer string) TokenConfig {
	return TokenConfig{
		Secret:          secret,
		Issuer:          issuer,
		DefaultDuration: time.Hour,
		Algorithm:       "HS256",
	}
}
