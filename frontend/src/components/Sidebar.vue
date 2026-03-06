<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Menu } from 'ant-design-vue'
import {
  DashboardOutlined,
  BlockOutlined,
  CloudServerOutlined,
  ApiOutlined,
  MessageOutlined,
  SettingOutlined
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

interface MenuItem {
  key: string
  icon: any
  label: string
  path: string
  requiredRole?: 'admin' | 'user'
}

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const menuItems: MenuItem[] = [
  {
    key: 'dashboard',
    icon: DashboardOutlined,
    label: 'Dashboard',
    path: '/'
  },
  {
    key: 'models',
    icon: BlockOutlined,
    label: 'Models',
    path: '/models'
  },
  {
    key: 'providers',
    icon: CloudServerOutlined,
    label: 'Providers',
    path: '/providers',
    requiredRole: 'admin'
  },
  {
    key: 'api-keys',
    icon: ApiOutlined,
    label: 'API Keys',
    path: '/api-keys'
  },
  {
    key: 'chat',
    icon: MessageOutlined,
    label: 'Chat',
    path: '/chat'
  },
  {
    key: 'settings',
    icon: SettingOutlined,
    label: 'Settings',
    path: '/settings'
  }
]

// 根据用户角色过滤菜单项
const filteredMenuItems = computed(() => {
  const userRole = authStore.userRole
  console.log('[DEBUG] authStore.userRole:', userRole)
  console.log('[DEBUG] authStore.user:', authStore.user)

  const filtered = menuItems.filter(item => {
    if (!item.requiredRole) return true
    return userRole === item.requiredRole
  })
  console.log('[DEBUG] Filtered menu items:', filtered.map(f => ({ key: f.key, label: f.label, requiredRole: f.requiredRole })))
  return filtered
})

const selectedKeys = computed(() => {
  const path = route.path
  if (path === '/') return ['dashboard']
  if (path.startsWith('/models')) return ['models']
  if (path.startsWith('/providers')) return ['providers']
  if (path.startsWith('/api-keys')) return ['api-keys']
  if (path.startsWith('/chat')) return ['chat']
  if (path.startsWith('/settings')) return ['settings']
  return []
})

const handleMenuClick = (key: string) => {
  const menuItem = menuItems.find(item => item.key === key)
  if (menuItem) {
    router.push(menuItem.path)
  }
}
</script>

<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <h1 class="sidebar-logo">Courier</h1>
    </div>
    <Menu
      v-model:selectedKeys="selectedKeys"
      mode="inline"
      class="sidebar-menu"
    >
      <Menu.Item
        v-for="item in filteredMenuItems"
        :key="item.key"
        @click="() => handleMenuClick(item.key)"
      >
        <template #icon>
          <component :is="item.icon" />
        </template>
        {{ item.label }}
      </Menu.Item>
    </Menu>
  </div>
</template>

<style scoped>
.sidebar {
  width: 240px;
  height: 100vh;
  background: #FFFFFF;
  border-right: 1px solid #E5E7EB;
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  border-bottom: 1px solid #E5E7EB;
}

.sidebar-logo {
  font-size: 20px;
  font-weight: 700;
  color: #10A37F;
  margin: 0;
}

.sidebar-menu {
  flex: 1;
  border: none;
  padding: 16px 0;
}

.sidebar-menu :deep(.ant-menu-item) {
  margin: 4px 12px;
  height: 40px;
  line-height: 40px;
  border-radius: 8px;
  color: #6B7280;
}

.sidebar-menu :deep(.ant-menu-item:hover) {
  background-color: #F3F4F6;
  color: #111827;
}

.sidebar-menu :deep(.ant-menu-item-selected) {
  background-color: #ECFDF5;
  color: #10A37F;
}

.sidebar-menu :deep(.ant-menu-item-selected .anticon) {
  color: #10A37F;
}

.sidebar-menu :deep(.anticon) {
  font-size: 16px;
}
</style>
