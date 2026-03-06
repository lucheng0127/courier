import type { ChatRequest } from '@/types'

// 获取 API Base URL
// 使用相对路径，由 Vite proxy (本地开发) 或 Nginx proxy (Docker) 处理
const getApiBaseUrl = () => {
  const envApiUrl = import.meta.env.VITE_API_BASE_URL
  if (envApiUrl && (envApiUrl.startsWith('http://') || envApiUrl.startsWith('https://'))) {
    return envApiUrl
  }
  // 默认使用相对路径（空字符串）
  return ''
}

// 发送流式 Chat 请求
// 支持两种认证方式：JWT Token 或 API Key
export const sendChatStream = async (
  authCredential: string, // 可以是 JWT Token 或 API Key
  data: Omit<ChatRequest, 'stream'>
): Promise<ReadableStream> => {
  const apiBaseUrl = getApiBaseUrl()
  const url = `${apiBaseUrl}/v1/chat/completions`

  console.log('[Chat API] Request URL:', url)

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authCredential}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ ...data, stream: true })
    })

    if (!response.ok) {
      let errorMessage = 'Chat request failed'
      try {
        const errorText = await response.text()
        if (errorText) {
          // 尝试解析 JSON 错误响应
          try {
            const errorJson = JSON.parse(errorText)
            errorMessage = errorJson.error?.message || errorJson.message || errorText
          } catch {
            errorMessage = errorText
          }
        }
      } catch {
        // 如果无法读取错误响应，使用状态码
        errorMessage = `请求失败，状态码：${response.status}`
      }
      throw new Error(errorMessage)
    }

    const body = response.body
    if (!body) {
      throw new Error('响应体为空')
    }

    return body
  } catch (error: any) {
    // 网络错误或其他 fetch 错误
    if (error.name === 'TypeError' && error.message.includes('fetch')) {
      throw new Error('网络连接失败，请检查网络或服务器状态')
    }
    throw error
  }
}

// 获取 Provider 的模型列表
export const getProviderModels = async (providerName: string) => {
  const apiBaseUrl = getApiBaseUrl()
  const token = localStorage.getItem('access_token')
  const url = `${apiBaseUrl}/api/v1/providers/${providerName}/models`

  console.log('[Provider API] Request URL:', url)

  const response = await fetch(url, {
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
