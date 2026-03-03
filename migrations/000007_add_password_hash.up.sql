-- 添加密码哈希字段到 users 表（用于登录）
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
