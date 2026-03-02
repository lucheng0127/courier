-- 删除密码哈希字段
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
