import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ProviderInfo, ProviderForm } from '@/types'
import {
  getProviders,
  createProvider,
  updateProvider,
  deleteProvider,
  reloadProvider,
  enableProvider,
  disableProvider
} from '@/api/providers'

export const useProviderStore = defineStore('providers', () => {
  // State
  const providers = ref<ProviderInfo[]>([])
  const loading = ref(false)

  // Actions
  const fetchProviders = async () => {
    loading.value = true
    try {
      providers.value = await getProviders()
    } finally {
      loading.value = false
    }
  }

  const addProvider = async (data: ProviderForm) => {
    await createProvider(data)
    // 重载 provider
    await reloadProvider(data.name)
    // 重新获取列表
    await fetchProviders()
  }

  const editProvider = async (name: string, data: Partial<ProviderForm>) => {
    await updateProvider(name, data)
    // 重载 provider
    await reloadProvider(name)
    // 重新获取列表
    await fetchProviders()
  }

  const removeProvider = async (name: string) => {
    await deleteProvider(name)
    await fetchProviders()
  }

  const enableProviderAction = async (name: string) => {
    await enableProvider(name)
    // 重新获取列表以确保数据格式一致
    await fetchProviders()
  }

  const disableProviderAction = async (name: string) => {
    await disableProvider(name)
    // 重新获取列表以确保数据格式一致
    await fetchProviders()
  }

  return {
    // State
    providers,
    loading,
    // Actions
    fetchProviders,
    addProvider,
    editProvider,
    removeProvider,
    enableProviderAction,
    disableProviderAction
  }
})
