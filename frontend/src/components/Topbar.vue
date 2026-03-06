<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from 'ant-design-vue'
import { LogoutOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="topbar">
    <div class="topbar-spacer"></div>
    <div class="topbar-right">
      <span class="user-email">{{ authStore.userEmail }}</span>
      <Button type="text" @click="handleLogout">
        <template #icon>
          <LogoutOutlined />
        </template>
        注销
      </Button>
    </div>
  </div>
</template>

<style scoped>
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  background: #FFFFFF;
  border-bottom: 1px solid #E5E7EB;
  position: fixed;
  top: 0;
  right: 0;
  left: 240px;
  z-index: 10;
}

.topbar-spacer {
  width: 24px;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-right: 24px;
}

.user-email {
  color: #6B7280;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .user-email {
    display: none;
  }

  .topbar-right {
    padding-right: 16px;
  }
}
</style>
