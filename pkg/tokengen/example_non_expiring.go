package tokengen

import (
	"fmt"
	"time"
)

// ExampleNonExpiringTokens demonstra como usar tokens não-expiráveis
func ExampleNonExpiringTokens() {

	// 1. Inicializar o gerador de tokens
	tokenGen := NewTokenGenerator(
		"your-secret-key-here",
		"lockari-app",
		time.Hour, // Duração padrão para tokens normais
	)

	// 2. Criar claims para um service account
	serviceAccountClaims := TokenClaims{
		UserID:   "service-account-monitoring",
		TenantID: "internal",
		Scope:    []string{"metrics", "logs", "alerts"},
		Metadata: map[string]interface{}{
			"type":        "service_account",
			"department":  "DevOps",
			"permissions": []string{"read", "write", "admin"},
		},
	}

	// 3. Gerar token SEM expiração
	nonExpiringToken, err := tokenGen.GenerateNonExpiring(serviceAccountClaims)
	if err != nil {
		fmt.Printf("Erro ao gerar token não-expirável: %v\n", err)
		return
	}

	fmt.Printf("Token não-expirável gerado: %s\n", nonExpiringToken[:50]+"...")

	// 4. Validar o token
	validatedClaims, err := tokenGen.Validate(nonExpiringToken)
	if err != nil {
		fmt.Printf("Erro ao validar token: %v\n", err)
		return
	}

	fmt.Printf("Token validado para usuário: %s\n", validatedClaims.UserID)
	fmt.Printf("Token é não-expirável: %v\n", validatedClaims.IsNonExpiring())
	fmt.Printf("Token expirou: %v\n", tokenGen.IsExpired(nonExpiringToken))

	// 5. Exemplo de token com expiração muito longa (10 anos)
	longLivedClaims := TokenClaims{
		UserID:   "mobile-app-user",
		TenantID: "tenant123",
		Scope:    []string{"profile", "data"},
		Metadata: map[string]interface{}{
			"device_type": "mobile",
			"app_version": "1.0.0",
		},
	}

	longLivedToken, err := tokenGen.GenerateLongLived(longLivedClaims, time.Hour*24*365*10)
	if err != nil {
		fmt.Printf("Erro ao gerar token de longa duração: %v\n", err)
		return
	}

	fmt.Printf("Token de longa duração gerado: %s\n", longLivedToken[:50]+"...")

	// 6. Validar token de longa duração
	longLivedValidated, err := tokenGen.Validate(longLivedToken)
	if err != nil {
		fmt.Printf("Erro ao validar token de longa duração: %v\n", err)
		return
	}

	fmt.Printf("Token de longa duração expira em: %v\n", longLivedValidated.ExpiresAt)
	fmt.Printf("Token de longa duração é não-expirável: %v\n", longLivedValidated.IsNonExpiring())
}

// ExampleServiceAccountFlow demonstra um fluxo completo para service accounts
func ExampleServiceAccountFlow() {
	// Para este exemplo, simule que o secret vem de uma variável de ambiente
	secret := "your-super-secret-key-for-production"

	tokenGen := NewTokenGenerator(secret, "lockari-backend", time.Hour)

	// Service account para monitoramento
	monitoringClaims := TokenClaims{
		UserID:   "monitoring-service",
		TenantID: "system",
		Scope:    []string{"metrics:read", "logs:read", "alerts:write"},
		Metadata: map[string]interface{}{
			"service_type": "monitoring",
			"version":      "2.1.0",
			"environment":  "production",
		},
	}

	// Gerar token que nunca expira
	serviceToken, err := tokenGen.GenerateNonExpiring(monitoringClaims)
	if err != nil {
		fmt.Printf("Falha ao gerar token do service account: %v\n", err)
		return
	}

	fmt.Printf("Token do service account gerado com sucesso!\n")
	fmt.Printf("Token: %s...\n", serviceToken[:30])

	// Simular validação em uma requisição
	fmt.Println("\n=== Simulando validação de requisição ===")

	claims, err := tokenGen.Validate(serviceToken)
	if err != nil {
		fmt.Printf("Token inválido: %v\n", err)
		return
	}

	// Verificar permissões
	hasMetricsPermission := false
	for _, scope := range claims.Scope {
		if scope == "metrics:read" {
			hasMetricsPermission = true
			break
		}
	}

	if hasMetricsPermission {
		fmt.Printf("✅ Service account '%s' autorizado para ler métricas\n", claims.UserID)
	} else {
		fmt.Printf("❌ Service account '%s' não tem permissão para ler métricas\n", claims.UserID)
	}

	// Verificar se o token está expirado
	if tokenGen.IsExpired(serviceToken) {
		fmt.Println("❌ Token expirado!")
	} else {
		fmt.Println("✅ Token válido e não expirado")
	}
}

// ExampleMobileAppFlow demonstra fluxo para apps mobile com tokens de longa duração
func ExampleMobileAppFlow() {
	tokenGen := NewTokenGenerator("mobile-app-secret", "lockari-mobile", time.Hour*24)

	// Claims para usuário mobile
	userClaims := TokenClaims{
		UserID:   "user-mobile-123",
		TenantID: "tenant-abc",
		Scope:    []string{"profile", "data", "sync"},
		Metadata: map[string]interface{}{
			"device_id":   "iPhone-12-Pro-ABC123",
			"app_version": "2.1.0",
			"platform":    "iOS",
		},
	}

	// Token que dura 1 ano para apps mobile
	mobileToken, err := tokenGen.GenerateLongLived(userClaims, time.Hour*24*365)
	if err != nil {
		fmt.Printf("Erro ao gerar token mobile: %v\n", err)
		return
	}

	fmt.Printf("Token mobile gerado (válido por 1 ano): %s...\n", mobileToken[:30])

	// Validar token
	validatedClaims, err := tokenGen.Validate(mobileToken)
	if err != nil {
		fmt.Printf("Token mobile inválido: %v\n", err)
		return
	}

	fmt.Printf("Token mobile válido para usuário: %s\n", validatedClaims.UserID)
	fmt.Printf("Expira em: %v\n", validatedClaims.ExpiresAt)
	fmt.Printf("Dispositivo: %v\n", validatedClaims.Metadata["device_id"])
}
