import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

// 获取 API Base URL
// 使用相对路径，由 Vite proxy (本地开发) 或 Nginx proxy (Docker) 处理
const getBaseURL = () => {
  // 如果环境变量设置了完整的 URL，使用它（用于特殊情况）
  const envApiUrl = import.meta.env.VITE_API_BASE_URL
  if (envApiUrl && (envApiUrl.startsWith('http://') || envApiUrl.startsWith('https://'))) {
    return `${envApiUrl}/api/v1`
  }

  // 默认使用相对路径
  return '/api/v1'
}

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token.access_token}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    return response.data
  },
  async (error: AxiosError) => {
    const authStore = useAuthStore()
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    // 处理 401 未授权错误
    if (error.response?.status === 401 && !originalRequest._retry) {
      // 如果有 refresh_token，尝试刷新
      if (authStore.token?.refresh_token) {
        originalRequest._retry = true
        try {
          await authStore.refreshToken()
          // 重试原请求
          if (authStore.token) {
            originalRequest.headers.Authorization = `Bearer ${authStore.token.access_token}`
          }
          return request(originalRequest)
        } catch (refreshError) {
          // 刷新失败，清除认证信息并跳转登录
          authStore.logout()
          router.push('/login')
          message.error('会话已过期，请重新登录')
          return Promise.reject(refreshError)
        }
      } else {
        // 没有 refresh_token，直接跳转登录
        authStore.logout()
        router.push('/login')
        message.error('请先登录')
      }
    }

    // 处理其他错误
    let errorMessage = '请求失败'
    if (error.response) {
      const data = error.response.data as any
      errorMessage = data?.message || data?.error?.message || errorMessage
    } else if (error.request) {
      errorMessage = '网络错误，请检查网络连接'
    } else {
      errorMessage = error.message || '未知错误'
    }

    message.error(errorMessage)
    return Promise.reject(error)
  }
)

export default request
