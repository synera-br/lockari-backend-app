#!/bin/bash

# Script para configurar o ambiente de deploy do OpenFGA no Google Cloud
# Este script deve ser executado uma única vez para configurar o projeto

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para print colorido
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar se o gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    print_error "Google Cloud SDK não encontrado. Por favor, instale: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Configurações
PROJECT_ID=""
REGION="us-central1"
SERVICE_ACCOUNT_NAME="openfga-deployer"
REPOSITORY_NAME="openfga-images"

# Função para solicitar input do usuário
get_input() {
    local prompt="$1"
    local var_name="$2"
    local default="$3"
    
    if [ -n "$default" ]; then
        read -p "$prompt [$default]: " input
        eval $var_name="\${input:-$default}"
    else
        read -p "$prompt: " input
        eval $var_name="$input"
    fi
}

# Obter informações do projeto
print_status "Configurando projeto Google Cloud..."
get_input "Digite o ID do projeto Google Cloud" PROJECT_ID
get_input "Digite a região" REGION "us-central1"

# Verificar se o projeto existe
if ! gcloud projects describe $PROJECT_ID &> /dev/null; then
    print_error "Projeto $PROJECT_ID não encontrado ou sem acesso"
    exit 1
fi

# Definir projeto padrão
gcloud config set project $PROJECT_ID
print_status "Projeto definido como: $PROJECT_ID"

# Habilitar APIs necessárias
print_status "Habilitando APIs necessárias..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable cloudsql.googleapis.com
gcloud services enable compute.googleapis.com

# Criar Artifact Registry repository
print_status "Criando repositório Artifact Registry..."
if ! gcloud artifacts repositories describe $REPOSITORY_NAME --location=$REGION &> /dev/null; then
    gcloud artifacts repositories create $REPOSITORY_NAME \
        --repository-format=docker \
        --location=$REGION \
        --description="Repositório para imagens do OpenFGA"
    print_status "Repositório criado: $REPOSITORY_NAME"
else
    print_warning "Repositório $REPOSITORY_NAME já existe"
fi

# Criar service account
print_status "Criando service account..."
if ! gcloud iam service-accounts describe $SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com &> /dev/null; then
    gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME \
        --display-name="OpenFGA Deployer" \
        --description="Service account para deploy do OpenFGA"
    print_status "Service account criada: $SERVICE_ACCOUNT_NAME"
else
    print_warning "Service account $SERVICE_ACCOUNT_NAME já existe"
fi

# Adicionar roles necessários
print_status "Adicionando roles ao service account..."
roles=(
    "roles/run.admin"
    "roles/storage.admin"
    "roles/artifactregistry.admin"
    "roles/iam.serviceAccountUser"
    "roles/cloudsql.client"
)

for role in "${roles[@]}"; do
    gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
        --role="$role"
    print_status "Role adicionado: $role"
done

# Gerar chave da service account
print_status "Gerando chave da service account..."
KEY_FILE="$SERVICE_ACCOUNT_NAME-key.json"
gcloud iam service-accounts keys create $KEY_FILE \
    --iam-account=$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com

print_status "Chave salva em: $KEY_FILE"

# Criar instância Cloud SQL (opcional)
read -p "Deseja criar uma instância Cloud SQL PostgreSQL? (y/n): " create_sql
if [ "$create_sql" == "y" ]; then
    print_status "Criando instância Cloud SQL..."
    
    get_input "Nome da instância Cloud SQL" SQL_INSTANCE_NAME "openfga-db"
    get_input "Versão do PostgreSQL" SQL_VERSION "POSTGRES_14"
    get_input "Tier da instância" SQL_TIER "db-f1-micro"
    
    if ! gcloud sql instances describe $SQL_INSTANCE_NAME &> /dev/null; then
        gcloud sql instances create $SQL_INSTANCE_NAME \
            --database-version=$SQL_VERSION \
            --tier=$SQL_TIER \
            --region=$REGION \
            --storage-auto-increase \
            --storage-size=10GB \
            --storage-type=SSD
        
        print_status "Instância Cloud SQL criada: $SQL_INSTANCE_NAME"
        
        # Criar databases
        print_status "Criando databases..."
        gcloud sql databases create openfga_dev --instance=$SQL_INSTANCE_NAME
        gcloud sql databases create openfga_prod --instance=$SQL_INSTANCE_NAME
        
        # Criar usuários
        print_status "Criando usuários..."
        gcloud sql users create openfga_dev --instance=$SQL_INSTANCE_NAME --password="$(openssl rand -base64 32)"
        gcloud sql users create openfga_prod --instance=$SQL_INSTANCE_NAME --password="$(openssl rand -base64 32)"
        
        print_warning "Anote as senhas geradas para os usuários!"
    else
        print_warning "Instância $SQL_INSTANCE_NAME já existe"
    fi
fi

# Exibir informações finais
print_status "Configuração concluída!"
echo ""
echo "=== INFORMAÇÕES PARA CONFIGURAR NO GITHUB ==="
echo "PROJECT_ID: $PROJECT_ID"
echo "REGION: $REGION"
echo "SERVICE_ACCOUNT_KEY: (conteúdo do arquivo $KEY_FILE)"
echo ""
echo "=== PRÓXIMOS PASSOS ==="
echo "1. Adicione os secrets no GitHub:"
echo "   - GCP_PROJECT_ID: $PROJECT_ID"
echo "   - GCP_SA_KEY: (conteúdo do arquivo $KEY_FILE)"
echo "   - DATABASE_URL: postgresql://user:pass@host:port/db"
echo ""
echo "2. Configure as URLs dos bancos de dados nos secrets:"
echo "   - Para desenvolvimento: postgresql://openfga_dev:PASSWORD@/openfga_dev?host=/cloudsql/$PROJECT_ID:$REGION:$SQL_INSTANCE_NAME"
echo "   - Para produção: postgresql://openfga_prod:PASSWORD@/openfga_prod?host=/cloudsql/$PROJECT_ID:$REGION:$SQL_INSTANCE_NAME"
echo ""
echo "3. O arquivo $KEY_FILE contém credenciais sensíveis. Guarde-o com segurança e delete após uso."
echo ""
print_warning "IMPORTANTE: Não commite o arquivo $KEY_FILE no repositório!"
