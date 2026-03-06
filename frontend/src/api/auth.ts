import request from '@/utils/request'
import type { RegisterRequest, LoginRequest, Token, User } from '@/types'

// 用户注册
export const register = (data: RegisterRequest) => {
  return request<any, User>({
    url: '/auth/register',
    method: 'POST',
    data
  })
}

// 用户登录
export const login = (data: LoginRequest) => {
  return request<any, Token>({
    url: '/auth/login',
    method: 'POST',
    data
  })
}

// 刷新 Token
export const refreshToken = (refreshToken: string) => {
  return request<any, Token>({
    url: '/auth/refresh',
    method: 'POST',
    data: { refresh_token: refreshToken }
  })
}

// 获取当前用户信息
export const getCurrentUser = (userId: number) => {
  return request<any, User>({
    url: `/users/${userId}`,
    method: 'GET'
  })
}
