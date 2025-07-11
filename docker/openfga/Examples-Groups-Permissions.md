# Exemplos Práticos - Grupos e Permissões Customizadas

## 🎯 Cenários de Uso

### 1. Configuração de Grupos

```bash
# Criar grupos
echo '{"writes": [
  {"user": "user:alice", "relation": "owner", "object": "group:developers"},
  {"user": "user:bob", "relation": "member", "object": "group:developers"},
  {"user": "user:charlie", "relation": "member", "object": "group:developers"},
  {"user": "user:david", "relation": "admin", "object": "group:devops"},
  {"user": "user:eve", "relation": "member", "object": "group:devops"}
]}' > setup-groups.json

# Aplicar no OpenFGA
fga tuple write --file setup-groups.json
```

### 2. Permissões de Vault com Grupos

```bash
# Configurar vault com grupos
echo '{"writes": [
  {"user": "user:alice", "relation": "owner", "object": "vault:api-secrets"},
  {"user": "group:developers", "relation": "writer", "object": "vault:api-secrets"},
  {"user": "group:devops", "relation": "admin", "object": "vault:api-secrets"},
  {"user": "tenant:company", "relation": "tenant", "object": "vault:api-secrets"}
]}' > vault-permissions.json

fga tuple write --file vault-permissions.json
```

### 3. Permissões Customizadas

```bash
# Dar permissões específicas
echo '{"writes": [
  {"user": "user:alice", "relation": "backup", "object": "vault:api-secrets"},
  {"user": "user:bob", "relation": "copy", "object": "vault:api-secrets"},
  {"user": "group:devops", "relation": "export", "object": "vault:api-secrets"},
  {"user": "user:charlie", "relation": "monitor", "object": "vault:api-secrets"}
]}' > custom-permissions.json

fga tuple write --file custom-permissions.json
```

## 🧪 Testes de Permissões

### Testar Permissões de Grupo

```bash
# Verificar se bob (membro do grupo developers) pode escrever
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:bob",
      "relation": "can_write",
      "object": "vault:api-secrets"
    }
  }'
# Resultado: {"allowed": true} - via grupo developers
```

### Testar Permissões Customizadas

```bash
# Verificar se bob pode copiar
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:bob",
      "relation": "can_copy",
      "object": "vault:api-secrets"
    }
  }'
# Resultado: {"allowed": true} - permissão direta de copy
```

```bash
# Verificar se eve pode fazer backup
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:eve",
      "relation": "can_backup",
      "object": "vault:api-secrets"
    }
  }'
# Resultado: {"allowed": true} - via grupo devops (admin)
```

## 🔍 Cenários Complexos

### 1. Secret Sensível

```bash
# Marcar secret como sensível
echo '{"writes": [
  {"user": "vault:api-secrets", "relation": "vault", "object": "secret:db-password"},
  {"user": "user:alice", "relation": "owner", "object": "secret:db-password"},
  {"user": "user:alice", "relation": "sensitive", "object": "secret:db-password"}
]}' > sensitive-secret.json

fga tuple write --file sensitive-secret.json

# Testar acesso sensível
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:alice",
      "relation": "can_read_sensitive",
      "object": "secret:db-password"
    }
  }'
# Resultado: {"allowed": true} - alice é owner E marcou como sensitive
```

### 2. Certificado Expirando

```bash
# Configurar certificado expirando
echo '{"writes": [
  {"user": "vault:api-secrets", "relation": "vault", "object": "certificate:ssl-cert"},
  {"user": "user:david", "relation": "owner", "object": "certificate:ssl-cert"},
  {"user": "user:charlie", "relation": "expires_soon", "object": "certificate:ssl-cert"}
]}' > expiring-cert.json

fga tuple write --file expiring-cert.json

# Testar alerta de expiração
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:charlie",
      "relation": "can_alert_expiry",
      "object": "certificate:ssl-cert"
    }
  }'
# Resultado: {"allowed": true} - charlie marcou como expires_soon E pode monitorar
```

### 3. Chave SSH de Produção

```bash
# Configurar chave SSH de produção
echo '{"writes": [
  {"user": "vault:api-secrets", "relation": "vault", "object": "ssh_key:prod-deploy"},
  {"user": "user:david", "relation": "owner", "object": "ssh_key:prod-deploy"},
  {"user": "user:david", "relation": "production", "object": "ssh_key:prod-deploy"}
]}' > prod-ssh-key.json

fga tuple write --file prod-ssh-key.json

# Testar deploy em produção
curl -X POST http://localhost:8080/stores/$STORE_ID/check \
  -H 'Content-Type: application/json' \
  -d '{
    "tuple_key": {
      "user": "user:david",
      "relation": "can_deploy",
      "object": "ssh_key:prod-deploy"
    }
  }'
# Resultado: {"allowed": true} - david é owner E marcou como production
```

## 📊 Relatórios e Auditoria

### Listar Todos os Vaults que um Usuário Pode Acessar

```bash
# Vaults que bob pode ler
curl -X POST http://localhost:8080/stores/$STORE_ID/list-objects \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user:bob",
    "relation": "can_read",
    "type": "vault"
  }'
```

### Listar Vaults que um Usuário Pode Fazer Backup

```bash
# Vaults que david pode fazer backup
curl -X POST http://localhost:8080/stores/$STORE_ID/list-objects \
  -H 'Content-Type: application/json' \
  -d '{
    "user": "user:david",
    "relation": "can_backup",
    "type": "vault"
  }'
```

### Listar Usuários que Podem Acessar um Vault

```bash
# Quem pode ler api-secrets
curl -X POST http://localhost:8080/stores/$STORE_ID/list-users \
  -H 'Content-Type: application/json' \
  -d '{
    "object": "vault:api-secrets",
    "relation": "can_read",
    "user_filters": [{"type": "user"}]
  }'
```

## 🚀 Automatização com Scripts

### Script para Configurar Equipe

```bash
#!/bin/bash

STORE_ID="your-store-id"
API_URL="http://localhost:8080"

# Função para criar grupo
create_group() {
    local group_name=$1
    local owner=$2
    
    curl -X POST "$API_URL/stores/$STORE_ID/write" \
        -H 'Content-Type: application/json' \
        -d "{
            \"writes\": [
                {\"user\": \"user:$owner\", \"relation\": \"owner\", \"object\": \"group:$group_name\"}
            ]
        }"
}

# Função para adicionar membro ao grupo
add_to_group() {
    local user=$1
    local group=$2
    local relation=${3:-member}
    
    curl -X POST "$API_URL/stores/$STORE_ID/write" \
        -H 'Content-Type: application/json' \
        -d "{
            \"writes\": [
                {\"user\": \"user:$user\", \"relation\": \"$relation\", \"object\": \"group:$group\"}
            ]
        }"
}

# Função para dar permissão customizada
grant_permission() {
    local user_or_group=$1
    local permission=$2
    local object=$3
    
    curl -X POST "$API_URL/stores/$STORE_ID/write" \
        -H 'Content-Type: application/json' \
        -d "{
            \"writes\": [
                {\"user\": \"$user_or_group\", \"relation\": \"$permission\", \"object\": \"$object\"}
            ]
        }"
}

# Configurar equipe
create_group "developers" "alice"
add_to_group "bob" "developers"
add_to_group "charlie" "developers"

create_group "devops" "david"
add_to_group "eve" "devops"

# Dar permissões
grant_permission "group:developers" "writer" "vault:api-secrets"
grant_permission "group:devops" "admin" "vault:api-secrets"
grant_permission "user:alice" "backup" "vault:api-secrets"
grant_permission "user:david" "restore" "vault:api-secrets"

echo "Equipe configurada com sucesso!"
```

## 💡 Melhores Práticas

### 1. Estrutura de Grupos
```
company
├── departments
│   ├── engineering
│   │   ├── frontend-team
│   │   ├── backend-team
│   │   └── devops-team
│   ├── marketing
│   └── sales
├── roles
│   ├── admin
│   ├── manager
│   └── intern
└── projects
    ├── project-alpha
    └── project-beta
```

### 2. Nomenclatura Consistente
```bash
# Grupos por departamento
group:engineering-frontend
group:engineering-backend
group:engineering-devops

# Grupos por projeto
group:project-alpha-team
group:project-beta-team

# Grupos por função
group:security-team
group:compliance-team
```

### 3. Permissões Granulares
```bash
# ✅ Específicas e claras
can_read_production_secrets
can_deploy_to_staging
can_backup_critical_data

# ❌ Muito genéricas
can_access
can_do_stuff
admin_access
```

### 4. Auditoria e Monitoramento
```bash
# Criar logs de auditoria
grant_permission "user:audit-system" "audit" "vault:api-secrets"
grant_permission "user:monitoring" "monitor" "vault:api-secrets"

# Verificar acessos periódicos
./scripts/audit-permissions.sh
```

## 🎯 Resumo Final

1. **Grupos são poderosos**: Use `[user, group]` para flexibilidade máxima
2. **Permissões customizadas**: Crie quantas precisar (`copy`, `backup`, `monitor`, etc.)
3. **Herança funciona**: `from` para herdar permissões de objetos pai
4. **Combine operadores**: `or`, `and`, `but not` para lógica complexa
5. **Teste sempre**: Valide cada permissão antes de usar em produção
6. **Documente**: Mantenha registro de quem tem acesso ao quê

O OpenFGA permite modelar qualquer sistema de permissões, por mais complexo que seja!
