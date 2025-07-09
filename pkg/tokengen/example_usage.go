package tokengen

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Exemplo prático de uso de tokens não-expiráveis
func ExampleUsage() {
	// 1. Configurar o gerador de tokens
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "demo-secret-key-for-testing-only"
		log.Println("Warning: Using demo secret. Set JWT_SECRET environment variable for production.")
	}

	tokenGen := NewTokenGenerator(secret, "lockari-demo", time.Hour)

	// 2. Exemplo de Service Account Token (não expira)
	fmt.Println("=== Service Account Token (Non-Expiring) ===")
	serviceAccountToken := generateServiceAccountToken(tokenGen)
	validateAndUseToken(tokenGen, serviceAccountToken, "Service Account")

	// 3. Exemplo de Mobile App Token (1 ano)
	fmt.Println("\n=== Mobile App Token (1 Year) ===")
	mobileToken := generateMobileAppToken(tokenGen)
	validateAndUseToken(tokenGen, mobileToken, "Mobile App")

	// 4. Exemplo de Token Normal (1 hora)
	fmt.Println("\n=== Normal Token (1 Hour) ===")
	normalToken := generateNormalToken(tokenGen)
	validateAndUseToken(tokenGen, normalToken, "Normal User")

	// 5. Demonstrar diferenças
	fmt.Println("\n=== Comparison ===")
	compareTokens(tokenGen, serviceAccountToken, mobileToken, normalToken)
}

func generateServiceAccountToken(tokenGen TokenGenerator) string {
	claims := TokenClaims{
		UserID:   "service-monitoring",
		TenantID: "system",
		Scope:    []string{"metrics:read", "logs:read", "alerts:write", "admin"},
		Metadata: map[string]interface{}{
			"service_type": "monitoring",
			"version":      "2.1.0",
			"environment":  "production",
			"department":   "DevOps",
		},
	}

	token, err := tokenGen.GenerateNonExpiring(claims)
	if err != nil {
		log.Fatalf("Failed to generate service account token: %v", err)
	}

	fmt.Printf("Service Account Token: %s...\n", token[:50])
	return token
}

func generateMobileAppToken(tokenGen TokenGenerator) string {
	claims := TokenClaims{
		UserID:   "user-mobile-456",
		TenantID: "tenant-xyz",
		Scope:    []string{"profile", "data", "sync", "notifications"},
		Metadata: map[string]interface{}{
			"device_id":   "iPhone-13-XYZ789",
			"app_version": "3.2.1",
			"platform":    "iOS",
			"device_type": "mobile",
		},
	}

	// Token que dura 1 ano
	token, err := tokenGen.GenerateLongLived(claims, time.Hour*24*365)
	if err != nil {
		log.Fatalf("Failed to generate mobile app token: %v", err)
	}

	fmt.Printf("Mobile App Token: %s...\n", token[:50])
	return token
}

func generateNormalToken(tokenGen TokenGenerator) string {
	claims := TokenClaims{
		UserID:   "user-web-789",
		TenantID: "tenant-abc",
		Scope:    []string{"profile", "data"},
		Metadata: map[string]interface{}{
			"session_id":  "session-123456",
			"user_agent":  "Mozilla/5.0...",
			"device_type": "web",
		},
	}

	token, err := tokenGen.Generate(claims)
	if err != nil {
		log.Fatalf("Failed to generate normal token: %v", err)
	}

	fmt.Printf("Normal Token: %s...\n", token[:50])
	return token
}

func validateAndUseToken(tokenGen TokenGenerator, token, tokenType string) {
	// Validar token
	claims, err := tokenGen.Validate(token)
	if err != nil {
		log.Printf("❌ %s token validation failed: %v", tokenType, err)
		return
	}

	fmt.Printf("✅ %s token is valid\n", tokenType)
	fmt.Printf("   User ID: %s\n", claims.UserID)
	fmt.Printf("   Tenant ID: %s\n", claims.TenantID)
	fmt.Printf("   Scope: %v\n", claims.Scope)
	fmt.Printf("   Non-Expiring: %v\n", claims.IsNonExpiring())

	if !claims.ExpiresAt.IsZero() {
		fmt.Printf("   Expires At: %v\n", claims.ExpiresAt)
	} else {
		fmt.Printf("   Expires At: NEVER\n")
	}

	fmt.Printf("   Is Expired: %v\n", tokenGen.IsExpired(token))

	// Exemplo de uso das claims
	if deviceType, exists := claims.Metadata["device_type"]; exists {
		fmt.Printf("   Device Type: %s\n", deviceType)
	}

	// Verificar permissões
	hasAdminScope := false
	for _, scope := range claims.Scope {
		if scope == "admin" {
			hasAdminScope = true
			break
		}
	}

	if hasAdminScope {
		fmt.Printf("   🔐 Has admin permissions\n")
	}
}

func compareTokens(tokenGen TokenGenerator, serviceToken, mobileToken, normalToken string) {
	tokens := map[string]string{
		"Service Account": serviceToken,
		"Mobile App":      mobileToken,
		"Normal User":     normalToken,
	}

	fmt.Printf("%-15s | %-12s | %-20s | %-10s\n", "Token Type", "Non-Expiring", "Expires At", "Is Expired")
	fmt.Println(strings.Repeat("-", 70))

	for tokenType, token := range tokens {
		claims, err := tokenGen.Validate(token)
		if err != nil {
			fmt.Printf("%-15s | %-12s | %-20s | %-10s\n", tokenType, "ERROR", "ERROR", "ERROR")
			continue
		}

		nonExpiring := "No"
		if claims.IsNonExpiring() {
			nonExpiring = "Yes"
		}

		expiresAt := "NEVER"
		if !claims.ExpiresAt.IsZero() {
			expiresAt = claims.ExpiresAt.Format("2006-01-02 15:04")
		}

		isExpired := "No"
		if tokenGen.IsExpired(token) {
			isExpired = "Yes"
		}

		fmt.Printf("%-15s | %-12s | %-20s | %-10s\n", tokenType, nonExpiring, expiresAt, isExpired)
	}
}

// Exemplo de uso em um handler HTTP
func exampleHTTPHandler() {
	/*
		func AuthenticatedHandler(tokenGen TokenGenerator) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				// Extrair token do header
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
					return
				}

				token := strings.TrimPrefix(authHeader, "Bearer ")

				// Validar token
				claims, err := tokenGen.Validate(token)
				if err != nil {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}

				// Log para tokens não-expiráveis
				if claims.IsNonExpiring() {
					log.Printf("Non-expiring token used by user: %s", claims.UserID)
				}

				// Processar requisição com as claims
				response := map[string]interface{}{
					"message":       "Success",
					"user_id":       claims.UserID,
					"tenant_id":     claims.TenantID,
					"scope":         claims.Scope,
					"non_expiring":  claims.IsNonExpiring(),
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	*/
}
