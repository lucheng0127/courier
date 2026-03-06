// 用户类型
export interface User {
  id: number
  name: string
  email: string
  role: 'admin' | 'user'
  status: 'active' | 'disabled'
  created_at: string
}

// Token 类型
export interface Token {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

// API Key 类型
export interface ApiKey {
  id: number
  key_prefix: string
  name: string
  status: 'active' | 'disabled' | 'revoked'
  created_at: string
  last_used_at?: string
}

// Provider 类型（后端 ProviderInfo 结构）
export interface ProviderInfo {
  provider: Provider
  is_running: boolean
}

// Provider 类型（数据库结构）
export interface Provider {
  id: number
  name: string
  type: string
  base_url: string
  timeout: number
  enabled: boolean
  fallback_models: string[]
  created_at: string
  updated_at?: string
  api_key?: string
}

// Provider 表单类型
export interface ProviderForm {
  name: string
  type: string
  base_url: string
  timeout: number
  api_key?: string
  enabled: boolean
  fallback_models: string[]
}

// Model 信息类型
export interface ModelInfo {
  name: string
  id: string
  enabled: boolean
}

// 使用记录类型
export interface UsageRecord {
  id: number
  user_id: number
  model: string
  provider_name: string
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
  latency_ms: number
  status: 'success' | 'error'
  timestamp: string
}

// Dashboard 统计类型
export interface DashboardStats {
  totalRequests: number
  totalTokens: number
  successRate: number
  avgLatency: number
}

// 注册请求类型
export interface RegisterRequest {
  name: string
  email: string
  password: string
}

// 登录请求类型
export interface LoginRequest {
  email: string
  password: string
}

// 使用统计查询参数
export interface UsageQueryParams {
  user_id?: number
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

// 使用统计响应类型
export interface UsageResponse {
  records: UsageRecord[]
  total: number
  page: number
  page_size: number
}

// Chat 消息类型
export interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
}

// Chat 请求类型
export interface ChatRequest {
  model: string
  messages: ChatMessage[]
  temperature?: number
  max_tokens?: number
  stream: boolean
}

// Chat 响应类型
export interface ChatResponse {
  id: string
  object: string
  created: number
  model: string
  choices: Array<{
    index: number
    message: ChatMessage
    finish_reason: string
  }>
  usage: {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
  }
}

// Chat 流式响应类型
export interface ChatStreamChunk {
  id: string
  object: string
  created: number
  model: string
  choices: Array<{
    index: number
    delta: {
      content?: string
      role?: string
    }
  }>
}

// Provider 模型类型（用于获取 Provider 的模型列表）
export interface ProviderModel {
  name: string
  id: string
  enabled: boolean
}
