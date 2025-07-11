# UX Design para Permissões Granulares - Lockari

## Visão Geral

Este documento apresenta sugestões de UX/UI para implementar as permissões granulares no Lockari, focando na usabilidade e clareza para diferentes tipos de usuários.

## Componentes de Interface

### 1. Indicadores de Permissão

#### Badge de Permissão
```
┌─────────────────────┐
│ 🔍 VIEW ONLY       │  - Cinza
│ 👀 READ ACCESS     │  - Azul
│ 📋 COPY ENABLED    │  - Verde
│ 💾 DOWNLOAD OK     │  - Laranja
│ ✏️  WRITE ACCESS    │  - Azul escuro
│ 🗑️  DELETE ACCESS   │  - Vermelho
│ 👥 SHARE ACCESS    │  - Roxo
│ ⚙️  MANAGE ACCESS   │  - Dourado
└─────────────────────┘
```

#### Status de Segredo
```
┌─────────────────────────────────────┐
│ 🔒 SSL Certificate                   │
│ expires-soon.com                     │
│ ──────────────────────────────────── │
│ 🔍 VIEW  👀 READ  📋 COPY  💾 DOWNLOAD │
│ ────────────────────────────────────  │
│ 🟢 Active • 🟡 Expires in 30 days    │
│ 👤 Owner: Alice                      │
└─────────────────────────────────────┘
```

### 2. Botões Condicionais

#### Botões de Ação por Permissão
```
VIEW ONLY:
┌─────────────────────────────────────┐
│ [👁️  View Details]                   │
│ [🚫 Read] [🚫 Copy] [🚫 Download]    │
└─────────────────────────────────────┘

READ ACCESS:
┌─────────────────────────────────────┐
│ [👁️  View Details] [👀 Show Secret] │
│ [🚫 Copy] [🚫 Download]              │
└─────────────────────────────────────┘

COPY ENABLED:
┌─────────────────────────────────────┐
│ [👁️  View] [👀 Show] [📋 Copy]       │
│ [🚫 Download]                        │
└─────────────────────────────────────┘

FULL ACCESS:
┌─────────────────────────────────────┐
│ [👁️  View] [👀 Show] [📋 Copy] [💾 Download] │
│ [✏️  Edit] [🗑️  Delete] [👥 Share]          │
└─────────────────────────────────────┘
```

### 3. Tooltips Explicativos

#### Tooltip para Permissões
```
┌─────────────────────────────────────┐
│ 🔍 VIEW ONLY                         │
│ ─────────────────────────────────── │
│ You can see this secret exists but  │
│ cannot access its content.          │
│                                     │
│ Available actions:                  │
│ • View metadata                     │
│ • See creation date                 │
│ • View tags and labels              │
│                                     │
│ Contact admin for more access.      │
└─────────────────────────────────────┘
```

### 4. Modais de Confirmação

#### Modal de Operação Perigosa
```
┌───────────────────────────────────────┐
│ ⚠️  Delete Production Secret?          │
│ ────────────────────────────────────  │
│ You are about to permanently delete   │
│ the secret "database-password".       │
│                                       │
│ This action cannot be undone and will │
│ affect production systems.            │
│                                       │
│ Type "DELETE" to confirm:             │
│ [________________]                    │
│                                       │
│ [Cancel] [🗑️  Delete Forever]          │
└───────────────────────────────────────┘
```

#### Modal de Compartilhamento Externo
```
┌───────────────────────────────────────┐
│ 👥 Share with External Tenant         │
│ ────────────────────────────────────  │
│ Target: startup-tech                  │
│ Permissions: VIEW + READ              │
│ Duration: 30 days                     │
│                                       │
│ ⚠️  Enterprise Feature                 │
│ This requires approval from:          │
│ • Your tenant admin                   │
│ • Target tenant admin                 │
│                                       │
│ [Cancel] [🚀 Request Approval]        │
└───────────────────────────────────────┘
```

## Layouts por Tipo de Usuário

### 1. Visualizador (VIEW only)
```
┌─────────────────────────────────────────┐
│ 🏢 Production Secrets                   │
│ ───────────────────────────────────────  │
│ 📊 Overview                             │
│ • 12 secrets available                  │
│ • 3 certificates expiring soon          │
│ • 5 API keys active                     │
│                                         │
│ 🔍 Secret List                          │
│ ┌─────────────────────────────────────┐ │
│ │ 🔒 SSL Cert    [View Details]       │ │
│ │ 🔑 API Key     [View Details]       │ │
│ │ 🗄️  Database    [View Details]       │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ 💡 Need more access? Contact admin     │
└─────────────────────────────────────────┘
```

### 2. Desenvolvedor (READ + COPY)
```
┌─────────────────────────────────────────┐
│ 🏢 Production Secrets                   │
│ ───────────────────────────────────────  │
│ 🔍 Quick Actions                        │
│ [🔍 Search] [🏷️  Filter] [⭐ Favorites] │
│                                         │
│ 📋 Recent Copied                        │
│ • API Key (2 min ago)                   │
│ • Database URL (1 hour ago)             │
│                                         │
│ 🔐 Secrets                              │
│ ┌─────────────────────────────────────┐ │
│ │ 🔑 Stripe API Key                   │ │
│ │ [👀 Show] [📋 Copy] [⭐ Favorite]    │ │
│ │ Last used: 2 hours ago               │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### 3. DevOps (FULL ACCESS)
```
┌─────────────────────────────────────────┐
│ 🏢 Production Secrets                   │
│ ───────────────────────────────────────  │
│ 🚀 Quick Actions                        │
│ [➕ Add Secret] [📊 Analytics] [⚙️ Settings] │
│                                         │
│ ⚠️  Alerts                               │
│ • SSL cert expires in 7 days            │
│ • API key rate limit reached            │
│                                         │
│ 📁 Secrets Management                   │
│ ┌─────────────────────────────────────┐ │
│ │ 🔒 SSL Certificate                  │ │
│ │ [👀 Show] [📋 Copy] [💾 Download]    │ │
│ │ [✏️  Edit] [🔄 Renew] [👥 Share]     │ │
│ │ Status: 🟡 Expires soon              │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### 4. Cliente Externo (LIMITED ACCESS)
```
┌─────────────────────────────────────────┐
│ 🌐 Shared Secrets - ACME Corp          │
│ ───────────────────────────────────────  │
│ 🔒 External Access                      │
│ Expires: 25 days remaining              │
│                                         │
│ 📋 Available Secrets                    │
│ ┌─────────────────────────────────────┐ │
│ │ 🔑 API Integration Key              │ │
│ │ [👀 Show] [📋 Copy]                 │ │
│ │ 🚫 Download disabled                 │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ ⚠️  Limited Access                       │
│ This access is temporary and monitored  │
└─────────────────────────────────────────┘
```

## Fluxos de Interação

### 1. Fluxo de Visualização de Segredo

#### Usuário com VIEW apenas
```
1. [Click] Secret Card
   ↓
2. [Show] Modal com metadados
   ↓
3. [Display] "Content hidden - Contact admin for access"
```

#### Usuário com READ
```
1. [Click] "Show Secret" 
   ↓
2. [Verify] MFA if sensitive
   ↓
3. [Display] Content with blur effect
   ↓
4. [Click] "Reveal" button
   ↓
5. [Show] Actual content
```

### 2. Fluxo de Cópia

#### Usuário com COPY
```
1. [Click] "Copy" button
   ↓
2. [Show] Loading spinner
   ↓
3. [Copy] To clipboard
   ↓
4. [Show] "Copied!" feedback
   ↓
5. [Auto-clear] Clipboard after 30s
   ↓
6. [Log] Audit event
```

### 3. Fluxo de Download

#### Usuário com DOWNLOAD
```
1. [Click] "Download" button
   ↓
2. [Show] Format selection modal
   ↓
3. [Select] Format (JSON, ENV, P12, etc.)
   ↓
4. [Generate] Encrypted file
   ↓
5. [Download] File
   ↓
6. [Show] Security warning
   ↓
7. [Log] Audit event
```

### 4. Fluxo de Compartilhamento Externo

#### Usuário Enterprise
```
1. [Click] "Share" button
   ↓
2. [Show] Sharing modal
   ↓
3. [Select] Target tenant
   ↓
4. [Choose] Permissions
   ↓
5. [Set] Duration
   ↓
6. [Submit] Request
   ↓
7. [Wait] Dual approval
   ↓
8. [Notify] All parties
   ↓
9. [Grant] Access
```

## Componentes Reutilizáveis

### 1. PermissionBadge Component
```jsx
<PermissionBadge 
  permission="read"
  granted={true}
  tooltip="You can view secret content"
/>
```

### 2. SecretCard Component
```jsx
<SecretCard
  secret={secret}
  permissions={userPermissions}
  onCopy={handleCopy}
  onDownload={handleDownload}
  onShare={handleShare}
/>
```

### 3. PermissionGate Component
```jsx
<PermissionGate 
  permission="can_download"
  resource={secret}
  fallback={<DisabledButton />}
>
  <DownloadButton />
</PermissionGate>
```

### 4. AuditTrail Component
```jsx
<AuditTrail
  resource={secret}
  actions={['view', 'copy', 'download']}
  user={currentUser}
/>
```

## Estados de Interface

### 1. Estados de Botão
```
Normal:     [📋 Copy]
Loading:    [⏳ Copying...]
Success:    [✅ Copied!]
Error:      [❌ Failed]
Disabled:   [🚫 Copy]
```

### 2. Estados de Permissão
```
Granted:    🟢 Allowed
Denied:     🔴 Denied
Pending:    🟡 Approval needed
Expired:    ⚫ Access expired
```

### 3. Estados de Segredo
```
Active:     🟢 Ready
Expiring:   🟡 Expires soon
Expired:    🔴 Expired
Sensitive:  🔒 MFA required
```

## Feedback Visual

### 1. Animações
- **Fade in**: Ao revelar conteúdo sensível
- **Shake**: Ao tentar ação não permitida
- **Pulse**: Em botões de ação crítica
- **Slide**: Em modais de compartilhamento

### 2. Cores por Tipo de Permissão
```
VIEW:     #6B7280 (Gray)
READ:     #3B82F6 (Blue)
COPY:     #10B981 (Green)
DOWNLOAD: #F59E0B (Orange)
WRITE:    #6366F1 (Indigo)
DELETE:   #EF4444 (Red)
SHARE:    #8B5CF6 (Purple)
MANAGE:   #F59E0B (Amber)
```

### 3. Ícones por Tipo de Segredo
```
🔒 SSL Certificate
🔑 API Key
🗄️  Database Connection
🔐 SSH Key
⚙️  Environment Variables
📝 Generic Secret
💳 Payment Token
🎫 Access Token
```

## Responsividade

### Desktop
- Cards em grid 3x3
- Sidebar com filtros
- Modais centralizados

### Tablet
- Cards em grid 2x2
- Sidebar colapsável
- Modais adaptados

### Mobile
- Cards em lista vertical
- Menu hambúrguer
- Modais full-screen

## Acessibilidade

### 1. Keyboard Navigation
- Tab order lógico
- Escape para fechar modais
- Enter para confirmar ações

### 2. Screen Readers
- ARIA labels descritivos
- Anúncios de mudanças de estado
- Descrições de permissões

### 3. Contraste
- Texto legível em todos os fundos
- Indicadores coloridos com texto
- Modo escuro disponível

## Testes de Usabilidade

### 1. Cenários de Teste
- Usuário encontra segredo sem permissão
- Desenvolvedor copia API key
- Admin compartilha vault
- Cliente externo acessa segredo

### 2. Métricas
- Tempo para encontrar segredo
- Taxa de erro em permissões
- Satisfação com feedback
- Compreensão de limitações

## Implementação Técnica

### 1. Hooks React
```jsx
const usePermissions = (resource) => {
  const { user } = useAuth();
  const [permissions, setPermissions] = useState({});
  
  useEffect(() => {
    checkPermissions(user, resource)
      .then(setPermissions);
  }, [user, resource]);
  
  return permissions;
};
```

### 2. Context para Permissões
```jsx
const PermissionContext = createContext();

const PermissionProvider = ({ children }) => {
  const [userPermissions, setUserPermissions] = useState({});
  
  return (
    <PermissionContext.Provider value={userPermissions}>
      {children}
    </PermissionContext.Provider>
  );
};
```

### 3. Middleware de Autorização
```javascript
const requirePermission = (permission) => (req, res, next) => {
  const { user, resource } = req;
  
  if (checkPermission(user, permission, resource)) {
    next();
  } else {
    res.status(403).json({ error: 'Insufficient permissions' });
  }
};
```

## Próximos Passos

1. **Protótipo**: Criar wireframes interativos
2. **Componentes**: Desenvolver biblioteca de componentes
3. **Testes**: Implementar testes automatizados
4. **Feedback**: Colher feedback de usuários
5. **Iteração**: Refinar baseado no uso real
6. **Documentação**: Criar guia de uso para desenvolvedores
