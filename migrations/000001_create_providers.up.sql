-- 创建 providers 表
CREATE TABLE IF NOT EXISTS providers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    base_url VARCHAR(2048) NOT NULL,
    timeout INTEGER NOT NULL DEFAULT 300,
    api_key VARCHAR(2048),
    extra_config JSONB,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_providers_name ON providers(name);
CREATE INDEX idx_providers_type ON providers(type);
CREATE INDEX idx_providers_enabled ON providers(enabled);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_providers_updated_at
    BEFORE UPDATE ON providers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
