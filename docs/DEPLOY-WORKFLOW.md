# Documentação do Workflow de Deploy OpenFGA

## Visão Geral

Este documento descreve o workflow automatizado para deploy do OpenFGA no Google Cloud Run, configurado para ambientes de desenvolvimento e produção.

## Estrutura do Workflow

### Gatilhos (Triggers)

1. **Push para branches:**
   - `main` → Deploy para produção
   - `develop` → Deploy para desenvolvimento

2. **Pull Request para main:**
   - Executa validação sem deploy

3. **Dispatch manual:**
   - Permite deploy manual para qualquer ambiente

### Ambientes

#### Desenvolvimento (`develop`)
- **Service Name:** `openfga-dev`
- **Branch:** `develop`
- **Recursos:** 1 CPU, 512Mi RAM
- **Instâncias:** 0-2 (auto-scaling)
- **Playground:** Habilitado
- **Database:** PostgreSQL de desenvolvimento

#### Produção (`main`)
- **Service Name:** `openfga-prod`
- **Branch:** `main`
- **Recursos:** 2 CPU, 1Gi RAM
- **Instâncias:** 1-10 (auto-scaling)
- **Playground:** Desabilitado
- **Database:** PostgreSQL de produção

## Configuração Necessária

### 1. Secrets do GitHub

Configure os seguintes secrets no repositório:

```bash
# Configurações do Google Cloud
GCP_PROJECT_ID=seu-project-id
GCP_SA_KEY=chave-da-service-account-em-json

# URLs dos bancos de dados
DATABASE_URL=postgresql://user:pass@host:port/db
```

### 2. Service Account do Google Cloud

Crie uma service account com as seguintes permissões:

```bash
# Criar service account
gcloud iam service-accounts create openfga-deployer \
  --display-name="OpenFGA Deployer"

# Adicionar roles necessários
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:openfga-deployer@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:openfga-deployer@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:openfga-deployer@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# Gerar chave JSON
gcloud iam service-accounts keys create key.json \
  --iam-account=openfga-deployer@$PROJECT_ID.iam.gserviceaccount.com
```

### 3. Configuração do PostgreSQL

#### Para Desenvolvimento:
```sql
-- Criar database para desenvolvimento
CREATE DATABASE openfga_dev;
CREATE USER openfga_dev WITH ENCRYPTED PASSWORD 'sua_senha_dev';
GRANT ALL PRIVILEGES ON DATABASE openfga_dev TO openfga_dev;
```

#### Para Produção:
```sql
-- Criar database para produção
CREATE DATABASE openfga_prod;
CREATE USER openfga_prod WITH ENCRYPTED PASSWORD 'sua_senha_prod';
GRANT ALL PRIVILEGES ON DATABASE openfga_prod TO openfga_prod;
```

## Processo de Deploy

### 1. Fluxo Automático

```mermaid
graph TD
    A[Push para branch] --> B{Branch?}
    B -->|develop| C[Deploy para DEV]
    B -->|main| D[Deploy para PROD]
    C --> E[Build Image]
    D --> E
    E --> F[Push para GCR]
    F --> G[Deploy Cloud Run]
    G --> H[Health Check]
    H --> I[Update Traffic]
    I --> J[Cleanup Old Revisions]
```

### 2. Etapas Detalhadas

1. **Checkout:** Baixa o código do repositório
2. **Build:** Constrói a imagem Docker do OpenFGA
3. **Push:** Envia a imagem para o Google Container Registry
4. **Deploy:** Atualiza o serviço no Cloud Run
5. **Health Check:** Verifica se o serviço está respondendo
6. **Traffic Update:** Direciona 100% do tráfego para a nova versão
7. **Cleanup:** Remove revisões antigas (mantém apenas 5)

## Configurações dos Serviços

### Variáveis de Ambiente

```bash
# Configurações comuns
OPENFGA_LOG_LEVEL=info
OPENFGA_DATASTORE_ENGINE=postgres
OPENFGA_AUTHN_METHOD=none
OPENFGA_HTTP_ADDR=0.0.0.0:8080
OPENFGA_GRPC_ADDR=0.0.0.0:8081

# Específicas por ambiente
OPENFGA_DATASTORE_URI=postgresql://...
OPENFGA_PLAYGROUND_ENABLED=true|false
```

### Recursos por Ambiente

| Ambiente | CPU | Memória | Min Instâncias | Max Instâncias |
|----------|-----|---------|----------------|----------------|
| DEV      | 1000m | 512Mi | 0 | 2 |
| PROD     | 2000m | 1Gi   | 1 | 10 |

## Monitoramento e Logs

### 1. Health Check

O workflow inclui verificação de saúde:

```bash
curl -f $SERVICE_URL/healthz
```

### 2. Logs do Cloud Run

```bash
# Ver logs do desenvolvimento
gcloud logs read "resource.type=cloud_run_revision resource.labels.service_name=openfga-dev" --limit 50

# Ver logs da produção
gcloud logs read "resource.type=cloud_run_revision resource.labels.service_name=openfga-prod" --limit 50
```

### 3. Métricas

Monitore as seguintes métricas no Cloud Monitoring:

- **Latência:** Tempo de resposta das requisições
- **Throughput:** Número de requisições por minuto
- **Erros:** Taxa de erro 4xx/5xx
- **CPU/Memória:** Utilização de recursos

## Rollback e Recuperação

### 1. Rollback Automático

Se o health check falhar, o deployment será cancelado automaticamente.

### 2. Rollback Manual

```bash
# Listar revisões
gcloud run revisions list --service=openfga-prod --region=us-central1

# Fazer rollback para revisão específica
gcloud run services update-traffic openfga-prod \
  --to-revisions=REVISION_NAME=100 \
  --region=us-central1
```

## Troubleshooting

### Problemas Comuns

1. **Falha no Health Check:**
   - Verificar se o banco está acessível
   - Conferir variáveis de ambiente
   - Analisar logs do serviço

2. **Falha no Build:**
   - Verificar sintaxe do Dockerfile
   - Confirmar que a imagem base está disponível

3. **Falha no Deploy:**
   - Verificar permissões da service account
   - Confirmar que o projeto está correto
   - Verificar quotas do Google Cloud

### Comandos Úteis

```bash
# Verificar status do serviço
gcloud run services describe openfga-dev --region=us-central1

# Testar conectividade
curl -I https://SERVICE_URL/healthz

# Ver logs em tempo real
gcloud logs tail "resource.type=cloud_run_revision resource.labels.service_name=openfga-dev"
```

## Segurança

### 1. Configurações de Segurança

- **Autenticação:** Desabilitada para simplicidade inicial
- **HTTPS:** Habilitado automaticamente pelo Cloud Run
- **Firewall:** Controlado pelas configurações do Cloud Run
- **Secrets:** Gerenciados pelo GitHub Secrets

### 2. Recomendações

1. **Habilitar autenticação** em produção
2. **Usar VPC** para isolamento de rede
3. **Implementar WAF** para proteção adicional
4. **Configurar alertas** para monitoramento proativo

## Custos Estimados

### Desenvolvimento
- **Compute:** ~$5-10/mês (baixo uso)
- **Storage:** ~$1-2/mês
- **Network:** ~$1-2/mês

### Produção
- **Compute:** ~$20-50/mês
- **Storage:** ~$5-10/mês
- **Network:** ~$5-10/mês

## Próximos Passos

1. **Implementar autenticação** via JWT/OAuth2
2. **Configurar VPC** para isolamento de rede
3. **Adicionar alertas** no Cloud Monitoring
4. **Implementar backup** automático do banco
5. **Configurar CDN** para otimização de performance

## Referências

- [OpenFGA Documentation](https://openfga.dev/)
- [Google Cloud Run Documentation](https://cloud.google.com/run/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
