-- 删除索引
DROP INDEX IF EXISTS idx_usage_records_request_id;
DROP INDEX IF EXISTS idx_usage_records_user_timestamp;
DROP INDEX IF EXISTS idx_usage_records_timestamp;
DROP INDEX IF EXISTS idx_usage_records_api_key_id;
DROP INDEX IF EXISTS idx_usage_records_user_id;

-- 删除表
DROP TABLE IF EXISTS usage_records;
