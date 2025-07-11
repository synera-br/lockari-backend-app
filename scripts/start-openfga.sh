#!/bin/bash

set -e

echo "🚀 Iniciando OpenFGA para Lockari..."

# Cores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%H:%M:%S')] $1${NC}"
}

# Parar containers existentes
log "Parando containers existentes..."
docker-compose down -v

# Iniciar serviços básicos
log "Iniciando PostgreSQL e Redis..."
docker-compose up -d postgres-openfga redis

# Aguardar PostgreSQL
log "Aguardando PostgreSQL..."
sleep 10

# Executar migrações
log "Executando migrações do banco..."
docker-compose up openfga-migrate

# Iniciar OpenFGA
log "Iniciando OpenFGA..."
docker-compose up -d openfga

# Aguardar OpenFGA
log "Aguardando OpenFGA inicializar..."
sleep 15

# Verificar se está funcionando
if curl -s http://localhost:8080/healthz > /dev/null; then
    log "✅ OpenFGA está funcionando!"
else
    error "❌ OpenFGA não está respondendo"
    echo "Verificando logs..."
    docker-compose logs openfga
    exit 1
fi

# Iniciar CLI
log "Iniciando CLI..."
docker-compose up -d openfga-cli

log "✅ Setup concluído!"
echo ""
echo "Serviços disponíveis:"
echo "  • OpenFGA API: http://localhost:8080"
echo "  • OpenFGA gRPC: http://localhost:8081"
echo "  • PostgreSQL: localhost:5433"
echo "  • Redis: localhost:6379"
echo ""
echo "Para configurar o modelo e dados:"
echo "  docker-compose exec openfga-cli bash"
echo "  # Dentro do container:"
echo "  # fga store create --name lockari --api-url http://openfga:8080"
echo ""
