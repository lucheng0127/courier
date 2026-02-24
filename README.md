# Courier

Courier 是一个 AI API 网关服务，提供统一的接口来调用多个上游 AI 模型。系统支持用户管理、API Key 认证、模型对话和完整的请求日志记录。

## 功能特性

- **用户管理**：创建用户、查询用户列表和详情
- **API Key 管理**：为用户生成 API Key、查询列表、删除和禁用
- **多模型支持**：通过 YAML 配置多个上游 AI 模型（如 Qwen、DeepSeek 等）
- **模型对话**：支持非流式和流式（SSE）两种对话模式
- **API Key 认证**：基于 Bearer Token 的请求认证机制
- **请求日志**：完整记录所有对话请求和响应，支持 Token 统计

## 技术栈

- **语言**：Go 1.23+
- **Web 框架**：Gin
- **ORM**：GORM
- **数据库**：SQLite
- **日志**：Zap

## 项目结构

```
courier/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── client/              # 上游模型客户端
│   ├── handler/             # HTTP 处理器
│   ├── middleware/          # 中间件（认证等）
│   ├── model/               # 数据模型
│   ├── repository/          # 数据仓储
│   ├── router/              # 路由配置
│   └── service/             # 业务逻辑
├── pkg/
│   ├── config/              # 配置管理
│   └── logger/              # 日志初始化
├── config.yaml              # 配置文件
└── go.mod
```

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/lucheng0127/courier.git
cd courier
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置文件

复制并编辑 `config.yaml`：

```yaml
server:
  port: "8080"

db:
  data_source_name: "courier.db"

models:
  - name: qwen-turbo
    provider: qwen
    base_url: https://dashscope.aliyuncs.com/compatible-mode/v1
    api_key: ${QWEN_API_KEY}  # 或直接填写 API Key

  - name: deepseek-chat
    provider: deepseek
    base_url: https://api.deepseek.com/v1
    api_key: ${DEEPSEEK_API_KEY}
```

**配置说明：**

- `name`: 模型唯一标识，用于 API 调用时指定模型
- `provider`: 模型提供商名称
- `base_url`: 上游模型的 API 地址
- `api_key`: 上游模型的 API Key（支持环境变量格式 `${VAR_NAME}` 或直接填写）

### 4. 编译并运行

```bash
# 编译
go build -o bin/server ./cmd/server/main.go

# 运行
./bin/server
```

或直接使用 `go run`：

```bash
go run ./cmd/server/main.go
```

服务启动后监听在 `http://localhost:8080`

## API 文档

### 用户管理

#### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "alice",
    "email": "alice@example.com"
  }'
```

**响应示例：**

```json
{
  "id": 1,
  "name": "alice",
  "email": "alice@example.com",
  "created_at": "2026-02-24T17:00:00+08:00",
  "updated_at": "2026-02-24T17:00:00+08:00"
}
```

#### 查询用户列表

```bash
curl http://localhost:8080/api/v1/users
```

#### 查询用户详情

```bash
curl http://localhost:8080/api/v1/users/1
```

### API Key 管理

#### 生成 API Key

```bash
curl -X POST http://localhost:8080/api/v1/users/1/apikeys
```

**响应示例：**

```json
{
  "id": 1,
  "user_id": 1,
  "key": "ck_7320c056b0b9b62553ab9924c3f0fedf055cebce271417b7",
  "status": "active",
  "last_used_at": "",
  "created_at": "2026-02-24T17:00:00+08:00",
  "updated_at": "2026-02-24T17:00:00+08:00"
}
```

#### 查询用户的 API Key 列表

```bash
curl http://localhost:8080/api/v1/users/1/apikeys
```

#### 删除 API Key

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1/apikeys/1
```

#### 禁用 API Key

```bash
curl -X PUT http://localhost:8080/api/v1/users/1/apikeys/1/disable
```

### 模型管理

#### 查询可用模型列表

```bash
curl http://localhost:8080/api/v1/models
```

**响应示例：**

```json
[
  {
    "name": "qwen-turbo",
    "provider": "qwen"
  },
  {
    "name": "deepseek-chat",
    "provider": "deepseek"
  }
]
```

### 模型对话

#### 非流式对话

```bash
curl -X POST http://localhost:8080/api/v1/models/qwen-turbo/chat \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "你好，请介绍一下你自己"}
    ]
  }'
```

**响应示例（OpenAI 兼容格式）：**

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1771926519,
  "model": "qwen-turbo",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "你好！我是通义千问..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 13,
    "completion_tokens": 20,
    "total_tokens": 33
  }
}
```

#### 流式对话

```bash
curl -X POST http://localhost:8080/api/v1/models/qwen-turbo/chat \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "数到10"}
    ],
    "stream": true
  }'
```

**响应格式（Server-Sent Events）：**

```
data:{"id":"chatcmpl-123","object":"chat.completion.chunk",...}

data:{"id":"chatcmpl-123","object":"chat.completion.chunk",...}

data:[DONE]
```

#### 多轮对话

```bash
curl -X POST http://localhost:8080/api/v1/models/qwen-turbo/chat \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "我叫小明"},
      {"role": "assistant", "content": "你好小明！"},
      {"role": "user", "content": "我叫什么名字？"}
    ]
  }'
```

## 认证说明

所有模型对话 API 都需要通过 API Key 认证。在请求头中添加：

```
Authorization: Bearer ck_your_api_key_here
```

**认证错误响应：**

```json
{
  "error": "未提供认证信息"
}
```

或

```json
{
  "error": "API Key 无效"
}
```

## 数据库

系统使用 SQLite 数据库存储数据，默认文件为 `courier.db`。

**表结构：**

- `users` - 用户表
- `api_keys` - API Key 表
- `request_logs` - 请求日志表

**查看请求日志：**

```bash
sqlite3 courier.db "SELECT * FROM request_logs ORDER BY created_at DESC LIMIT 10;"
```

## 环境变量

支持在配置文件中使用环境变量：

```yaml
models:
  - name: qwen-turbo
    api_key: ${QWEN_API_KEY}
```

然后设置环境变量：

```bash
export QWEN_API_KEY=your_actual_api_key
./bin/server
```

## 开发

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
gofmt -w .
```

### 静态检查

```bash
go vet ./...
```

## License

MIT License
