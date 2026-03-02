-- 删除触发器
DROP TRIGGER IF EXISTS update_providers_updated_at ON providers;

-- 删除函数
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 删除表
DROP TABLE IF EXISTS providers;
