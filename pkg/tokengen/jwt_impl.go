package tokengen

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// jwtGenerator implementação da interface usando JWT
type jwtGenerator struct {
	config TokenConfig
}

// jwtClaims estrutura para claims JWT
type jwtClaims struct {
	UserID   string                 `json:"user_id"`
	TenantID string                 `json:"tenant_id,omitempty"`
	Scope    []string               `json:"scope,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

// newJWTGenerator cria nova instância do gerador JWT
func newJWTGenerator(secret, issuer string, defaultDuration time.Duration) TokenGenerator {
	if secret == "" {
		panic("secret cannot be empty")
	}
	if issuer == "" {
		panic("issuer cannot be empty")
	}

	return &jwtGenerator{
		config: TokenConfig{
			Secret:          secret,
			Issuer:          issuer,
			DefaultDuration: defaultDuration,
			Algorithm:       "HS256",
		},
	}
}

// Generate implementa TokenGenerator.Generate
func (jg *jwtGenerator) Generate(claims TokenClaims) (string, error) {
	now := time.Now()
	expiresAt := claims.ExpiresAt

	// Se não especificada, usa duração padrão
	if expiresAt.IsZero() && !claims.NonExpiring {
		expiresAt = now.Add(jg.config.DefaultDuration)
	}

	jwtClaims := &jwtClaims{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Scope:    claims.Scope,
		Metadata: claims.Metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jg.config.Issuer,
			Subject:   claims.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Só define expiração se não for non-expiring
	if !claims.NonExpiring && !expiresAt.IsZero() {
		jwtClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tokenString, err := token.SignedString([]byte(jg.config.Secret))
	if err != nil {
		return "", NewTokenError("GENERATION_FAILED", "Failed to sign token", err)
	}

	return tokenString, nil
}

// Validate implementa TokenGenerator.Validate
func (jg *jwtGenerator) Validate(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, NewTokenError("EMPTY_TOKEN", "Token is empty", ErrInvalidToken)
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jg.config.Secret), nil
	})

	if err != nil {
		// Verificar se é erro de expiração
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, NewTokenError("EXPIRED", "Token has expired", ErrExpiredToken)
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, NewTokenError("MALFORMED", "Token is malformed", ErrMalformedToken)
			}
		}
		return nil, NewTokenError("INVALID", "Token validation failed", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, NewTokenError("INVALID_CLAIMS", "Invalid token claims", ErrInvalidClaims)
	}

	// Verificar se é um token não-expirável
	isNonExpiring := claims.ExpiresAt == nil || claims.ExpiresAt.Time.IsZero()

	expiresAt := time.Time{}
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	return &TokenClaims{
		UserID:      claims.UserID,
		TenantID:    claims.TenantID,
		Scope:       claims.Scope,
		Metadata:    claims.Metadata,
		IssuedAt:    claims.IssuedAt.Time,
		ExpiresAt:   expiresAt,
		NonExpiring: isNonExpiring,
	}, nil
}

// Refresh implementa TokenGenerator.Refresh
func (jg *jwtGenerator) Refresh(tokenString string, newDuration time.Duration) (string, error) {
	claims, err := jg.Validate(tokenString)
	if err != nil {
		return "", NewTokenError("REFRESH_FAILED", "Cannot refresh invalid token", err)
	}

	// Criar novo token com nova expiração
	newClaims := TokenClaims{
		UserID:    claims.UserID,
		TenantID:  claims.TenantID,
		Scope:     claims.Scope,
		Metadata:  claims.Metadata,
		ExpiresAt: time.Now().Add(newDuration),
	}

	return jg.Generate(newClaims)
}

// IsExpired implementa TokenGenerator.IsExpired
func (jg *jwtGenerator) IsExpired(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jg.config.Secret), nil
	})

	if err != nil {
		return true // Se não conseguir parsear, considera expirado
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return true
	}

	// Se não tem ExpiresAt ou é zero, token não expira
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.IsZero() {
		return false
	}

	return time.Now().After(claims.ExpiresAt.Time)
}

// GenerateNonExpiring implementa TokenGenerator.GenerateNonExpiring
func (jg *jwtGenerator) GenerateNonExpiring(claims TokenClaims) (string, error) {
	claims.NonExpiring = true
	claims.ExpiresAt = time.Time{} // Remove expiração
	return jg.Generate(claims)
}

// GenerateLongLived implementa TokenGenerator.GenerateLongLived
func (jg *jwtGenerator) GenerateLongLived(claims TokenClaims, duration time.Duration) (string, error) {
	claims.NonExpiring = false
	claims.ExpiresAt = time.Now().Add(duration)
	return jg.Generate(claims)
}
