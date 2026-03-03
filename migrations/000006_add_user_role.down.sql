-- 删除角色索引
DROP INDEX IF EXISTS idx_users_role;

-- 删除角色约束和字段
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_role;
ALTER TABLE users DROP COLUMN IF EXISTS role;
