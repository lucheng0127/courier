-- 删除 fallback_models 字段
ALTER TABLE providers DROP COLUMN IF EXISTS fallback_models;
