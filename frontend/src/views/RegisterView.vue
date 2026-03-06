<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

interface RegisterForm {
  name: string
  email: string
  password: string
  confirmPassword: string
}

const formRef = ref()
const loading = ref(false)
const formData = ref<RegisterForm>({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  name: [
    { required: true, message: '请输入姓名' }
  ],
  email: [
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '邮箱格式不正确' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 8, message: '密码至少需要 8 个字符' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    {
      validator: (_rule: any, value: string) => {
        if (value !== formData.value.password) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    await authStore.register(formData.value.name, formData.value.email, formData.value.password)
    message.success('注册成功，请登录')
    router.push('/login')
  } catch (error: any) {
    if (error.errorFields) {
      // 表单验证错误
      console.log('Validation failed:', error)
    } else {
      // 注册失败（已在 request.ts 中处理错误提示）
      console.error('Register failed:', error)
    }
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h1>Courier LLM Gateway</h1>
        <p>创建新账户</p>
      </div>

      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <a-form-item label="姓名" name="name">
          <a-input
            v-model:value="formData.name"
            placeholder="请输入姓名"
            size="large"
          />
        </a-form-item>

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
            placeholder="请输入密码（至少 8 个字符）"
            size="large"
          />
        </a-form-item>

        <a-form-item label="确认密码" name="confirmPassword">
          <a-input-password
            v-model:value="formData.confirmPassword"
            placeholder="请再次输入密码"
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
            注册
          </a-button>
        </a-form-item>
      </a-form>

      <div class="register-footer">
        已有账户？
        <a-button type="link" @click="goToLogin">立即登录</a-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #F9FAFB;
}

.register-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: #FFFFFF;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 8px;
}

.register-header p {
  font-size: 14px;
  color: #6B7280;
}

.register-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
  color: #6B7280;
}
</style>
