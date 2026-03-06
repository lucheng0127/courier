import request from '@/utils/request'
import type { UsageQueryParams, UsageResponse, Provider } from '@/types'

// 获取使用统计
export const getUsageStats = (params: UsageQueryParams) => {
  return request<any, UsageResponse>({
    url: '/usage',
    method: 'GET',
    params
  })
}

// 获取 Provider 列表
export const getProviders = () => {
  return request<any, { providers: Array<{ provider: Provider; is_running: boolean }> }>({
    url: '/providers',
    method: 'GET'
  })
}
