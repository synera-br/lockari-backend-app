# Modelo de Autorização Lockari - Permissões Granulares

## Visão Geral

Este documento detalha o modelo de autorização refinado do Lockari, explicando as diferenças entre as permissões granulares e como elas se aplicam em cenários reais.

## Hierarquia de Permissões

### 1. VIEW (Visualizar Metadados)
- **O que permite**: Ver informações básicas dos segredos
  - Nome do segredo
  - Tipo (certificado, SSH key, API key, etc.)
  - Tags e labels
  - Data de criação e última modificação
  - Status (ativo, expirado, etc.)
- **O que NÃO permite**: Ver o conteúdo real dos segredos
- **Casos de uso**: 
  - Usuários que precisam saber quais segredos existem
  - Auditoria de inventário
  - Navegação e busca

### 2. READ (Ler Conteúdo)
- **O que permite**: 
  - Tudo do VIEW +
  - Ver o conteúdo completo dos segredos
  - Acessar valores sensíveis
- **Casos de uso**:
  - Desenvolvedores que precisam usar os segredos
  - Operações de CI/CD
  - Troubleshooting e debugging

### 3. COPY (Copiar para Clipboard)
- **O que permite**:
  - Copiar segredos para área de transferência
  - Uso temporário em outras aplicações
  - Não cria arquivos persistentes
- **Segurança**: 
  - Clipboard é limpo automaticamente após X segundos
  - Sem rastro em arquivos do sistema
- **Casos de uso**:
  - Login em sistemas
  - Preenchimento de formulários
  - Uso pontual de credenciais

### 4. DOWNLOAD (Baixar/Exportar)
- **O que permite**:
  - Exportar segredos em arquivos
  - Diferentes formatos (JSON, ENV, YAML, P12, PEM)
  - Criação de arquivos persistentes
- **Segurança**:
  - Mais restritivo que COPY
  - Cria trilha de auditoria
  - Arquivos podem ser criptografados
- **Casos de uso**:
  - Backup local
  - Migração entre sistemas
  - Integração com ferramentas externas

### 5. WRITE (Escrever/Editar)
- **O que permite**:
  - Criar novos segredos
  - Editar segredos existentes
  - Atualizar metadados
  - Rotacionar chaves
- **Casos de uso**:
  - Administração de segredos
  - Automação de rotação
  - Atualizações de configuração

### 6. DELETE (Deletar)
- **O que permite**:
  - Remover segredos permanentemente
  - Revogar certificados/chaves
- **Segurança**:
  - Operação irreversível
  - Requer confirmação adicional
  - Log de auditoria obrigatório

### 7. SHARE (Compartilhar)
- **O que permite**:
  - Conceder acesso a outros usuários
  - Criar links de compartilhamento
  - Configurar permissões granulares

### 8. MANAGE (Gerenciar)
- **O que permite**:
  - Configurar vault
  - Gerenciar permissões
  - Configurar políticas de backup
  - Definir regras de auditoria

## Diferenças Práticas: COPY vs DOWNLOAD

### COPY (Clipboard)
```
✅ Uso temporário
✅ Sem arquivos criados
✅ Auto-limpeza do clipboard
✅ Menor risco de vazamento
❌ Não persistente
❌ Limitado a texto simples
```

### DOWNLOAD (Arquivo)
```
✅ Backup permanente
✅ Múltiplos formatos
✅ Integração com ferramentas
✅ Transferência entre sistemas
❌ Arquivos podem ser esquecidos
❌ Risco de vazamento maior
❌ Requer gerenciamento manual
```

## Cenários de Uso

### Desenvolvedor Frontend
```
Permissões: VIEW + READ + COPY
- Pode ver quais APIs estão disponíveis
- Pode ler as chaves de API
- Pode copiar para usar em testes
- NÃO pode baixar arquivos (proteção)
```

### DevOps Engineer
```
Permissões: VIEW + READ + COPY + DOWNLOAD + WRITE
- Acesso completo para automação
- Pode exportar para CI/CD
- Pode rotacionar chaves
- Pode criar novos segredos
```

### Auditor/Compliance
```
Permissões: VIEW apenas
- Pode ver inventário de segredos
- Pode verificar conformidade
- NÃO pode acessar conteúdo sensível
- Acesso aos logs de auditoria
```

### Cliente Externo (Enterprise)
```
Permissões: VIEW + READ limitado
- Acesso apenas aos segredos compartilhados
- Pode usar mas não exportar
- Tempo de acesso limitado
- Aprovação dupla necessária
```

## Compartilhamento Entre Tenants

### Plano Free
- ❌ Sem compartilhamento externo
- ✅ Compartilhamento interno ilimitado
- ✅ Grupos internos

### Plano Enterprise
- ✅ Compartilhamento externo com aprovação dupla
- ✅ Controle granular de permissões
- ✅ Auditoria completa
- ✅ TTL nos acessos externos

### Fluxo de Aprovação Dupla
1. **Solicitação**: Usuário solicita compartilhamento
2. **Primeira Aprovação**: Admin do tenant de origem
3. **Segunda Aprovação**: Admin do tenant de destino
4. **Configuração**: Definição de permissões específicas
5. **TTL**: Definição de tempo de acesso
6. **Auditoria**: Log completo do processo

## Tipos de Segredos e Permissões Específicas

### Certificados Digitais
```
Permissões específicas:
- can_renew: Renovar certificado
- can_revoke: Revogar certificado
- can_export_p12: Exportar em formato P12
- can_export_pem: Exportar em formato PEM
- can_install: Instalar em sistema
- can_validate: Validar certificado
```

### Chaves SSH
```
Permissões específicas:
- can_use: Usar para conexão
- can_deploy_prod: Deploy em produção
- can_add_to_agent: Adicionar ao SSH agent
- can_generate_public: Gerar chave pública
- can_rotate: Rotacionar chave
```

### API Keys
```
Permissões específicas:
- can_rotate: Rotacionar chave
- can_revoke: Revogar chave
- can_test: Testar conectividade
```

### Database Connections
```
Permissões específicas:
- can_test_connection: Testar conexão
- can_connect_prod: Conectar em produção
- can_export_schema: Exportar schema
```

## Segurança e Auditoria

### Logs de Auditoria
- Todas as operações são logadas
- Rastreamento de acesso por usuário
- Detecção de atividades suspeitas
- Relatórios de compliance

### Notificações
- Acesso a segredos sensíveis
- Certificados próximos ao vencimento
- Tentativas de acesso suspeitas
- Compartilhamentos externos

### Sessões
- Controle de dispositivos
- Verificação de localização
- MFA obrigatório para operações sensíveis
- Detecção de atividade suspeita

## Implementação Recomendada

### Frontend (React/Next.js)
- Botões condicionais baseados em permissões
- UI diferente para cada nível de acesso
- Confirmações para operações perigosas
- Indicadores visuais de segurança

### Backend (Go)
- Middleware de autorização
- Cache de permissões
- Rate limiting por operação
- Validação dupla para operações críticas

### Integração OpenFGA
- Cache local de permissões
- Batch queries para performance
- Fallback para modo offline
- Sync automático de mudanças

## Próximos Passos

1. **Validar modelo no OpenFGA Playground**
2. **Criar queries de exemplo**
3. **Implementar no backend Go**
4. **Criar componentes React**
5. **Testes de integração**
6. **Documentação de API**
