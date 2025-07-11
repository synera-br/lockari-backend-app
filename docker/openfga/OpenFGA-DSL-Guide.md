# OpenFGA DSL - Guia Completo e Detalhado

## 📋 Índice
1. [Conceitos Fundamentais](#conceitos-fundamentais)
2. [Estrutura Básica](#estrutura-básica)
3. [Tipos de Relações](#tipos-de-relações)
4. [Operadores Lógicos](#operadores-lógicos)
5. [Grupos e Roles](#grupos-e-roles)
6. [Permissões Customizadas](#permissões-customizadas)
7. [Exemplos Práticos](#exemplos-práticos)
8. [Casos de Uso Avançados](#casos-de-uso-avançados)

---

## 🧠 Conceitos Fundamentais

### O que é OpenFGA?
OpenFGA é um sistema de autorização baseado em **relacionamentos** (Relationship-Based Access Control - ReBAC). Ele funciona com:

- **Tuplas**: Relacionamentos entre usuários e objetos
- **Tipos**: Definições de objetos (user, vault, secret, etc.)
- **Relações**: Como os objetos se relacionam (owner, viewer, etc.)

### Estrutura de uma Tupla
```
user:alice#owner@vault:marketing-secrets
 ↑      ↑       ↑        ↑
user   ID    relação   objeto
```

---

## 🏗️ Estrutura Básica

### Modelo Mínimo
```fga
model
  schema 1.1

type user

type document
  relations
    define owner: [user]
    define viewer: [user]
    define can_read: owner or viewer
```

### Componentes Explicados

#### 1. **Schema Version**
```fga
model
  schema 1.1  # Versão do schema OpenFGA
```

#### 2. **Types** (Tipos)
```fga
type user        # Tipo base para usuários
type group       # Tipo para grupos
type vault       # Tipo para vaults
type secret      # Tipo para secrets
```

#### 3. **Relations** (Relações)
```fga
type vault
  relations
    define owner: [user]           # Relação direta
    define admin: [user]           # Relação direta
    define member: [user, group]   # Aceita usuários E grupos
```

---

## 🔗 Tipos de Relações

### 1. **Relações Diretas**
```fga
define owner: [user]              # Apenas usuários
define admin: [user, group]       # Usuários ou grupos
define viewers: [user:*]          # Qualquer usuário (wildcard)
```

### 2. **Relações Computadas**
```fga
define can_read: owner or admin   # União de relações
define can_write: owner           # Apenas owners
define can_delete: owner          # Apenas owners
```

### 3. **Relações Herdadas**
```fga
type secret
  relations
    define vault: [vault]
    define can_read: can_read from vault  # Herda do vault
```

---

## ⚡ Operadores Lógicos

### 1. **OR (União)**
```fga
define can_read: owner or admin or viewer
# Qualquer um dos três pode ler
```

### 2. **AND (Interseção)**
```fga
define can_edit: writer and approved_user
# Precisa ser writer E approved_user
```

### 3. **FROM (Herança)**
```fga
define can_read: can_read from parent
# Herda permissão do objeto pai
```

### 4. **BUT NOT (Exclusão)**
```fga
define can_read: member but not blocked_user
# Membros podem ler, exceto bloqueados
```

---

## 👥 Grupos e Roles

### Modelo com Grupos
```fga
model
  schema 1.1

type user

type group
  relations
    define owner: [user]
    define admin: [user]
    define member: [user]

type role
  relations
    define assignee: [user, group]

type vault
  relations
    define owner: [user, group]
    define admin: [user, group]
    define writer: [user, group]
    define viewer: [user, group]
    define tenant: [tenant]
    
    # Permissões que incluem grupos
    define can_read: owner or admin or writer or viewer or member from tenant
    define can_write: owner or admin or writer
    define can_delete: owner or admin
    define can_share: owner or admin
```

### Exemplo de Uso com Grupos
```bash
# Criar grupo
user:alice#owner@group:developers

# Adicionar usuários ao grupo
user:bob#member@group:developers
user:charlie#member@group:developers

# Dar permissão ao grupo
group:developers#writer@vault:api-secrets

# Agora bob e charlie podem escrever no vault via grupo
```

---

## 🎯 Permissões Customizadas

### Você pode criar QUALQUER permissão personalizada:

```fga
type vault
  relations
    define owner: [user, group]
    define admin: [user, group]
    define writer: [user, group]
    define viewer: [user, group]
    
    # Permissões customizadas
    define copy: [user, group]
    define read: [user, group]
    define download: [user, group]
    define export: [user, group]
    define backup: [user, group]
    define restore: [user, group]
    define audit: [user, group]
    define monitor: [user, group]
    
    # Permissões computadas
    define can_copy: copy or admin or owner
    define can_read: read or viewer or writer or admin or owner
    define can_download: download or writer or admin or owner
    define can_export: export or admin or owner
    define can_backup: backup or admin or owner
    define can_restore: restore or owner
    define can_audit: audit or admin or owner
    define can_monitor: monitor or admin or owner
    
    # Permissões condicionais
    define can_read_sensitive: (read or viewer) and approved_user
    define can_write_prod: writer and prod_access and approved_user
    define can_delete_critical: owner and two_factor_enabled
```

---

## 📚 Exemplos Práticos

### 1. **Sistema de Arquivos**
```fga
model
  schema 1.1

type user

type group
  relations
    define member: [user]
    define admin: [user]

type folder
  relations
    define owner: [user, group]
    define editor: [user, group]
    define viewer: [user, group]
    define parent: [folder]
    
    define can_read: owner or editor or viewer or can_read from parent
    define can_write: owner or editor or can_write from parent
    define can_delete: owner or can_delete from parent
    define can_create_subfolder: owner or editor

type file
  relations
    define folder: [folder]
    define owner: [user, group]
    define editor: [user, group]
    
    define can_read: owner or editor or can_read from folder
    define can_write: owner or editor or can_write from folder
    define can_delete: owner or can_delete from folder
```

### 2. **Sistema de Blog**
```fga
model
  schema 1.1

type user

type blog
  relations
    define owner: [user]
    define editor: [user]
    define contributor: [user]
    define subscriber: [user]
    
    define can_read: owner or editor or contributor or subscriber
    define can_write: owner or editor
    define can_publish: owner or editor
    define can_moderate: owner

type post
  relations
    define blog: [blog]
    define author: [user]
    define editor: [user]
    
    define can_read: can_read from blog
    define can_write: author or editor or can_write from blog
    define can_publish: can_publish from blog
    define can_comment: can_read from blog
    define can_moderate: can_moderate from blog

type comment
  relations
    define post: [post]
    define author: [user]
    
    define can_read: can_read from post
    define can_write: author
    define can_moderate: can_moderate from post
```

### 3. **Sistema Empresarial com Departamentos**
```fga
model
  schema 1.1

type user

type department
  relations
    define manager: [user]
    define employee: [user]
    define intern: [user]
    
    define can_read_internal: manager or employee
    define can_write_internal: manager or employee
    define can_approve: manager

type project
  relations
    define department: [department]
    define owner: [user]
    define contributor: [user]
    define viewer: [user]
    
    define can_read: owner or contributor or viewer or can_read_internal from department
    define can_write: owner or contributor or can_write_internal from department
    define can_approve: owner or can_approve from department

type document
  relations
    define project: [project]
    define author: [user]
    define reviewer: [user]
    
    define can_read: can_read from project
    define can_write: author or can_write from project
    define can_review: reviewer or can_approve from project
```

---

## 🚀 Casos de Uso Avançados

### 1. **Multi-Tenancy com Isolamento**
```fga
type tenant
  relations
    define owner: [user]
    define admin: [user]
    define member: [user]
    define banned: [user]
    
    define active_member: member but not banned

type vault
  relations
    define tenant: [tenant]
    define owner: [user]
    define admin: [user]
    define writer: [user]
    define viewer: [user]
    
    # Apenas membros ativos do tenant podem acessar
    define can_read: (owner or admin or writer or viewer) and active_member from tenant
    define can_write: (owner or admin or writer) and active_member from tenant
```

### 2. **Sistema de Aprovação**
```fga
type approval_workflow
  relations
    define approver: [user]
    define reviewer: [user]
    define submitter: [user]
    
    define can_submit: submitter
    define can_review: reviewer or approver
    define can_approve: approver
    define can_reject: approver

type document
  relations
    define workflow: [approval_workflow]
    define author: [user]
    define status: [status]
    
    define can_read: author or can_review from workflow
    define can_edit: author and not (approved from status)
    define can_submit: author and can_submit from workflow
```

### 3. **Sistema de Recursos Compartilhados**
```fga
type resource_pool
  relations
    define owner: [user]
    define admin: [user]
    define user: [user]
    
    define can_allocate: owner or admin
    define can_use: user or admin or owner

type resource
  relations
    define pool: [resource_pool]
    define allocated_to: [user]
    define reserved_by: [user]
    
    define can_use: allocated_to or (can_use from pool and not reserved_by)
    define can_reserve: can_use from pool
    define can_release: allocated_to or reserved_by
```

---

## 🔧 Comandos Úteis

### Validar Modelo
```bash
fga model validate --file model.fga
```

### Converter DSL para JSON
```bash
fga model transform --input-format dsl --output-format json --file model.fga
```

### Testar Permissões
```bash
# Verificar se user:alice pode ler vault:secrets
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "can_read",
      "object": "vault:secrets"
    }
  }'
```

### Listar Objetos que Usuário pode Acessar
```bash
curl -X POST http://localhost:8080/stores/$STORE_ID/list-objects \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user:alice",
    "relation": "can_read",
    "type": "vault"
  }'
```

---

## 💡 Dicas Importantes

### 1. **Nomeação de Relações**
- Use nomes descritivos: `can_read`, `can_write`, `can_delete`
- Seja consistente: sempre use o mesmo padrão
- Evite ambiguidade: `editor` vs `can_edit`

### 2. **Performance**
- Evite relações muito complexas
- Use índices adequados no banco
- Monitore queries lentas

### 3. **Debugging**
- Use o playground do OpenFGA
- Teste cada relação individualmente
- Documente casos de uso complexos

### 4. **Grupos vs Usuários**
```fga
# ✅ Bom: Flexibilidade
define viewer: [user, group]

# ❌ Limitado: Apenas usuários
define viewer: [user]
```

### 5. **Herança de Permissões**
```fga
# ✅ Bom: Herança clara
define can_read: can_read from parent_folder

# ❌ Confuso: Lógica duplicada
define can_read: owner or viewer or (owner from parent_folder) or (viewer from parent_folder)
```

---

## 🎯 Resumo Final

1. **Você pode criar QUALQUER permissão customizada** (`copy`, `read`, `download`, etc.)
2. **Grupos são totalmente suportados** - adicione `[user, group]` nas relações
3. **Combine relações com operadores**: `or`, `and`, `but not`, `from`
4. **Herança é poderosa**: use `from` para herdar permissões
5. **Teste sempre**: Use o playground e comandos curl para validar

O OpenFGA é extremamente flexível - você pode modelar praticamente qualquer sistema de permissões!
