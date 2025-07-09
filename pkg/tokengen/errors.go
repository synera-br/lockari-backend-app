package tokengen

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidToken    = errors.New("token is invalid")
	ErrExpiredToken    = errors.New("token has expired")
	ErrMalformedToken  = errors.New("token is malformed")
	ErrEmptySecret     = errors.New("secret cannot be empty")
	ErrEmptyIssuer     = errors.New("issuer cannot be empty")
	ErrInvalidClaims   = errors.New("invalid token claims")
	ErrTokenGeneration = errors.New("failed to generate token")
)

// TokenError representa um erro específico de token
type TokenError struct {
	Err     error
	Message string
	Code    string
}

func (te *TokenError) Error() string {
	return fmt.Sprintf("token error [%s]: %s - %v", te.Code, te.Message, te.Err)
}

// NewTokenError cria um novo erro de token
func NewTokenError(code, message string, err error) *TokenError {
	return &TokenError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// IsTokenError verifica se um erro é do tipo TokenError
func IsTokenError(err error) bool {
	_, ok := err.(*TokenError)
	return ok
}

// GetTokenErrorCode retorna o código do erro se for TokenError
func GetTokenErrorCode(err error) string {
	if te, ok := err.(*TokenError); ok {
		return te.Code
	}
	return ""
}
