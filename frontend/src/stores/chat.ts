import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ChatMessage, ProviderModel } from '@/types'
import { sendChatStream, getProviderModels } from '@/api/chat'

export const useChatStore = defineStore('chat', () => {
  // State
  const messages = ref<ChatMessage[]>([])
  const loading = ref(false)
  const streamingContent = ref('')
  const selectedProvider = ref<string | null>(null)
  const selectedModel = ref<string | null>(null)
  const providerModels = ref<ProviderModel[]>([])

  // Computed
  const fullModelIdentifier = computed(() => {
    if (!selectedProvider.value || !selectedModel.value) {
      return ''
    }
    return `${selectedProvider.value}/${selectedModel.value}`
  })

  const hasMessages = computed(() => messages.value.length > 0)

  // Actions
  const selectProvider = async (providerName: string) => {
    selectedProvider.value = providerName
    selectedModel.value = null
    providerModels.value = []

    // 获取该 Provider 的模型列表
    try {
      const data = await getProviderModels(providerName)
      const models = data.models || data || []

      // 如果模型是字符串数组，转换为 ProviderModel 格式
      if (Array.isArray(models) && models.length > 0) {
        if (typeof models[0] === 'string') {
          providerModels.value = (models as string[]).map(m => ({
            name: m,
            id: m,
            enabled: true
          }))
        } else {
          providerModels.value = models as ProviderModel[]
        }
      }
    } catch (error) {
      console.error('Failed to fetch provider models:', error)
      providerModels.value = []
    }
  }

  const selectModel = (modelName: string) => {
    selectedModel.value = modelName
  }

  const addMessage = (message: ChatMessage) => {
    messages.value.push(message)
  }

  const updateLastAssistantMessage = (content: string) => {
    const lastMessage = messages.value[messages.value.length - 1]
    if (lastMessage && lastMessage.role === 'assistant') {
      lastMessage.content = content
    }
  }

  const sendMessage = async (apiKey: string, userMessage: string) => {
    if (!fullModelIdentifier.value) {
      throw new Error('Please select provider and model')
    }

    loading.value = true
    streamingContent.value = ''

    // 添加用户消息
    const userMsg: ChatMessage = { role: 'user', content: userMessage }
    messages.value.push(userMsg)

    // 添加空的助手消息
    const assistantMsg: ChatMessage = { role: 'assistant', content: '' }
    messages.value.push(assistantMsg)

    try {
      const stream = await sendChatStream(apiKey, {
        model: fullModelIdentifier.value,
        messages: messages.value.slice(0, -1), // 不包含刚添加的空消息
        temperature: 0.7,
        max_tokens: 2000
      })

      const reader = stream.getReader()
      const decoder = new TextDecoder()

      while (true) {
        const { done, value } = await reader.read()

        if (done) break

        const chunk = decoder.decode(value)
        const lines = chunk.split('\n')

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)

            if (data === '[DONE]') {
              break
            }

            try {
              const parsed = JSON.parse(data)
              const content = parsed.choices?.[0]?.delta?.content

              if (content) {
                streamingContent.value += content
                updateLastAssistantMessage(streamingContent.value)
              }
            } catch (e) {
              // 忽略解析错误
            }
          }
        }
      }
    } catch (error: any) {
      // 发生错误，更新最后一条消息为错误提示
      updateLastAssistantMessage(`错误：${error.message || '发送消息失败'}`)
      throw error
    } finally {
      loading.value = false
      streamingContent.value = ''
    }
  }

  const clearMessages = () => {
    messages.value = []
  }

  const resetSelection = () => {
    selectedProvider.value = null
    selectedModel.value = null
    providerModels.value = []
  }

  return {
    // State
    messages,
    loading,
    streamingContent,
    selectedProvider,
    selectedModel,
    providerModels,
    // Computed
    fullModelIdentifier,
    hasMessages,
    // Actions
    selectProvider,
    selectModel,
    addMessage,
    sendMessage,
    clearMessages,
    resetSelection
  }
})
