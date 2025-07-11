-- Inicialização do banco PostgreSQL para OpenFGA
-- Este script é executado automaticamente quando o container inicia

-- Criar database e usuário para OpenFGA
CREATE DATABASE openfga;
CREATE USER openfga WITH PASSWORD 'dev_password_123';
GRANT ALL PRIVILEGES ON DATABASE openfga TO openfga;

-- Configurações de performance para desenvolvimento
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
ALTER SYSTEM SET max_connections = 100;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;

-- Conectar ao database openfga e criar extensões necessárias
\c openfga;

-- Criar extensões úteis
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Criar índices para performance (OpenFGA criará as tabelas automaticamente)
-- Estes comandos serão executados após o OpenFGA criar as tabelas
-- Por agora, deixamos comentado para não dar erro na inicialização

-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tuple_user ON tuple (user_object_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tuple_relation ON tuple (relation);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tuple_object ON tuple (object_object_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tuple_user_relation ON tuple (user_object_id, relation);

-- Log da inicialização
\echo 'OpenFGA database initialized successfully!'
\echo 'Database: openfga'
\echo 'User: openfga'
\echo 'Extensions: uuid-ossp, pg_stat_statements'
