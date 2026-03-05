import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, Token } from '@/types'
import { register as registerApi, login as loginApi, refreshToken as refreshTokenApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<Token | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const userEmail = computed(() => user.value?.email || '')
  const userName = computed(() => user.value?.name || '')
  const userRole = computed(() => user.value?.role || 'user')

  // Actions
  const setToken = (newToken: Token) => {
    token.value = newToken
    localStorage.setItem('access_token', newToken.access_token)
    localStorage.setItem('refresh_token', newToken.refresh_token)
  }

  const register = async (name: string, email: string, password: string) => {
    await registerApi({ name, email, password })
    // 注册成功后不自动登录，需要用户手动登录
    // 不设置 user 和 token 状态
  }

  const login = async (email: string, password: string) => {
    const tokenData = await loginApi({ email, password })
    setToken(tokenData)
    // 设置临时用户信息（因为 API 不返回用户信息）
    // 后续可以通过其他接口获取完整用户信息
    const emailName = email.split('@')[0] || email
    user.value = {
      id: 0, // 临时 ID
      name: emailName, // 从邮箱提取用户名
      email: email,
      role: 'user', // 默认角色
      status: 'active',
      created_at: new Date().toISOString()
    }
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  }

  const refreshToken = async () => {
    if (!token.value?.refresh_token) {
      throw new Error('No refresh token available')
    }
    const tokenData = await refreshTokenApi(token.value.refresh_token)
    setToken(tokenData)
  }

  // 从 localStorage 恢复状态
  const restoreState = () => {
    const savedToken = localStorage.getItem('access_token')
    const savedRefreshToken = localStorage.getItem('refresh_token')
    const savedUser = localStorage.getItem('user')

    if (savedToken && savedRefreshToken) {
      token.value = {
        access_token: savedToken,
        refresh_token: savedRefreshToken,
        token_type: 'Bearer',
        expires_in: 900
      }
    }

    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (e) {
        console.error('Failed to parse saved user:', e)
      }
    }
  }

  return {
    // State
    user,
    token,
    // Getters
    isAuthenticated,
    userEmail,
    userName,
    userRole,
    // Actions
    register,
    login,
    logout,
    refreshToken,
    restoreState
  }
})
