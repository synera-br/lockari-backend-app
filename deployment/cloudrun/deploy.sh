#!/bin/bash

# Cloud Run Deployment Script
# This script deploys the Lockari Platform to Google Cloud Run (production environment)

set -e

# Configuration
PROJECT_ID="lockari-project"
REGION="us-central1"
DEPLOYMENT_DIR="$(dirname "$0")"
CLOUDRUN_DIR="${DEPLOYMENT_DIR}/cloudrun"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if gcloud is installed
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud is not installed. Please install Google Cloud SDK first."
        exit 1
    fi
    
    # Check if terraform is installed
    if ! command -v terraform &> /dev/null; then
        log_error "terraform is not installed. Please install terraform first."
        exit 1
    fi
    
    # Check if docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "docker is not installed. Please install docker first."
        exit 1
    fi
    
    # Check authentication
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
        log_error "Not authenticated with Google Cloud. Please run 'gcloud auth login'."
        exit 1
    fi
    
    # Set project
    gcloud config set project ${PROJECT_ID}
    
    log_info "Prerequisites check passed."
}

setup_infrastructure() {
    log_info "Setting up cloud infrastructure..."
    
    cd "${CLOUDRUN_DIR}"
    
    # Initialize Terraform
    terraform init
    
    # Plan infrastructure
    terraform plan -out=tfplan
    
    # Apply infrastructure
    terraform apply tfplan
    
    log_info "Infrastructure setup completed."
}

build_images() {
    log_info "Building and pushing container images..."
    
    # Configure Docker for GCR
    gcloud auth configure-docker gcr.io
    
    # Build Backend Image
    log_info "Building backend image..."
    cd "${DEPLOYMENT_DIR}/.."
    docker build -t gcr.io/${PROJECT_ID}/lockari-backend:latest -f build/Dockerfile .
    docker push gcr.io/${PROJECT_ID}/lockari-backend:latest
    
    # Build Frontend Image
    log_info "Building frontend image..."
    cd "${DEPLOYMENT_DIR}/../../lockari-frontend-app"
    docker build -t gcr.io/${PROJECT_ID}/lockari-frontend:latest -f Dockerfile .
    docker push gcr.io/${PROJECT_ID}/lockari-frontend:latest
    
    log_info "Container images built and pushed successfully."
}

deploy_secrets() {
    log_info "Deploying secrets..."
    
    # Create secrets using gcloud (since we need to set actual values)
    gcloud secrets create lockari-database-url --data-file=- <<EOF
postgres://lockari:${POSTGRES_PASSWORD}@/lockari_prod?host=/cloudsql/${PROJECT_ID}:${REGION}:lockari-prod
EOF
    
    gcloud secrets create lockari-openfga-database-url --data-file=- <<EOF
postgres://openfga:${OPENFGA_POSTGRES_PASSWORD}@/openfga_prod?host=/cloudsql/${PROJECT_ID}:${REGION}:lockari-openfga-prod
EOF
    
    gcloud secrets create lockari-jwt-secret --data-file=- <<EOF
${JWT_SECRET}
EOF
    
    gcloud secrets create lockari-jwt-refresh-secret --data-file=- <<EOF
${JWT_REFRESH_SECRET}
EOF
    
    gcloud secrets create lockari-encryption-key --data-file=- <<EOF
${ENCRYPTION_KEY}
EOF
    
    gcloud secrets create lockari-openfga-preshared-key --data-file=- <<EOF
${OPENFGA_PRESHARED_KEY}
EOF
    
    # Set IAM permissions
    gcloud secrets add-iam-policy-binding lockari-database-url \
        --member="serviceAccount:lockari-backend@${PROJECT_ID}.iam.gserviceaccount.com" \
        --role="roles/secretmanager.secretAccessor"
    
    gcloud secrets add-iam-policy-binding lockari-openfga-database-url \
        --member="serviceAccount:lockari-openfga@${PROJECT_ID}.iam.gserviceaccount.com" \
        --role="roles/secretmanager.secretAccessor"
    
    log_info "Secrets deployed successfully."
}

deploy_openfga() {
    log_info "Deploying OpenFGA service..."
    
    # Deploy OpenFGA
    gcloud run services replace ${CLOUDRUN_DIR}/openfga-service.yaml \
        --region=${REGION}
    
    # Wait for deployment to complete
    gcloud run services describe lockari-openfga \
        --region=${REGION} \
        --format="value(status.conditions[0].status)" | grep -q "True"
    
    log_info "OpenFGA service deployed successfully."
}

deploy_backend() {
    log_info "Deploying backend service..."
    
    # Deploy Backend
    gcloud run services replace ${CLOUDRUN_DIR}/backend-service.yaml \
        --region=${REGION}
    
    # Wait for deployment to complete
    gcloud run services describe lockari-backend \
        --region=${REGION} \
        --format="value(status.conditions[0].status)" | grep -q "True"
    
    log_info "Backend service deployed successfully."
}

deploy_frontend() {
    log_info "Deploying frontend service..."
    
    # Deploy Frontend
    gcloud run services replace ${CLOUDRUN_DIR}/frontend-service.yaml \
        --region=${REGION}
    
    # Wait for deployment to complete
    gcloud run services describe lockari-frontend \
        --region=${REGION} \
        --format="value(status.conditions[0].status)" | grep -q "True"
    
    log_info "Frontend service deployed successfully."
}

setup_load_balancer() {
    log_info "Setting up load balancer..."
    
    # Apply load balancer configuration
    kubectl apply -f ${CLOUDRUN_DIR}/load-balancer.yaml
    
    # Wait for load balancer to be ready
    kubectl wait --for=condition=Ready ingress/lockari-ingress --timeout=600s
    
    log_info "Load balancer setup completed."
}

migrate_database() {
    log_info "Running database migrations..."
    
    # Run migrations using a Cloud Run job
    gcloud run jobs create lockari-migrate \
        --image=gcr.io/${PROJECT_ID}/lockari-backend:latest \
        --region=${REGION} \
        --service-account=lockari-backend@${PROJECT_ID}.iam.gserviceaccount.com \
        --vpc-connector=lockari-vpc-connector \
        --set-env-vars="APP_ENV=production" \
        --set-secrets="DATABASE_URL=lockari-database-url:latest" \
        --command="./migrate" \
        --args="up"
    
    # Execute migration job
    gcloud run jobs execute lockari-migrate --region=${REGION} --wait
    
    log_info "Database migrations completed."
}

verify_deployment() {
    log_info "Verifying deployment..."
    
    # Get service URLs
    BACKEND_URL=$(gcloud run services describe lockari-backend --region=${REGION} --format="value(status.url)")
    FRONTEND_URL=$(gcloud run services describe lockari-frontend --region=${REGION} --format="value(status.url)")
    OPENFGA_URL=$(gcloud run services describe lockari-openfga --region=${REGION} --format="value(status.url)")
    
    # Test health endpoints
    log_info "Testing health endpoints..."
    
    if curl -f "${BACKEND_URL}/health" > /dev/null 2>&1; then
        log_info "Backend health check passed."
    else
        log_warn "Backend health check failed."
    fi
    
    if curl -f "${FRONTEND_URL}/api/health" > /dev/null 2>&1; then
        log_info "Frontend health check passed."
    else
        log_warn "Frontend health check failed."
    fi
    
    if curl -f "${OPENFGA_URL}/healthz" > /dev/null 2>&1; then
        log_info "OpenFGA health check passed."
    else
        log_warn "OpenFGA health check failed."
    fi
    
    log_info "Deployment verification completed."
}

print_access_info() {
    log_info "Deployment completed successfully!"
    echo
    echo "Service URLs:"
    echo "  Frontend: https://app.lockari.com"
    echo "  Backend API: https://api.lockari.com"
    echo "  OpenFGA: https://openfga.lockari.com"
    echo
    echo "Direct Cloud Run URLs:"
    echo "  Backend: $(gcloud run services describe lockari-backend --region=${REGION} --format="value(status.url)")"
    echo "  Frontend: $(gcloud run services describe lockari-frontend --region=${REGION} --format="value(status.url)")"
    echo "  OpenFGA: $(gcloud run services describe lockari-openfga --region=${REGION} --format="value(status.url)")"
    echo
    echo "To check logs:"
    echo "  gcloud run logs read --service=lockari-backend --region=${REGION}"
    echo "  gcloud run logs read --service=lockari-frontend --region=${REGION}"
    echo "  gcloud run logs read --service=lockari-openfga --region=${REGION}"
    echo
    echo "To monitor services:"
    echo "  gcloud run services list --region=${REGION}"
    echo "  gcloud run revisions list --service=lockari-backend --region=${REGION}"
}

rollback_deployment() {
    log_warn "Rolling back deployment..."
    
    # Get previous revisions
    BACKEND_PREV=$(gcloud run revisions list --service=lockari-backend --region=${REGION} --format="value(metadata.name)" --limit=2 | tail -1)
    FRONTEND_PREV=$(gcloud run revisions list --service=lockari-frontend --region=${REGION} --format="value(metadata.name)" --limit=2 | tail -1)
    OPENFGA_PREV=$(gcloud run revisions list --service=lockari-openfga --region=${REGION} --format="value(metadata.name)" --limit=2 | tail -1)
    
    # Rollback services
    gcloud run services update-traffic lockari-backend --to-revisions=${BACKEND_PREV}=100 --region=${REGION}
    gcloud run services update-traffic lockari-frontend --to-revisions=${FRONTEND_PREV}=100 --region=${REGION}
    gcloud run services update-traffic lockari-openfga --to-revisions=${OPENFGA_PREV}=100 --region=${REGION}
    
    log_info "Rollback completed."
}

cleanup_deployment() {
    log_warn "Cleaning up deployment..."
    
    # Delete Cloud Run services
    gcloud run services delete lockari-backend --region=${REGION} --quiet
    gcloud run services delete lockari-frontend --region=${REGION} --quiet
    gcloud run services delete lockari-openfga --region=${REGION} --quiet
    
    # Delete infrastructure
    cd "${CLOUDRUN_DIR}"
    terraform destroy -auto-approve
    
    log_info "Cleanup completed."
}

# Main execution
main() {
    case "${1:-deploy}" in
        deploy)
            log_info "Starting Cloud Run deployment..."
            check_prerequisites
            setup_infrastructure
            build_images
            deploy_secrets
            deploy_openfga
            deploy_backend
            deploy_frontend
            setup_load_balancer
            migrate_database
            verify_deployment
            print_access_info
            ;;
        rollback)
            log_info "Rolling back deployment..."
            check_prerequisites
            rollback_deployment
            ;;
        cleanup)
            log_info "Cleaning up deployment..."
            check_prerequisites
            cleanup_deployment
            ;;
        verify)
            log_info "Verifying deployment..."
            check_prerequisites
            verify_deployment
            ;;
        build)
            log_info "Building container images..."
            check_prerequisites
            build_images
            ;;
        migrate)
            log_info "Running database migrations..."
            check_prerequisites
            migrate_database
            ;;
        *)
            echo "Usage: $0 [deploy|rollback|cleanup|verify|build|migrate]"
            echo "  deploy   - Deploy the application (default)"
            echo "  rollback - Rollback the deployment"
            echo "  cleanup  - Clean up all resources"
            echo "  verify   - Verify the deployment"
            echo "  build    - Build and push container images"
            echo "  migrate  - Run database migrations"
            exit 1
            ;;
    esac
}

# Check for required environment variables
if [[ "${1:-deploy}" == "deploy" ]]; then
    required_vars=(
        "POSTGRES_PASSWORD"
        "OPENFGA_POSTGRES_PASSWORD"
        "JWT_SECRET"
        "JWT_REFRESH_SECRET"
        "ENCRYPTION_KEY"
        "OPENFGA_PRESHARED_KEY"
    )
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            log_error "Environment variable ${var} is required but not set."
            exit 1
        fi
    done
fi

# Run main function
main "$@"
