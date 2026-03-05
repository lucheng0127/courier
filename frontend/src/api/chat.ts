import type { ChatRequest } from '@/types'

// 发送流式 Chat 请求
export const sendChatStream = async (
  apiKey: string,
  data: Omit<ChatRequest, 'stream'>
): Promise<ReadableStream> => {
  const response = await fetch('/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${apiKey}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ ...data, stream: true })
  })

  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || 'Chat request failed')
  }

  return response.body!
}

// 获取 Provider 的模型列表
export const getProviderModels = async (providerName: string) => {
  const token = localStorage.getItem('access_token')
  const response = await fetch(`/api/v1/providers/${providerName}/models`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  })

  if (!response.ok) {
    throw new Error('Failed to fetch provider models')
  }

  return response.json()
}
