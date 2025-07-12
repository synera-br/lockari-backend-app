# Plan Features Constants Reference

## Location
The `PlanFeatureVaultLimit` constant and related plan limit constants are defined in:
- **File**: `pkg/authorization/types.go`
- **Section**: Plan Features constants (around line 167)

## Defined Constants

### Limit-related Plan Features
```go
const (
    PlanFeatureVaultLimit          PlanFeature = "vault_limit"
    PlanFeatureUserLimit           PlanFeature = "user_limit"
    PlanFeatureUnlimitedVaults     PlanFeature = "unlimited_vaults"
    PlanFeatureUnlimitedUsers      PlanFeature = "unlimited_users"
)
```

### All Plan Features
```go
const (
    PlanFeatureBasic               PlanFeature = "basic"
    PlanFeatureAdvancedPermissions PlanFeature = "advanced_permissions"
    PlanFeatureCrossTenantSharing  PlanFeature = "cross_tenant_sharing"
    PlanFeatureAuditLogs           PlanFeature = "audit_logs"
    PlanFeatureBackup              PlanFeature = "backup"
    PlanFeatureExternalSharing     PlanFeature = "external_sharing"
    PlanFeatureVaultLimit          PlanFeature = "vault_limit"
    PlanFeatureUserLimit           PlanFeature = "user_limit"
    PlanFeatureUnlimitedVaults     PlanFeature = "unlimited_vaults"
    PlanFeatureUnlimitedUsers      PlanFeature = "unlimited_users"
)
```

## Usage in Code

The constants are used in:
1. **lockari_service.go** - For automatic plan setup during onboarding
2. **SETUP_PERMISSIONS_GUIDE.md** - For documentation and examples

### Example Usage
```go
// Free plan (3 vaults limit)
features := []authorization.PlanFeature{
    authorization.PlanFeatureBasic,
    authorization.PlanFeatureVaultLimit,   // Limite: 3 vaults
}

// Pro plan (50 vaults limit)
features := []authorization.PlanFeature{
    authorization.PlanFeatureBasic,
    authorization.PlanFeatureAdvancedPermissions,
    authorization.PlanFeatureVaultLimit,      // Limite: 50 vaults
    authorization.PlanFeatureUserLimit,       // Limite: 10 usuários
}

// Enterprise plan (unlimited)
features := []authorization.PlanFeature{
    authorization.PlanFeatureBasic,
    authorization.PlanFeatureAdvancedPermissions,
    authorization.PlanFeatureCrossTenantSharing,
    authorization.PlanFeatureAuditLogs,
    authorization.PlanFeatureBackup,
    authorization.PlanFeatureExternalSharing,
    authorization.PlanFeatureUnlimitedVaults,  // Vaults ilimitados
    authorization.PlanFeatureUnlimitedUsers,   // Usuários ilimitados
}
```

## Plan Limits Summary

| Plan | Vault Limit | User Limit | Features |
|------|-------------|------------|----------|
| Free | 3 | 1 | Basic features only |
| Pro | 50 | 10 | Advanced permissions, audit logs |
| Enterprise | Unlimited | Unlimited | All features including cross-tenant sharing, backup, external sharing |

## Implementation Notes

1. **Feature Detection**: The service checks for `PlanFeatureVaultLimit` vs `PlanFeatureUnlimitedVaults` to determine limits
2. **Enforcement**: Actual limit values (3, 50, unlimited) are implemented in the service logic
3. **Validation**: All constants are type-safe and validated through the `IsValid()` method
4. **Extensibility**: New plan features can be added by extending the constants and updating `AllPlanFeatures()`
