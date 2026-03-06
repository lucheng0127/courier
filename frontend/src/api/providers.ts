import request from '@/utils/request'
import type { ProviderInfo, ProviderForm, ModelInfo } from '@/types'

// 获取 Provider 列表
export const getProviders = () => {
  return request<any, { providers: any[] }>({
    url: '/providers',
    method: 'GET'
  }).then(res => {
    // 判断数据格式：管理员返回嵌套格式，普通用户返回扁平格式
    const isAdminFormat = res.providers.length > 0 && 'provider' in res.providers[0]

    if (isAdminFormat) {
      // 管理员格式：{ provider: {...}, is_running: boolean }
      return res.providers.map((p: any) => {
        const newProvider = { ...p.provider }
        // fallback_models 已经是数组（从后端 JSON 字段转换而来）
        return {
          provider: newProvider,
          is_running: p.is_running
        }
      })
    } else {
      // 普通用户格式：扁平的 { name, type, base_url, enabled, fallback_models }
      return res.providers.map((p: any) => ({
        provider: {
          id: 0,
          name: p.name,
          type: p.type,
          base_url: p.base_url,
          timeout: 30,
          enabled: p.enabled,
          fallback_models: p.fallback_models || [],
          created_at: ''
        },
        is_running: p.enabled // 普通用户看不到 is_running，假设与 enabled 相同
      }))
    }
  })
}

// 创建 Provider
export const createProvider = (data: ProviderForm) => {
  return request<any, ProviderInfo>({
    url: '/providers',
    method: 'POST',
    data
  })
}

// 更新 Provider
export const updateProvider = (name: string, data: Partial<ProviderForm>) => {
  return request<any, ProviderInfo>({
    url: `/providers/${name}`,
    method: 'PUT',
    data
  })
}

// 删除 Provider
export const deleteProvider = (name: string) => {
  return request<any, void>({
    url: `/providers/${name}`,
    method: 'DELETE'
  })
}

// 重载 Provider
export const reloadProvider = (name: string) => {
  return request<any, ProviderInfo>({
    url: `/admin/providers/${name}/reload`,
    method: 'POST'
  })
}

// 启用 Provider
export const enableProvider = (name: string) => {
  return request<any, ProviderInfo>({
    url: `/admin/providers/${name}/enable`,
    method: 'POST'
  })
}

// 禁用 Provider
export const disableProvider = (name: string) => {
  return request<any, ProviderInfo>({
    url: `/admin/providers/${name}/disable`,
    method: 'POST'
  })
}

// 获取 Provider 模型列表
export const getProviderModels = (name: string) => {
  return request<any, ModelInfo[]>({
    url: `/providers/${name}/models`,
    method: 'GET'
  })
}
