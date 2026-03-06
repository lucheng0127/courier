<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

interface LoginForm {
  email: string
  password: string
}

const formRef = ref()
const loading = ref(false)
const formData = ref<LoginForm>({
  email: '',
  password: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '邮箱格式不正确' }
  ],
  password: [
    { required: true, message: '请输入密码' }
  ]
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    await authStore.login(formData.value.email, formData.value.password)
    message.success('登录成功')

    // 跳转到目标页面或 Dashboard
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (error: any) {
    if (error.errorFields) {
      // 表单验证错误
      console.log('Validation failed:', error)
    } else {
      // 登录失败（已在 request.ts 中处理错误提示）
      console.error('Login failed:', error)
    }
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>Courier LLM Gateway</h1>
        <p>登录到您的账户</p>
      </div>

      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <a-form-item label="邮箱" name="email">
          <a-input
            v-model:value="formData.email"
            placeholder="请输入邮箱"
            size="large"
          />
        </a-form-item>

        <a-form-item label="密码" name="password">
          <a-input-password
            v-model:value="formData.password"
            placeholder="请输入密码"
            size="large"
          />
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            block
            size="large"
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>

      <div class="login-footer">
        还没有账户？
        <a-button type="link" @click="goToRegister">立即注册</a-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #F9FAFB;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: #FFFFFF;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #6B7280;
}

.login-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
  color: #6B7280;
}
</style>
