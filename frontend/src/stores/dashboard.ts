import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DashboardStats, Provider, UsageRecord } from '@/types'
import { getUsageStats, getProviders } from '@/api/dashboard'

export const useDashboardStore = defineStore('dashboard', () => {
  // State
  const stats = ref<DashboardStats>({
    totalRequests: 0,
    totalTokens: 0,
    successRate: 0,
    avgLatency: 0
  })
  const providers = ref<Array<Provider & { is_running: boolean }>>([])
  const recentActivity = ref<UsageRecord[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  const fetchDashboardData = async () => {
    loading.value = true
    error.value = null
    try {
      // 获取使用统计
      const usageData = await getUsageStats({ page: 1, page_size: 100 })

      // 计算统计数据
      const totalRequests = usageData.total || 0
      const successRecords = usageData.records.filter(r => r.status === 'success')
      const successRate = totalRequests > 0
        ? (successRecords.length / totalRequests) * 100
        : 0

      const totalTokens = usageData.records.reduce((sum, r) => sum + r.total_tokens, 0)
      const avgLatency = usageData.records.length > 0
        ? usageData.records.reduce((sum, r) => sum + r.latency_ms, 0) / usageData.records.length
        : 0

      stats.value = {
        totalRequests,
        totalTokens,
        successRate,
        avgLatency
      }

      // 获取最近活动（前 10 条）
      recentActivity.value = usageData.records.slice(0, 10)

      // 获取 Provider 状态
      const providersData = await getProviders()
      providers.value = providersData.providers.map(p => ({
        ...p.provider,
        is_running: p.is_running
      }))
    } catch (err: any) {
      error.value = err.message || '加载数据失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const refreshDashboardData = () => {
    return fetchDashboardData()
  }

  return {
    // State
    stats,
    providers,
    recentActivity,
    loading,
    error,
    // Actions
    fetchDashboardData,
    refreshDashboardData
  }
})
