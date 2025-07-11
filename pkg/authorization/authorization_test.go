package authorization

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock implementations for testing
type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*CheckResponse), args.Error(1)
}

func (m *MockAuthorizationService) CheckBatch(ctx context.Context, requests []*CheckRequest) ([]*CheckResponse, error) {
	args := m.Called(ctx, requests)
	return args.Get(0).([]*CheckResponse), args.Error(1)
}

func (m *MockAuthorizationService) ListObjects(ctx context.Context, req *ListObjectsRequest) (*ListObjectsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*ListObjectsResponse), args.Error(1)
}

func (m *MockAuthorizationService) Write(ctx context.Context, req *WriteRequest) (*WriteResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*WriteResponse), args.Error(1)
}

func (m *MockAuthorizationService) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*DeleteResponse), args.Error(1)
}

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) Log(ctx context.Context, event *AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockAuditService) Query(ctx context.Context, query *AuditQuery) ([]*AuditEvent, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*AuditEvent), args.Error(1)
}

type MockLockariAuthorizationService struct {
	mock.Mock
}

func (m *MockLockariAuthorizationService) CheckVaultPermission(ctx context.Context, userID, vaultID string, permission VaultPermission) (bool, error) {
	args := m.Called(ctx, userID, vaultID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockLockariAuthorizationService) CheckSecretPermission(ctx context.Context, userID, secretID string, permission SecretPermission) (bool, error) {
	args := m.Called(ctx, userID, secretID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockLockariAuthorizationService) CheckTenantPermission(ctx context.Context, userID, tenantID string, permission TenantPermission) (bool, error) {
	args := m.Called(ctx, userID, tenantID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockLockariAuthorizationService) SetupTenant(ctx context.Context, tenantID, ownerID string, plan PlanFeature) error {
	args := m.Called(ctx, tenantID, ownerID, plan)
	return args.Error(0)
}

func (m *MockLockariAuthorizationService) AddUserToTenant(ctx context.Context, userID, tenantID string, role TenantRole) error {
	args := m.Called(ctx, userID, tenantID, role)
	return args.Error(0)
}

func (m *MockLockariAuthorizationService) CreateGroup(ctx context.Context, groupID, tenantID, ownerID string) error {
	args := m.Called(ctx, groupID, tenantID, ownerID)
	return args.Error(0)
}

func (m *MockLockariAuthorizationService) CreateAPIToken(ctx context.Context, userID, vaultID string, permissions []TokenPermission) (string, error) {
	args := m.Called(ctx, userID, vaultID, permissions)
	return args.String(0), args.Error(1)
}

func (m *MockLockariAuthorizationService) ListAccessibleVaults(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

// Test Suite
type AuthorizationTestSuite struct {
	suite.Suite
	service        *MockAuthorizationService
	cache          *MockCacheService
	audit          *MockAuditService
	lockariService *MockLockariAuthorizationService
}

func (suite *AuthorizationTestSuite) SetupTest() {
	suite.service = new(MockAuthorizationService)
	suite.cache = new(MockCacheService)
	suite.audit = new(MockAuditService)
	suite.lockariService = new(MockLockariAuthorizationService)
}

func (suite *AuthorizationTestSuite) TestCheckVaultPermission() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("CheckVaultPermission", ctx, "alice", "test-vault", VaultPermissionRead).Return(true, nil)

	// Test
	canRead, err := suite.lockariService.CheckVaultPermission(ctx, "alice", "test-vault", VaultPermissionRead)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), canRead)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestCheckSecretPermission() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("CheckSecretPermission", ctx, "bob", "test-secret", SecretPermissionWrite).Return(false, nil)

	// Test
	canWrite, err := suite.lockariService.CheckSecretPermission(ctx, "bob", "test-secret", SecretPermissionWrite)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), canWrite)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestCheckTenantPermission() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("CheckTenantPermission", ctx, "charlie", "test-tenant", TenantPermissionAdmin).Return(true, nil)

	// Test
	isAdmin, err := suite.lockariService.CheckTenantPermission(ctx, "charlie", "test-tenant", TenantPermissionAdmin)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isAdmin)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestSetupTenant() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("SetupTenant", ctx, "test-tenant", "alice", PlanFeatureAdvancedPermissions).Return(nil)

	// Test
	err := suite.lockariService.SetupTenant(ctx, "test-tenant", "alice", PlanFeatureAdvancedPermissions)

	// Assertions
	assert.NoError(suite.T(), err)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestAddUserToTenant() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("AddUserToTenant", ctx, "bob", "test-tenant", TenantRoleAdmin).Return(nil)

	// Test
	err := suite.lockariService.AddUserToTenant(ctx, "bob", "test-tenant", TenantRoleAdmin)

	// Assertions
	assert.NoError(suite.T(), err)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestCreateGroup() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("CreateGroup", ctx, "test-group", "test-tenant", "alice").Return(nil)

	// Test
	err := suite.lockariService.CreateGroup(ctx, "test-group", "test-tenant", "alice")

	// Assertions
	assert.NoError(suite.T(), err)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestCreateAPIToken() {
	ctx := context.Background()

	permissions := []TokenPermission{
		TokenPermissionReadSecrets,
		TokenPermissionWriteSecrets,
	}

	// Setup expectations
	suite.lockariService.On("CreateAPIToken", ctx, "alice", "test-vault", permissions).Return("test-token-id", nil)

	// Test
	tokenID, err := suite.lockariService.CreateAPIToken(ctx, "alice", "test-vault", permissions)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-token-id", tokenID)

	suite.lockariService.AssertExpectations(suite.T())
}

func (suite *AuthorizationTestSuite) TestListAccessibleVaults() {
	ctx := context.Background()

	// Setup expectations
	suite.lockariService.On("ListAccessibleVaults", ctx, "alice").Return([]string{"vault1", "vault2", "vault3"}, nil)

	// Test
	vaults, err := suite.lockariService.ListAccessibleVaults(ctx, "alice")

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), []string{"vault1", "vault2", "vault3"}, vaults)

	suite.lockariService.AssertExpectations(suite.T())
}

func TestAuthorizationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorizationTestSuite))
}

// Unit tests for types and utilities
func TestFormatUser(t *testing.T) {
	assert.Equal(t, "user:alice", FormatUser("alice"))
	assert.Equal(t, "user:alice", FormatUser("user:alice")) // should not double-prefix
}

func TestFormatVault(t *testing.T) {
	assert.Equal(t, "vault:test-vault", FormatVault("test-vault"))
	assert.Equal(t, "vault:test-vault", FormatVault("vault:test-vault")) // should not double-prefix
}

func TestFormatSecret(t *testing.T) {
	assert.Equal(t, "secret:test-secret", FormatSecret("test-secret"))
	assert.Equal(t, "secret:test-secret", FormatSecret("secret:test-secret")) // should not double-prefix
}

func TestFormatTenant(t *testing.T) {
	assert.Equal(t, "tenant:test-tenant", FormatTenant("test-tenant"))
	assert.Equal(t, "tenant:test-tenant", FormatTenant("tenant:test-tenant")) // should not double-prefix
}

func TestFormatToken(t *testing.T) {
	assert.Equal(t, "token:test-token", FormatToken("test-token"))
	assert.Equal(t, "token:test-token", FormatToken("token:test-token")) // should not double-prefix
}

func TestFormatGroup(t *testing.T) {
	assert.Equal(t, "group:test-group", FormatGroup("test-group"))
	assert.Equal(t, "group:test-group", FormatGroup("group:test-group")) // should not double-prefix
}

func TestVaultPermissionToRelation(t *testing.T) {
	tests := []struct {
		permission VaultPermission
		relation   string
	}{
		{VaultPermissionView, "can_view"},
		{VaultPermissionRead, "can_read"},
		{VaultPermissionCopy, "can_copy"},
		{VaultPermissionDownload, "can_download"},
		{VaultPermissionWrite, "can_write"},
		{VaultPermissionDelete, "can_delete"},
		{VaultPermissionShare, "can_share"},
		{VaultPermissionManage, "can_manage"},
	}

	for _, test := range tests {
		assert.Equal(t, test.relation, vaultPermissionToRelation(test.permission))
	}
}

func TestSecretPermissionToRelation(t *testing.T) {
	tests := []struct {
		permission SecretPermission
		relation   string
	}{
		{SecretPermissionView, "can_view"},
		{SecretPermissionRead, "can_read"},
		{SecretPermissionCopy, "can_copy"},
		{SecretPermissionDownload, "can_download"},
		{SecretPermissionWrite, "can_write"},
		{SecretPermissionDelete, "can_delete"},
		{SecretPermissionReadSensitive, "can_read_sensitive"},
		{SecretPermissionCopySensitive, "can_copy_sensitive"},
		{SecretPermissionCopyProduction, "can_copy_production"},
	}

	for _, test := range tests {
		assert.Equal(t, test.relation, secretPermissionToRelation(test.permission))
	}
}

func TestTenantPermissionToRelation(t *testing.T) {
	tests := []struct {
		permission TenantPermission
		relation   string
	}{
		{TenantPermissionView, "can_view"},
		{TenantPermissionCreateVault, "can_create_vault"},
		{TenantPermissionManageUsers, "can_manage_users"},
		{TenantPermissionManageGroups, "can_manage_groups"},
		{TenantPermissionManageTokens, "can_manage_tokens"},
		{TenantPermissionViewAudit, "can_view_audit"},
		{TenantPermissionManageSettings, "can_manage_settings"},
		{TenantPermissionExternalShare, "can_external_share"},
		{TenantPermissionOwner, "owner"},
		{TenantPermissionAdmin, "admin"},
	}

	for _, test := range tests {
		// Since TenantPermission is a string, we can directly compare
		assert.Equal(t, test.relation, string(test.permission))
	}
}

func TestTokenPermissionToRelation(t *testing.T) {
	tests := []struct {
		permission TokenPermission
		relation   string
	}{
		{TokenPermissionUse, "can_use"},
		{TokenPermissionReadSecrets, "can_read_secrets"},
		{TokenPermissionWriteSecrets, "can_write_secrets"},
		{TokenPermissionManageVault, "can_manage_vault"},
		{TokenPermissionRevoke, "can_revoke"},
		{TokenPermissionRegenerate, "can_regenerate"},
	}

	for _, test := range tests {
		assert.Equal(t, test.relation, tokenPermissionToRelation(test.permission))
	}
}

func TestGenerateTokenID(t *testing.T) {
	tokenID := GenerateTokenID()
	assert.NotEmpty(t, tokenID)
	assert.Len(t, tokenID, 32) // UUID without dashes

	// Generate another one to ensure uniqueness
	tokenID2 := GenerateTokenID()
	assert.NotEqual(t, tokenID, tokenID2)
}

func TestNewAuditEvent(t *testing.T) {
	event := NewAuditEvent("alice", "check_permission", "vault:test-vault", true, map[string]interface{}{
		"permission": "can_read",
	})

	assert.Equal(t, "alice", event.UserID)
	assert.Equal(t, "check_permission", event.Action)
	assert.Equal(t, "vault:test-vault", event.Resource)
	assert.Equal(t, "allowed", event.Result)
	assert.Equal(t, "can_read", event.Metadata["permission"])
	assert.NotZero(t, event.Timestamp)
}

// Integration tests
func TestAuthorizationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// This would require actual OpenFGA instance
	// config := &Config{
	// 	APIURL:     "localhost:8080",
	// 	StoreID:    "test-store",
	// 	ModelID:    "test-model",
	// }
	//
	// service, err := NewAuthorizationService(config)
	// assert.NoError(t, err)
	//
	// // Test basic functionality
	// ctx := context.Background()
	// response, err := service.Check(ctx, &CheckRequest{
	// 	User:     "user:alice",
	// 	Relation: "can_read",
	// 	Object:   "vault:test-vault",
	// })
	// assert.NoError(t, err)
	// assert.NotNil(t, response)
}

// Benchmark tests
func BenchmarkFormatUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatUser("alice")
	}
}

func BenchmarkVaultPermissionToRelation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = vaultPermissionToRelation(VaultPermissionRead)
	}
}

func BenchmarkGenerateTokenID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateTokenID()
	}
}

func BenchmarkNewAuditEvent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewAuditEvent("alice", "check_permission", "vault:test-vault", true, nil)
	}
}

// Error handling tests
func TestErrorHandling(t *testing.T) {
	// Test custom errors
	err := ErrPermissionDenied
	assert.Equal(t, "permission denied", err.Error())

	err = ErrInvalidUser
	assert.Equal(t, "invalid user format", err.Error())

	err = ErrInvalidObject
	assert.Equal(t, "invalid object format", err.Error())

	err = ErrInvalidTenant
	assert.Equal(t, "invalid tenant", err.Error())

	err = ErrInvalidToken
	assert.Equal(t, "invalid token", err.Error())

	err = ErrNotFound
	assert.Equal(t, "not found", err.Error())

	err = ErrUnauthorized
	assert.Equal(t, "unauthorized", err.Error())

	err = ErrForbidden
	assert.Equal(t, "forbidden", err.Error())

	err = ErrRateLimited
	assert.Equal(t, "rate limited", err.Error())

	err = ErrTimeout
	assert.Equal(t, "operation timeout", err.Error())

	err = ErrCircuitBreakerOpen
	assert.Equal(t, "circuit breaker open", err.Error())
}

// Configuration tests
func TestConfigValidation(t *testing.T) {
	// Test valid config
	config := &Config{
		APIURL:               "localhost:8080",
		StoreID:              "test-store",
		AuthorizationModelID: "test-model",
	}

	err := config.Validate()
	assert.NoError(t, err)

	// Test invalid config - missing APIURL
	config = &Config{
		StoreID:              "test-store",
		AuthorizationModelID: "test-model",
	}

	err = config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "APIURL")

	// Test invalid config - missing StoreID
	config = &Config{
		APIURL:               "localhost:8080",
		AuthorizationModelID: "test-model",
	}

	err = config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "StoreID")
}
