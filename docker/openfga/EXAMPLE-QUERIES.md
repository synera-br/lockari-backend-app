# Queries de Exemplo para Teste do Modelo Lockari

## Queries Básicas de Permissões

### 1. Verificar se Alice pode ler segredos no vault de produção
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_read",
    "object": "vault:production-secrets"
  }
}
```

### 2. Verificar se Charlie pode fazer download de segredos
```json
{
  "tuple_key": {
    "user": "user:charlie",
    "relation": "can_download",
    "object": "vault:production-secrets"
  }
}
```

### 3. Verificar se Frank (apenas viewer) pode ler conteúdo
```json
{
  "tuple_key": {
    "user": "user:frank",
    "relation": "can_read",
    "object": "vault:production-secrets"
  }
}
```

### 4. Verificar se Grace (copier) pode copiar segredos
```json
{
  "tuple_key": {
    "user": "user:grace",
    "relation": "can_copy",
    "object": "vault:production-secrets"
  }
}
```

### 5. Verificar se Henry (downloader) pode fazer download
```json
{
  "tuple_key": {
    "user": "user:henry",
    "relation": "can_download",
    "object": "vault:production-secrets"
  }
}
```

## Queries de Permissões Específicas por Tipo

### 6. Verificar se Bob pode renovar certificado SSL
```json
{
  "tuple_key": {
    "user": "user:bob",
    "relation": "can_renew",
    "object": "certificate:ssl-cert"
  }
}
```

### 7. Verificar se Charlie pode fazer deploy em produção com SSH key
```json
{
  "tuple_key": {
    "user": "user:charlie",
    "relation": "can_deploy_prod",
    "object": "ssh_key:deploy-key"
  }
}
```

### 8. Verificar se Alice pode rotacionar API key
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_rotate",
    "object": "api_key:stripe-key"
  }
}
```

### 9. Verificar se Bob pode exportar variáveis de ambiente
```json
{
  "tuple_key": {
    "user": "user:bob",
    "relation": "can_export_env",
    "object": "key_value:env-vars"
  }
}
```

### 10. Verificar se Alice pode conectar no banco de produção
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_connect_prod",
    "object": "database_connection:main-db"
  }
}
```

## Queries de Segredos Sensíveis

### 11. Verificar se Alice pode ler senha sensível do banco
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_read_sensitive",
    "object": "secret:database-password"
  }
}
```

### 12. Verificar se Bob pode copiar senha sensível
```json
{
  "tuple_key": {
    "user": "user:bob",
    "relation": "can_copy_sensitive",
    "object": "secret:database-password"
  }
}
```

## Queries de Compartilhamento Externo

### 13. Verificar se cliente externo pode ver metadados
```json
{
  "tuple_key": {
    "user": "user:external-client",
    "relation": "can_view_external",
    "object": "vault:production-secrets"
  }
}
```

### 14. Verificar se cliente externo pode ler conteúdo
```json
{
  "tuple_key": {
    "user": "user:external-client",
    "relation": "can_read_external",
    "object": "vault:production-secrets"
  }
}
```

### 15. Verificar se Alice pode compartilhar entre tenants
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_share_cross_tenant",
    "object": "vault:production-secrets"
  }
}
```

## Queries de Auditoria e Backup

### 16. Verificar se compliance officer pode ler logs de auditoria
```json
{
  "tuple_key": {
    "user": "user:ciso",
    "relation": "can_read",
    "object": "audit_log:vault-audit"
  }
}
```

### 17. Verificar se Alice pode criar backup
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_create",
    "object": "backup:daily-backup"
  }
}
```

### 18. Verificar se sistema pode fazer backup
```json
{
  "tuple_key": {
    "user": "user:system",
    "relation": "can_create",
    "object": "backup:daily-backup"
  }
}
```

## Queries de Sessão e Notificações

### 19. Verificar se Alice pode usar sua sessão
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_use",
    "object": "session:alice-session"
  }
}
```

### 20. Verificar se Alice pode ler notificação de expiração
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_read",
    "object": "notification:cert-expiry"
  }
}
```

## Queries de Listagem (List Objects)

### 21. Listar todos os vaults que Alice pode gerenciar
```json
{
  "user": "user:alice",
  "relation": "can_manage",
  "type": "vault"
}
```

### 22. Listar todos os segredos que Bob pode ler
```json
{
  "user": "user:bob",
  "relation": "can_read",
  "type": "secret"
}
```

### 23. Listar todos os certificados que Charlie pode renovar
```json
{
  "user": "user:charlie",
  "relation": "can_renew",
  "type": "certificate"
}
```

### 24. Listar todos os backups que o sistema pode criar
```json
{
  "user": "user:system",
  "relation": "can_create",
  "type": "backup"
}
```

## Queries de Usuários (List Users)

### 25. Listar todos os usuários que podem ler um vault específico
```json
{
  "object": "vault:production-secrets",
  "relation": "can_read",
  "user_filters": [
    {
      "type": "user"
    }
  ]
}
```

### 26. Listar todos os usuários que podem gerenciar certificados
```json
{
  "object": "certificate:ssl-cert",
  "relation": "can_manage",
  "user_filters": [
    {
      "type": "user"
    }
  ]
}
```

## Queries de Diferenciação de Planos

### 27. Verificar se Diana (free plan) pode compartilhar externamente
```json
{
  "tuple_key": {
    "user": "user:diana",
    "relation": "can_share_cross_tenant",
    "object": "vault:startup-configs"
  }
}
```

### 28. Verificar se Alice (enterprise) pode compartilhar externamente
```json
{
  "tuple_key": {
    "user": "user:alice",
    "relation": "can_share_cross_tenant",
    "object": "vault:production-secrets"
  }
}
```

## Queries de Grupos

### 29. Verificar se membro do grupo dev-team pode ler
```json
{
  "tuple_key": {
    "user": "user:charlie",
    "relation": "can_read",
    "object": "vault:production-secrets"
  }
}
```

### 30. Verificar se grupo dev-team pode ler vault
```json
{
  "tuple_key": {
    "user": "group:dev-team",
    "relation": "can_read",
    "object": "vault:production-secrets"
  }
}
```

## Resultados Esperados

### Permissões Básicas
- Alice (owner): ✅ Todas as permissões
- Bob (admin): ✅ Todas exceto algumas de owner
- Charlie (writer): ✅ Read, write, copy, mas não download
- Frank (viewer): ✅ Apenas view, não read
- Grace (copier): ✅ Copy mas não download
- Henry (downloader): ✅ Download mas não copy

### Permissões Específicas
- Bob pode renovar certificados (auto_renew)
- Charlie pode fazer deploy prod (owner + production)
- Alice pode rotacionar API keys (owner)
- Bob pode exportar env vars (owner)
- Alice pode conectar prod DB (owner + production)

### Segredos Sensíveis
- Alice pode ler/copiar (owner + sensitive)
- Bob não pode copiar senha sensível (não é owner)

### Compartilhamento Externo
- Cliente externo pode apenas ver metadados
- Alice pode compartilhar entre tenants (enterprise + aprovação)
- Diana não pode compartilhar externamente (free plan)

### Auditoria
- CISO pode ler logs (compliance_officer)
- Sistema pode fazer backup (system)
- Alice pode usar sessão (MFA verificado)

## Como Testar no OpenFGA Playground

1. **Carregue o modelo**: Cole o conteúdo de `lockari-refined-model.fga`
2. **Carregue os dados**: Cole o conteúdo de `refined-playground-data.json`
3. **Execute as queries**: Use as queries acima na seção "Check"
4. **Teste list objects**: Use na seção "List Objects"
5. **Teste list users**: Use na seção "List Users"

## Debugging

Se alguma query não retornar o resultado esperado:
1. Verifique se todos os tuple_keys necessários estão nos dados
2. Confirme se as relações estão definidas corretamente no modelo
3. Trace o caminho de permissões usando as relações intermediárias
4. Verifique se há conflitos entre permissões (como `but not banned`)

## Próximos Passos

1. Validar todas as queries no playground
2. Criar queries para casos edge
3. Testar performance com dados maiores
4. Implementar no backend Go
5. Criar testes automatizados
