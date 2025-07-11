# Lockari Platform - Deployment Guide

## Overview

Este diretório contém os arquivos de deployment para a Lockari Platform em dois ambientes:

- **Kubernetes**: Ambiente de testes e desenvolvimento
- **Google Cloud Run**: Ambiente de produção

## Arquitetura de Deployment

### Ambiente de Testes (Kubernetes)
```
┌─────────────────────────────────────────────────────────────────┐
│                        GKE Cluster                              │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  Frontend   │  │   Backend   │  │   OpenFGA   │            │
│  │   (Next.js) │  │    (Go)     │  │    (FGA)    │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ PostgreSQL  │  │    Redis    │  │  RabbitMQ   │            │
│  │ (Database)  │  │   (Cache)   │  │ (Messages)  │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   Ingress Controller                    │   │
│  │            (nginx + Let's Encrypt SSL)                 │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Ambiente de Produção (Google Cloud Run)
```
┌─────────────────────────────────────────────────────────────────┐
│                     Google Cloud Platform                       │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  Frontend   │  │   Backend   │  │   OpenFGA   │            │
│  │ (Cloud Run) │  │ (Cloud Run) │  │ (Cloud Run) │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ Cloud SQL   │  │ Memorystore │  │   Pub/Sub   │            │
│  │(PostgreSQL) │  │   (Redis)   │  │ (Messages)  │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  Load Balancer                          │   │
│  │              (Global SSL + CDN)                        │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Estrutura de Arquivos

```
deployment/
├── kubernetes/                 # Ambiente de testes
│   ├── 00-namespace.yaml      # Namespace e ConfigMap
│   ├── 01-secrets.yaml        # Secrets do Kubernetes
│   ├── 02-postgres.yaml       # PostgreSQL deployment
│   ├── 03-redis.yaml          # Redis deployment
│   ├── 04-rabbitmq.yaml       # RabbitMQ deployment
│   ├── 05-openfga.yaml        # OpenFGA deployment
│   ├── 06-backend.yaml        # Backend deployment
│   ├── 07-frontend.yaml       # Frontend deployment
│   ├── 08-ingress.yaml        # Ingress e SSL
│   ├── 09-monitoring.yaml     # Monitoring stack
│   └── deploy.sh              # Script de deploy
│
├── cloudrun/                   # Ambiente de produção
│   ├── backend-service.yaml   # Backend Cloud Run service
│   ├── frontend-service.yaml  # Frontend Cloud Run service
│   ├── openfga-service.yaml   # OpenFGA Cloud Run service
│   ├── secrets.yaml           # Secret Manager config
│   ├── load-balancer.yaml     # Load balancer config
│   ├── terraform.tf           # Infraestrutura como código
│   └── deploy.sh              # Script de deploy
│
└── README.md                   # Este arquivo
```

## Pré-requisitos

### Para Kubernetes (Ambiente de Testes)
- `kubectl` instalado e configurado
- `helm` instalado
- `gcloud` CLI instalado
- Cluster GKE configurado
- Certificados SSL configurados (Let's Encrypt)

### Para Cloud Run (Ambiente de Produção)
- `gcloud` CLI instalado e autenticado
- `terraform` instalado
- `docker` instalado
- Projeto Google Cloud configurado
- Permissões de IAM adequadas

## Deployment - Kubernetes (Testes)

### 1. Configuração Inicial

```bash
# Configurar context do Kubernetes
kubectl config use-context gke_lockari-project_us-central1-a_lockari-test-cluster

# Verificar conectividade
kubectl cluster-info
```

### 2. Deploy Completo

```bash
# Executar deploy completo
cd deployment/kubernetes
chmod +x deploy.sh
./deploy.sh deploy
```

### 3. Comandos Úteis

```bash
# Verificar status do deployment
./deploy.sh verify

# Rollback para versão anterior
./deploy.sh rollback

# Limpar todos os recursos
./deploy.sh cleanup

# Port forwarding para acesso local
kubectl port-forward service/frontend-service 3000:3000 -n lockari-test
kubectl port-forward service/backend-service 8080:8080 -n lockari-test
```

### 4. Acesso aos Serviços

- **Frontend**: https://test.lockari.com
- **Backend API**: https://api-test.lockari.com
- **OpenFGA**: https://openfga-test.lockari.com
- **Grafana**: Port forward para 3000

## Deployment - Cloud Run (Produção)

### 1. Configuração Inicial

```bash
# Configurar projeto Google Cloud
gcloud config set project lockari-project
gcloud auth login
gcloud auth application-default login

# Configurar variáveis de ambiente
export POSTGRES_PASSWORD="your-secure-password"
export OPENFGA_POSTGRES_PASSWORD="your-openfga-password"
export JWT_SECRET="your-jwt-secret"
export JWT_REFRESH_SECRET="your-jwt-refresh-secret"
export ENCRYPTION_KEY="your-encryption-key"
export OPENFGA_PRESHARED_KEY="your-openfga-key"
```

### 2. Deploy Completo

```bash
# Executar deploy completo
cd deployment/cloudrun
chmod +x deploy.sh
./deploy.sh deploy
```

### 3. Comandos Úteis

```bash
# Apenas build das imagens
./deploy.sh build

# Apenas migração do banco
./deploy.sh migrate

# Verificar deployment
./deploy.sh verify

# Rollback para versão anterior
./deploy.sh rollback

# Limpar todos os recursos
./deploy.sh cleanup
```

### 4. Acesso aos Serviços

- **Frontend**: https://app.lockari.com
- **Backend API**: https://api.lockari.com
- **OpenFGA**: https://openfga.lockari.com (interno)

## Configuração de Domínios

### Para Kubernetes (Testes)
```bash
# Configurar DNS para apontar para o LoadBalancer IP
test.lockari.com → [EXTERNAL-IP]
api-test.lockari.com → [EXTERNAL-IP]
openfga-test.lockari.com → [EXTERNAL-IP]
```

### Para Cloud Run (Produção)
```bash
# Configurar DNS para apontar para o Load Balancer IP
lockari.com → [GLOBAL-IP]
app.lockari.com → [GLOBAL-IP]
api.lockari.com → [GLOBAL-IP]
openfga.lockari.com → [GLOBAL-IP]
```

## Monitoramento

### Kubernetes
- **Grafana**: Port forward para 3000
- **Prometheus**: Integrado via ServiceMonitor
- **Logs**: `kubectl logs -f deployment/[service] -n lockari-test`

### Cloud Run
- **Cloud Monitoring**: Console GCP
- **Cloud Logging**: Console GCP
- **Cloud Trace**: Console GCP
- **Logs**: `gcloud run logs read --service=[service] --region=us-central1`

## Backup e Recuperação

### Kubernetes
```bash
# Backup dos volumes persistentes
kubectl get pvc -n lockari-test
kubectl get pv

# Backup das configurações
kubectl get all -n lockari-test -o yaml > backup.yaml
```

### Cloud Run
```bash
# Backup do Cloud SQL
gcloud sql backups create --instance=lockari-prod

# Backup do Memorystore
gcloud redis instances backup lockari-cache --region=us-central1
```

## Troubleshooting

### Kubernetes
```bash
# Verificar pods com problemas
kubectl get pods -n lockari-test
kubectl describe pod [pod-name] -n lockari-test
kubectl logs [pod-name] -n lockari-test

# Verificar serviços
kubectl get services -n lockari-test
kubectl describe service [service-name] -n lockari-test

# Verificar ingress
kubectl get ingress -n lockari-test
kubectl describe ingress lockari-ingress -n lockari-test
```

### Cloud Run
```bash
# Verificar serviços
gcloud run services list --region=us-central1
gcloud run services describe [service-name] --region=us-central1

# Verificar logs
gcloud run logs read --service=[service-name] --region=us-central1 --limit=100

# Verificar revisões
gcloud run revisions list --service=[service-name] --region=us-central1
```

## Segurança

### Kubernetes
- **Network Policies**: Isolamento de rede entre pods
- **RBAC**: Controle de acesso granular
- **Secrets**: Encrypted at rest
- **Pod Security Standards**: Restricted security context

### Cloud Run
- **IAM**: Service accounts com least privilege
- **VPC**: Isolamento de rede
- **Secret Manager**: Encrypted secrets
- **Security Headers**: Implementados no load balancer

## Configuração de CI/CD

### GitHub Actions (Kubernetes)
```yaml
name: Deploy to Kubernetes
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Deploy to GKE
      run: |
        cd deployment/kubernetes
        ./deploy.sh deploy
```

### Cloud Build (Cloud Run)
```yaml
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/lockari-backend', '.']
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', 'gcr.io/$PROJECT_ID/lockari-backend']
- name: 'gcr.io/cloud-builders/gcloud'
  args: ['run', 'deploy', 'lockari-backend', '--image', 'gcr.io/$PROJECT_ID/lockari-backend']
```

## Custos Estimados

### Kubernetes (Testes)
- **GKE Cluster**: $70/mês (3 nodes e1-standard-2)
- **Persistent Volumes**: $20/mês (100GB)
- **Load Balancer**: $18/mês
- **Total**: ~$108/mês

### Cloud Run (Produção)
- **Cloud Run**: $150/mês (baseado em uso)
- **Cloud SQL**: $200/mês (instância HA)
- **Memorystore**: $50/mês
- **Load Balancer**: $18/mês
- **Total**: ~$418/mês

## Suporte

Para suporte técnico:
- **Email**: devops@lockari.com
- **Slack**: #devops-support
- **Documentation**: https://docs.lockari.com/deployment
- **Runbooks**: https://runbooks.lockari.com

---

**Lockari Platform** - Deployment Guide v1.0
*Última atualização: Dezembro 2023*
