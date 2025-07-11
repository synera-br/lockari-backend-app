# Firebase Collections Structure - Lockari

## Visão Geral da Arquitetura

Este documento detalha a estrutura das collections do Firebase para o sistema Lockari, considerando o modelo de autorização multi-tenant com OpenFGA.

## 🏗️ Estrutura Principal das Collections

### 1. **users** (Collection)
Usuários do sistema com informações básicas de perfil.

```
users/
├── {userId}/
│   ├── email: string
│   ├── displayName: string
│   ├── photoURL: string
│   ├── providerId: string (google, email, etc.)
│   ├── emailVerified: boolean
│   ├── createdAt: timestamp
│   ├── updatedAt: timestamp
│   ├── lastLoginAt: timestamp
│   ├── isActive: boolean
│   ├── preferences: object
│   │   ├── theme: string
│   │   ├── language: string
│   │   └── notifications: object
│   └── metadata: object
│       ├── ipAddress: string
│       ├── userAgent: string
│       └── location: string
```

### 2. **tenants** (Collection)
Organizações/empresas que usam o sistema.

```
tenants/
├── {tenantId}/
│   ├── name: string
│   ├── slug: string (unique)
│   ├── description: string
│   ├── domain: string
│   ├── plan: string (free, enterprise)
│   ├── planExpiry: timestamp
│   ├── maxUsers: number
│   ├── maxVaults: number
│   ├── maxSecrets: number
│   ├── features: array<string>
│   ├── billing: object
│   │   ├── customerId: string
│   │   ├── subscriptionId: string
│   │   └── status: string
│   ├── settings: object
│   │   ├── allowExternalSharing: boolean
│   │   ├── requireMFA: boolean
│   │   ├── sessionTimeout: number
│   │   └── auditRetention: number
│   ├── createdAt: timestamp
│   ├── updatedAt: timestamp
│   ├── createdBy: string (userId)
│   └── isActive: boolean
```

### 3. **tenant_members** (Collection)
Relacionamento entre usuários e tenants com roles.

```
tenant_members/
├── {tenantId}_{userId}/
│   ├── tenantId: string
│   ├── userId: string
│   ├── role: string (owner, admin, member)
│   ├── permissions: array<string>
│   ├── joinedAt: timestamp
│   ├── invitedBy: string (userId)
│   ├── status: string (active, invited, suspended)
│   ├── lastAccessAt: timestamp
│   └── metadata: object
```

### 4. **groups** (Collection)
Grupos de usuários dentro de tenants.

```
groups/
├── {groupId}/
│   ├── name: string
│   ├── description: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── members: array<string> (userIds)
│   ├── admins: array<string> (userIds)
│   ├── permissions: array<string>
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 5. **vaults** (Collection)
Cofres que contêm os segredos.

```
vaults/
├── {vaultId}/
│   ├── name: string
│   ├── description: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── icon: string
│   ├── color: string
│   ├── tags: array<string>
│   ├── secretsCount: number
│   ├── permissions: object
│   │   ├── owners: array<string>
│   │   ├── admins: array<string>
│   │   ├── writers: array<string>
│   │   ├── readers: array<string>
│   │   ├── viewers: array<string>
│   │   ├── copiers: array<string>
│   │   └── downloaders: array<string>
│   ├── sharing: object
│   │   ├── isPublic: boolean
│   │   ├── allowExternalSharing: boolean
│   │   └── externalShares: array<object>
│   ├── backup: object
│   │   ├── enabled: boolean
│   │   ├── frequency: string
│   │   └── lastBackupAt: timestamp
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 6. **secrets** (Collection)
Segredos genéricos (senhas, tokens, etc.).

```
secrets/
├── {secretId}/
│   ├── name: string
│   ├── description: string
│   ├── type: string (password, token, api_key, etc.)
│   ├── vaultId: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── encryptedValue: string
│   ├── metadata: object
│   │   ├── url: string
│   │   ├── username: string
│   │   ├── notes: string
│   │   └── customFields: object
│   ├── tags: array<string>
│   ├── isSensitive: boolean
│   ├── isProduction: boolean
│   ├── expiresAt: timestamp
│   ├── lastUsedAt: timestamp
│   ├── usageCount: number
│   ├── version: number
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 7. **certificates** (Collection)
Certificados SSL/TLS.

```
certificates/
├── {certificateId}/
│   ├── name: string
│   ├── description: string
│   ├── vaultId: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── certificateData: object
│   │   ├── encryptedCert: string
│   │   ├── encryptedKey: string
│   │   ├── encryptedChain: string
│   │   └── passphrase: string (encrypted)
│   ├── details: object
│   │   ├── commonName: string
│   │   ├── subjectAltNames: array<string>
│   │   ├── issuer: string
│   │   ├── serialNumber: string
│   │   ├── fingerprint: string
│   │   └── algorithm: string
│   ├── validity: object
│   │   ├── notBefore: timestamp
│   │   ├── notAfter: timestamp
│   │   └── isExpired: boolean
│   ├── renewal: object
│   │   ├── autoRenew: boolean
│   │   ├── renewalDays: number
│   │   └── provider: string
│   ├── tags: array<string>
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 8. **ssh_keys** (Collection)
Chaves SSH para acesso e deploy.

```
ssh_keys/
├── {sshKeyId}/
│   ├── name: string
│   ├── description: string
│   ├── vaultId: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── keyData: object
│   │   ├── encryptedPrivateKey: string
│   │   ├── publicKey: string
│   │   ├── passphrase: string (encrypted)
│   │   └── keyType: string (rsa, ed25519, etc.)
│   ├── usage: object
│   │   ├── isProduction: boolean
│   │   ├── authorizedHosts: array<string>
│   │   ├── allowedUsers: array<string>
│   │   └── restrictions: object
│   ├── metadata: object
│   │   ├── keySize: number
│   │   ├── fingerprint: string
│   │   └── comment: string
│   ├── tags: array<string>
│   ├── lastUsedAt: timestamp
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 9. **key_values** (Collection)
Variáveis de ambiente e configurações.

```
key_values/
├── {keyValueId}/
│   ├── name: string
│   ├── description: string
│   ├── vaultId: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── data: object
│   │   ├── encryptedValues: object
│   │   ├── environment: string (dev, staging, prod)
│   │   └── format: string (json, env, yaml)
│   ├── metadata: object
│   │   ├── totalKeys: number
│   │   ├── lastSyncAt: timestamp
│   │   └── syncSource: string
│   ├── tags: array<string>
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 10. **database_connections** (Collection)
Conexões de banco de dados.

```
database_connections/
├── {connectionId}/
│   ├── name: string
│   ├── description: string
│   ├── vaultId: string
│   ├── tenantId: string
│   ├── createdBy: string (userId)
│   ├── connection: object
│   │   ├── type: string (mysql, postgres, mongo, etc.)
│   │   ├── host: string
│   │   ├── port: number
│   │   ├── database: string
│   │   ├── encryptedUsername: string
│   │   ├── encryptedPassword: string
│   │   └── encryptedConnectionString: string
│   ├── ssl: object
│   │   ├── enabled: boolean
│   │   ├── certificate: string
│   │   └── verifyCA: boolean
│   ├── pools: object
│   │   ├── maxConnections: number
│   │   ├── minConnections: number
│   │   └── timeout: number
│   ├── isProduction: boolean
│   ├── isReadOnly: boolean
│   ├── tags: array<string>
│   ├── lastTestedAt: timestamp
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 11. **audit_logs** (Collection)
Logs de auditoria de todas as ações.

```
audit_logs/
├── {logId}/
│   ├── tenantId: string
│   ├── userId: string
│   ├── action: string
│   ├── resourceType: string
│   ├── resourceId: string
│   ├── resourceName: string
│   ├── details: object
│   │   ├── oldValues: object
│   │   ├── newValues: object
│   │   ├── ipAddress: string
│   │   ├── userAgent: string
│   │   └── location: string
│   ├── result: string (success, failure, denied)
│   ├── risk: string (low, medium, high)
│   ├── session: object
│   │   ├── sessionId: string
│   │   ├── deviceId: string
│   │   └── mfaVerified: boolean
│   ├── timestamp: timestamp
│   └── expiresAt: timestamp
```

### 12. **activity_logs** (Collection)
Logs de atividades dos usuários.

```
activity_logs/
├── {activityId}/
│   ├── tenantId: string
│   ├── userId: string
│   ├── action: string
│   ├── description: string
│   ├── metadata: object
│   │   ├── ipAddress: string
│   │   ├── userAgent: string
│   │   ├── location: string
│   │   └── duration: number
│   ├── timestamp: timestamp
│   └── expiresAt: timestamp
```

### 13. **sessions** (Collection)
Sessões ativas dos usuários.

```
sessions/
├── {sessionId}/
│   ├── userId: string
│   ├── tenantId: string
│   ├── deviceId: string
│   ├── deviceInfo: object
│   │   ├── type: string
│   │   ├── os: string
│   │   ├── browser: string
│   │   └── version: string
│   ├── location: object
│   │   ├── country: string
│   │   ├── region: string
│   │   ├── city: string
│   │   └── ipAddress: string
│   ├── security: object
│   │   ├── mfaVerified: boolean
│   │   ├── mfaAt: timestamp
│   │   ├── riskScore: number
│   │   └── isSuspicious: boolean
│   ├── isActive: boolean
│   ├── createdAt: timestamp
│   ├── lastAccessAt: timestamp
│   └── expiresAt: timestamp
```

### 14. **external_share_requests** (Collection)
Solicitações de compartilhamento externo.

```
external_share_requests/
├── {requestId}/
│   ├── fromTenantId: string
│   ├── toTenantId: string
│   ├── requesterId: string (userId)
│   ├── vaultId: string
│   ├── permissions: array<string>
│   ├── duration: number (days)
│   ├── message: string
│   ├── approvals: array<object>
│   │   ├── userId: string
│   │   ├── tenantId: string
│   │   ├── status: string
│   │   ├── approvedAt: timestamp
│   │   └── comment: string
│   ├── status: string (pending, approved, rejected, expired)
│   ├── expiresAt: timestamp
│   ├── createdAt: timestamp
│   └── updatedAt: timestamp
```

### 15. **notifications** (Collection)
Notificações do sistema.

```
notifications/
├── {notificationId}/
│   ├── tenantId: string
│   ├── userId: string
│   ├── type: string
│   ├── title: string
│   ├── message: string
│   ├── data: object
│   ├── channels: array<string>
│   ├── priority: string
│   ├── isRead: boolean
│   ├── readAt: timestamp
│   ├── createdAt: timestamp
│   └── expiresAt: timestamp
```

### 16. **backups** (Collection)
Backups dos vaults.

```
backups/
├── {backupId}/
│   ├── tenantId: string
│   ├── vaultId: string
│   ├── createdBy: string (userId)
│   ├── type: string (manual, scheduled)
│   ├── status: string (in_progress, completed, failed)
│   ├── metadata: object
│   │   ├── totalSecrets: number
│   │   ├── totalSize: number
│   │   ├── encryptionKey: string
│   │   └── checksum: string
│   ├── storage: object
│   │   ├── provider: string
│   │   ├── location: string
│   │   └── expiresAt: timestamp
│   ├── createdAt: timestamp
│   └── completedAt: timestamp
```

## 🔄 Índices Recomendados

### Índices Compostos Essenciais
```javascript
// Firestore Indexes
tenants: [['isActive', 'plan'], ['createdAt', 'desc']]
tenant_members: [['tenantId', 'status'], ['userId', 'role']]
vaults: [['tenantId', 'isActive'], ['createdBy', 'updatedAt']]
secrets: [['vaultId', 'isActive'], ['tenantId', 'type']]
audit_logs: [['tenantId', 'timestamp'], ['userId', 'action']]
sessions: [['userId', 'isActive'], ['tenantId', 'lastAccessAt']]
```

## 🔐 Regras de Segurança

### Exemplo de Security Rules
```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Usuários só podem acessar seus próprios dados
    match /users/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    // Tenants - apenas membros podem acessar
    match /tenants/{tenantId} {
      allow read, write: if request.auth != null && 
        exists(/databases/$(database)/documents/tenant_members/$(tenantId + '_' + request.auth.uid));
    }
    
    // Vaults - verificar permissões via OpenFGA
    match /vaults/{vaultId} {
      allow read, write: if request.auth != null && 
        resource.data.tenantId in getUserTenants(request.auth.uid);
    }
  }
}
```

## 📊 Considerações de Performance

### 1. **Particionamento**
- Particionar por `tenantId` para melhor distribuição
- Usar TTL para logs e sessões
- Implementar cleanup automático

### 2. **Caching**
- Cache frequente de tenant_members
- Cache de permissões por usuário
- Cache de sessões ativas

### 3. **Batch Operations**
- Usar batch writes para operações relacionadas
- Implementar retry logic para operações críticas
- Monitorar limites de quota

## 🔍 Queries Comuns

### Buscar vaults do usuário
```javascript
vaults.where('tenantId', '==', tenantId)
      .where('isActive', '==', true)
      .orderBy('updatedAt', 'desc')
```

### Logs de auditoria do tenant
```javascript
audit_logs.where('tenantId', '==', tenantId)
          .where('timestamp', '>=', startDate)
          .orderBy('timestamp', 'desc')
          .limit(100)
```

### Sessões ativas do usuário
```javascript
sessions.where('userId', '==', userId)
        .where('isActive', '==', true)
        .orderBy('lastAccessAt', 'desc')
```

## 🎯 Próximos Passos

1. **Implementar Security Rules detalhadas**
2. **Criar índices otimizados**
3. **Configurar TTL para logs**
4. **Implementar cleanup automático**
5. **Monitorar performance e custos**
6. **Implementar backup strategy**

Esta estrutura fornece uma base sólida para o sistema Lockari, com foco em segurança, performance e escalabilidade.
