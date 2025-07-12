# 🚀 Guia Rápido: Configurando Permissões para Novos Usuários/Tenants

## 📋 Resumo Executivo

**Para criar um novo usuário/tenant com permissões no Lockari:**

1. **Criar Tenant** → Define o plano e recursos disponíveis
2. **Adicionar Usuário** → Associa usuário ao tenant com papel específico
3. **Configurar Vault** → Cria vault inicial com permissões básicas

## 🛠️ Implementação Prática

### 1. **Fluxo Completo de Criação**

```go
func CreateNewUserTenant(ctx context.Context, authService authorization.LockariAuthorizationService, 
                        userID, tenantID string, plan PlanType) error {
    
    // 1. Configurar tenant com plano específico
    features := getPlanFeatures(plan)
    err := authService.SetupTenant(ctx, tenantID, userID, features)
    if err != nil {
        return fmt.Errorf("failed to setup tenant: %w", err)
    }
    
    // 2. Adicionar usuário como owner do tenant
    err = authService.AddUserToTenant(ctx, userID, tenantID, authorization.TenantRoleOwner)
    if err != nil {
        return fmt.Errorf("failed to add user to tenant: %w", err)
    }
    
    // 3. Criar vault inicial (se permitido pelo plano)
    if plan != PlanFree || isPlanFeatureEnabled(features, authorization.PlanFeatureVaultCreation) {
        vaultID := fmt.Sprintf("vault-%s-default", userID)
        err = authService.SetupVault(ctx, vaultID, tenantID, userID)
        if err != nil {
            return fmt.Errorf("failed to setup initial vault: %w", err)
        }
    }
    
    return nil
}
```

### 2. **Diferenças Entre Planos**

```go
type PlanType string

const (
    PlanFree       PlanType = "free"
    PlanPro        PlanType = "pro"
    PlanEnterprise PlanType = "enterprise"
)

func getPlanFeatures(plan PlanType) []authorization.PlanFeature {
    switch plan {
    case PlanFree:
        return []authorization.PlanFeature{
            authorization.PlanFeatureVaultLimit,      // Limite: 3 vaults
            authorization.PlanFeatureUserLimit,       // Limite: 1 usuário
            authorization.PlanFeatureBasicSharing,    // Compartilhamento básico
        }
        
    case PlanPro:
        return []authorization.PlanFeature{
            authorization.PlanFeatureVaultLimit,      // Limite: 50 vaults
            authorization.PlanFeatureUserLimit,       // Limite: 10 usuários
            authorization.PlanFeatureAdvancedSharing, // Compartilhamento avançado
            authorization.PlanFeatureAPIAccess,       // Acesso à API
            authorization.PlanFeatureGroupManagement, // Gerenciamento de grupos
            authorization.PlanFeatureAuditLogs,       // Logs de auditoria
        }
        
    case PlanEnterprise:
        return []authorization.PlanFeature{
            authorization.PlanFeatureUnlimitedVaults,  // Vaults ilimitados
            authorization.PlanFeatureUnlimitedUsers,   // Usuários ilimitados
            authorization.PlanFeatureExternalSharing,  // Compartilhamento externo
            authorization.PlanFeatureAPIAccess,        // Acesso à API
            authorization.PlanFeatureGroupManagement,  // Gerenciamento de grupos
            authorization.PlanFeatureAuditLogs,        // Logs de auditoria
            authorization.PlanFeatureSSO,              // Single Sign-On
            authorization.PlanFeatureAdvancedSecurity, // Segurança avançada
        }
        
    default:
        return []authorization.PlanFeature{}
    }
}
```

### 3. **Configuração Específica por Plano**

#### **🆓 Plano Free**
```go
func setupFreePlan(ctx context.Context, authService authorization.LockariAuthorizationService, 
                   userID, tenantID string) error {
    // 1. Tenant com recursos limitados
    features := []authorization.PlanFeature{
        authorization.PlanFeatureVaultLimit,    // Máximo 3 vaults
        authorization.PlanFeatureBasicSharing,  // Compartilhamento básico
    }
    
    err := authService.SetupTenant(ctx, tenantID, userID, features)
    if err != nil {
        return err
    }
    
    // 2. Usuário como owner (único usuário permitido)
    err = authService.AddUserToTenant(ctx, userID, tenantID, authorization.TenantRoleOwner)
    if err != nil {
        return err
    }
    
    // 3. Vault inicial gratuito
    defaultVaultID := fmt.Sprintf("vault-%s-personal", userID)
    return authService.SetupVault(ctx, defaultVaultID, tenantID, userID)
}
```

#### **💼 Plano Pro**
```go
func setupProPlan(ctx context.Context, authService authorization.LockariAuthorizationService, 
                  userID, tenantID string) error {
    // 1. Tenant com recursos expandidos
    features := []authorization.PlanFeature{
        authorization.PlanFeatureVaultLimit,        // Até 50 vaults
        authorization.PlanFeatureUserLimit,         // Até 10 usuários
        authorization.PlanFeatureAdvancedSharing,   // Compartilhamento avançado
        authorization.PlanFeatureAPIAccess,         // Tokens de API
        authorization.PlanFeatureGroupManagement,   // Grupos
        authorization.PlanFeatureAuditLogs,         // Auditoria
    }
    
    err := authService.SetupTenant(ctx, tenantID, userID, features)
    if err != nil {
        return err
    }
    
    // 2. Usuário como owner
    err = authService.AddUserToTenant(ctx, userID, tenantID, authorization.TenantRoleOwner)
    if err != nil {
        return err
    }
    
    // 3. Múltiplos vaults iniciais
    vaultIDs := []string{
        fmt.Sprintf("vault-%s-personal", userID),
        fmt.Sprintf("vault-%s-business", userID),
    }
    
    for _, vaultID := range vaultIDs {
        err = authService.SetupVault(ctx, vaultID, tenantID, userID)
        if err != nil {
            return err
        }
    }
    
    // 4. Grupo padrão para colaboração
    groupID := fmt.Sprintf("group-%s-team", tenantID)
    err = authService.CreateGroup(ctx, groupID, tenantID, userID)
    if err != nil {
        return err
    }
    
    return authService.AddUserToGroup(ctx, userID, groupID, authorization.GroupRoleOwner)
}
```

#### **🏢 Plano Enterprise**
```go
func setupEnterprisePlan(ctx context.Context, authService authorization.LockariAuthorizationService, 
                         userID, tenantID string) error {
    // 1. Tenant com recursos ilimitados
    features := []authorization.PlanFeature{
        authorization.PlanFeatureUnlimitedVaults,   // Vaults ilimitados
        authorization.PlanFeatureUnlimitedUsers,    // Usuários ilimitados
        authorization.PlanFeatureExternalSharing,   // Compartilhamento externo
        authorization.PlanFeatureAPIAccess,         // API completa
        authorization.PlanFeatureGroupManagement,   // Grupos avançados
        authorization.PlanFeatureAuditLogs,         // Auditoria completa
        authorization.PlanFeatureSSO,               // SSO
        authorization.PlanFeatureAdvancedSecurity,  // Segurança avançada
    }
    
    err := authService.SetupTenant(ctx, tenantID, userID, features)
    if err != nil {
        return err
    }
    
    // 2. Usuário como owner
    err = authService.AddUserToTenant(ctx, userID, tenantID, authorization.TenantRoleOwner)
    if err != nil {
        return err
    }
    
    // 3. Estrutura organizacional completa
    vaultIDs := []string{
        fmt.Sprintf("vault-%s-executive", userID),
        fmt.Sprintf("vault-%s-hr", userID),
        fmt.Sprintf("vault-%s-finance", userID),
        fmt.Sprintf("vault-%s-it", userID),
    }
    
    for _, vaultID := range vaultIDs {
        err = authService.SetupVault(ctx, vaultID, tenantID, userID)
        if err != nil {
            return err
        }
    }
    
    // 4. Grupos departamentais
    groupIDs := []string{
        fmt.Sprintf("group-%s-executives", tenantID),
        fmt.Sprintf("group-%s-hr", tenantID),
        fmt.Sprintf("group-%s-finance", tenantID),
        fmt.Sprintf("group-%s-it", tenantID),
    }
    
    for _, groupID := range groupIDs {
        err = authService.CreateGroup(ctx, groupID, tenantID, userID)
        if err != nil {
            return err
        }
        
        err = authService.AddUserToGroup(ctx, userID, groupID, authorization.GroupRoleOwner)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

## 📊 Tabela de Comparação

| Recurso | Free | Pro | Enterprise |
|---------|------|-----|------------|
| **Vaults** | 3 | 50 | ♾️ Ilimitado |
| **Usuários** | 1 | 10 | ♾️ Ilimitado |
| **Compartilhamento** | Básico | Avançado | Externo |
| **API** | ❌ | ✅ | ✅ |
| **Grupos** | ❌ | ✅ | ✅ |
| **Auditoria** | ❌ | ✅ | ✅ |
| **SSO** | ❌ | ❌ | ✅ |

## 🎯 Exemplo de Uso Prático

```go
// No seu handler de cadastro
func SignUpHandler(c *gin.Context) {
    var req SignUpRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Criar usuário no sistema de autenticação
    userID, err := createUser(req.Email, req.Password)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to create user"})
        return
    }
    
    // Criar tenant com plano escolhido
    tenantID := fmt.Sprintf("tenant-%s", userID)
    err = CreateNewUserTenant(c.Request.Context(), authService, userID, tenantID, req.Plan)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to setup tenant"})
        return
    }
    
    c.JSON(201, gin.H{
        "message": "User created successfully",
        "user_id": userID,
        "tenant_id": tenantID,
        "plan": req.Plan,
    })
}
```

## 💡 Resumo das Diferenças

### **🔑 Principais Diferenças:**

1. **Recursos Limitados vs Ilimitados**
   - Free: 3 vaults, 1 usuário
   - Pro: 50 vaults, 10 usuários  
   - Enterprise: Ilimitado

2. **Funcionalidades Avançadas**
   - Free: Apenas básico
   - Pro: API, Grupos, Auditoria
   - Enterprise: SSO, Compartilhamento externo

3. **Configuração Inicial**
   - Free: 1 vault pessoal
   - Pro: 2 vaults + 1 grupo
   - Enterprise: 4 vaults + 4 grupos

**Use o método `CreateNewUserTenant()` passando o plano desejado e o sistema configurará automaticamente as permissões corretas!**
