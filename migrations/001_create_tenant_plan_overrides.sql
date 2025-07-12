-- Migration: Create tenant_plan_overrides table
-- Description: Tabela para armazenar customizações de limites por tenant/plano

CREATE TABLE tenant_plan_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    plan_type VARCHAR(50) NOT NULL,
    limits JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT unique_tenant_plan UNIQUE (tenant_id, plan_type),
    CONSTRAINT valid_plan_type CHECK (plan_type IN ('free', 'pro', 'enterprise'))
);

-- Índices para performance
CREATE INDEX idx_tenant_plan_overrides_tenant_id ON tenant_plan_overrides(tenant_id);
CREATE INDEX idx_tenant_plan_overrides_plan_type ON tenant_plan_overrides(plan_type);
CREATE INDEX idx_tenant_plan_overrides_updated_at ON tenant_plan_overrides(updated_at);

-- Exemplo de estrutura do campo limits (JSONB)
/*
{
  "vault_limit": 100,
  "user_limit": 25,
  "is_unlimited": false,
  "features": [
    "basic",
    "vault_limit",
    "user_limit",
    "api_access",
    "audit_logs",
    "group_management",
    "sso"
  ],
  "custom_features": {
    "custom_branding": true,
    "dedicated_support": true,
    "custom_integrations": ["slack", "teams"]
  },
  "description": "Plano Pro customizado para cliente VIP"
}
*/

-- Função para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para atualizar updated_at
CREATE TRIGGER update_tenant_plan_overrides_updated_at
    BEFORE UPDATE ON tenant_plan_overrides
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentários para documentação
COMMENT ON TABLE tenant_plan_overrides IS 'Customizações de limites e funcionalidades por tenant/plano';
COMMENT ON COLUMN tenant_plan_overrides.tenant_id IS 'ID do tenant que possui as customizações';
COMMENT ON COLUMN tenant_plan_overrides.plan_type IS 'Tipo do plano (free, pro, enterprise)';
COMMENT ON COLUMN tenant_plan_overrides.limits IS 'Limites e funcionalidades customizadas em formato JSON';
COMMENT ON COLUMN tenant_plan_overrides.created_at IS 'Data de criação da customização';
COMMENT ON COLUMN tenant_plan_overrides.updated_at IS 'Data da última atualização';

-- Exemplos de uso
-- 1. Inserir customização para um tenant específico
/*
INSERT INTO tenant_plan_overrides (tenant_id, plan_type, limits)
VALUES (
    'tenant-vip-123',
    'pro',
    '{
        "vault_limit": 100,
        "user_limit": 25,
        "is_unlimited": false,
        "features": [
            "basic",
            "vault_limit",
            "user_limit",
            "api_access",
            "audit_logs",
            "group_management",
            "sso"
        ],
        "description": "Plano Pro customizado para cliente VIP"
    }'::jsonb
);
*/

-- 2. Buscar customizações de um tenant
/*
SELECT * FROM tenant_plan_overrides 
WHERE tenant_id = 'tenant-vip-123' AND plan_type = 'pro';
*/

-- 3. Atualizar limites
/*
UPDATE tenant_plan_overrides 
SET limits = jsonb_set(limits, '{vault_limit}', '200')
WHERE tenant_id = 'tenant-vip-123' AND plan_type = 'pro';
*/

-- 4. Verificar se tenant tem funcionalidade específica
/*
SELECT EXISTS (
    SELECT 1 FROM tenant_plan_overrides 
    WHERE tenant_id = 'tenant-vip-123' 
    AND plan_type = 'pro'
    AND limits->'features' ? 'sso'
);
*/

-- 5. Listar tenants com SSO customizado
/*
SELECT tenant_id, plan_type, limits->>'description' as description
FROM tenant_plan_overrides
WHERE limits->'features' ? 'sso';
*/
