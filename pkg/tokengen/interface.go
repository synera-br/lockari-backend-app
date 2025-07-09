package tokengen

import (
	"time"
)

// TokenGenerator interface principal para operações com tokens
type TokenGenerator interface {
	// Generate cria um novo token com as claims fornecidas
	Generate(claims TokenClaims) (string, error)

	// GenerateNonExpiring cria um token sem expiração
	GenerateNonExpiring(claims TokenClaims) (string, error)

	// GenerateLongLived cria um token com expiração muito longa (ex: 100 anos)
	GenerateLongLived(claims TokenClaims, duration time.Duration) (string, error)

	// Validate valida um token e retorna as claims se válido
	Validate(token string) (*TokenClaims, error)

	// Refresh cria um novo token baseado em um token válido existente
	Refresh(token string, newDuration time.Duration) (string, error)

	// IsExpired verifica se um token está expirado sem validar a assinatura
	IsExpired(token string) bool
}

// NewTokenGenerator cria uma nova instância do gerador de tokens
// secret: string secreta para assinar tokens (vem de variável de ambiente)
// issuer: identificador da aplicação (ex: "lockari-app")
// defaultDuration: duração padrão dos tokens (ex: 1 hora)
func NewTokenGenerator(secret, issuer string, defaultDuration time.Duration) TokenGenerator {
	return newJWTGenerator(secret, issuer, defaultDuration)
}
