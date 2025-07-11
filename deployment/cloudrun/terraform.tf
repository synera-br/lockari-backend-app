# Terraform configuration for Google Cloud resources
terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }
  
  backend "gcs" {
    bucket = "lockari-terraform-state"
    prefix = "production"
  }
}

# Provider configuration
provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

# Variables
variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
  default     = "lockari-project"
}

variable "region" {
  description = "Google Cloud Region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "Google Cloud Zone"
  type        = string
  default     = "us-central1-a"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

# Enable required APIs
resource "google_project_service" "apis" {
  for_each = toset([
    "cloudsql.googleapis.com",
    "redis.googleapis.com",
    "pubsub.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "compute.googleapis.com",
    "vpcaccess.googleapis.com",
    "servicenetworking.googleapis.com",
    "monitoring.googleapis.com",
    "logging.googleapis.com",
    "cloudtrace.googleapis.com",
    "cloudbuild.googleapis.com",
    "artifactregistry.googleapis.com"
  ])
  
  service = each.value
  
  disable_dependent_services = true
}

# VPC Network
resource "google_compute_network" "vpc" {
  name                    = "lockari-vpc"
  auto_create_subnetworks = false
  routing_mode           = "REGIONAL"
}

resource "google_compute_subnetwork" "subnet" {
  name          = "lockari-subnet"
  ip_cidr_range = "10.0.0.0/24"
  region        = var.region
  network       = google_compute_network.vpc.id
  
  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.1.0.0/24"
  }
  
  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.2.0.0/16"
  }
}

# VPC Connector for Cloud Run
resource "google_vpc_access_connector" "connector" {
  name          = "lockari-vpc-connector"
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.vpc.name
  region        = var.region
  
  depends_on = [google_project_service.apis]
}

# Private Service Connection for Cloud SQL
resource "google_compute_global_address" "private_ip_address" {
  name          = "lockari-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

# Cloud SQL Instance for Main Database
resource "google_sql_database_instance" "main" {
  name             = "lockari-prod"
  database_version = "POSTGRES_15"
  region           = var.region
  deletion_protection = true
  
  depends_on = [google_service_networking_connection.private_vpc_connection]
  
  settings {
    tier                        = "db-custom-2-4096"
    availability_type          = "REGIONAL"
    disk_type                  = "PD_SSD"
    disk_size                  = 100
    disk_autoresize            = true
    disk_autoresize_limit      = 500
    
    backup_configuration {
      enabled                        = true
      start_time                    = "02:00"
      location                      = var.region
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 30
        retention_unit   = "COUNT"
      }
    }
    
    maintenance_window {
      day          = 7
      hour         = 3
      update_track = "stable"
    }
    
    database_flags {
      name  = "log_statement"
      value = "all"
    }
    
    database_flags {
      name  = "log_min_duration_statement"
      value = "1000"
    }
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.vpc.id
      require_ssl     = true
    }
    
    insights_config {
      query_insights_enabled  = true
      record_application_tags = true
      record_client_address   = true
    }
  }
}

# Cloud SQL Database
resource "google_sql_database" "main" {
  name     = "lockari_prod"
  instance = google_sql_database_instance.main.name
}

# Cloud SQL User
resource "google_sql_user" "main" {
  name     = "lockari"
  instance = google_sql_database_instance.main.name
  password = var.postgres_password
}

# Cloud SQL Instance for OpenFGA
resource "google_sql_database_instance" "openfga" {
  name             = "lockari-openfga-prod"
  database_version = "POSTGRES_15"
  region           = var.region
  deletion_protection = true
  
  depends_on = [google_service_networking_connection.private_vpc_connection]
  
  settings {
    tier              = "db-custom-1-2048"
    availability_type = "REGIONAL"
    disk_type         = "PD_SSD"
    disk_size         = 50
    disk_autoresize   = true
    
    backup_configuration {
      enabled                        = true
      start_time                    = "03:00"
      location                      = var.region
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 30
        retention_unit   = "COUNT"
      }
    }
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.vpc.id
      require_ssl     = true
    }
  }
}

# Cloud SQL Database for OpenFGA
resource "google_sql_database" "openfga" {
  name     = "openfga_prod"
  instance = google_sql_database_instance.openfga.name
}

# Cloud SQL User for OpenFGA
resource "google_sql_user" "openfga" {
  name     = "openfga"
  instance = google_sql_database_instance.openfga.name
  password = var.openfga_postgres_password
}

# Memorystore Redis Instance
resource "google_redis_instance" "cache" {
  name           = "lockari-cache"
  tier           = "STANDARD_HA"
  memory_size_gb = 1
  region         = var.region
  
  location_id             = var.zone
  alternative_location_id = "${substr(var.region, 0, length(var.region)-1)}b"
  
  authorized_network = google_compute_network.vpc.id
  
  redis_version     = "REDIS_7_0"
  display_name      = "Lockari Cache"
  
  auth_enabled = true
  
  maintenance_policy {
    weekly_maintenance_window {
      day = "SUNDAY"
      start_time {
        hours   = 3
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
  }
}

# Pub/Sub Topics
resource "google_pubsub_topic" "audit_logs" {
  name = "lockari-audit-logs"
  
  message_retention_duration = "604800s"  # 7 days
}

resource "google_pubsub_topic" "notifications" {
  name = "lockari-notifications"
  
  message_retention_duration = "86400s"  # 1 day
}

resource "google_pubsub_topic" "webhooks" {
  name = "lockari-webhooks"
  
  message_retention_duration = "86400s"  # 1 day
}

# Pub/Sub Subscriptions
resource "google_pubsub_subscription" "audit_logs" {
  name  = "lockari-audit-logs-sub"
  topic = google_pubsub_topic.audit_logs.name
  
  message_retention_duration = "604800s"  # 7 days
  retain_acked_messages      = false
  ack_deadline_seconds       = 20
  
  expiration_policy {
    ttl = "300000.5s"
  }
  
  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }
}

# Service Accounts
resource "google_service_account" "backend" {
  account_id   = "lockari-backend"
  display_name = "Lockari Backend Service Account"
  description  = "Service account for Lockari backend Cloud Run service"
}

resource "google_service_account" "frontend" {
  account_id   = "lockari-frontend"
  display_name = "Lockari Frontend Service Account"
  description  = "Service account for Lockari frontend Cloud Run service"
}

resource "google_service_account" "openfga" {
  account_id   = "lockari-openfga"
  display_name = "Lockari OpenFGA Service Account"
  description  = "Service account for Lockari OpenFGA Cloud Run service"
}

# IAM Bindings for Backend
resource "google_project_iam_member" "backend_sql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_redis_editor" {
  project = var.project_id
  role    = "roles/redis.editor"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_pubsub_editor" {
  project = var.project_id
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_trace_agent" {
  project = var.project_id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_monitoring_writer" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# IAM Bindings for Frontend
resource "google_project_iam_member" "frontend_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.frontend.email}"
}

# IAM Bindings for OpenFGA
resource "google_project_iam_member" "openfga_sql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.openfga.email}"
}

resource "google_project_iam_member" "openfga_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.openfga.email}"
}

# Static IP for Load Balancer
resource "google_compute_global_address" "lb_ip" {
  name = "lockari-ip"
}

# Outputs
output "vpc_network_id" {
  value = google_compute_network.vpc.id
}

output "vpc_connector_id" {
  value = google_vpc_access_connector.connector.id
}

output "database_connection_name" {
  value = google_sql_database_instance.main.connection_name
}

output "openfga_database_connection_name" {
  value = google_sql_database_instance.openfga.connection_name
}

output "redis_host" {
  value = google_redis_instance.cache.host
}

output "redis_port" {
  value = google_redis_instance.cache.port
}

output "load_balancer_ip" {
  value = google_compute_global_address.lb_ip.address
}

output "backend_service_account_email" {
  value = google_service_account.backend.email
}

output "frontend_service_account_email" {
  value = google_service_account.frontend.email
}

output "openfga_service_account_email" {
  value = google_service_account.openfga.email
}
