## 1. 项目初始化
- [x] 1.1 初始化 Go 模块 `github.com/lucheng0127/courier`
- [x] 1.2 创建项目目录结构（`cmd/`、`internal/`、`pkg/`）
- [x] 1.3 添加依赖项（Gin、GORM、SQLite、Zap）

## 2. 数据库层实现
- [x] 2.1 定义 User 数据模型
- [x] 2.2 定义 APIKey 数据模型
- [x] 2.3 实现 UserRepository 接口和 SQLite 实现
- [x] 2.4 实现 APIKeyRepository 接口和 SQLite 实现
- [x] 2.5 创建数据库连接和迁移逻辑

## 3. 业务逻辑层实现
- [x] 3.1 实现 UserService（创建用户、查询用户）
- [x] 3.2 实现 APIKeyService（生成、查询、删除、禁用 API Key）
- [x] 3.3 实现 API Key 生成逻辑（使用 crypto/rand）

## 4. HTTP 层实现
- [x] 4.1 定义 RESTful API 路由
- [x] 4.2 实现 UserHandler（创建用户、查询用户列表/详情）
- [x] 4.3 实现 APIKeyHandler（生成、查询、删除、禁用 API Key）
- [x] 4.4 实现请求和响应 DTO

## 5. 中间件和配置
- [x] 5.1 实现配置加载（YAML 格式）
- [x] 5.2 初始化 Zap 日志器
- [x] 5.3 创建 main.go 入口文件
- [x] 5.4 实现依赖注入组装

## 6. 验证和测试
- [x] 6.1 确保所有 API 接口可正常调用
- [x] 6.2 确认数据库正确持久化数据
- [x] 6.3 验证 API Key 生成格式正确
