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

### 17. **Token** - Token de API para Automação

```go
// internal/core/entity/token.go
package entity

import (
    "errors"
    "time"
)

type Token struct {
    ID            string                 `json:"id" firestore:"id"`
    Name          string                 `json:"name" firestore:"name" validate:"required,min=2,max=100"`
    Description   string                 `json:"description,omitempty" firestore:"description"`
    TenantID      string                 `json:"tenantId" firestore:"tenantId" validate:"required"`
    CreatedBy     string                 `json:"createdBy" firestore:"createdBy" validate:"required"`
    VaultIDs      []string               `json:"vaultIds" firestore:"vaultIds" validate:"required,min=1"`
    Permissions   TokenPermissions       `json:"permissions" firestore:"permissions"`
    Restrictions  TokenRestrictions      `json:"restrictions" firestore:"restrictions"`
    Usage         TokenUsage             `json:"usage" firestore:"usage"`
    Security      TokenSecurity          `json:"security" firestore:"security"`
    HashedToken   string                 `json:"hashedToken" firestore:"hashedToken" validate:"required"`
    LastUsedAt    *time.Time             `json:"lastUsedAt,omitempty" firestore:"lastUsedAt"`
    UsageCount    int                    `json:"usageCount" firestore:"usageCount"`
    IsActive      bool                   `json:"isActive" firestore:"isActive"`
    IsRevoked     bool                   `json:"isRevoked" firestore:"isRevoked"`
    CreatedAt     time.Time              `json:"createdAt" firestore:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt" firestore:"updatedAt"`
    ExpiresAt     *time.Time             `json:"expiresAt,omitempty" firestore:"expiresAt"`
}

type TokenPermissions struct {
    CanRead       bool                   `json:"canRead" firestore:"canRead"`
    CanWrite      bool                   `json:"canWrite" firestore:"canWrite"`
    CanDelete     bool                   `json:"canDelete" firestore:"canDelete"`
    CanManage     bool                   `json:"canManage" firestore:"canManage"`
    CanBackup     bool                   `json:"canBackup" firestore:"canBackup"`
    CanAudit      bool                   `json:"canAudit" firestore:"canAudit"`
    ReadOnly      bool                   `json:"readOnly" firestore:"readOnly"`
    Scopes        []string               `json:"scopes" firestore:"scopes"` // specific operations
}

type TokenRestrictions struct {
    IPWhitelist   []string               `json:"ipWhitelist" firestore:"ipWhitelist"`
    UserAgent     string                 `json:"userAgent,omitempty" firestore:"userAgent"`
    RateLimit     TokenRateLimit         `json:"rateLimit" firestore:"rateLimit"`
    TimeWindow    TokenTimeWindow        `json:"timeWindow" firestore:"timeWindow"`
    Environment   []string               `json:"environment" firestore:"environment"` // prod, dev, staging
}

type TokenRateLimit struct {
    Enabled       bool                   `json:"enabled" firestore:"enabled"`
    RequestsPerHour int                  `json:"requestsPerHour" firestore:"requestsPerHour"`
    RequestsPerDay  int                  `json:"requestsPerDay" firestore:"requestsPerDay"`
    BurstLimit    int                    `json:"burstLimit" firestore:"burstLimit"`
}

type TokenTimeWindow struct {
    Enabled       bool                   `json:"enabled" firestore:"enabled"`
    StartTime     string                 `json:"startTime,omitempty" firestore:"startTime"` // HH:MM
    EndTime       string                 `json:"endTime,omitempty" firestore:"endTime"`     // HH:MM
    DaysOfWeek    []string               `json:"daysOfWeek,omitempty" firestore:"daysOfWeek"` // mon, tue, etc
    Timezone      string                 `json:"timezone,omitempty" firestore:"timezone"`
}

type TokenUsage struct {
    TotalRequests    int                 `json:"totalRequests" firestore:"totalRequests"`
    LastRequestAt    *time.Time          `json:"lastRequestAt,omitempty" firestore:"lastRequestAt"`
    LastRequestIP    string              `json:"lastRequestIP,omitempty" firestore:"lastRequestIP"`
    LastUserAgent    string              `json:"lastUserAgent,omitempty" firestore:"lastUserAgent"`
    SuccessfulReqs   int                 `json:"successfulReqs" firestore:"successfulReqs"`
    FailedReqs       int                 `json:"failedReqs" firestore:"failedReqs"`
    RateLimitHits    int                 `json:"rateLimitHits" firestore:"rateLimitHits"`
}

type TokenSecurity struct {
    RequireHTTPS     bool                `json:"requireHTTPS" firestore:"requireHTTPS"`
    AllowedOrigins   []string            `json:"allowedOrigins,omitempty" firestore:"allowedOrigins"`
    RequireUserAgent bool                `json:"requireUserAgent" firestore:"requireUserAgent"`
    LogAllRequests   bool                `json:"logAllRequests" firestore:"logAllRequests"`
    AlertOnSuspicious bool               `json:"alertOnSuspicious" firestore:"alertOnSuspicious"`
}

// Métodos de validação e controle
func (t *Token) IsValid() bool {
    if t.IsRevoked || !t.IsActive {
        return false
    }
    
    if t.ExpiresAt != nil && t.ExpiresAt.Before(time.Now()) {
        return false
    }
    
    return true
}

func (t *Token) IsExpired() bool {
    return t.ExpiresAt != nil && t.ExpiresAt.Before(time.Now())
}

func (t *Token) CanAccessVault(vaultID string) bool {
    for _, id := range t.VaultIDs {
        if id == vaultID {
            return true
        }
    }
    return false
}

func (t *Token) Revoke() {
    t.IsRevoked = true
    t.UpdatedAt = time.Now()
}

func (t *Token) Activate() {
    t.IsActive = true
    t.UpdatedAt = time.Now()
}

func (t *Token) Deactivate() {
    t.IsActive = false
    t.UpdatedAt = time.Now()
}

func (t *Token) RecordUsage(ip, userAgent string, success bool) {
    t.Usage.TotalRequests++
    t.Usage.LastRequestAt = &time.Time{}
    *t.Usage.LastRequestAt = time.Now()
    t.Usage.LastRequestIP = ip
    t.Usage.LastUserAgent = userAgent
    t.LastUsedAt = t.Usage.LastRequestAt
    t.UsageCount++
    
    if success {
        t.Usage.SuccessfulReqs++
    } else {
        t.Usage.FailedReqs++
    }
    
    t.UpdatedAt = time.Now()
}

func (t *Token) CheckRateLimit() error {
    if !t.Restrictions.RateLimit.Enabled {
        return nil
    }
    
    // Implementar lógica de rate limiting
    // Esta seria uma implementação simplificada
    if t.Usage.TotalRequests >= t.Restrictions.RateLimit.RequestsPerDay {
        t.Usage.RateLimitHits++
        return errors.New("daily rate limit exceeded")
    }
    
    return nil
}

func (t *Token) CheckTimeWindow() error {
    if !t.Restrictions.TimeWindow.Enabled {
        return nil
    }
    
    now := time.Now()
    
    // Verificar dia da semana
    if len(t.Restrictions.TimeWindow.DaysOfWeek) > 0 {
        currentDay := now.Weekday().String()[:3] // Mon, Tue, etc
        allowed := false
        for _, day := range t.Restrictions.TimeWindow.DaysOfWeek {
            if strings.ToLower(day) == strings.ToLower(currentDay) {
                allowed = true
                break
            }
        }
        if !allowed {
            return errors.New("token not allowed on this day of week")
        }
    }
    
    // Verificar horário (implementação simplificada)
    if t.Restrictions.TimeWindow.StartTime != "" && t.Restrictions.TimeWindow.EndTime != "" {
        // Lógica de verificação de horário seria implementada aqui
        return nil
    }
    
    return nil
}

func (t *Token) CheckIPWhitelist(ip string) error {
    if len(t.Restrictions.IPWhitelist) == 0 {
        return nil
    }
    
    for _, allowedIP := range t.Restrictions.IPWhitelist {
        if allowedIP == ip {
            return nil
        }
    }
    
    return errors.New("IP not in whitelist")
}

func (t *Token) GetOpenFGAID() string {
    return "token:" + t.ID
}

func (t *Token) HasPermission(permission string) bool {
    switch permission {
    case "read":
        return t.Permissions.CanRead
    case "write":
        return t.Permissions.CanWrite && !t.Permissions.ReadOnly
    case "delete":
        return t.Permissions.CanDelete && !t.Permissions.ReadOnly
    case "manage":
        return t.Permissions.CanManage && !t.Permissions.ReadOnly
    case "backup":
        return t.Permissions.CanBackup
    case "audit":
        return t.Permissions.CanAudit
    default:
        return false
    }
}

func (t *Token) HasScope(scope string) bool {
    for _, s := range t.Permissions.Scopes {
        if s == scope {
            return true
        }
    }
    return false
}
```
```go
// internal/core/repository/token_repository.go
package repository

import (
    "context"
    "github.com/yourusername/yourproject/internal/core/entity"
)

type TokenRepository interface {
    Create(ctx context.Context, token *entity.Token) error
    GetByID(ctx context.Context, id string) (*entity.Token, error)
    GetByHash(ctx context.Context, hashedToken string) (*entity.Token, error)
    Update(ctx context.Context, token *entity.Token) error
    Delete(ctx context.Context, id string) error
    ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.Token, error)
    ListByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.Token, error)
    ListByVault(ctx context.Context, vaultID string, limit, offset int) ([]*entity.Token, error)
    RevokeByUser(ctx context.Context, userID string) error
    RevokeByVault(ctx context.Context, vaultID string) error
    GetActiveTokens(ctx context.Context, tenantID string) ([]*entity.Token, error)
    GetExpiredTokens(ctx context.Context) ([]*entity.Token, error)
    RecordUsage(ctx context.Context, tokenID, ip, userAgent string, success bool) error
}
```
```go
// internal/core/services/token.go
package services

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"
    
    "github.com/lockari/internal/core/entity"
    "github.com/lockari/internal/core/repository"
)

type TokenService struct {
    tokenRepo  repository.TokenRepository
    vaultRepo  repository.VaultRepository
    auditRepo  repository.AuditLogRepository
    authzService *AuthorizationService
}

func NewTokenService(
    tokenRepo repository.TokenRepository,
    vaultRepo repository.VaultRepository,
    auditRepo repository.AuditLogRepository,
    authzService *AuthorizationService,
) *TokenService {
    return &TokenService{
        tokenRepo:    tokenRepo,
        vaultRepo:    vaultRepo,
        auditRepo:    auditRepo,
        authzService: authzService,
    }
}

func (s *TokenService) CreateToken(ctx context.Context, userID, tenantID string, req *CreateTokenRequest) (*entity.Token, string, error) {
    // Verificar se o usuário pode gerenciar os vaults especificados
    for _, vaultID := range req.VaultIDs {
        allowed, err := s.authzService.CheckPermission(ctx, userID, "can_manage", "vault", vaultID)
        if err != nil {
            return nil, "", err
        }
        if !allowed {
            return nil, "", errors.New("permission denied for vault: " + vaultID)
        }
    }
    
    // Gerar token único
    tokenValue, hashedToken, err := s.generateToken()
    if err != nil {
        return nil, "", err
    }
    
    // Criar entidade do token
    token := &entity.Token{
        ID:           generateID(),
        Name:         req.Name,
        Description:  req.Description,
        TenantID:     tenantID,
        CreatedBy:    userID,
        VaultIDs:     req.VaultIDs,
        Permissions:  req.Permissions,
        Restrictions: req.Restrictions,
        Security:     req.Security,
        HashedToken:  hashedToken,
        IsActive:     true,
        IsRevoked:    false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        ExpiresAt:    req.ExpiresAt,
    }
    
    // Salvar no banco
    err = s.tokenRepo.Create(ctx, token)
    if err != nil {
        return nil, "", err
    }
    
    // Criar relacionamentos no OpenFGA
    err = s.createTokenRelationships(ctx, token)
    if err != nil {
        // Rollback se falhar
        s.tokenRepo.Delete(ctx, token.ID)
        return nil, "", err
    }
    
    // Log de auditoria
    s.auditRepo.Create(ctx, &entity.AuditLog{
        TenantID:     tenantID,
        UserID:       userID,
        Action:       "token_created",
        ResourceType: "token",
        ResourceID:   token.ID,
        ResourceName: token.Name,
        Result:       entity.AuditResultSuccess,
        Risk:         entity.RiskLow,
        Timestamp:    time.Now(),
    })
    
    return token, tokenValue, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, tokenValue string) (*entity.Token, error) {
    hashedToken := s.hashToken(tokenValue)
    
    token, err := s.tokenRepo.GetByHash(ctx, hashedToken)
    if err != nil {
        return nil, errors.New("invalid token")
    }
    
    if !token.IsValid() {
        return nil, errors.New("token is invalid, revoked or expired")
    }
    
    return token, nil
}

func (s *TokenService) AuthorizeTokenAction(ctx context.Context, token *entity.Token, action, resourceType, resourceID string) error {
    // Verificar se o token pode acessar o vault
    if resourceType == "vault" {
        if !token.CanAccessVault(resourceID) {
            return errors.New("token cannot access this vault")
        }
    }
    
    // Verificar permissões específicas
    if !token.HasPermission(action) {
        return errors.New("token does not have permission for this action")
    }
    
    // Verificar restrições de IP
    // (IP seria passado via context ou parâmetro adicional)
    
    // Verificar rate limiting
    err := token.CheckRateLimit()
    if err != nil {
        return err
    }
    
    // Verificar janela de tempo
    err = token.CheckTimeWindow()
    if err != nil {
        return err
    }
    
    // Verificar via OpenFGA
    allowed, err := s.authzService.CheckPermission(ctx, 
        fmt.Sprintf("token:%s", token.ID), 
        fmt.Sprintf("can_%s_with_token", action), 
        resourceType, 
        resourceID)
    if err != nil {
        return err
    }
    
    if !allowed {
        return errors.New("token authorization failed")
    }
    
    return nil
}

func (s *TokenService) RevokeToken(ctx context.Context, userID, tokenID string) error {
    token, err := s.tokenRepo.GetByID(ctx, tokenID)
    if err != nil {
        return err
    }
    
    // Verificar se o usuário pode revogar este token
    if token.CreatedBy != userID {
        // Verificar se é admin do tenant
        allowed, err := s.authzService.CheckPermission(ctx, userID, "admin", "tenant", token.TenantID)
        if err != nil {
            return err
        }
        if !allowed {
            return errors.New("permission denied")
        }
    }
    
    // Revogar token
    token.Revoke()
    err = s.tokenRepo.Update(ctx, token)
    if err != nil {
        return err
    }
    
    // Remover relacionamentos do OpenFGA
    err = s.removeTokenRelationships(ctx, token)
    if err != nil {
        return err
    }
    
    // Log de auditoria
    s.auditRepo.Create(ctx, &entity.AuditLog{
        TenantID:     token.TenantID,
        UserID:       userID,
        Action:       "token_revoked",
        ResourceType: "token",
        ResourceID:   token.ID,
        ResourceName: token.Name,
        Result:       entity.AuditResultSuccess,
        Risk:         entity.RiskMedium,
        Timestamp:    time.Now(),
    })
    
    return nil
}

func (s *TokenService) generateToken() (string, string, error) {
    // Gerar token de 32 bytes
    tokenBytes := make([]byte, 32)
    _, err := rand.Read(tokenBytes)
    if err != nil {
        return "", "", err
    }
    
    tokenValue := hex.EncodeToString(tokenBytes)
    hashedToken := s.hashToken(tokenValue)
    
    return tokenValue, hashedToken, nil
}

func (s *TokenService) hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}

func (s *TokenService) createTokenRelationships(ctx context.Context, token *entity.Token) error {
    // Criar relacionamentos no OpenFGA
    // Implementação específica do OpenFGA seria feita aqui
    return nil
}

func (s *TokenService) removeTokenRelationships(ctx context.Context, token *entity.Token) error {
    // Remover relacionamentos no OpenFGA
    // Implementação específica do OpenFGA seria feita aqui
    return nil
}

type CreateTokenRequest struct {
    Name         string                      `json:"name" validate:"required,min=2,max=100"`
    Description  string                      `json:"description,omitempty"`
    VaultIDs     []string                    `json:"vaultIds" validate:"required,min=1"`
    Permissions  entity.TokenPermissions     `json:"permissions" validate:"required"`
    Restrictions entity.TokenRestrictions    `json:"restrictions"`
    Security     entity.TokenSecurity        `json:"security"`
    ExpiresAt    *time.Time                  `json:"expiresAt,omitempty"`
}

type UpdateTokenRequest struct {
    Name         string                      `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Description  string                      `json:"description,omitempty"`
    Permissions  entity.TokenPermissions     `json:"permissions,omitempty"`
    Restrictions entity.TokenRestrictions    `json:"restrictions,omitempty"`
    Security     entity.TokenSecurity        `json:"security,omitempty"`
    IsActive     *bool                       `json:"isActive,omitempty"`
}

type TokenResponse struct {
    ID           string                      `json:"id"`
    Name         string                      `json:"name"`
    Description  string                      `json:"description,omitempty"`
    VaultIDs     []string                    `json:"vaultIds"`
    Permissions  entity.TokenPermissions     `json:"permissions"`
    Restrictions entity.TokenRestrictions    `json:"restrictions"`
    Security     entity.TokenSecurity        `json:"security"`
    Usage        entity.TokenUsage           `json:"usage"`
    IsActive     bool                        `json:"isActive"`
    IsRevoked    bool                        `json:"isRevoked"`
    CreatedAt    time.Time                   `json:"createdAt"`
    UpdatedAt    time.Time                   `json:"updatedAt"`
    ExpiresAt    *time.Time                  `json:"expiresAt,omitempty"`
    LastUsedAt   *time.Time                  `json:"lastUsedAt,omitempty"`
}

type TokenCreateResponse struct {
    Token     TokenResponse                 `json:"token"`
    TokenValue string                       `json:"tokenValue"` // Retornado apenas na criação
    Warning   string                        `json:"warning,omitempty"`
}

type TokenUsageResponse struct {
    TokenID         string                  `json:"tokenId"`
    TotalRequests   int                     `json:"totalRequests"`
    SuccessfulReqs  int                     `json:"successfulReqs"`
    FailedReqs      int                     `json:"failedReqs"`
    RateLimitHits   int                     `json:"rateLimitHits"`
    LastUsedAt      *time.Time              `json:"lastUsedAt,omitempty"`
    LastRequestIP   string                  `json:"lastRequestIP,omitempty"`
    IsActive        bool                    `json:"isActive"`
    IsExpired       bool                    `json:"isExpired"`
}
````markdown
## 🔑 Sistema de Tokens para Automação

### Características

1. **Tokens com Permissões Específicas** - Cada token pode ter permissões granulares (read, write, manage, etc.)
2. **Associação a Vaults** - Tokens são associados a vaults específicos
3. **Restrições de Acesso** - IP whitelist, rate limiting, janelas de tempo
4. **Controle de Expiração** - Tokens podem ter data de expiração
5. **Auditoria Completa** - Todas as operações com tokens são auditadas
6. **Revogação Instantânea** - Tokens podem ser revogados imediatamente

### Tipos de Tokens

```go
// Permissões de Token
type TokenPermissions struct {
    CanRead    bool     // Ler segredos
    CanWrite   bool     // Escrever/editar segredos
    CanDelete  bool     // Deletar segredos
    CanManage  bool     // Gerenciar vault
    CanBackup  bool     // Fazer backup
    CanAudit   bool     // Acessar logs de auditoria
    ReadOnly   bool     // Força modo somente leitura
    Scopes     []string // Operações específicas permitidas
}

// Restrições de Token
type TokenRestrictions struct {
    IPWhitelist   []string        // IPs permitidos
    RateLimit     TokenRateLimit  // Limite de requisições
    TimeWindow    TokenTimeWindow // Janela de tempo permitida
    Environment   []string        // Ambientes permitidos
}
```

### Fluxo de Uso

1. **Criação de Token**:
   - Usuário cria token com permissões específicas
   - Sistema gera token único e hash
   - Relacionamentos são criados no OpenFGA

2. **Validação de Token**:
   - Sistema valida hash do token
   - Verifica se não está revogado/expirado
   - Aplica restrições (IP, rate limit, etc.)

3. **Autorização**:
   - Verifica permissões específicas do token
   - Consulta OpenFGA para autorização
   - Registra uso para auditoria

4. **Revogação**:
   - Token pode ser revogado pelo dono ou admin
   - Relacionamentos são removidos do OpenFGA
   - Ação é auditada

### Exemplos de Uso

#### Token de CI/CD (Somente Leitura)
```go
token := &entity.Token{
    Name: "CI/CD Pipeline",
    Permissions: entity.TokenPermissions{
        CanRead:  true,
        ReadOnly: true,
    },
    Restrictions: entity.TokenRestrictions{
        IPWhitelist: []string{"192.168.1.10", "10.0.0.5"},
        RateLimit: entity.TokenRateLimit{
            Enabled: true,
            RequestsPerHour: 100,
            RequestsPerDay: 1000,
        },
        Environment: []string{"prod", "staging"},
    },
    ExpiresAt: &time.Time{}, // 90 dias
}
```

#### Token de Backup
```go
token := &entity.Token{
    Name: "Backup Service",
    Permissions: entity.TokenPermissions{
        CanRead:   true,
        CanBackup: true,
    },
    Restrictions: entity.TokenRestrictions{
        IPWhitelist: []string{"10.0.0.100"},
        TimeWindow: entity.TokenTimeWindow{
            Enabled: true,
            StartTime: "02:00",
            EndTime: "06:00",
            DaysOfWeek: []string{"sun", "wed"},
        },
    },
}
```

#### Token de Desenvolvimento
```go
token := &entity.Token{
    Name: "Development Team",
    Permissions: entity.TokenPermissions{
        CanRead:  true,
        CanWrite: true,
    },
    Restrictions: entity.TokenRestrictions{
        RateLimit: entity.TokenRateLimit{
            Enabled: true,
            RequestsPerHour: 50,
            BurstLimit: 10,
        },
        Environment: []string{"dev", "staging"},
    },
}
```

### Monitoramento e Auditoria

#### Métricas de Token
- Total de requisições
- Requisições bem-sucedidas vs falhadas
- Rate limit hits
- Último uso
- IPs de origem

#### Logs de Auditoria
- Criação/revogação de tokens
- Uso de tokens
- Violações de restrições
- Tentativas de acesso negadas

### Segurança

#### Boas Práticas
1. **Princípio do Menor Privilégio** - Tokens devem ter apenas as permissões necessárias
2. **Rotação Regular** - Tokens devem ser rotacionados periodicamente
3. **Monitoramento Ativo** - Alertas para uso suspeito
4. **Revogação Imediata** - Capacidade de revogar tokens comprometidos
5. **Auditoria Completa** - Logs de todas as operações

#### Controles de Segurança
- Hash SHA-256 para armazenamento
- Tokens únicos de 32 bytes
- Verificação de IP whitelist
- Rate limiting configurable
- Janelas de tempo restritas
- Controle de ambiente (dev/prod)

### Integração com OpenFGA

#### Relacionamentos de Token
```
token:ci-cd-token-123
├── owner: user:alice
├── vault: vault:prod-secrets
├── can_read_secrets: user:system
├── read_only: user:system
└── production: user:system
```

#### Verificações de Autorização
```bash
# Verificar se token pode ler vault
curl -X POST /stores/{store_id}/check \
  -d '{
    "tuple_key": {
      "user": "token:ci-cd-token-123",
      "relation": "can_read_via_token",
      "object": "vault:prod-secrets"
    }
  }'
```

### API de Tokens

#### Endpoints Principais
- `POST /tokens` - Criar novo token
- `GET /tokens` - Listar tokens do usuário
- `GET /tokens/{id}` - Obter detalhes do token
- `PUT /tokens/{id}` - Atualizar token
- `DELETE /tokens/{id}` - Revogar token
- `POST /tokens/{id}/revoke` - Revogar token
- `GET /tokens/{id}/usage` - Estatísticas de uso

#### Autenticação de Token
```bash
# Usar token para acessar API
curl -H "Authorization: Bearer lockari_token_abc123..." \
  -X GET /api/v1/vaults/vault-123/secrets
```

Este sistema de tokens fornece uma base sólida para automação, mantendo a segurança e controle de acesso granular necessários para um sistema de gerenciamento de segredos.
