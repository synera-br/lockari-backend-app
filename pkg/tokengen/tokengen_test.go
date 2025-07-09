package tokengen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenGenerator_Generate(t *testing.T) {
	tokenGen := NewTokenGenerator("test-secret", "test-issuer", time.Hour)

	t.Run("successful token generation", func(t *testing.T) {
		claims := TokenClaims{
			UserID:   "user123",
			TenantID: "tenant456",
			Scope:    []string{"read", "write"},
			Metadata: map[string]interface{}{
				"role": "user",
				"plan": "free",
			},
		}

		token, err := tokenGen.Generate(claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("token with custom expiration", func(t *testing.T) {
		claims := TokenClaims{
			UserID:    "user123",
			TenantID:  "tenant456",
			ExpiresAt: time.Now().Add(2 * time.Hour),
		}

		token, err := tokenGen.Generate(claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestTokenGenerator_Validate(t *testing.T) {
	tokenGen := NewTokenGenerator("test-secret", "test-issuer", time.Hour)

	t.Run("successful token validation", func(t *testing.T) {
		originalClaims := TokenClaims{
			UserID:   "user123",
			TenantID: "tenant456",
			Scope:    []string{"read", "write"},
			Metadata: map[string]interface{}{
				"role": "user",
				"plan": "free",
			},
		}

		token, err := tokenGen.Generate(originalClaims)
		require.NoError(t, err)

		validatedClaims, err := tokenGen.Validate(token)
		require.NoError(t, err)
		assert.Equal(t, originalClaims.UserID, validatedClaims.UserID)
		assert.Equal(t, originalClaims.TenantID, validatedClaims.TenantID)
		assert.Equal(t, originalClaims.Scope, validatedClaims.Scope)
		assert.Equal(t, originalClaims.Metadata["role"], validatedClaims.Metadata["role"])
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := tokenGen.Validate("")
		assert.Error(t, err)
		assert.True(t, IsTokenError(err))
		assert.Equal(t, "EMPTY_TOKEN", GetTokenErrorCode(err))
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := tokenGen.Validate("invalid-token")
		assert.Error(t, err)
		assert.True(t, IsTokenError(err))
	})

	t.Run("expired token", func(t *testing.T) {
		claims := TokenClaims{
			UserID:    "user123",
			ExpiresAt: time.Now().Add(-time.Hour), // Token expirado
		}

		token, err := tokenGen.Generate(claims)
		require.NoError(t, err)

		_, err = tokenGen.Validate(token)
		assert.Error(t, err)
		assert.True(t, IsTokenError(err))
		assert.Equal(t, "EXPIRED", GetTokenErrorCode(err))
	})

	t.Run("wrong secret", func(t *testing.T) {
		// Gerar token com um secret
		tokenGen1 := NewTokenGenerator("secret1", "test-issuer", time.Hour)
		claims := TokenClaims{UserID: "user123"}
		token, err := tokenGen1.Generate(claims)
		require.NoError(t, err)

		// Tentar validar com outro secret
		tokenGen2 := NewTokenGenerator("secret2", "test-issuer", time.Hour)
		_, err = tokenGen2.Validate(token)
		assert.Error(t, err)
		assert.True(t, IsTokenError(err))
	})
}

func TestTokenGenerator_Refresh(t *testing.T) {
	tokenGen := NewTokenGenerator("test-secret", "test-issuer", time.Hour)

	t.Run("successful token refresh", func(t *testing.T) {
		claims := TokenClaims{
			UserID:   "user123",
			TenantID: "tenant456",
			Scope:    []string{"read", "write"},
		}

		originalToken, err := tokenGen.Generate(claims)
		require.NoError(t, err)

		newToken, err := tokenGen.Refresh(originalToken, 2*time.Hour)
		require.NoError(t, err)
		assert.NotEmpty(t, newToken)
		assert.NotEqual(t, originalToken, newToken)

		// Validar que o novo token tem as mesmas claims
		newClaims, err := tokenGen.Validate(newToken)
		require.NoError(t, err)
		assert.Equal(t, claims.UserID, newClaims.UserID)
		assert.Equal(t, claims.TenantID, newClaims.TenantID)
		assert.Equal(t, claims.Scope, newClaims.Scope)
	})

	t.Run("refresh invalid token", func(t *testing.T) {
		_, err := tokenGen.Refresh("invalid-token", time.Hour)
		assert.Error(t, err)
		assert.True(t, IsTokenError(err))
		assert.Equal(t, "REFRESH_FAILED", GetTokenErrorCode(err))
	})
}

func TestTokenGenerator_IsExpired(t *testing.T) {
	tokenGen := NewTokenGenerator("test-secret", "test-issuer", time.Hour)

	t.Run("valid token not expired", func(t *testing.T) {
		claims := TokenClaims{
			UserID:    "user123",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		token, err := tokenGen.Generate(claims)
		require.NoError(t, err)

		assert.False(t, tokenGen.IsExpired(token))
	})

	t.Run("expired token", func(t *testing.T) {
		claims := TokenClaims{
			UserID:    "user123",
			ExpiresAt: time.Now().Add(-time.Hour),
		}

		token, err := tokenGen.Generate(claims)
		require.NoError(t, err)

		assert.True(t, tokenGen.IsExpired(token))
	})

	t.Run("invalid token", func(t *testing.T) {
		assert.True(t, tokenGen.IsExpired("invalid-token"))
	})
}

func TestTokenConfig_DefaultConfig(t *testing.T) {
	config := DefaultConfig("test-secret", "test-issuer")

	assert.Equal(t, "test-secret", config.Secret)
	assert.Equal(t, "test-issuer", config.Issuer)
	assert.Equal(t, time.Hour, config.DefaultDuration)
	assert.Equal(t, "HS256", config.Algorithm)
}

func TestTokenGenerator_PanicOnEmptySecret(t *testing.T) {
	assert.Panics(t, func() {
		NewTokenGenerator("", "test-issuer", time.Hour)
	})
}

func TestTokenGenerator_PanicOnEmptyIssuer(t *testing.T) {
	assert.Panics(t, func() {
		NewTokenGenerator("test-secret", "", time.Hour)
	})
}

func TestTokenClaims_SerializationDeserialization(t *testing.T) {
	tokenGen := NewTokenGenerator("test-secret", "test-issuer", time.Hour)

	originalClaims := TokenClaims{
		UserID:   "user123",
		TenantID: "tenant456",
		Scope:    []string{"read", "write", "signup"},
		Metadata: map[string]interface{}{
			"role":        "admin",
			"plan":        "premium",
			"permissions": []string{"create", "update", "delete"},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "pt-BR",
			},
		},
	}

	token, err := tokenGen.Generate(originalClaims)
	require.NoError(t, err)

	validatedClaims, err := tokenGen.Validate(token)
	require.NoError(t, err)

	assert.Equal(t, originalClaims.UserID, validatedClaims.UserID)
	assert.Equal(t, originalClaims.TenantID, validatedClaims.TenantID)
	assert.Equal(t, originalClaims.Scope, validatedClaims.Scope)
	assert.Equal(t, originalClaims.Metadata["role"], validatedClaims.Metadata["role"])
	assert.Equal(t, originalClaims.Metadata["plan"], validatedClaims.Metadata["plan"])
}
