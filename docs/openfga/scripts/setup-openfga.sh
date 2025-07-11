#!/bin/bash

# Setup OpenFGA para Lockari Backend
# Este script configura o OpenFGA com modelo de autorização e dados iniciais

set -e

echo "🚀 Setting up OpenFGA for Lockari Backend..."

# Definir ambiente
ENV=${1:-dev}
echo "📁 Environment: $ENV"

# Verificar se docker-compose está disponível
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose not found. Please install it first."
    exit 1
fi

# Verificar se curl está disponível
if ! command -v curl &> /dev/null; then
    echo "❌ curl not found. Please install it first."
    exit 1
fi

# Verificar se jq está disponível
if ! command -v jq &> /dev/null; then
    echo "❌ jq not found. Please install it first."
    echo "   Ubuntu/Debian: sudo apt-get install jq"
    echo "   macOS: brew install jq"
    exit 1
fi

# Carregar variáveis de ambiente
if [ -f ".env.$ENV" ]; then
    export $(cat .env.$ENV | grep -v '#' | grep -v '^$' | xargs)
    echo "✅ Loaded .env.$ENV"
else
    echo "❌ .env.$ENV file not found"
    exit 1
fi

# Parar serviços existentes se estiverem rodando
echo "🛑 Stopping existing services..."
docker-compose down -v 2>/dev/null || true

# Limpar volumes antigos se existirem
echo "🧹 Cleaning up old volumes..."
docker volume rm lockari-backend-app_postgres_openfga_data 2>/dev/null || true
docker volume rm lockari-backend-app_redis_data 2>/dev/null || true

# Iniciar serviços
echo "🐳 Starting services..."
docker-compose up -d postgres-openfga

# Aguardar PostgreSQL estar pronto
echo "⏳ Waiting for PostgreSQL to be ready..."
timeout=60
while [ $timeout -gt 0 ]; do
    if docker-compose exec -T postgres-openfga pg_isready -U openfga -d openfga > /dev/null 2>&1; then
        echo "✅ PostgreSQL is ready!"
        break
    fi
    sleep 2
    timeout=$((timeout-2))
done

if [ $timeout -le 0 ]; then
    echo "❌ PostgreSQL failed to start within 60 seconds"
    echo "📋 Logs:"
    docker-compose logs postgres-openfga
    exit 1
fi

# Iniciar OpenFGA
echo "🔧 Starting OpenFGA..."
docker-compose up -d openfga

# Aguardar OpenFGA estar pronto
echo "⏳ Waiting for OpenFGA to be ready..."
timeout=120
while [ $timeout -gt 0 ]; do
    if curl -f -s http://localhost:8080/healthz > /dev/null 2>&1; then
        echo "✅ OpenFGA is ready!"
        break
    fi
    sleep 3
    timeout=$((timeout-3))
done

if [ $timeout -le 0 ]; then
    echo "❌ OpenFGA failed to start within 120 seconds"
    echo "📋 Logs:"
    docker-compose logs openfga
    exit 1
fi

# Iniciar CLI container
echo "🔧 Starting OpenFGA CLI..."
docker-compose up -d openfga-cli

# Criar store
echo "🏪 Creating OpenFGA store..."
STORE_RESPONSE=$(curl -s -X POST http://localhost:8080/stores \
    -H "Content-Type: application/json" \
    -d '{"name": "lockari-vault"}')

if [ $? -ne 0 ]; then
    echo "❌ Failed to create store - curl command failed"
    exit 1
fi

STORE_ID=$(echo $STORE_RESPONSE | jq -r '.id // empty')

if [ -z "$STORE_ID" ] || [ "$STORE_ID" = "null" ]; then
    echo "❌ Failed to create OpenFGA store"
    echo "Response: $STORE_RESPONSE"
    exit 1
fi

echo "✅ Store created with ID: $STORE_ID"

# Atualizar .env com STORE_ID
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/OPENFGA_STORE_ID=.*/OPENFGA_STORE_ID=$STORE_ID/" .env.$ENV
else
    # Linux
    sed -i "s/OPENFGA_STORE_ID=.*/OPENFGA_STORE_ID=$STORE_ID/" .env.$ENV
fi

# Carregar modelo de autorização
echo "📋 Loading authorization model..."
if [ -f "docker/openfga/model.json" ]; then
    MODEL_RESPONSE=$(curl -s -X POST "http://localhost:8080/stores/$STORE_ID/authorization-models" \
        -H "Content-Type: application/json" \
        -d @docker/openfga/model.json)
    
    if [ $? -ne 0 ]; then
        echo "❌ Failed to load model - curl command failed"
        exit 1
    fi
    
    MODEL_ID=$(echo $MODEL_RESPONSE | jq -r '.authorization_model_id // empty')
    
    if [ -z "$MODEL_ID" ] || [ "$MODEL_ID" = "null" ]; then
        echo "❌ Failed to load authorization model"
        echo "Response: $MODEL_RESPONSE"
        exit 1
    fi
    
    echo "✅ Model loaded with ID: $MODEL_ID"
    
    # Atualizar .env com MODEL_ID
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/OPENFGA_MODEL_ID=.*/OPENFGA_MODEL_ID=$MODEL_ID/" .env.$ENV
    else
        # Linux
        sed -i "s/OPENFGA_MODEL_ID=.*/OPENFGA_MODEL_ID=$MODEL_ID/" .env.$ENV
    fi
else
    echo "❌ Authorization model file not found: docker/openfga/model.json"
    exit 1
fi

# Carregar dados iniciais se existirem
if [ -f "docker/openfga/init-data.json" ]; then
    echo "📊 Loading initial demo data..."
    WRITE_RESPONSE=$(curl -s -X POST "http://localhost:8080/stores/$STORE_ID/write" \
        -H "Content-Type: application/json" \
        -d @docker/openfga/init-data.json)
    
    if [ $? -ne 0 ]; then
        echo "❌ Failed to load initial data - curl command failed"
        exit 1
    fi
    
    # Verificar se houve erro na resposta
    if echo $WRITE_RESPONSE | jq -e '.code' > /dev/null; then
        echo "❌ Failed to load initial data:"
        echo $WRITE_RESPONSE | jq '.'
        exit 1
    fi
    
    echo "✅ Initial demo data loaded"
else
    echo "⚠️  No initial data file found (docker/openfga/init-data.json)"
fi

# Iniciar Redis
echo "🔧 Starting Redis..."
docker-compose up -d redis

# Aguardar Redis estar pronto
echo "⏳ Waiting for Redis to be ready..."
timeout=30
while [ $timeout -gt 0 ]; do
    if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
        echo "✅ Redis is ready!"
        break
    fi
    sleep 2
    timeout=$((timeout-2))
done

if [ $timeout -le 0 ]; then
    echo "❌ Redis failed to start within 30 seconds"
    echo "📋 Logs:"
    docker-compose logs redis
    exit 1
fi

echo ""
echo "🎉 OpenFGA setup completed successfully!"
echo ""
echo "📋 Configuration:"
echo "   Store ID: $STORE_ID"
echo "   Model ID: $MODEL_ID"
echo "   API URL:  http://localhost:8080"
echo "   Redis:    redis://localhost:6379"
echo ""
echo "🔧 Services running:"
docker-compose ps
echo ""
echo "🔧 Next steps:"
echo "   1. Test the setup: ./scripts/test-openfga.sh"
echo "   2. View OpenFGA logs: docker-compose logs -f openfga"
echo "   3. View all logs: docker-compose logs -f"
echo "   4. Start your Golang backend with the updated .env.$ENV"
echo ""
echo "🧪 Demo data available:"
echo "   - Tenant: demo-company"
echo "   - Users: demo-alice (owner), demo-bob (member/viewer)"
echo "   - Vault: demo-marketing-secrets"
echo "   - Items: demo-api-key-prod, demo-ssl-cert, demo-deploy-key"
echo ""
echo "💡 Use OpenFGA CLI:"
echo "   docker exec -it lockari-openfga-cli fga model get --store-id=$STORE_ID"
echo "   docker exec -it lockari-openfga-cli fga tuple read --store-id=$STORE_ID"
echo ""
