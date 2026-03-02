-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- 删除索引
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;

-- 删除表
DROP TABLE IF EXISTS users;
