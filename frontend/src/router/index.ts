import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  // 恢复认证状态
  if (!authStore.isAuthenticated) {
    authStore.restoreState()
  }

  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !authStore.isAuthenticated) {
    // 未登录访问受保护路由，跳转到登录页
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (!requiresAuth && authStore.isAuthenticated && (to.name === 'Login' || to.name === 'Register')) {
    // 已登录访问登录/注册页，跳转到 Dashboard
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
