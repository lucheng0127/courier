import request from '@/utils/request'
import type { ApiKey } from '@/types'

// 创建 API Key
export const createApiKey = (userId: number, data: { name: string }) => {
  return request<any, { id: number; key: string; key_prefix: string; name: string; status: string; created_at: string }>({
    url: `/users/${userId}/api-keys`,
    method: 'POST',
    data
  })
}

// 获取 API Key 列表
export const getApiKeys = (userId: number) => {
  return request<any, { api_keys: ApiKey[] }>({
    url: `/users/${userId}/api-keys`,
    method: 'GET'
  }).then(res => res.api_keys)
}

// 删除 API Key
export const deleteApiKey = (userId: number, keyId: number) => {
  return request<any, void>({
    url: `/users/${userId}/api-keys/${keyId}`,
    method: 'DELETE'
  })
}
