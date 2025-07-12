package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/synera-br/lockari-backend-app/pkg/authenticator"
	"github.com/synera-br/lockari-backend-app/pkg/logger"
	"github.com/synera-br/lockari-backend-app/pkg/tokengen"
)

// contextKey is a type for context keys.
type contextKey string

const (
	// UserIDKey is the context key for the user ID.
	UserIDKey = contextKey("userID")
	// TenantIDKey is the context key for the tenant ID.
	TenantIDKey = contextKey("tenantID")
	// ClaimsKey is the context key for the claims.
	ClaimsKey = contextKey("claims")
)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain creates a new middleware chain.
func Chain(middlewares ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}

// Auth validates the token and adds the user information to the context.
func Auth(auth authenticator.Authenticator, log logger.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("X-AUTHORIZATION")
			if authHeader == "" {
				log.Error("Authorization header is missing")
				http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
				return
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			if idToken == authHeader {
				log.Error("Bearer token not found in Authorization header")
				http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateToken(r.Context(), idToken)
			if err != nil {
				log.Errorf("Token validation failed: %v", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			uid, ok := claims["sub"].(string)
			if !ok || uid == "" {
				log.Error("UID not found in token claims")
				http.Error(w, "User identifier not found in token", http.StatusUnauthorized)
				return
			}

			tenantID, ok := claims["tenantId"].(string)
			if !ok || tenantID == "" {
				log.Errorf("Custom claim 'tenantId' not found for user %s", uid)
				http.Error(w, "User is not associated with a tenant", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			ctx = context.WithValue(ctx, TenantIDKey, tenantID)
			ctx = context.WithValue(ctx, ClaimsKey, claims)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthJWT validates the JWT and adds the claims to the context.
func AuthJWT(token tokengen.TokenGenerator, log logger.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("X-TOKEN")
			if authHeader == "" {
				log.Error("X-TOKEN header is missing")
				http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
				return
			}

			claims, err := token.Validate(authHeader)
			if err != nil {
				log.Errorf("Token validation failed: %v", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
