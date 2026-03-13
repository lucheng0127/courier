# 套餐配额和使用统计 - 设计文档

## 背景

平台需要实现套餐功能，允许管理员创建包含Provider Token配额的套餐产品，用户可以购买套餐并在有效期内使用配额。每次Chat API调用需要检查用户是否有足够的配额。

## 目标 / 非目标

### 目标
- 实现套餐产品管理（创建、上下架、查询）
- 实现用户套餐购买（含模拟支付）
- 实现Token配额检测中间件
- 实现套餐配额使用统计
- 支持套餐叠加（一个用户可购买多个套餐）
- 支持套餐有效期管理

### 非目标
- 真实支付集成（本次使用模拟支付）
- 套餐修改功能（套餐创建后不可修改）
- 套餐删除功能（只能下架）
- 复杂的计费策略
- 套餐推荐功能

## 决策

### 1. 数据模型设计

#### 1.1 Package（套餐产品）

```go
type Package struct {
    ID          int64     // 套餐ID
    Name        string    // 套餐名称
    Description string    // 套餐描述
    Price       int       // 价格（分为单位）
    Status      string    // 状态：draft（草稿）、online（上架）、offline（下架）
    ValidityDays int      // 有效期（天数）
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**设计理由**：
- 套餐创建后不可修改，确保已购买用户权益
- 状态管理：draft → online → offline，状态流转单向
- 价格使用分为单位，避免浮点数精度问题
- ValidityDays 支持灵活的有效期配置

#### 1.2 PackageQuota（套餐配额）

```go
type PackageQuota struct {
    ID          int64  // 配额ID
    PackageID   int64  // 套餐ID
    ProviderName string // Provider名称（支持*表示所有Provider）
    TokenLimit  int64  // Token配额限制
}
```

**设计理由**：
- 一个套餐可包含多个Provider的配额
- 支持 `*` 作为Provider名称，表示套餐适用于所有Provider
- TokenLimit 为0表示不限量

#### 1.3 UserPackage（用户套餐）

```go
type UserPackage struct {
    ID              int64      // 用户套餐ID
    UserID          int64      // 用户ID
    PackageID       int64      // 套餐ID
    Status          string     // 状态：active（激活）、expired（过期）
    ExpiresAt       time.Time  // 过期时间
    ActivatedAt     time.Time  // 激活时间
    CreatedAt       time.Time
}
```

**设计理由**：
- 记录用户购买套餐的激活时间
- ExpiresAt 计算方式：购买时间 + 套餐有效期天数
- 支持用户购买多个相同套餐（配额叠加）

#### 1.4 PackageUsage（套餐使用记录）

```go
type PackageUsage struct {
    ID               int64      // 使用记录ID
    UserPackageID    int64      // 用户套餐ID
    UserID           int64      // 用户ID
    PackageID        int64      // 套餐ID
    ProviderName     string     // 使用的Provider
    TokensUsed       int64      // 使用的Token数
    RequestID        string     // 关联的请求ID
    UsageRecordID    int64      // 关联的使用记录ID
    CreatedAt        time.Time
}
```

**设计理由**：
- 关联到 UserPackage，支持按套餐统计使用量
- 关联到 UsageRecord，保持数据一致性
- TokensUsed 记录实际扣减的Token数

### 2. 配额检测中间件设计

#### 2.1 检测流程

```
1. 从上下文获取 user_id 和 model 参数
2. 解析 model 获取 provider_name
3. 查询用户所有激活且未过期的套餐
4. 筛选包含该 provider 配额的套餐（或包含 * 的套餐）
5. 按过期时间升序排序（优先使用即将过期的套餐）
6. 检查是否有足够的剩余配额
7. 如果配额不足，返回 429 Too Many Requests
8. 将可用的套餐配额信息注入到上下文（供后续扣减使用）
```

#### 2.2 剩余配额计算

```
剩余配额 = 套餐配额限制 - 已使用配额

注意：
- 如果套餐配额限制为0，表示不限量
- 已使用配额从 package_usage 表聚合计算
- 优先使用即将过期的套餐配额
```

#### 2.3 配额扣减时机

- 方案1：请求前预扣，请求后调整（实现复杂，暂不采用）
- **方案2：请求后扣减**（本次采用）

请求后扣减流程：
1. Chat请求完成后，获取实际使用的token数
2. 从上下文获取之前检测时选择的套餐列表
3. 按优先级扣减配额
4. 记录 PackageUsage

### 3. 套餐购买流程（含模拟支付）

```
1. 用户查询可购买的套餐列表（status=online）
2. 用户选择套餐并发起购买请求
3. 系统检查套餐是否可购买（online状态）
4. 模拟支付成功（直接通过，不进行真实支付）
5. 创建 UserPackage 记录
6. 计算 ExpiresAt = 当前时间 + 套餐有效期天数
7. 返回购买成功
```

### 4. 套餐状态流转

```
draft（草稿）→ online（上架）→ offline（下架）
           ↑                        ↓
           └────────────────────────┘
           （管理员可以重新上架已下架的套餐）
```

**状态说明**：
- `draft`：创建套餐时的初始状态，不可购买
- `online`：上架状态，用户可以购买
- `offline`：下架状态，新用户无法购买，已购买的用户继续使用

### 5. API 路由设计

```
# 管理员接口
POST   /api/v1/admin/packages                    创建套餐
PUT    /api/v1/admin/packages/:id                更新套餐（仅限草稿状态）
POST   /api/v1/admin/packages/:id/online         上架套餐
POST   /api/v1/admin/packages/:id/offline        下架套餐
GET    /api/v1/admin/packages                    查询套餐列表
GET    /api/v1/admin/packages/:id                查询套餐详情
DELETE /api/v1/admin/packages/:id                删除套餐（仅限草稿状态）

# 用户接口
GET    /api/v1/packages                          查询可购买套餐列表
GET    /api/v1/packages/:id                      查询套餐详情
POST   /api/v1/user/packages                     购买套餐
GET    /api/v1/user/packages                     查询我的套餐
GET    /api/v1/user/packages/:id                 查询套餐详情
GET    /api/v1/user/packages/:id/usage           查询套餐使用统计
```

## 替代方案考虑

### 配额扣减时机

| 方案 | 优点 | 缺点 | 选择 |
|------|------|------|------|
| 请求前预扣 | 配额保证充足 | 实际使用可能与预估不符，需要回滚机制 | ❌ |
| 请求后扣减 | 实现简单，数据准确 | 可能短时间超量使用 | ✅ |

**选择理由**：请求后扣减实现简单，且Chat API的实际token使用量与请求的max_tokens参数差异通常不大，短暂超量在可接受范围内。

### 套餐修改

| 方案 | 优点 | 缺点 | 选择 |
|------|------|------|------|
| 允许修改 | 灵活 | 影响已购买用户权益 | ❌ |
| 不允许修改 | 保护用户权益 | 需要创建新套餐 | ✅ |

**选择理由**：套餐修改可能影响已购买用户的预期权益（如配额减少），不允许修改更符合业务逻辑。需要调整时创建新套餐即可。

## 风险 / 权衡

### 风险1：套餐配额用完后请求被拒绝

**影响**：用户体验下降
**缓解措施**：
- 提供配额即将用尽的预警
- 返回明确的错误信息，引导用户购买新套餐
- 支持套餐叠加，避免单一套餐配额不足

### 风险2：配额扣减失败导致数据不一致

**影响**：用户配额与实际使用不符
**缓解措施**：
- 使用数据库事务保证原子性
- 记录失败日志，支持人工对账
- 提供配额同步修复工具

### 风险3：大量并发请求下配额检测性能

**影响**：API响应变慢
**缓解措施**：
- 对用户套餐配额进行缓存（Redis）
- 缓存TTL设置为5分钟
- 定时任务同步缓存

### 风险4：过期套餐的配额占用存储空间

**影响**：数据库膨胀
**缓解措施**：
- 定期归档过期套餐的使用记录
- 提供数据清理工具

## 迁移计划

### 阶段1：数据库迁移
1. 创建新表（packages, package_quotas, user_packages, package_usage）
2. 运行数据库迁移脚本

### 阶段2：后端实现
1. 实现套餐管理功能
2. 实现套餐购买功能
3. 实现配额检测中间件
4. 实现套餐使用统计

### 阶段3：测试验证
1. 单元测试
2. 集成测试
3. 性能测试

### 阶段4：上线部署
1. 灰度发布
2. 监控观察
3. 全量发布

**回滚计划**：
- 移除配额检测中间件
- 保留数据表，不影响现有功能

## 开放问题

1. **套餐推荐**：是否需要根据用户使用情况推荐套餐？（暂不实现）
2. **套餐赠送**：管理员是否需要给用户赠送套餐的功能？（可后续添加）
3. **套餐转让**：用户是否可以将套餐转让给其他用户？（暂不支持）
4. **套餐退款**：购买后是否支持退款？（暂不支持）
