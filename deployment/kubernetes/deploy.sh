#!/bin/bash

# Kubernetes Deployment Script
# This script deploys the Lockari Platform to Kubernetes (test environment)

set -e

# Configuration
NAMESPACE="lockari-test"
KUBECTL_CONTEXT="gke_lockari-project_us-central1-a_lockari-test-cluster"
DEPLOYMENT_DIR="$(dirname "$0")"
KUBERNETES_DIR="${DEPLOYMENT_DIR}/kubernetes"

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
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed. Please install kubectl first."
        exit 1
    fi
    
    # Check if helm is installed
    if ! command -v helm &> /dev/null; then
        log_error "helm is not installed. Please install helm first."
        exit 1
    fi
    
    # Check if gcloud is installed
    if ! command -v gcloud &> /dev/null; then
        log_error "gcloud is not installed. Please install Google Cloud SDK first."
        exit 1
    fi
    
    # Check if context exists
    if ! kubectl config get-contexts | grep -q "${KUBECTL_CONTEXT}"; then
        log_error "Kubernetes context ${KUBECTL_CONTEXT} not found. Please configure your cluster access."
        exit 1
    fi
    
    log_info "Prerequisites check passed."
}

setup_context() {
    log_info "Setting up Kubernetes context..."
    kubectl config use-context "${KUBECTL_CONTEXT}"
    
    # Verify connection
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Please check your configuration."
        exit 1
    fi
    
    log_info "Connected to Kubernetes cluster successfully."
}

deploy_infrastructure() {
    log_info "Deploying infrastructure components..."
    
    # Deploy namespace and config
    kubectl apply -f "${KUBERNETES_DIR}/00-namespace.yaml"
    kubectl apply -f "${KUBERNETES_DIR}/01-secrets.yaml"
    
    # Wait for namespace to be ready
    kubectl wait --for=condition=Active namespace/${NAMESPACE} --timeout=60s
    
    # Deploy databases
    log_info "Deploying PostgreSQL..."
    kubectl apply -f "${KUBERNETES_DIR}/02-postgres.yaml"
    
    log_info "Deploying Redis..."
    kubectl apply -f "${KUBERNETES_DIR}/03-redis.yaml"
    
    log_info "Deploying RabbitMQ..."
    kubectl apply -f "${KUBERNETES_DIR}/04-rabbitmq.yaml"
    
    # Wait for databases to be ready
    log_info "Waiting for databases to be ready..."
    kubectl wait --for=condition=Ready pod -l app=postgres -n ${NAMESPACE} --timeout=300s
    kubectl wait --for=condition=Ready pod -l app=redis -n ${NAMESPACE} --timeout=300s
    kubectl wait --for=condition=Ready pod -l app=rabbitmq -n ${NAMESPACE} --timeout=300s
    
    log_info "Infrastructure components deployed successfully."
}

deploy_services() {
    log_info "Deploying application services..."
    
    # Deploy OpenFGA
    log_info "Deploying OpenFGA..."
    kubectl apply -f "${KUBERNETES_DIR}/05-openfga.yaml"
    
    # Wait for OpenFGA to be ready
    kubectl wait --for=condition=Ready pod -l app=openfga -n ${NAMESPACE} --timeout=300s
    
    # Deploy Backend
    log_info "Deploying Backend..."
    kubectl apply -f "${KUBERNETES_DIR}/06-backend.yaml"
    
    # Wait for Backend to be ready
    kubectl wait --for=condition=Ready pod -l app=backend -n ${NAMESPACE} --timeout=300s
    
    # Deploy Frontend
    log_info "Deploying Frontend..."
    kubectl apply -f "${KUBERNETES_DIR}/07-frontend.yaml"
    
    # Wait for Frontend to be ready
    kubectl wait --for=condition=Ready pod -l app=frontend -n ${NAMESPACE} --timeout=300s
    
    log_info "Application services deployed successfully."
}

deploy_networking() {
    log_info "Deploying networking components..."
    
    # Deploy Ingress
    kubectl apply -f "${KUBERNETES_DIR}/08-ingress.yaml"
    
    # Wait for ingress to be ready
    kubectl wait --for=condition=Ready ingress/lockari-ingress -n ${NAMESPACE} --timeout=300s
    
    log_info "Networking components deployed successfully."
}

deploy_monitoring() {
    log_info "Deploying monitoring components..."
    
    # Deploy monitoring
    kubectl apply -f "${KUBERNETES_DIR}/09-monitoring.yaml"
    
    # Wait for monitoring to be ready
    kubectl wait --for=condition=Ready pod -l app=grafana -n ${NAMESPACE} --timeout=300s
    
    log_info "Monitoring components deployed successfully."
}

verify_deployment() {
    log_info "Verifying deployment..."
    
    # Check all pods are running
    log_info "Checking pod status..."
    kubectl get pods -n ${NAMESPACE}
    
    # Check services
    log_info "Checking service status..."
    kubectl get services -n ${NAMESPACE}
    
    # Check ingress
    log_info "Checking ingress status..."
    kubectl get ingress -n ${NAMESPACE}
    
    # Test health endpoints
    log_info "Testing health endpoints..."
    
    # Port forward to test locally
    kubectl port-forward service/backend-service 8080:8080 -n ${NAMESPACE} &
    BACKEND_PID=$!
    
    sleep 5
    
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_info "Backend health check passed."
    else
        log_warn "Backend health check failed."
    fi
    
    kill $BACKEND_PID
    
    log_info "Deployment verification completed."
}

print_access_info() {
    log_info "Deployment completed successfully!"
    echo
    echo "Access URLs:"
    echo "  Frontend: https://test.lockari.com"
    echo "  Backend API: https://api-test.lockari.com"
    echo "  OpenFGA: https://openfga-test.lockari.com"
    echo "  Grafana: https://grafana-test.lockari.com"
    echo
    echo "Local port forwarding commands:"
    echo "  Frontend: kubectl port-forward service/frontend-service 3000:3000 -n ${NAMESPACE}"
    echo "  Backend: kubectl port-forward service/backend-service 8080:8080 -n ${NAMESPACE}"
    echo "  OpenFGA: kubectl port-forward service/openfga-service 8081:8080 -n ${NAMESPACE}"
    echo "  Grafana: kubectl port-forward service/grafana-service 3001:3000 -n ${NAMESPACE}"
    echo
    echo "To check logs:"
    echo "  kubectl logs -f deployment/backend -n ${NAMESPACE}"
    echo "  kubectl logs -f deployment/frontend -n ${NAMESPACE}"
    echo "  kubectl logs -f deployment/openfga -n ${NAMESPACE}"
}

rollback_deployment() {
    log_warn "Rolling back deployment..."
    
    # Rollback in reverse order
    kubectl rollout undo deployment/frontend -n ${NAMESPACE}
    kubectl rollout undo deployment/backend -n ${NAMESPACE}
    kubectl rollout undo deployment/openfga -n ${NAMESPACE}
    
    log_info "Rollback completed."
}

cleanup_deployment() {
    log_warn "Cleaning up deployment..."
    
    # Delete all resources
    kubectl delete namespace ${NAMESPACE}
    
    log_info "Cleanup completed."
}

# Main execution
main() {
    case "${1:-deploy}" in
        deploy)
            log_info "Starting Kubernetes deployment..."
            check_prerequisites
            setup_context
            deploy_infrastructure
            deploy_services
            deploy_networking
            deploy_monitoring
            verify_deployment
            print_access_info
            ;;
        rollback)
            log_info "Rolling back deployment..."
            setup_context
            rollback_deployment
            ;;
        cleanup)
            log_info "Cleaning up deployment..."
            setup_context
            cleanup_deployment
            ;;
        verify)
            log_info "Verifying deployment..."
            setup_context
            verify_deployment
            ;;
        *)
            echo "Usage: $0 [deploy|rollback|cleanup|verify]"
            echo "  deploy   - Deploy the application (default)"
            echo "  rollback - Rollback the deployment"
            echo "  cleanup  - Clean up all resources"
            echo "  verify   - Verify the deployment"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
