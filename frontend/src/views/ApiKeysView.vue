<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined, CopyOutlined, CheckCircleOutlined, StopOutlined } from '@ant-design/icons-vue'
import { useApiKeyStore } from '@/stores/api-keys'
import { useAuthStore } from '@/stores/auth'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import type { ApiKey } from '@/types'

const apiKeyStore = useApiKeyStore()
const authStore = useAuthStore()

// 表格数据
const columns = [
  { title: '名称', dataIndex: 'name', key: 'name', width: 200 },
  { title: 'Key 前缀', dataIndex: 'key_prefix', key: 'key_prefix', width: 150 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '最后使用时间', dataIndex: 'last_used_at', key: 'last_used_at', width: 180 },
  { title: '操作', key: 'actions', width: 200, fixed: 'right' }
]

// 创建表单相关
const createModalVisible = ref(false)
const createFormRef = ref()
const createFormData = ref({
  name: ''
})
const createFormRules = {
  name: [{ required: true, message: '请输入 API Key 名称', trigger: 'blur' }]
}

// 显示完整 Key 的模态框
const showKeyModalVisible = ref(false)
const fullKey = ref('')
const copied = ref(false)

// 加载数据
const loadData = async () => {
  await apiKeyStore.fetchApiKeys(authStore.user!.id)
}

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

// 打开创建模态框
const openCreateModal = () => {
  createFormData.value = { name: '' }
  createModalVisible.value = true
}

// 提交创建表单
const handleSubmit = async () => {
  try {
    await createFormRef.value.validate()
    const result = await apiKeyStore.addApiKey(authStore.user!.id, { name: createFormData.value.name })
    createModalVisible.value = false
    fullKey.value = result.key
    showKeyModalVisible.value = true
    message.success('API Key 创建成功')
  } catch (error: any) {
    console.error('创建 API Key 失败:', error)
  }
}

// 启用 API Key
const handleEnable = async (record: ApiKey) => {
  try {
    await apiKeyStore.enableKey(authStore.user!.id, record.id)
    message.success('API Key 已启用')
  } catch (error: any) {
    message.error('启用失败: ' + (error.message || '未知错误'))
  }
}

// 禁用 API Key
const handleDisable = (record: ApiKey) => {
  Modal.confirm({
    title: '确认禁用',
    content: `确定要禁用 API Key "${record.name}" 吗？禁用后将无法使用此 Key 进行 API 调用。`,
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      try {
        await apiKeyStore.disableKey(authStore.user!.id, record.id)
        message.success('API Key 已禁用')
      } catch (error: any) {
        message.error('禁用失败: ' + (error.message || '未知错误'))
      }
    }
  })
}

// 删除 API Key
const handleDelete = (record: ApiKey) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除 API Key "${record.name}" 吗？删除后无法恢复。`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await apiKeyStore.removeApiKey(authStore.user!.id, record.id)
        message.success('API Key 已删除')
      } catch (error: any) {
        message.error('删除失败: ' + (error.message || '未知错误'))
      }
    }
  })
}

// 复制完整 Key
const copyFullKey = async () => {
  try {
    await navigator.clipboard.writeText(fullKey.value)
    copied.value = true
    message.success('已复制到剪贴板')
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (error) {
    message.error('复制失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="api-keys-layout">
    <Sidebar />
    <Topbar />

    <div class="main-content">
      <div class="page-header">
        <h2 class="page-title">API Keys 管理</h2>
        <a-button type="primary" @click="openCreateModal">
          <template #icon>
            <PlusOutlined />
          </template>
          创建 API Key
        </a-button>
      </div>

      <a-card class="table-card">
        <a-table
          :columns="columns"
          :data-source="apiKeyStore.apiKeys"
          :loading="apiKeyStore.loading"
          :pagination="{ pageSize: 10 }"
          :scroll="{ x: 1000 }"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="record.status === 'active' ? 'success' : 'default'">
                {{ record.status === 'active' ? '已启用' : record.status === 'disabled' ? '已禁用' : '已撤销' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatTime(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'last_used_at'">
              {{ formatTime(record.last_used_at!) }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-space>
                <!-- 启用按钮（仅在禁用状态显示） -->
                <a-button
                  v-if="record.status === 'disabled'"
                  type="link"
                  size="small"
                  @click="handleEnable(record)"
                >
                  <template #icon>
                    <CheckCircleOutlined />
                  </template>
                  启用
                </a-button>
                <!-- 禁用按钮（仅在启用状态显示） -->
                <a-button
                  v-if="record.status === 'active'"
                  type="link"
                  size="small"
                  @click="handleDisable(record)"
                >
                  <template #icon>
                    <StopOutlined />
                  </template>
                  禁用
                </a-button>
                <!-- 删除按钮（任何状态都显示） -->
                <a-button
                  type="link"
                  size="small"
                  danger
                  @click="handleDelete(record)"
                >
                  删除
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>

      <!-- 创建 API Key 表单模态框 -->
      <a-modal
        v-model:open="createModalVisible"
        title="创建 API Key"
        width="500px"
        @ok="handleSubmit"
        @cancel="createModalVisible = false"
      >
        <a-form
          ref="createFormRef"
          :model="createFormData"
          :rules="createFormRules"
          layout="vertical"
          style="margin-top: 24px"
        >
          <a-form-item label="名称" name="name">
            <a-input
              v-model:value="createFormData.name"
              placeholder="例如：生产环境 Key"
              :maxlength="50"
            />
            <div class="form-tip">请为您的 API Key 输入一个便于识别的名称</div>
          </a-form-item>
        </a-form>
      </a-modal>

      <!-- 显示完整 Key 模态框 -->
      <a-modal
        v-model:open="showKeyModalVisible"
        title="API Key 创建成功"
        :closable="false"
        :mask-closable="false"
        :footer="null"
        width="600px"
      >
        <div class="key-display">
          <div class="key-warning">
            <a-alert
              type="warning"
              message="请妥善保存您的 API Key"
              description="关闭此窗口后，您将无法再次查看完整的 API Key。请立即复制并保存在安全的地方。"
              show-icon
            />
          </div>
          <div class="key-content">
            <div class="key-label">完整的 API Key：</div>
            <div class="key-value">{{ fullKey }}</div>
            <a-button
              type="primary"
              @click="copyFullKey"
            >
              <template v-if="!copied" #icon>
                <CopyOutlined />
              </template>
              {{ copied ? '已复制' : '复制完整 Key' }}
            </a-button>
          </div>
        </div>
        <div class="key-modal-footer">
          <a-button @click="showKeyModalVisible = false">我已保存，关闭</a-button>
        </div>
      </a-modal>
    </div>
  </div>
</template>

<style scoped>
.api-keys-layout {
  min-height: 100vh;
  background-color: #F9FAFB;
}

.main-content {
  padding: 24px;
  margin-left: 240px;
  margin-top: 60px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.table-card {
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.form-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}

.key-display {
  padding: 16px 0;
}

.key-warning {
  margin-bottom: 24px;
}

.key-content {
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}

.key-label {
  font-size: 14px;
  color: #6B7280;
  margin-bottom: 12px;
}

.key-value {
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 14px;
  color: #10A37F;
  background: #FFFFFF;
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  padding: 12px 16px;
  word-break: break-all;
  margin-bottom: 16px;
  line-height: 1.6;
}

.key-modal-footer {
  margin-top: 24px;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .main-content {
    padding: 16px;
    margin-left: 0;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
}
</style>

