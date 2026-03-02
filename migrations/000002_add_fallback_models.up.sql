-- 添加 fallback_models 字段到 providers 表
ALTER TABLE providers ADD COLUMN IF NOT EXISTS fallback_models JSONB;
