package authorization

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestIntegration tests the full integration flow
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This is a basic integration test that would work with a real OpenFGA instance
	// For now, we'll test the configuration and basic setup

	config := &Config{
		APIURL:                "localhost:8080",
		StoreID:               "test-store",
		AuthorizationModelID:  "test-model",
		Timeout:               30 * time.Second,
		RetryAttempts:         3,
		RetryDelay:            time.Second,
		MaxRetryDelay:         10 * time.Second,
		CacheEnabled:          true,
		CacheTTL:              15 * time.Minute,
		CacheMaxSize:          10000,
		CacheCleanupInterval:  5 * time.Minute,
		AuditEnabled:          true,
		AuditLevel:            "info",
		MaxConcurrentRequests: 100,
		BatchSize:             50,
		ConnectionPoolSize:    10,
		HealthCheckEnabled:    true,
		HealthCheckInterval:   30 * time.Second,
		HealthCheckTimeout:    5 * time.Second,
	}

	// Test configuration validation
	err := config.Validate()
	assert.NoError(t, err)

	// Test basic configuration properties
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 3, config.RetryAttempts)
	assert.True(t, config.CacheEnabled)
	assert.Equal(t, 15*time.Minute, config.CacheTTL)
	assert.Equal(t, 10000, config.CacheMaxSize)

	// Test helper functions
	userID := "alice"
	vaultID := "test-vault"
	secretID := "test-secret"
	tenantID := "test-tenant"
	tokenID := "test-token"
	groupID := "test-group"

	assert.Equal(t, "user:alice", FormatUser(userID))
	assert.Equal(t, "vault:test-vault", FormatVault(vaultID))
	assert.Equal(t, "secret:test-secret", FormatSecret(secretID))
	assert.Equal(t, "tenant:test-tenant", FormatTenant(tenantID))
	assert.Equal(t, "token:test-token", FormatToken(tokenID))
	assert.Equal(t, "group:test-group", FormatGroup(groupID))

	// Test parsing
	objectType, objectID, err := ParseObject("vault:test-vault")
	assert.NoError(t, err)
	assert.Equal(t, "vault", objectType)
	assert.Equal(t, "test-vault", objectID)

	parsedUserID, err := ParseUser("user:alice")
	assert.NoError(t, err)
	assert.Equal(t, "alice", parsedUserID)

	// Test audit event creation
	event := NewAuditEvent(userID, "check_permission", "vault:test-vault", true, map[string]interface{}{
		"permission": "can_read",
		"tenant":     tenantID,
	})

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, "check_permission", event.Action)
	assert.Equal(t, "vault:test-vault", event.Resource)
	assert.Equal(t, "allowed", event.Result)
	assert.Equal(t, "can_read", event.Metadata["permission"])
	assert.Equal(t, tenantID, event.Metadata["tenant"])
	assert.NotZero(t, event.Timestamp)
	assert.NotEmpty(t, event.ID)

	// Test token generation
	token1 := GenerateTokenID()
	token2 := GenerateTokenID()

	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
	assert.Len(t, token1, 32) // UUID without dashes
	assert.Len(t, token2, 32)

	// Test permission validation
	assert.True(t, VaultPermissionRead.IsValid())
	assert.True(t, SecretPermissionWrite.IsValid())
	assert.True(t, TenantRoleAdmin.IsValid())
	assert.True(t, GroupRoleOwner.IsValid())
	assert.True(t, PlanFeatureAdvancedPermissions.IsValid())
	assert.True(t, TokenPermissionUse.IsValid())
	assert.True(t, TenantPermissionAdmin.IsValid())

	// Test request validation
	checkReq := &CheckRequest{
		User:     FormatUser("alice"),
		Relation: "can_read",
		Object:   FormatVault("test-vault"),
	}

	err = checkReq.Validate()
	assert.NoError(t, err)

	// Test invalid request
	invalidReq := &CheckRequest{
		User:     "",
		Relation: "can_read",
		Object:   FormatVault("test-vault"),
	}

	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user cannot be empty")

	// Test tuple validation
	tuple := &Tuple{
		User:     FormatUser("alice"),
		Relation: "can_read",
		Object:   FormatVault("test-vault"),
	}

	err = tuple.Validate()
	assert.NoError(t, err)

	// Test tuple string representation
	tupleStr := tuple.String()
	assert.Equal(t, "user:alice#can_read@vault:test-vault", tupleStr)

	// Test cache entry
	cacheEntry := &CacheEntry{
		Key:       "test-key",
		Value:     true,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}

	assert.False(t, cacheEntry.IsExpired())

	// Test expired cache entry
	expiredEntry := &CacheEntry{
		Key:       "test-key",
		Value:     true,
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}

	assert.True(t, expiredEntry.IsExpired())

	// Test write request validation
	writeReq := &WriteRequest{
		Tuples: []Tuple{*tuple},
	}

	err = writeReq.Validate()
	assert.NoError(t, err)

	// Test delete request validation
	deleteReq := &DeleteRequest{
		Tuples: []Tuple{*tuple},
	}

	err = deleteReq.Validate()
	assert.NoError(t, err)

	// Test list objects request validation
	listReq := &ListObjectsRequest{
		User:     FormatUser("alice"),
		Relation: "can_read",
		Type:     "vault",
		Context:  map[string]interface{}{"tenant": tenantID},
	}

	err = listReq.Validate()
	assert.NoError(t, err)

	t.Log("All integration tests passed!")
}

// TestPermissionHierarchy tests the permission hierarchy logic
func TestPermissionHierarchy(t *testing.T) {
	// Test vault permission hierarchy
	vaultPermissions := AllVaultPermissions()
	assert.Contains(t, vaultPermissions, VaultPermissionView)
	assert.Contains(t, vaultPermissions, VaultPermissionRead)
	assert.Contains(t, vaultPermissions, VaultPermissionWrite)
	assert.Contains(t, vaultPermissions, VaultPermissionManage)

	// Test secret permission hierarchy
	secretPermissions := AllSecretPermissions()
	assert.Contains(t, secretPermissions, SecretPermissionView)
	assert.Contains(t, secretPermissions, SecretPermissionRead)
	assert.Contains(t, secretPermissions, SecretPermissionWrite)
	assert.Contains(t, secretPermissions, SecretPermissionReadSensitive)
	assert.Contains(t, secretPermissions, SecretPermissionCopyProduction)

	// Test tenant roles
	tenantRoles := AllTenantRoles()
	assert.Contains(t, tenantRoles, TenantRoleOwner)
	assert.Contains(t, tenantRoles, TenantRoleAdmin)
	assert.Contains(t, tenantRoles, TenantRoleMember)
	assert.Contains(t, tenantRoles, TenantRoleGuest)

	// Test group roles
	groupRoles := AllGroupRoles()
	assert.Contains(t, groupRoles, GroupRoleOwner)
	assert.Contains(t, groupRoles, GroupRoleAdmin)
	assert.Contains(t, groupRoles, GroupRoleMember)

	// Test plan features
	planFeatures := AllPlanFeatures()
	assert.Contains(t, planFeatures, PlanFeatureBasic)
	assert.Contains(t, planFeatures, PlanFeatureAdvancedPermissions)
	assert.Contains(t, planFeatures, PlanFeatureCrossTenantSharing)
	assert.Contains(t, planFeatures, PlanFeatureExternalSharing)

	// Test token permissions
	tokenPermissions := AllTokenPermissions()
	assert.Contains(t, tokenPermissions, TokenPermissionUse)
	assert.Contains(t, tokenPermissions, TokenPermissionReadSecrets)
	assert.Contains(t, tokenPermissions, TokenPermissionWriteSecrets)
	assert.Contains(t, tokenPermissions, TokenPermissionRevoke)

	// Test tenant permissions
	tenantPermissions := AllTenantPermissions()
	assert.Contains(t, tenantPermissions, TenantPermissionView)
	assert.Contains(t, tenantPermissions, TenantPermissionCreateVault)
	assert.Contains(t, tenantPermissions, TenantPermissionManageUsers)
	assert.Contains(t, tenantPermissions, TenantPermissionOwner)
	assert.Contains(t, tenantPermissions, TenantPermissionAdmin)
}

// TestErrorScenarios tests various error scenarios
func TestErrorScenarios(t *testing.T) {
	// Test invalid object parsing
	_, _, err := ParseObject("invalid-object")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid object format")

	// Test invalid user parsing
	_, err = ParseUser("invalid-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user format")

	// Test invalid check request
	invalidReq := &CheckRequest{
		User:     "invalid-user",
		Relation: "can_read",
		Object:   "vault:test-vault",
	}

	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user must start with 'user:' or 'token:'")

	// Test invalid check request with malformed object
	invalidReq2 := &CheckRequest{
		User:     "user:alice",
		Relation: "can_read",
		Object:   "invalid-object",
	}

	err = invalidReq2.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "object must be in format 'type:id'")

	// Test empty write request
	emptyWriteReq := &WriteRequest{
		Tuples: []Tuple{},
	}

	err = emptyWriteReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tuples cannot be empty")

	// Test invalid tuple in write request
	invalidWriteReq := &WriteRequest{
		Tuples: []Tuple{
			{
				User:     "",
				Relation: "can_read",
				Object:   "vault:test-vault",
			},
		},
	}

	err = invalidWriteReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user cannot be empty")

	// Test invalid list objects request
	invalidListReq := &ListObjectsRequest{
		User:     "invalid-user",
		Relation: "can_read",
		Type:     "vault",
	}

	err = invalidListReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user must start with 'user:' or 'token:'")
}

// TestStringRepresentations tests string representations of types
func TestStringRepresentations(t *testing.T) {
	// Test vault permission strings
	assert.Equal(t, "can_read", VaultPermissionRead.String())
	assert.Equal(t, "can_write", VaultPermissionWrite.String())
	assert.Equal(t, "can_manage", VaultPermissionManage.String())

	// Test secret permission strings
	assert.Equal(t, "can_read", SecretPermissionRead.String())
	assert.Equal(t, "can_write", SecretPermissionWrite.String())
	assert.Equal(t, "can_read_sensitive", SecretPermissionReadSensitive.String())

	// Test tenant role strings
	assert.Equal(t, "owner", TenantRoleOwner.String())
	assert.Equal(t, "admin", TenantRoleAdmin.String())
	assert.Equal(t, "member", TenantRoleMember.String())

	// Test group role strings
	assert.Equal(t, "owner", GroupRoleOwner.String())
	assert.Equal(t, "admin", GroupRoleAdmin.String())
	assert.Equal(t, "member", GroupRoleMember.String())

	// Test plan feature strings
	assert.Equal(t, "basic", PlanFeatureBasic.String())
	assert.Equal(t, "advanced_permissions", PlanFeatureAdvancedPermissions.String())
	assert.Equal(t, "external_sharing", PlanFeatureExternalSharing.String())

	// Test token permission strings
	assert.Equal(t, "can_use", TokenPermissionUse.String())
	assert.Equal(t, "can_read_secrets", TokenPermissionReadSecrets.String())
	assert.Equal(t, "can_write_secrets", TokenPermissionWriteSecrets.String())

	// Test tenant permission strings
	assert.Equal(t, "can_view", TenantPermissionView.String())
	assert.Equal(t, "can_create_vault", TenantPermissionCreateVault.String())
	assert.Equal(t, "owner", TenantPermissionOwner.String())
	assert.Equal(t, "admin", TenantPermissionAdmin.String())

	// Test health state strings
	assert.Equal(t, "healthy", HealthStateHealthy.String())
	assert.Equal(t, "unhealthy", HealthStateUnhealthy.String())
	assert.Equal(t, "degraded", HealthStateDegraded.String())
}

// BenchmarkIntegration runs performance benchmarks for key operations
func BenchmarkIntegration(b *testing.B) {
	// Benchmark format operations
	b.Run("FormatUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FormatUser("alice")
		}
	})

	b.Run("FormatVault", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FormatVault("test-vault")
		}
	})

	b.Run("ParseObject", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ParseObject("vault:test-vault")
		}
	})

	b.Run("ParseUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ParseUser("user:alice")
		}
	})

	b.Run("GenerateTokenID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GenerateTokenID()
		}
	})

	b.Run("NewAuditEvent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewAuditEvent("alice", "check_permission", "vault:test-vault", true, nil)
		}
	})

	b.Run("TupleValidation", func(b *testing.B) {
		tuple := &Tuple{
			User:     "user:alice",
			Relation: "can_read",
			Object:   "vault:test-vault",
		}

		for i := 0; i < b.N; i++ {
			tuple.Validate()
		}
	})

	b.Run("CheckRequestValidation", func(b *testing.B) {
		req := &CheckRequest{
			User:     "user:alice",
			Relation: "can_read",
			Object:   "vault:test-vault",
		}

		for i := 0; i < b.N; i++ {
			req.Validate()
		}
	})
}
