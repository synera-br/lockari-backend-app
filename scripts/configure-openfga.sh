#!/bin/bash

set -e

echo "🔧 Configurando modelo e dados no OpenFGA..."

# Cores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%H:%M:%S')] $1${NC}"
}

# Verificar se OpenFGA está rodando
if ! curl -s http://localhost:8080/healthz > /dev/null; then
    echo -e "${RED}❌ OpenFGA não está rodando. Execute primeiro:${NC}"
    echo "  ./scripts/start-openfga.sh"
    exit 1
fi

log "Criando store..."
STORE_OUTPUT=$(docker-compose exec -T openfga-cli fga store create --name lockari-store --api-url http://openfga:8080 --format json)
STORE_ID=$(echo "$STORE_OUTPUT" | jq -r '.id')

if [ "$STORE_ID" == "null" ] || [ -z "$STORE_ID" ]; then
    echo -e "${RED}❌ Erro ao criar store${NC}"
    echo "$STORE_OUTPUT"
    exit 1
fi

log "Store criado: $STORE_ID"

log "Carregando modelo..."
# Primeiro converte DSL para JSON se necessário
if [ -f docker/openfga/model.fga ]; then
    log "Convertendo DSL para JSON..."
    docker-compose exec -T openfga-cli fga model transform --input-format dsl --output-format json --file model.fga > docker/openfga/model.json.tmp
    if [ -s docker/openfga/model.json.tmp ]; then
        mv docker/openfga/model.json.tmp docker/openfga/model.json
        log "DSL convertida com sucesso!"
    else
        rm -f docker/openfga/model.json.tmp
        log "Usando modelo JSON existente..."
    fi
fi

MODEL_OUTPUT=$(docker-compose exec -T openfga-cli fga model write --api-url http://openfga:8080 --store-id "$STORE_ID" --file model.json --format json)
MODEL_ID=$(echo "$MODEL_OUTPUT" | jq -r '.authorization_model_id')

if [ "$MODEL_ID" == "null" ] || [ -z "$MODEL_ID" ]; then
    echo -e "${RED}❌ Erro ao carregar modelo${NC}"
    echo "$MODEL_OUTPUT"
    exit 1
fi

log "Modelo carregado: $MODEL_ID"

log "Carregando dados de exemplo..."
docker-compose exec -T openfga-cli fga tuple write --api-url http://openfga:8080 --store-id "$STORE_ID" --file init-data.json

log "Atualizando .env.dev..."
if [ -f .env.dev ]; then
    sed -i "s/OPENFGA_STORE_ID=.*/OPENFGA_STORE_ID=$STORE_ID/" .env.dev
    sed -i "s/OPENFGA_AUTH_MODEL_ID=.*/OPENFGA_AUTH_MODEL_ID=$MODEL_ID/" .env.dev
else
    echo "OPENFGA_STORE_ID=$STORE_ID" >> .env.dev
    echo "OPENFGA_AUTH_MODEL_ID=$MODEL_ID" >> .env.dev
fi

echo ""
log "✅ Configuração concluída!"
echo ""
echo -e "${YELLOW}Informações:${NC}"
echo "  Store ID: $STORE_ID"
echo "  Model ID: $MODEL_ID"
echo ""
echo -e "${YELLOW}Teste uma consulta:${NC}"
echo "  curl -X POST http://localhost:8080/stores/$STORE_ID/check \\"
echo "    -H 'Content-Type: application/json' \\"
echo "    -d '{\"tuple_key\":{\"user\":\"user:demo-user-1\",\"relation\":\"read\",\"object\":\"vault:demo-vault-1\"}}'"
echo ""
