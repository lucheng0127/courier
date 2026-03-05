import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ApiKey } from '@/types'
import {
  getApiKeys,
  createApiKey,
  deleteApiKey
} from '@/api/api-keys'

export const useApiKeyStore = defineStore('apiKeys', () => {
  // State
  const apiKeys = ref<ApiKey[]>([])
  const loading = ref(false)

  // Actions
  const fetchApiKeys = async (userId: number) => {
    loading.value = true
    try {
      apiKeys.value = await getApiKeys(userId)
    } finally {
      loading.value = false
    }
  }

  const addApiKey = async (userId: number, data: { name: string }) => {
    const result = await createApiKey(userId, data)
    // 先刷新列表以获取完整的 API Key 信息（包括 id）
    await fetchApiKeys(userId)
    // 找到刚创建的 API Key（通过 key_prefix 匹配）
    const newKey = apiKeys.value.find(k => k.key_prefix === result.key_prefix)
    if (newKey && result.key) {
      // 将完整的 key 保存到 localStorage，供聊天功能使用
      localStorage.setItem(`api_key_${newKey.id}`, result.key)
    }
    return result
  }

  const removeApiKey = async (userId: number, keyId: number) => {
    await deleteApiKey(userId, keyId)
    // 清除 localStorage 中的完整 key
    localStorage.removeItem(`api_key_${keyId}`)
    await fetchApiKeys(userId)
  }

  // 获取第一个可用的 API Key
  const getFirstActiveApiKey = (): ApiKey | null => {
    return apiKeys.value.find(k => k.status === 'active') || null
  }

  // 检查是否有可用的 API Key
  const hasActiveApiKey = (): boolean => {
    return apiKeys.value.some(k => k.status === 'active')
  }

  return {
    // State
    apiKeys,
    loading,
    // Actions
    fetchApiKeys,
    addApiKey,
    removeApiKey,
    getFirstActiveApiKey,
    hasActiveApiKey
  }
})
