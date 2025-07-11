#!/bin/bash

# Script para testar o OpenFGA localmente
# Este script configura e testa o OpenFGA com Docker

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar se o Docker está instalado
if ! command -v docker &> /dev/null; then
    print_error "Docker não encontrado. Por favor, instale o Docker."
    exit 1
fi

# Verificar se o Docker Compose está instalado
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose não encontrado. Por favor, instale o Docker Compose."
    exit 1
fi

# Diretório do projeto
PROJECT_DIR="$(dirname "$(dirname "$(readlink -f "$0")")")"
DOCKER_DIR="$PROJECT_DIR/docker/openfga"

print_status "Iniciando teste local do OpenFGA..."

# Criar docker-compose.yml para teste
cat > "$DOCKER_DIR/docker-compose.yml" << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: openfga
      POSTGRES_USER: openfga
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U openfga"]
      interval: 5s
      timeout: 5s
      retries: 5

  openfga:
    build: .
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      OPENFGA_DATASTORE_ENGINE: postgres
      OPENFGA_DATASTORE_URI: postgresql://openfga:password@postgres:5432/openfga?sslmode=disable
      OPENFGA_HTTP_ADDR: 0.0.0.0:8080
      OPENFGA_GRPC_ADDR: 0.0.0.0:8081
      OPENFGA_LOG_LEVEL: info
      OPENFGA_PLAYGROUND_ENABLED: true
      OPENFGA_AUTHN_METHOD: none
    ports:
      - "8080:8080"
      - "8081:8081"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/healthz"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 30s

volumes:
  postgres_data:
EOF

# Navegar para o diretório do Docker
cd "$DOCKER_DIR"

# Parar containers existentes
print_status "Parando containers existentes..."
docker-compose down -v

# Construir e iniciar os serviços
print_status "Construindo e iniciando serviços..."
docker-compose up -d --build

# Aguardar os serviços ficarem prontos
print_status "Aguardando serviços ficarem prontos..."
sleep 10

# Verificar se os serviços estão rodando
if ! docker-compose ps | grep -q "Up"; then
    print_error "Falha ao iniciar os serviços"
    docker-compose logs
    exit 1
fi

# Teste de conectividade
print_status "Testando conectividade..."
for i in {1..30}; do
    if curl -f http://localhost:8080/healthz &> /dev/null; then
        print_status "OpenFGA está rodando!"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "Timeout aguardando OpenFGA ficar pronto"
        docker-compose logs openfga
        exit 1
    fi
    echo "Aguardando OpenFGA... ($i/30)"
    sleep 2
done

# Executar testes básicos
print_status "Executando testes básicos..."

# Teste 1: Verificar API de stores
print_status "Teste 1: Listando stores..."
STORES_RESPONSE=$(curl -s -X GET http://localhost:8080/stores)
echo "Resposta: $STORES_RESPONSE"

# Teste 2: Criar um store
print_status "Teste 2: Criando store de teste..."
CREATE_STORE_RESPONSE=$(curl -s -X POST http://localhost:8080/stores \
    -H "Content-Type: application/json" \
    -d '{
        "name": "test-store"
    }')
echo "Resposta: $CREATE_STORE_RESPONSE"

# Extrair store ID da resposta
STORE_ID=$(echo $CREATE_STORE_RESPONSE | grep -o '"id":"[^"]*' | grep -o '[^"]*$' || echo "")

if [ -n "$STORE_ID" ]; then
    print_status "Store criado com ID: $STORE_ID"
    
    # Teste 3: Definir modelo de autorização
    print_status "Teste 3: Definindo modelo de autorização..."
    MODEL_RESPONSE=$(curl -s -X POST "http://localhost:8080/stores/$STORE_ID/authorization-models" \
        -H "Content-Type: application/json" \
        -d '{
            "schema_version": "1.1",
            "type_definitions": [
                {
                    "type": "user",
                    "relations": {}
                },
                {
                    "type": "organization",
                    "relations": {
                        "member": {
                            "this": {}
                        }
                    },
                    "metadata": {
                        "relations": {
                            "member": {
                                "directly_related_user_types": [
                                    {
                                        "type": "user"
                                    }
                                ]
                            }
                        }
                    }
                },
                {
                    "type": "vault",
                    "relations": {
                        "owner": {
                            "this": {}
                        },
                        "viewer": {
                            "this": {}
                        }
                    },
                    "metadata": {
                        "relations": {
                            "owner": {
                                "directly_related_user_types": [
                                    {
                                        "type": "user"
                                    }
                                ]
                            },
                            "viewer": {
                                "directly_related_user_types": [
                                    {
                                        "type": "user"
                                    }
                                ]
                            }
                        }
                    }
                }
            ]
        }')
    echo "Resposta do modelo: $MODEL_RESPONSE"
    
    # Teste 4: Verificar playground
    print_status "Teste 4: Verificando playground..."
    if curl -f http://localhost:8080/playground &> /dev/null; then
        print_status "Playground está acessível!"
    else
        print_warning "Playground não está acessível"
    fi
    
else
    print_error "Falha ao criar store"
fi

# Exibir informações finais
echo ""
print_status "=== TESTE CONCLUÍDO ==="
echo "OpenFGA está rodando em: http://localhost:8080"
echo "Playground disponível em: http://localhost:8080/playground"
echo "API gRPC disponível em: localhost:8081"
echo ""
echo "Para parar os serviços, execute:"
echo "  cd $DOCKER_DIR && docker-compose down"
echo ""
echo "Para ver os logs, execute:"
echo "  cd $DOCKER_DIR && docker-compose logs -f"
echo ""
echo "Para acessar o banco de dados:"
echo "  docker-compose exec postgres psql -U openfga -d openfga"
