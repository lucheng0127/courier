-- 添加角色字段到 users 表
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user';
ALTER TABLE users ADD CONSTRAINT check_role CHECK (role IN ('user', 'admin'));

-- 为角色字段创建索引
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 为现有用户设置默认角色
UPDATE users SET role = 'user' WHERE role IS NULL;
