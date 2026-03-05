<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { message } from 'ant-design-vue'
import { SendOutlined, RobotOutlined, UserOutlined } from '@ant-design/icons-vue'
import { useChatStore } from '@/stores/chat'
import { useApiKeyStore } from '@/stores/api-keys'
import { useAuthStore } from '@/stores/auth'
import { getProviders } from '@/api/providers'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import type { ProviderInfo } from '@/types'

const chatStore = useChatStore()
const apiKeyStore = useApiKeyStore()
const authStore = useAuthStore()

// Provider 和 Model 列表
const providers = ref<ProviderInfo[]>([])
const loadingProviders = ref(false)
const hasApiKey = ref(false)
const loadingApiKeyCheck = ref(true)

// 输入相关
const inputMessage = ref('')
const messagesContainerRef = ref<HTMLElement>()

// 检查 API Key
const checkApiKey = async () => {
  loadingApiKeyCheck.value = true
  try {
    await apiKeyStore.fetchApiKeys(authStore.user!.id)
    hasApiKey.value = apiKeyStore.hasActiveApiKey()
  } catch (error) {
    console.error('Failed to check API keys:', error)
  } finally {
    loadingApiKeyCheck.value = false
  }
}

// 加载 Provider 列表
const loadProviders = async () => {
  loadingProviders.value = true
  try {
    providers.value = await getProviders()
    // 只显示已启用的 Provider
    providers.value = providers.value.filter(p => p.provider.enabled)
  } catch (error) {
    console.error('Failed to load providers:', error)
  } finally {
    loadingProviders.value = false
  }
}

// 处理发送消息
const handleSend = async () => {
  const content = inputMessage.value.trim()
  if (!content) {
    return
  }

  if (!chatStore.fullModelIdentifier) {
    message.warning('请先选择 Provider 和 Model')
    return
  }

  if (!hasApiKey.value) {
    message.warning('请先创建 API Key')
    return
  }

  // 获取第一个可用的 API Key
  const activeApiKey = apiKeyStore.getFirstActiveApiKey()
  if (!activeApiKey) {
    message.error('没有可用的 API Key')
    return
  }

  // 尝试从 localStorage 获取完整的 API Key
  const fullApiKey = localStorage.getItem(`api_key_${activeApiKey.id}`)

  if (!fullApiKey) {
    message.warning({
      content: '此 API Key 是在更新前创建的，无法用于聊天。请删除旧 Key 并重新创建一个新的 API Key。',
      duration: 5
    })
    return
  }

  inputMessage.value = ''

  try {
    await chatStore.sendMessage(fullApiKey, content)

    // 滚动到底部
    await nextTick()
    scrollToBottom()
  } catch (error: any) {
    message.error('发送失败：' + (error.message || '未知错误'))
  }
}

// 处理键盘事件
const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainerRef.value) {
    messagesContainerRef.value.scrollTop = messagesContainerRef.value.scrollHeight
  }
}

// 监听消息变化自动滚动
watch(() => chatStore.messages, async () => {
  await nextTick()
  scrollToBottom()
}, { deep: true })

// 监听 Provider 选择变化
const handleProviderChange = (value: string) => {
  chatStore.selectProvider(value)
}

// 监听 Model 选择变化
const handleModelChange = (value: string) => {
  chatStore.selectModel(value)
}

// 跳转到 API Keys 页面
const goToApiKeys = () => {
  window.location.href = '/api-keys'
}

// 清空对话
const handleClear = () => {
  chatStore.clearMessages()
}

onMounted(async () => {
  await checkApiKey()
  await loadProviders()
})
</script>

<template>
  <div class="chat-layout">
    <Sidebar />
    <Topbar />

    <div class="main-content">
      <!-- 无 API Key 提示 -->
      <div v-if="!loadingApiKeyCheck && !hasApiKey" class="no-api-key-container">
        <a-card class="no-api-key-card">
          <div class="no-api-key-content">
            <div class="no-api-key-icon">
              <RobotOutlined :style="{ fontSize: '64px', color: '#9CA3AF' }" />
            </div>
            <h2 class="no-api-key-title">您还没有创建 API Key</h2>
            <p class="no-api-key-description">请先创建 API Key 再进行对话</p>
            <a-button type="primary" size="large" @click="goToApiKeys">
              创建 API Key
            </a-button>
          </div>
        </a-card>
      </div>

      <!-- 聊天主界面 -->
      <div v-else class="chat-container">
        <!-- 顶部选择区域 -->
        <div class="chat-header">
          <div class="selector-group">
            <div class="selector-item">
              <label class="selector-label">Provider:</label>
              <a-select
                v-model:value="chatStore.selectedProvider"
                placeholder="选择 Provider"
                :loading="loadingProviders"
                style="width: 200px"
                @change="handleProviderChange"
              >
                <a-select-option
                  v-for="provider in providers"
                  :key="provider.provider.name"
                  :value="provider.provider.name"
                >
                  {{ provider.provider.name }} ({{ provider.provider.type }})
                </a-select-option>
              </a-select>
            </div>

            <div class="selector-item">
              <label class="selector-label">Model:</label>
              <a-select
                v-model:value="chatStore.selectedModel"
                placeholder="选择 Model"
                :disabled="!chatStore.selectedProvider"
                style="width: 200px"
                @change="handleModelChange"
              >
                <a-select-option
                  v-for="model in chatStore.providerModels"
                  :key="model.name || model.id"
                  :value="model.name || model.id"
                >
                  {{ model.name || model.id }}
                </a-select-option>
              </a-select>
            </div>

            <div v-if="chatStore.fullModelIdentifier" class="model-display">
              <a-tag color="green">{{ chatStore.fullModelIdentifier }}</a-tag>
            </div>

            <a-button
              v-if="chatStore.hasMessages"
              type="text"
              @click="handleClear"
            >
              清空对话
            </a-button>
          </div>
        </div>

        <!-- 消息区域 -->
        <div ref="messagesContainerRef" class="messages-container">
          <div v-if="!chatStore.hasMessages" class="empty-state">
            <RobotOutlined :style="{ fontSize: '48px', color: '#D1D5DB' }" />
            <p>选择 Provider 和 Model 开始对话</p>
          </div>

          <div v-else class="messages-list">
            <div
              v-for="(msg, index) in chatStore.messages"
              :key="index"
              :class="['message-item', msg.role === 'user' ? 'user-message' : 'assistant-message']"
            >
              <div class="message-avatar">
                <UserOutlined v-if="msg.role === 'user'" />
                <RobotOutlined v-else />
              </div>
              <div class="message-content">
                <div class="message-role">{{ msg.role === 'user' ? '用户' : 'AI' }}</div>
                <div class="message-text">{{ msg.content }}</div>
                <div v-if="msg.role === 'assistant' && chatStore.loading && index === chatStore.messages.length - 1" class="typing-indicator">
                  <span>•</span><span>•</span><span>•</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-area">
          <div class="input-container">
            <a-textarea
              ref="inputRef"
              v-model:value="inputMessage"
              placeholder="输入消息... (Enter 发送，Shift+Enter 换行)"
              :disabled="chatStore.loading || !chatStore.fullModelIdentifier"
              :auto-size="{ minRows: 1, maxRows: 6 }"
              @keydown="handleKeyDown"
            />
            <a-button
              type="primary"
              :loading="chatStore.loading"
              :disabled="!inputMessage.trim() || !chatStore.fullModelIdentifier"
              @click="handleSend"
            >
              <template #icon>
                <SendOutlined />
              </template>
              发送
            </a-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-layout {
  min-height: 100vh;
  background-color: #F9FAFB;
}

.main-content {
  height: 100vh;
  padding: 24px;
  margin-left: 240px;
  margin-top: 60px;
  display: flex;
  flex-direction: column;
}

/* 无 API Key 提示 */
.no-api-key-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.no-api-key-card {
  border-radius: 12px;
  text-align: center;
  max-width: 400px;
}

.no-api-key-content {
  padding: 24px;
}

.no-api-key-icon {
  margin-bottom: 24px;
}

.no-api-key-title {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12px;
}

.no-api-key-description {
  font-size: 14px;
  color: #6B7280;
  margin-bottom: 24px;
}

/* 聊天容器 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #FFFFFF;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

/* 聊天头部 */
.chat-header {
  padding: 16px 24px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.selector-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.selector-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.selector-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  white-space: nowrap;
}

.model-display {
  margin-left: auto;
}

/* 消息区域 */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #9CA3AF;
}

.empty-state p {
  margin-top: 16px;
  font-size: 14px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 12px;
}

.user-message {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
}

.user-message .message-avatar {
  background: #10A37F;
  color: #FFFFFF;
}

.assistant-message .message-avatar {
  background: #F3F4F6;
  color: #6B7280;
}

.message-content {
  max-width: 70%;
}

.message-role {
  font-size: 12px;
  color: #9CA3AF;
  margin-bottom: 4px;
}

.user-message .message-role {
  text-align: right;
}

.message-text {
  padding: 12px 16px;
  border-radius: 12px;
  line-height: 1.5;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.user-message .message-text {
  background: #10A37F;
  color: #FFFFFF;
  border-bottom-right-radius: 4px;
}

.assistant-message .message-text {
  background: #F3F4F6;
  color: #111827;
  border-bottom-left-radius: 4px;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  margin-top: 8px;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #D1D5DB;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-4px);
  }
}

/* 输入区域 */
.input-area {
  padding: 16px 24px;
  border-top: 1px solid #E5E7EB;
  background: #FFFFFF;
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.input-container :deep(.ant-input-textarea) {
  flex: 1;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .main-content {
    padding: 16px;
    margin-left: 0;
    margin-top: 60px;
  }

  .selector-group {
    flex-direction: column;
    align-items: stretch;
  }

  .selector-item {
    flex-direction: column;
    align-items: stretch;
  }

  .model-display {
    margin-left: 0;
  }

  .message-content {
    max-width: 85%;
  }
}
</style>
