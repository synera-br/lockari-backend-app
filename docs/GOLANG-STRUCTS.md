# Gola## 📦 Estrutura de Packages

```
internal/
├── core/
│   ├── entity/         # Entidades de domínio
│   ├── repository/     # Interfaces/contratos
│   └── services/       # Lógica de ### 4. **Vault** - Cofre de Dados

```go
// internal/core/entity/vault.go
package entity

import (
    "errors"
    "time"
)

type Vault struct 
├── adapters/
│   ├── firebase/       # Adaptador Firebase
│   ├── openfga/        # Adaptador OpenFGA
│   └── http/           # Handlers HTTP
└── config/             # Configurações
```

## Visão Geral

Este documento detalha as structs Go para o backend do Lockari, alinhadas com as collections do Firebase e o modelo de autorização OpenFGA.

## 📦 Estrutura de Packages

```
internal/
├── core/
│   ├── entity/         # Entidades de domínio
│   ├── repository/     # Interfaces/contratos de repositório
│   └── services/       # Lógica de negócio
├── adapters/
│   ├── firebase/       # Adaptador Firebase
│   ├── openfga/        # Adaptador OpenFGA
│   └── http/           # Handlers HTTP
└── config/             # Configurações
```

## 🏗️ Structs Principais

### 1. **User** - Usuário do Sistema

```go
// internal/core/entity/user.go
package entity

import (
    "time"
)

type User struct {
    ID            string                 `json:"id" firestore:"id"`
    Email         string                 `json:"email" firestore:"email" validate:"required,email"`
    DisplayName   string                 `json:"displayName" firestore:"displayName" validate:"required"`
    PhotoURL      string                 `json:"photoURL,omitempty" firestore:"photoURL"`
    ProviderID    string                 `json:"providerId" firestore:"providerId"`
    EmailVerified bool                   `json:"emailVerified" firestore:"emailVerified"`
    IsActive      bool                   `json:"isActive" firestore:"isActive"`
    Preferences   UserPreferences        `json:"preferences" firestore:"preferences"`
    Metadata      UserMetadata           `json:"metadata" firestore:"metadata"`
    CreatedAt     time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt" firestore:"updatedAt"`
    LastLoginAt   *time.Time             `json:"lastLoginAt,omitempty" firestore:"lastLoginAt"`
}

type UserPreferences struct {
    Theme         string                 `json:"theme" firestore:"theme" validate:"oneof=light dark system"`
    Language      string                 `json:"language" firestore:"language" validate:"required"`
    Notifications NotificationSettings   `json:"notifications" firestore:"notifications"`
}

type NotificationSettings struct {
    Email         bool                   `json:"email" firestore:"email"`
    Browser       bool                   `json:"browser" firestore:"browser"`
    Mobile        bool                   `json:"mobile" firestore:"mobile"`
    Security      bool                   `json:"security" firestore:"security"`
    Marketing     bool                   `json:"marketing" firestore:"marketing"`
}

type UserMetadata struct {
    IPAddress     string                 `json:"ipAddress,omitempty" firestore:"ipAddress"`
    UserAgent     string                 `json:"userAgent,omitempty" firestore:"userAgent"`
    Location      string                 `json:"location,omitempty" firestore:"location"`
}
```

### 2. **Tenant** - Organização/Empresa

```go
// internal/core/entity/tenant.go
package entity

type Tenant struct {
    ID          string                 `json:"id" firestore:"id"`
    Name        string                 `json:"name" firestore:"name" validate:"required,min=2,max=100"`
    Slug        string                 `json:"slug" firestore:"slug" validate:"required,alphanum,min=2,max=50"`
    Description string                 `json:"description,omitempty" firestore:"description"`
    Domain      string                 `json:"domain,omitempty" firestore:"domain"`
    Plan        PlanType               `json:"plan" firestore:"plan" validate:"required"`
    PlanExpiry  *time.Time             `json:"planExpiry,omitempty" firestore:"planExpiry"`
    Limits      TenantLimits           `json:"limits" firestore:"limits"`
    Features    []string               `json:"features" firestore:"features"`
    Billing     TenantBilling          `json:"billing" firestore:"billing"`
    Settings    TenantSettings         `json:"settings" firestore:"settings"`
    CreatedBy   string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    IsActive    bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt   time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt   time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

type PlanType string

const (
    PlanFree       PlanType = "free"
    PlanEnterprise PlanType = "enterprise"
)

type TenantLimits struct {
    MaxUsers    int                    `json:"maxUsers" firestore:"maxUsers"`
    MaxVaults   int                    `json:"maxVaults" firestore:"maxVaults"`
    MaxSecrets  int                    `json:"maxSecrets" firestore:"maxSecrets"`
}

type TenantBilling struct {
    CustomerID     string              `json:"customerId,omitempty" firestore:"customerId"`
    SubscriptionID string              `json:"subscriptionId,omitempty" firestore:"subscriptionId"`
    Status         string              `json:"status" firestore:"status"`
}

type TenantSettings struct {
    AllowExternalSharing bool           `json:"allowExternalSharing" firestore:"allowExternalSharing"`
    RequireMFA          bool           `json:"requireMFA" firestore:"requireMFA"`
    SessionTimeout      int            `json:"sessionTimeout" firestore:"sessionTimeout"` // minutes
    AuditRetention      int            `json:"auditRetention" firestore:"auditRetention"` // days
}
```

### 3. **TenantMember** - Relacionamento Usuário-Tenant

```go
// internal/core/entity/tenant_member.go
package entity

type TenantMember struct {
    ID           string                 `json:"id" firestore:"id"`
    TenantID     string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    UserID       string                 `json:"userId" firestore:"userId" validate:"required"`
    Role         TenantRole             `json:"role" firestore:"role" validate:"required"`
    Permissions  []string               `json:"permissions" firestore:"permissions"`
    Status       MemberStatus           `json:"status" firestore:"status"`
    InvitedBy    string                 `json:"invitedBy,omitempty" firestore:"invitedBy"`
    JoinedAt     time.Time              `json:"joinedAt" firestore:"joinedAt"`
    LastAccessAt *time.Time             `json:"lastAccessAt,omitempty" firestore:"lastAccessAt"`
    Metadata     map[string]interface{} `json:"metadata,omitempty" firestore:"metadata"`
}

type TenantRole string

const (
    RoleOwner  TenantRole = "owner"
    RoleAdmin  TenantRole = "admin"
    RoleMember TenantRole = "member"
)

type MemberStatus string

const (
    StatusActive    MemberStatus = "active"
    StatusInvited   MemberStatus = "invited"
    StatusSuspended MemberStatus = "suspended"
)
```

### 4. **Group** - Grupo de Usuários

```go
package domain

type Group struct {
    ID          string                 `json:"id" firestore:"id"`
    Name        string                 `json:"name" firestore:"name" validate:"required,min=2,max=100"`
    Description string                 `json:"description,omitempty" firestore:"description"`
    TenantID    string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy   string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    Members     []string               `json:"members" firestore:"members"`
    Admins      []string               `json:"admins" firestore:"admins"`
    Permissions []string               `json:"permissions" firestore:"permissions"`
    IsActive    bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt   time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt   time.Time              `json:"updatedAt" firestore:"updatedAt"`
}
```

### 5. **Vault** - Cofre de Segredos

```go
// internal/core/entity/vault.go
package entity

import (
    "errors"
    "time"
)

type Vault struct {
    ID           string                 `json:"id" firestore:"id"`
    Name         string                 `json:"name" firestore:"name" validate:"required,min=2,max=100"`
    Description  string                 `json:"description,omitempty" firestore:"description"`
    TenantID     string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy    string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    Icon         string                 `json:"icon,omitempty" firestore:"icon"`
    Color        string                 `json:"color,omitempty" firestore:"color" validate:"hexcolor"`
    Tags         []string               `json:"tags" firestore:"tags" validate:"max=5"`
    SecretsCount int                    `json:"secretsCount" firestore:"secretsCount"`
    Permissions  VaultPermissions       `json:"permissions" firestore:"permissions"`
    Sharing      VaultSharing           `json:"sharing" firestore:"sharing"`
    Backup       VaultBackup            `json:"backup" firestore:"backup"`
    IsActive     bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt    time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt    time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

// Validação de tags (máximo 5)
func (v *Vault) AddTag(tag string) error {
    if len(v.Tags) >= 5 {
        return errors.New("maximum of 5 tags allowed")
    }
    
    // Verificar se já existe
    for _, existingTag := range v.Tags {
        if existingTag == tag {
            return errors.New("tag already exists")
        }
    }
    
    v.Tags = append(v.Tags, tag)
    return nil
}

func (v *Vault) RemoveTag(tag string) {
    for i, existingTag := range v.Tags {
        if existingTag == tag {
            v.Tags = append(v.Tags[:i], v.Tags[i+1:]...)
            return
        }
    }
}

func (v *Vault) GetOpenFGAID() string {
    return "vault:" + v.ID
}

type VaultPermissions struct {
    Owners      []string               `json:"owners" firestore:"owners"`
    Admins      []string               `json:"admins" firestore:"admins"`
    Writers     []string               `json:"writers" firestore:"writers"`
    Readers     []string               `json:"readers" firestore:"readers"`
    Viewers     []string               `json:"viewers" firestore:"viewers"`
    Copiers     []string               `json:"copiers" firestore:"copiers"`
    Downloaders []string               `json:"downloaders" firestore:"downloaders"`
}

type VaultSharing struct {
    IsPublic             bool                   `json:"isPublic" firestore:"isPublic"`
    AllowExternalSharing bool                   `json:"allowExternalSharing" firestore:"allowExternalSharing"`
    ExternalShares       []ExternalShare        `json:"externalShares" firestore:"externalShares"`
}

type ExternalShare struct {
    TenantID    string                 `json:"tenantId" firestore:"tenantId"`
    UserID      string                 `json:"userId" firestore:"userId"`
    Permissions []string               `json:"permissions" firestore:"permissions"`
    ExpiresAt   *time.Time             `json:"expiresAt,omitempty" firestore:"expiresAt"`
    CreatedAt   time.Time              `json:"createdAt" firestore:"createdAt"`
}

type VaultBackup struct {
    Enabled      bool                   `json:"enabled" firestore:"enabled"`
    Frequency    string                 `json:"frequency" firestore:"frequency"`
    LastBackupAt *time.Time             `json:"lastBackupAt,omitempty" firestore:"lastBackupAt"`
}
```

### 6. **Secret** - Segredo Genérico

```go
// internal/core/entity/secret.go
package entity

import (
    "errors"
    "time"
)

type Secret struct {
    ID             string                 `json:"id" firestore:"id"`
    Name           string                 `json:"name" firestore:"name" validate:"required,min=1,max=200"`
    Description    string                 `json:"description,omitempty" firestore:"description"`
    Type           SecretType             `json:"type" firestore:"type" validate:"required"`
    VaultID        string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    TenantID       string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy      string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    EncryptedValue string                 `json:"encryptedValue" firestore:"encryptedValue" validate:"required"`
    Metadata       SecretMetadata         `json:"metadata" firestore:"metadata"`
    Tags           []string               `json:"tags" firestore:"tags" validate:"max=5"`
    IsSensitive    bool                   `json:"isSensitive" firestore:"isSensitive"`
    IsProduction   bool                   `json:"isProduction" firestore:"isProduction"`
    ExpiresAt      *time.Time             `json:"expiresAt,omitempty" firestore:"expiresAt"`
    LastUsedAt     *time.Time             `json:"lastUsedAt,omitempty" firestore:"lastUsedAt"`
    UsageCount     int                    `json:"usageCount" firestore:"usageCount"`
    Version        int                    `json:"version" firestore:"version"`
    IsActive       bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt      time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt      time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

// Validação de tags (máximo 5)
func (s *Secret) AddTag(tag string) error {
    if len(s.Tags) >= 5 {
        return errors.New("maximum of 5 tags allowed")
    }
    
    // Verificar se já existe
    for _, existingTag := range s.Tags {
        if existingTag == tag {
            return errors.New("tag already exists")
        }
    }
    
    s.Tags = append(s.Tags, tag)
    return nil
}

func (s *Secret) RemoveTag(tag string) {
    for i, existingTag := range s.Tags {
        if existingTag == tag {
            s.Tags = append(s.Tags[:i], s.Tags[i+1:]...)
            return
        }
    }
}

func (s *Secret) GetOpenFGAID() string {
    return "object:" + s.ID
}

func (s *Secret) GetVaultOpenFGAID() string {
    return "vault:" + s.VaultID
}

type SecretType string

const (
    SecretTypePassword SecretType = "password"
    SecretTypeToken    SecretType = "token"
    SecretTypeAPIKey   SecretType = "api_key"
    SecretTypeOAuth    SecretType = "oauth"
    SecretTypeGeneric  SecretType = "generic"
)

type SecretMetadata struct {
    URL          string                 `json:"url,omitempty" firestore:"url"`
    Username     string                 `json:"username,omitempty" firestore:"username"`
    Notes        string                 `json:"notes,omitempty" firestore:"notes"`
    CustomFields map[string]interface{} `json:"customFields,omitempty" firestore:"customFields"`
}
```

### 7. **Tag** - Tags Predefinidas

```go
// internal/core/entity/tag.go
package entity

import (
    "time"
)

type TagType string

const (
    TagTypeSystem TagType = "system"
    TagTypeCustom TagType = "custom"
)

type Tag struct {
    ID          string    `json:"id" firestore:"id"`
    TenantID    string    `json:"tenant_id" firestore:"tenant_id"`
    Name        string    `json:"name" firestore:"name" validate:"required,min=1,max=50"`
    Color       string    `json:"color" firestore:"color" validate:"hexcolor"`
    Type        TagType   `json:"type" firestore:"type" validate:"required"`
    IsActive    bool      `json:"is_active" firestore:"is_active"`
    CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
    
    // Configurações
    Description string `json:"description" firestore:"description"`
    Category    string `json:"category" firestore:"category"`
    
    // Estatísticas
    UsageCount  int `json:"usage_count" firestore:"usage_count"`
}

// Tags predefinidas do sistema
var SystemTags = []Tag{
    {ID: "tag-confidential", Name: "Confidential", Color: "#FF5722", Type: TagTypeSystem, Category: "security"},
    {ID: "tag-public", Name: "Public", Color: "#4CAF50", Type: TagTypeSystem, Category: "security"},
    {ID: "tag-internal", Name: "Internal", Color: "#FF9800", Type: TagTypeSystem, Category: "security"},
    {ID: "tag-draft", Name: "Draft", Color: "#9E9E9E", Type: TagTypeSystem, Category: "status"},
    {ID: "tag-approved", Name: "Approved", Color: "#4CAF50", Type: TagTypeSystem, Category: "status"},
    {ID: "tag-review", Name: "Under Review", Color: "#FF9800", Type: TagTypeSystem, Category: "status"},
    {ID: "tag-archived", Name: "Archived", Color: "#795548", Type: TagTypeSystem, Category: "status"},
    {ID: "tag-important", Name: "Important", Color: "#F44336", Type: TagTypeSystem, Category: "priority"},
    {ID: "tag-urgent", Name: "Urgent", Color: "#E91E63", Type: TagTypeSystem, Category: "priority"},
    {ID: "tag-financial", Name: "Financial", Color: "#009688", Type: TagTypeSystem, Category: "department"},
    {ID: "tag-legal", Name: "Legal", Color: "#673AB7", Type: TagTypeSystem, Category: "department"},
    {ID: "tag-hr", Name: "HR", Color: "#3F51B5", Type: TagTypeSystem, Category: "department"},
    {ID: "tag-marketing", Name: "Marketing", Color: "#E91E63", Type: TagTypeSystem, Category: "department"},
    {ID: "tag-technical", Name: "Technical", Color: "#607D8B", Type: TagTypeSystem, Category: "department"},
}

func (t *Tag) GetOpenFGAID() string {
    return "tag:" + t.ID
}

// Normalizar nome da tag para evitar duplicatas
func NormalizeTagName(name string) string {
    return strings.ToLower(strings.TrimSpace(name))
}
```

### 8. **Certificate** - Certificado Digital

```go
// internal/core/entity/certificate.go
package entity

import (
    "errors"
    "time"
)

type Certificate struct {
    ID              string                 `json:"id" firestore:"id"`
    Name            string                 `json:"name" firestore:"name" validate:"required,min=1,max=200"`
    Description     string                 `json:"description,omitempty" firestore:"description"`
    VaultID         string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    TenantID        string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy       string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    CertificateData CertificateData        `json:"certificateData" firestore:"certificateData"`
    Details         CertificateDetails     `json:"details" firestore:"details"`
    Validity        CertificateValidity    `json:"validity" firestore:"validity"`
    Renewal         CertificateRenewal     `json:"renewal" firestore:"renewal"`
    Tags            []string               `json:"tags" firestore:"tags" validate:"max=5"`
    IsActive        bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt       time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt       time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

// Validação de tags (máximo 5)
func (c *Certificate) AddTag(tag string) error {
    if len(c.Tags) >= 5 {
        return errors.New("maximum of 5 tags allowed")
    }
    
    // Verificar se já existe
    for _, existingTag := range c.Tags {
        if existingTag == tag {
            return errors.New("tag already exists")
        }
    }
    
    c.Tags = append(c.Tags, tag)
    return nil
}

func (c *Certificate) RemoveTag(tag string) {
    for i, existingTag := range c.Tags {
        if existingTag == tag {
            c.Tags = append(c.Tags[:i], c.Tags[i+1:]...)
            return
        }
    }
}

func (c *Certificate) GetOpenFGAID() string {
    return "object:" + c.ID
}

type CertificateData struct {
    EncryptedCert  string                 `json:"encryptedCert" firestore:"encryptedCert"`
    EncryptedKey   string                 `json:"encryptedKey" firestore:"encryptedKey"`
    EncryptedChain string                 `json:"encryptedChain,omitempty" firestore:"encryptedChain"`
    Passphrase     string                 `json:"passphrase,omitempty" firestore:"passphrase"`
}

type CertificateDetails struct {
    CommonName       string               `json:"commonName" firestore:"commonName"`
    SubjectAltNames  []string             `json:"subjectAltNames" firestore:"subjectAltNames"`
    Issuer          string               `json:"issuer" firestore:"issuer"`
    SerialNumber    string               `json:"serialNumber" firestore:"serialNumber"`
    Fingerprint     string               `json:"fingerprint" firestore:"fingerprint"`
    Algorithm       string               `json:"algorithm" firestore:"algorithm"`
}

type CertificateValidity struct {
    NotBefore time.Time              `json:"notBefore" firestore:"notBefore"`
    NotAfter  time.Time              `json:"notAfter" firestore:"notAfter"`
    IsExpired bool                   `json:"isExpired" firestore:"isExpired"`
}

type CertificateRenewal struct {
    AutoRenew    bool                   `json:"autoRenew" firestore:"autoRenew"`
    RenewalDays  int                    `json:"renewalDays" firestore:"renewalDays"`
    Provider     string                 `json:"provider,omitempty" firestore:"provider"`
}
```

### 9. **SSHKey** - Chave SSH

```go
package domain

type SSHKey struct {
    ID          string                 `json:"id" firestore:"id"`
    Name        string                 `json:"name" firestore:"name" validate:"required,min=1,max=200"`
    Description string                 `json:"description,omitempty" firestore:"description"`
    VaultID     string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    TenantID    string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy   string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    KeyData     SSHKeyData             `json:"keyData" firestore:"keyData"`
    Usage       SSHKeyUsage            `json:"usage" firestore:"usage"`
    Metadata    SSHKeyMetadata         `json:"metadata" firestore:"metadata"`
    Tags        []string               `json:"tags" firestore:"tags"`
    LastUsedAt  *time.Time             `json:"lastUsedAt,omitempty" firestore:"lastUsedAt"`
    IsActive    bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt   time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt   time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

type SSHKeyData struct {
    EncryptedPrivateKey string             `json:"encryptedPrivateKey" firestore:"encryptedPrivateKey"`
    PublicKey          string             `json:"publicKey" firestore:"publicKey"`
    Passphrase         string             `json:"passphrase,omitempty" firestore:"passphrase"`
    KeyType            string             `json:"keyType" firestore:"keyType"`
}

type SSHKeyUsage struct {
    IsProduction    bool                   `json:"isProduction" firestore:"isProduction"`
    AuthorizedHosts []string               `json:"authorizedHosts" firestore:"authorizedHosts"`
    AllowedUsers    []string               `json:"allowedUsers" firestore:"allowedUsers"`
    Restrictions    map[string]interface{} `json:"restrictions,omitempty" firestore:"restrictions"`
}

type SSHKeyMetadata struct {
    KeySize     int                    `json:"keySize" firestore:"keySize"`
    Fingerprint string                 `json:"fingerprint" firestore:"fingerprint"`
    Comment     string                 `json:"comment,omitempty" firestore:"comment"`
}
```

### 10. **KeyValue** - Variáveis de Ambiente

```go
package domain

type KeyValue struct {
    ID          string                 `json:"id" firestore:"id"`
    Name        string                 `json:"name" firestore:"name" validate:"required,min=1,max=200"`
    Description string                 `json:"description,omitempty" firestore:"description"`
    VaultID     string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    TenantID    string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy   string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    Data        KeyValueData           `json:"data" firestore:"data"`
    Metadata    KeyValueMetadata       `json:"metadata" firestore:"metadata"`
    Tags        []string               `json:"tags" firestore:"tags"`
    IsActive    bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt   time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt   time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

type KeyValueData struct {
    EncryptedValues map[string]string      `json:"encryptedValues" firestore:"encryptedValues"`
    Environment     string                 `json:"environment" firestore:"environment"`
    Format          string                 `json:"format" firestore:"format"`
}

type KeyValueMetadata struct {
    TotalKeys   int                    `json:"totalKeys" firestore:"totalKeys"`
    LastSyncAt  *time.Time             `json:"lastSyncAt,omitempty" firestore:"lastSyncAt"`
    SyncSource  string                 `json:"syncSource,omitempty" firestore:"syncSource"`
}
```

### 11. **DatabaseConnection** - Conexão de Banco

```go
package domain

type DatabaseConnection struct {
    ID           string                 `json:"id" firestore:"id"`
    Name         string                 `json:"name" firestore:"name" validate:"required,min=1,max=200"`
    Description  string                 `json:"description,omitempty" firestore:"description"`
    VaultID      string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    TenantID     string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy    string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    Connection   DBConnectionData       `json:"connection" firestore:"connection"`
    SSL          DBSSLConfig            `json:"ssl" firestore:"ssl"`
    Pools        DBPoolConfig           `json:"pools" firestore:"pools"`
    IsProduction bool                   `json:"isProduction" firestore:"isProduction"`
    IsReadOnly   bool                   `json:"isReadOnly" firestore:"isReadOnly"`
    Tags         []string               `json:"tags" firestore:"tags"`
    LastTestedAt *time.Time             `json:"lastTestedAt,omitempty" firestore:"lastTestedAt"`
    IsActive     bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt    time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt    time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

type DBConnectionData struct {
    Type                     string     `json:"type" firestore:"type" validate:"required"`
    Host                     string     `json:"host" firestore:"host" validate:"required"`
    Port                     int        `json:"port" firestore:"port" validate:"required,min=1,max=65535"`
    Database                 string     `json:"database" firestore:"database" validate:"required"`
    EncryptedUsername        string     `json:"encryptedUsername" firestore:"encryptedUsername"`
    EncryptedPassword        string     `json:"encryptedPassword" firestore:"encryptedPassword"`
    EncryptedConnectionString string     `json:"encryptedConnectionString,omitempty" firestore:"encryptedConnectionString"`
}

type DBSSLConfig struct {
    Enabled     bool                   `json:"enabled" firestore:"enabled"`
    Certificate string                 `json:"certificate,omitempty" firestore:"certificate"`
    VerifyCA    bool                   `json:"verifyCA" firestore:"verifyCA"`
}

type DBPoolConfig struct {
    MaxConnections int                `json:"maxConnections" firestore:"maxConnections"`
    MinConnections int                `json:"minConnections" firestore:"minConnections"`
    Timeout        int                `json:"timeout" firestore:"timeout"` // seconds
}
```

### 12. **AuditLog** - Log de Auditoria

```go
package domain

type AuditLog struct {
    ID           string                 `json:"id" firestore:"id"`
    TenantID     string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    UserID       string                 `json:"userId" firestore:"userId" validate:"required"`
    Action       string                 `json:"action" firestore:"action" validate:"required"`
    ResourceType string                 `json:"resourceType" firestore:"resourceType" validate:"required"`
    ResourceID   string                 `json:"resourceId" firestore:"resourceId" validate:"required"`
    ResourceName string                 `json:"resourceName" firestore:"resourceName"`
    Details      AuditDetails           `json:"details" firestore:"details"`
    Result       AuditResult            `json:"result" firestore:"result" validate:"required"`
    Risk         RiskLevel              `json:"risk" firestore:"risk" validate:"required"`
    Session      AuditSession           `json:"session" firestore:"session"`
    Timestamp    time.Time              `json:"timestamp" firestore:"timestamp"`
    ExpiresAt    time.Time              `json:"expiresAt" firestore:"expiresAt"`
}

type AuditDetails struct {
    OldValues map[string]interface{} `json:"oldValues,omitempty" firestore:"oldValues"`
    NewValues map[string]interface{} `json:"newValues,omitempty" firestore:"newValues"`
    IPAddress string                 `json:"ipAddress" firestore:"ipAddress"`
    UserAgent string                 `json:"userAgent" firestore:"userAgent"`
    Location  string                 `json:"location,omitempty" firestore:"location"`
}

type AuditResult string

const (
    AuditResultSuccess AuditResult = "success"
    AuditResultFailure AuditResult = "failure"
    AuditResultDenied  AuditResult = "denied"
)

type RiskLevel string

const (
    RiskLow    RiskLevel = "low"
    RiskMedium RiskLevel = "medium"
    RiskHigh   RiskLevel = "high"
)

type AuditSession struct {
    SessionID   string                 `json:"sessionId" firestore:"sessionId"`
    DeviceID    string                 `json:"deviceId" firestore:"deviceId"`
    MFAVerified bool                   `json:"mfaVerified" firestore:"mfaVerified"`
}
```

### 13. **Session** - Sessão do Usuário

```go
package domain

type Session struct {
    ID           string                 `json:"id" firestore:"id"`
    UserID       string                 `json:"userId" firestore:"userId" validate:"required"`
    TenantID     string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    DeviceID     string                 `json:"deviceId" firestore:"deviceId" validate:"required"`
    DeviceInfo   DeviceInfo             `json:"deviceInfo" firestore:"deviceInfo"`
    Location     LocationInfo           `json:"location" firestore:"location"`
    Security     SecurityInfo           `json:"security" firestore:"security"`
    IsActive     bool                   `json:"isActive" firestore:"isActive"`
    CreatedAt    time.Time              `json:"createdAt" firestore:"createdAt"`
    LastAccessAt time.Time              `json:"lastAccessAt" firestore:"lastAccessAt"`
    ExpiresAt    time.Time              `json:"expiresAt" firestore:"expiresAt"`
}

type DeviceInfo struct {
    Type    string                 `json:"type" firestore:"type"`
    OS      string                 `json:"os" firestore:"os"`
    Browser string                 `json:"browser" firestore:"browser"`
    Version string                 `json:"version" firestore:"version"`
}

type LocationInfo struct {
    Country   string                 `json:"country" firestore:"country"`
    Region    string                 `json:"region" firestore:"region"`
    City      string                 `json:"city" firestore:"city"`
    IPAddress string                 `json:"ipAddress" firestore:"ipAddress"`
}

type SecurityInfo struct {
    MFAVerified  bool                   `json:"mfaVerified" firestore:"mfaVerified"`
    MFAAt        *time.Time             `json:"mfaAt,omitempty" firestore:"mfaAt"`
    RiskScore    int                    `json:"riskScore" firestore:"riskScore"`
    IsSuspicious bool                   `json:"isSuspicious" firestore:"isSuspicious"`
}
```

### 14. **ExternalShareRequest** - Solicitação Compartilhamento

```go
package domain

type ExternalShareRequest struct {
    ID           string                 `json:"id" firestore:"id"`
    FromTenantID string                 `json:"fromTenantId" firestore:"fromTenantId" validate:"required"`
    ToTenantID   string                 `json:"toTenantId" firestore:"toTenantId" validate:"required"`
    RequesterID  string                 `json:"requesterId" firestore:"requesterId" validate:"required"`
    VaultID      string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    Permissions  []string               `json:"permissions" firestore:"permissions" validate:"required"`
    Duration     int                    `json:"duration" firestore:"duration" validate:"required,min=1"` // days
    Message      string                 `json:"message,omitempty" firestore:"message"`
    Approvals    []ShareApproval        `json:"approvals" firestore:"approvals"`
    Status       ShareStatus            `json:"status" firestore:"status"`
    ExpiresAt    time.Time              `json:"expiresAt" firestore:"expiresAt"`
    CreatedAt    time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt    time.Time              `json:"updatedAt" firestore:"updatedAt"`
}

type ShareApproval struct {
    UserID     string                 `json:"userId" firestore:"userId"`
    TenantID   string                 `json:"tenantId" firestore:"tenantId"`
    Status     ApprovalStatus         `json:"status" firestore:"status"`
    ApprovedAt *time.Time             `json:"approvedAt,omitempty" firestore:"approvedAt"`
    Comment    string                 `json:"comment,omitempty" firestore:"comment"`
}

type ShareStatus string

const (
    ShareStatusPending  ShareStatus = "pending"
    ShareStatusApproved ShareStatus = "approved"
    ShareStatusRejected ShareStatus = "rejected"
    ShareStatusExpired  ShareStatus = "expired"
)

type ApprovalStatus string

const (
    ApprovalStatusPending  ApprovalStatus = "pending"
    ApprovalStatusApproved ApprovalStatus = "approved"
    ApprovalStatusRejected ApprovalStatus = "rejected"
)
```

### 15. **Notification** - Notificação

```go
package domain

type Notification struct {
    ID        string                 `json:"id" firestore:"id"`
    TenantID  string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    UserID    string                 `json:"userId" firestore:"userId" validate:"required"`
    Type      NotificationType       `json:"type" firestore:"type" validate:"required"`
    Title     string                 `json:"title" firestore:"title" validate:"required"`
    Message   string                 `json:"message" firestore:"message" validate:"required"`
    Data      map[string]interface{} `json:"data,omitempty" firestore:"data"`
    Channels  []NotificationChannel  `json:"channels" firestore:"channels"`
    Priority  NotificationPriority   `json:"priority" firestore:"priority"`
    IsRead    bool                   `json:"isRead" firestore:"isRead"`
    ReadAt    *time.Time             `json:"readAt,omitempty" firestore:"readAt"`
    CreatedAt time.Time              `json:"createdAt" firestore:"createdAt"`
    ExpiresAt *time.Time             `json:"expiresAt,omitempty" firestore:"expiresAt"`
}

type NotificationType string

const (
    NotificationTypeSecurity    NotificationType = "security"
    NotificationTypeExpiration  NotificationType = "expiration"
    NotificationTypeSharing     NotificationType = "sharing"
    NotificationTypeSystem      NotificationType = "system"
    NotificationTypeMarketing   NotificationType = "marketing"
)

type NotificationChannel string

const (
    ChannelEmail   NotificationChannel = "email"
    ChannelBrowser NotificationChannel = "browser"
    ChannelMobile  NotificationChannel = "mobile"
    ChannelSlack   NotificationChannel = "slack"
)

type NotificationPriority string

const (
    PriorityLow      NotificationPriority = "low"
    PriorityNormal   NotificationPriority = "normal"
    PriorityHigh     NotificationPriority = "high"
    PriorityCritical NotificationPriority = "critical"
)
```

### 16. **Backup** - Backup do Sistema

```go
package domain

type Backup struct {
    ID          string                 `json:"id" firestore:"id"`
    TenantID    string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    VaultID     string                 `json:"vaultId" firestore:"vaultId" validate:"required"`
    CreatedBy   string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    Type        BackupType             `json:"type" firestore:"type" validate:"required"`
    Status      BackupStatus           `json:"status" firestore:"status"`
    Metadata    BackupMetadata         `json:"metadata" firestore:"metadata"`
    Storage     BackupStorage          `json:"storage" firestore:"storage"`
    CreatedAt   time.Time              `json:"createdAt" firestore:"createdAt"`
    CompletedAt *time.Time             `json:"completedAt,omitempty" firestore:"completedAt"`
}

type BackupType string

const (
    BackupTypeManual    BackupType = "manual"
    BackupTypeScheduled BackupType = "scheduled"
)

type BackupStatus string

const (
    BackupStatusInProgress BackupStatus = "in_progress"
    BackupStatusCompleted  BackupStatus = "completed"
    BackupStatusFailed     BackupStatus = "failed"
)

type BackupMetadata struct {
    TotalSecrets  int                    `json:"totalSecrets" firestore:"totalSecrets"`
    TotalSize     int64                  `json:"totalSize" firestore:"totalSize"` // bytes
    EncryptionKey string                 `json:"encryptionKey" firestore:"encryptionKey"`
    Checksum      string                 `json:"checksum" firestore:"checksum"`
}

type BackupStorage struct {
    Provider  string                 `json:"provider" firestore:"provider"`
    Location  string                 `json:"location" firestore:"location"`
    ExpiresAt *time.Time             `json:"expiresAt,omitempty" firestore:"expiresAt"`
}
```

## 🔧 Interfaces de Repositório

### Repository Interfaces

```go
// internal/core/repository/user.go
package repository

import (
    "context"
    "github.com/lockari/internal/core/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id string) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

type TenantRepository interface {
    Create(ctx context.Context, tenant *entity.Tenant) error
    GetByID(ctx context.Context, id string) (*entity.Tenant, error)
    GetBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
    Update(ctx context.Context, tenant *entity.Tenant) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*entity.Tenant, error)
    GetUserTenants(ctx context.Context, userID string) ([]*entity.Tenant, error)
}

type VaultRepository interface {
    Create(ctx context.Context, vault *entity.Vault) error
    GetByID(ctx context.Context, id string) (*entity.Vault, error)
    Update(ctx context.Context, vault *entity.Vault) error
    Delete(ctx context.Context, id string) error
    ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.Vault, error)
    ListByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.Vault, error)
}

type ObjectRepository interface {
    Create(ctx context.Context, object *entity.Object) error
    GetByID(ctx context.Context, id string) (*entity.Object, error)
    Update(ctx context.Context, object *entity.Object) error
    Delete(ctx context.Context, id string) error
    ListByVault(ctx context.Context, vaultID string, limit, offset int) ([]*entity.Object, error)
    Search(ctx context.Context, tenantID, query string, limit, offset int) ([]*entity.Object, error)
}

type TagRepository interface {
    Create(ctx context.Context, tag *entity.Tag) error
    GetByID(ctx context.Context, id string) (*entity.Tag, error)
    GetByName(ctx context.Context, tenantID, name string) (*entity.Tag, error)
    Update(ctx context.Context, tag *entity.Tag) error
    Delete(ctx context.Context, id string) error
    ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.Tag, error)
    ListByCategory(ctx context.Context, tenantID, category string, limit, offset int) ([]*entity.Tag, error)
    GetSystemTags() ([]*entity.Tag, error)
    IncrementUsage(ctx context.Context, tagID string) error
    DecrementUsage(ctx context.Context, tagID string) error
}

type AuditLogRepository interface {
    Create(ctx context.Context, log *entity.AuditLog) error
    List(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]*entity.AuditLog, error)
    GetByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.AuditLog, error)
    GetByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*entity.AuditLog, error)
}
```

## 🔐 Services de Autorização

### Authorization Service

```go
// internal/core/services/authorization.go
package services

import (
    "context"
    "fmt"
    "github.com/lockari/internal/core/entity"
    "github.com/lockari/internal/core/repository"
)

type AuthorizationService struct {
    openFGAClient OpenFGAClient
    userRepo      repository.UserRepository
    tenantRepo    repository.TenantRepository
}

func NewAuthorizationService(
    openFGAClient OpenFGAClient,
    userRepo repository.UserRepository,
    tenantRepo repository.TenantRepository,
) *AuthorizationService {
    return &AuthorizationService{
        openFGAClient: openFGAClient,
        userRepo:      userRepo,
        tenantRepo:    tenantRepo,
    }
}

func (s *AuthorizationService) CheckPermission(
    ctx context.Context,
    userID string,
    permission string,
    resourceType string,
    resourceID string,
) (bool, error) {
    return s.openFGAClient.Check(ctx, &CheckRequest{
        User:     fmt.Sprintf("user:%s", userID),
        Relation: permission,
        Object:   fmt.Sprintf("%s:%s", resourceType, resourceID),
    })
}

func (s *AuthorizationService) ListUserObjects(
    ctx context.Context,
    userID string,
    permission string,
    objectType string,
) ([]string, error) {
    return s.openFGAClient.ListObjects(ctx, &ListObjectsRequest{
        User:     fmt.Sprintf("user:%s", userID),
        Relation: permission,
        Type:     objectType,
    })
}

// Estruturas auxiliares para OpenFGA
type CheckRequest struct {
    User     string `json:"user"`
    Relation string `json:"relation"`
    Object   string `json:"object"`
}

type ListObjectsRequest struct {
    User     string `json:"user"`
    Relation string `json:"relation"`
    Type     string `json:"type"`
}

type OpenFGAClient interface {
    Check(ctx context.Context, req *CheckRequest) (bool, error)
    ListObjects(ctx context.Context, req *ListObjectsRequest) ([]string, error)
}
```

### Tag Service

```go
// internal/core/services/tag.go
package services

import (
    "context"
    "strings"
    "github.com/lockari/internal/core/entity"
    "github.com/lockari/internal/core/repository"
)

type TagService struct {
    tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) *TagService {
    return &TagService{tagRepo: tagRepo}
}

func (s *TagService) CreateTag(ctx context.Context, tag *entity.Tag) error {
    // Normalizar nome
    tag.Name = s.normalizeTagName(tag.Name)
    
    // Verificar se já existe
    existingTag, err := s.tagRepo.GetByName(ctx, tag.TenantID, tag.Name)
    if err == nil && existingTag != nil {
        return errors.New("tag already exists")
    }
    
    return s.tagRepo.Create(ctx, tag)
}

func (s *TagService) GetOrCreateTag(ctx context.Context, tenantID, tagName string) (*entity.Tag, error) {
    normalizedName := s.normalizeTagName(tagName)
    
    // Tentar buscar tag existente
    tag, err := s.tagRepo.GetByName(ctx, tenantID, normalizedName)
    if err == nil && tag != nil {
        return tag, nil
    }
    
    // Criar nova tag
    newTag := &entity.Tag{
        TenantID:  tenantID,
        Name:      normalizedName,
        Type:      entity.TagTypeCustom,
        IsActive:  true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    err = s.tagRepo.Create(ctx, newTag)
    if err != nil {
        return nil, err
    }
    
    return newTag, nil
}

func (s *TagService) GetSuggestedTags(ctx context.Context, tenantID string, query string) ([]*entity.Tag, error) {
    tags, err := s.tagRepo.ListByTenant(ctx, tenantID, 50, 0)
    if err != nil {
        return nil, err
    }
    
    var suggested []*entity.Tag
    query = strings.ToLower(query)
    
    for _, tag := range tags {
        if strings.Contains(strings.ToLower(tag.Name), query) {
            suggested = append(suggested, tag)
        }
    }
    
    return suggested, nil
}

func (s *TagService) normalizeTagName(name string) string {
    return strings.ToLower(strings.TrimSpace(name))
}
```

## 🎯 Validação e Transformação

### DTO Structs

```go
// internal/adapters/http/dto/vault.go
package dto

import (
    "time"
    "github.com/lockari/internal/core/entity"
)

type CreateVaultRequest struct {
    Name        string   `json:"name" validate:"required,min=2,max=100"`
    Description string   `json:"description,omitempty"`
    Icon        string   `json:"icon,omitempty"`
    Color       string   `json:"color,omitempty" validate:"omitempty,hexcolor"`
    Tags        []string `json:"tags,omitempty" validate:"max=5"`
}

type CreateSecretRequest struct {
    Name        string                 `json:"name" validate:"required,min=1,max=200"`
    Description string                 `json:"description,omitempty"`
    Type        entity.SecretType      `json:"type" validate:"required"`
    Value       string                 `json:"value" validate:"required"`
    Metadata    entity.SecretMetadata  `json:"metadata,omitempty"`
    Tags        []string               `json:"tags,omitempty" validate:"max=5"`
    IsSensitive bool                   `json:"isSensitive"`
}

type CreateTagRequest struct {
    Name        string   `json:"name" validate:"required,min=1,max=50"`
    Description string   `json:"description,omitempty"`
    Color       string   `json:"color,omitempty" validate:"omitempty,hexcolor"`
    Category    string   `json:"category,omitempty"`
}

type VaultResponse struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    Description  string                 `json:"description,omitempty"`
    Icon         string                 `json:"icon,omitempty"`
    Color        string                 `json:"color,omitempty"`
    Tags         []TagResponse          `json:"tags"`
    SecretsCount int                    `json:"secretsCount"`
    Permissions  []string               `json:"permissions"`
    CreatedAt    time.Time              `json:"createdAt"`
    UpdatedAt    time.Time              `json:"updatedAt"`
}

type TagResponse struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Color       string     `json:"color,omitempty"`
    Category    string     `json:"category,omitempty"`
    Type        string     `json:"type"`
    UsageCount  int        `json:"usage_count"`
}
```

## 📈 Considerações de Performance

### 1. **Índices de Banco**
```go
// Firestore indexes configuration
type FirestoreIndex struct {
    Collection string
    Fields     []IndexField
}

type IndexField struct {
    Name       string
    Direction  string // ASC, DESC
}

var RequiredIndexes = []FirestoreIndex{
    {
        Collection: "vaults",
        Fields: []IndexField{
            {Name: "tenantId", Direction: "ASC"},
            {Name: "isActive", Direction: "ASC"},
            {Name: "updatedAt", Direction: "DESC"},
        },
    },
    // ... more indexes
}
```

### 2. **Cache Strategy**
```go
type CacheService struct {
    redis        *redis.Client
    defaultTTL   time.Duration
}

func (c *CacheService) CacheUserPermissions(
    ctx context.Context,
    userID string,
    permissions map[string]bool,
) error {
    key := fmt.Sprintf("permissions:%s", userID)
    return c.redis.Set(ctx, key, permissions, c.defaultTTL).Err()
}
```

## 🔧 Configuração

### Config Struct

```go
package config

type Config struct {
    Server    ServerConfig    `yaml:"server"`
    Firebase  FirebaseConfig  `yaml:"firebase"`
    OpenFGA   OpenFGAConfig   `yaml:"openfga"`
    Redis     RedisConfig     `yaml:"redis"`
    Security  SecurityConfig  `yaml:"security"`
}

type ServerConfig struct {
    Port         int           `yaml:"port"`
    Host         string        `yaml:"host"`
    ReadTimeout  time.Duration `yaml:"readTimeout"`
    WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type FirebaseConfig struct {
    ProjectID        string `yaml:"projectId"`
    CredentialsPath  string `yaml:"credentialsPath"`
    DatabaseURL      string `yaml:"databaseUrl"`
}

type OpenFGAConfig struct {
    Host    string `yaml:"host"`
    Port    int    `yaml:"port"`
    StoreID string `yaml:"storeId"`
}

type SecurityConfig struct {
    JWTSecret       string        `yaml:"jwtSecret"`
    EncryptionKey   string        `yaml:"encryptionKey"`
    SessionTimeout  time.Duration `yaml:"sessionTimeout"`
    MFARequired     bool          `yaml:"mfaRequired"`
}
```

## 🏷️ Sistema de Tags

### Características

1. **Limite de 5 tags por vault/objeto** - Para facilitar organização e pesquisa
2. **Tags predefinidas** - Sistema fornece tags comuns para evitar duplicatas
3. **Normalização automática** - Nomes são normalizados (lowercase, trim) para consistência
4. **Contagem de uso** - Tracking de quantas vezes cada tag é usada
5. **Categorização** - Tags organizadas por categoria (security, status, priority, department)

### Tags Predefinidas do Sistema

```go
// Tags de Segurança
"confidential" - Informações confidenciais (#FF5722)
"public" - Informações públicas (#4CAF50)
"internal" - Uso interno apenas (#FF9800)

// Tags de Status
"draft" - Rascunho (#9E9E9E)
"approved" - Aprovado (#4CAF50)
"review" - Em revisão (#FF9800)
"archived" - Arquivado (#795548)

// Tags de Prioridade
"important" - Importante (#F44336)
"urgent" - Urgente (#E91E63)

// Tags de Departamento
"financial" - Financeiro (#009688)
"legal" - Jurídico (#673AB7)
"hr" - Recursos Humanos (#3F51B5)
"marketing" - Marketing (#E91E63)
"technical" - Técnico (#607D8B)
```

### Fluxo de Uso

1. **Criação de Vault/Objeto**:
   - Sistema sugere tags existentes baseado no nome/tipo
   - Usuário pode selecionar até 5 tags
   - Novas tags são criadas automaticamente se não existirem

2. **Busca e Filtro**:
   - Filtros por tag para encontrar rapidamente recursos
   - Autocomplete com tags existentes
   - Pesquisa por categoria de tag

3. **Gestão de Tags**:
   - Admins podem gerenciar tags personalizadas
   - Estatísticas de uso de tags
   - Limpeza automática de tags não utilizadas
