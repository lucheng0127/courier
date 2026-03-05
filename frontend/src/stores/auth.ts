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

    // 从 JWT token 中解析用户信息
    const payload = parseJWT(tokenData.access_token)
    console.log('[DEBUG] JWT Payload:', payload)

    // 设置用户信息
    const emailName = email.split('@')[0] || email
    user.value = {
      id: payload.user_id || 0,
      name: payload.user_email ? payload.user_email.split('@')[0] : emailName,
      email: payload.user_email || email,
      role: payload.user_role || 'user',
      status: 'active',
      created_at: new Date().toISOString()
    }
    console.log('[DEBUG] User after login:', user.value)
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  // 解析 JWT token
  const parseJWT = (token: string) => {
    try {
      const parts = token.split('.')
      if (parts.length < 2) return {}

      const base64Url = parts[1]
      if (!base64Url) return {}

      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(atob(base64).split('').map((c) => {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
      }).join(''))
      return JSON.parse(jsonPayload)
    } catch (e) {
      console.error('Failed to parse JWT:', e)
      return {}
    }
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

    console.log('[DEBUG] restoreState - savedToken:', !!savedToken)
    console.log('[DEBUG] restoreState - savedUser:', savedUser)

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
        console.log('[DEBUG] restoreState - parsed user:', user.value)
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
