# OpenFGA Deploy - Guia Completo

## Visão Geral

Este guia documenta a implementação completa do sistema de deploy automatizado do OpenFGA no Google Cloud Run, com ambientes separados para desenvolvimento e produção.

## Estrutura do Projeto

```
lockari-backend-app/
├── docker/
│   └── openfga/
│       ├── Dockerfile
│       └── docker-compose.yml
├── scripts/
│   ├── setup-gcp.sh
│   └── test-local.sh
├── docs/
│   └── DEPLOY-WORKFLOW.md
└── .github/
    └── workflows/
        └── deploy-openfga.yml
```

## Configuração Inicial

### 1. Executar Setup do Google Cloud

```bash
# Tornar o script executável
chmod +x lockari-backend-app/scripts/setup-gcp.sh

# Executar o setup
./lockari-backend-app/scripts/setup-gcp.sh
```

Este script irá:
- Criar service account
- Configurar permissões
- Criar repositório Artifact Registry
- Opcionalmente criar instância Cloud SQL
- Gerar arquivo de credenciais

### 2. Configurar Secrets no GitHub

Navegue até `Settings > Secrets and variables > Actions` no repositório e adicione:

```
GCP_PROJECT_ID: seu-project-id
GCP_SA_KEY: {conteúdo do arquivo JSON da service account}
DATABASE_URL_DEV: postgresql://user:pass@host:port/openfga_dev
DATABASE_URL_PROD: postgresql://user:pass@host:port/openfga_prod
```

### 3. Configurar Environments no GitHub

Vá para `Settings > Environments` e crie:

- **develop**: Para ambiente de desenvolvimento
- **production**: Para ambiente de produção (pode adicionar revisores)

## Teste Local

Para testar o OpenFGA localmente antes do deploy:

```bash
# Tornar o script executável
chmod +x lockari-backend-app/scripts/test-local.sh

# Executar o teste
./lockari-backend-app/scripts/test-local.sh
```

O teste irá:
- Subir PostgreSQL e OpenFGA via Docker Compose
- Executar testes básicos de API
- Validar conectividade
- Disponibilizar playground em http://localhost:8080/playground

## Deploy Automático

### Gatilhos

1. **Push para `develop`**: Deploy automático para ambiente de desenvolvimento
2. **Push para `main`**: Deploy automático para ambiente de produção
3. **Pull Request para `main`**: Validação sem deploy

### Configurações por Ambiente

#### Desenvolvimento
- **URL**: `openfga-dev-[hash].run.app`
- **Recursos**: 1 CPU, 512Mi RAM
- **Instâncias**: 0-2 (auto-scaling)
- **Playground**: Habilitado
- **Logs**: Detalhados

#### Produção
- **URL**: `openfga-prod-[hash].run.app`
- **Recursos**: 2 CPU, 1Gi RAM
- **Instâncias**: 1-10 (auto-scaling)
- **Playground**: Desabilitado
- **Logs**: Otimizados

## Comandos Úteis

### Verificar Status dos Serviços

```bash
# Desenvolvimento
gcloud run services describe openfga-dev --region=us-central1

# Produção
gcloud run services describe openfga-prod --region=us-central1
```

### Ver Logs

```bash
# Desenvolvimento
gcloud logs read "resource.type=cloud_run_revision resource.labels.service_name=openfga-dev" --limit=50

# Produção
gcloud logs read "resource.type=cloud_run_revision resource.labels.service_name=openfga-prod" --limit=50
```

### Fazer Rollback

```bash
# Listar revisões
gcloud run revisions list --service=openfga-prod --region=us-central1

# Rollback para revisão específica
gcloud run services update-traffic openfga-prod \
  --to-revisions=REVISION_NAME=100 \
  --region=us-central1
```

## Monitoramento

### Métricas Importantes

1. **Latência**: Tempo de resposta das requisições
2. **Throughput**: Número de requisições por minuto
3. **Erros**: Taxa de erro 4xx/5xx
4. **Recursos**: Utilização de CPU e memória

### Alertas Recomendados

```yaml
# Exemplo de alerta no Cloud Monitoring
displayName: "OpenFGA High Error Rate"
conditions:
  - displayName: "Error rate > 5%"
    conditionThreshold:
      filter: 'resource.type="cloud_run_revision" AND resource.labels.service_name="openfga-prod"'
      comparison: COMPARISON_GREATER_THAN
      thresholdValue: 0.05
```

## Segurança

### Configurações Aplicadas

1. **HTTPS**: Automaticamente habilitado pelo Cloud Run
2. **Firewall**: Controlado pelas configurações do Cloud Run
3. **Secrets**: Gerenciados pelo GitHub Secrets
4. **Service Account**: Permissões mínimas necessárias

### Recomendações Futuras

1. **Autenticação**: Implementar JWT/OAuth2
2. **VPC**: Configurar isolamento de rede
3. **WAF**: Adicionar proteção contra ataques
4. **Audit Logs**: Habilitar logs de auditoria

## Troubleshooting

### Problemas Comuns

1. **Falha no Health Check**
   - Verificar conectividade com o banco
   - Conferir variáveis de ambiente
   - Analisar logs do serviço

2. **Falha no Build**
   - Verificar sintaxe do Dockerfile
   - Confirmar disponibilidade da imagem base

3. **Falha no Deploy**
   - Verificar permissões da service account
   - Confirmar configurações do projeto
   - Verificar quotas do Google Cloud

### Logs de Debug

```bash
# Logs detalhados do workflow
# Disponíveis na aba "Actions" do GitHub

# Logs do Cloud Run
gcloud logs tail "resource.type=cloud_run_revision resource.labels.service_name=openfga-dev"
```

## Custos Estimados

### Desenvolvimento
- **Compute**: $5-10/mês
- **Storage**: $1-2/mês
- **Network**: $1-2/mês
- **Total**: ~$10-15/mês

### Produção
- **Compute**: $20-50/mês
- **Storage**: $5-10/mês
- **Network**: $5-10/mês
- **Total**: ~$30-70/mês

## Roadmap

### Fase 1 (Implementado)
- [x] Workflow de deploy automatizado
- [x] Ambientes separados (dev/prod)
- [x] Health checks
- [x] Cleanup automático de revisões
- [x] Documentação completa

### Fase 2 (Próximos passos)
- [ ] Implementar autenticação
- [ ] Configurar VPC
- [ ] Adicionar alertas
- [ ] Implementar backup automático
- [ ] Configurar CDN

### Fase 3 (Futuro)
- [ ] Multi-region deployment
- [ ] Blue-green deployment
- [ ] Canary releases
- [ ] Disaster recovery

## Suporte

Para questões técnicas ou problemas:

1. **Logs**: Verificar logs do Cloud Run e GitHub Actions
2. **Documentação**: Consultar `docs/DEPLOY-WORKFLOW.md`
3. **Testes**: Executar `scripts/test-local.sh` para debug local
4. **Monitoramento**: Verificar métricas no Cloud Monitoring

## Referências

- [OpenFGA Documentation](https://openfga.dev/)
- [Google Cloud Run Documentation](https://cloud.google.com/run/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Documentation](https://docs.docker.com/)

---

**Nota**: Este sistema foi projetado para ser robusto, escalável e fácil de manter. Siga as melhores práticas de segurança e monitore regularmente o desempenho dos serviços.
