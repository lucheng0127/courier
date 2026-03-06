<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  EditOutlined,
  CheckCircleOutlined,
  StopOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import type { ProviderInfo, ProviderForm } from '@/types'
import { useProviderStore } from '@/stores/providers'
import { useAuthStore } from '@/stores/auth'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'

const providerStore = useProviderStore()
const authStore = useAuthStore()

// 权限检查
if (authStore.userRole !== 'admin') {
  message.error('权限不足')
  // 路由守卫会处理跳转
}

// 计算属性：从 ProviderInfo 中提取 Provider 数据用于显示
const providersForDisplay = computed(() => {
  return providerStore.providers.map((p) => ({
    ...p.provider,
    is_running: p.is_running,
    _providerInfo: p  // 保存原始 ProviderInfo 引用
  }))
})

// 表格数据
const columns = [
  { title: '名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', key: 'type', width: 120 },
  { title: 'Base URL', dataIndex: 'base_url', key: 'base_url', ellipsis: true },
  { title: '状态', dataIndex: 'enabled', key: 'enabled', width: 100 },
  { title: '运行状态', dataIndex: 'is_running', key: 'is_running', width: 120 },
  { title: 'Fallback 模型', dataIndex: 'fallback_models', key: 'fallback_models', width: 200 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'actions', width: 250, fixed: 'right' }
]

// 表单相关
const formModalVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formRef = ref()
const formData = ref<ProviderForm>({
  name: '',
  type: 'openai',
  base_url: '',
  timeout: 30,
  api_key: '',
  enabled: true,
  fallback_models: []
})

const providerTypeOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Azure OpenAI', value: 'azure_openai' },
  { label: 'Custom', value: 'custom' }
]

// 删除模型抽屉相关，不再需要

// 加载数据
const loadData = async () => {
  await providerStore.fetchProviders()
}

// 格式化时间
const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

// 打开创建表单
const openCreateModal = () => {
  formMode.value = 'create'
  formData.value = {
    name: '',
    type: 'openai',
    base_url: '',
    timeout: 30,
    api_key: '',
    enabled: true,
    fallback_models: []
  }
  formModalVisible.value = true
}

// 打开编辑表单
const openEditModal = (record: any) => {
  formMode.value = 'edit'
  formData.value = {
    name: record.name,
    type: record.type,
    base_url: record.base_url,
    timeout: record.timeout,
    api_key: record.api_key || '',
    enabled: record.enabled,
    fallback_models: [...record.fallback_models]
  }
  formModalVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    if (formMode.value === 'create') {
      await providerStore.addProvider(formData.value)
      message.success('创建成功')
    } else {
      await providerStore.editProvider(formData.value.name, formData.value)
      message.success('更新成功')
    }
    formModalVisible.value = false
  } catch (error: any) {
    if (error.errorFields) {
      message.error('请检查表单填写')
    } else {
      message.error('操作失败: ' + (error.message || '未知错误'))
    }
  }
}

// 启用/禁用 Provider
const handleToggleEnabled = async (record: ProviderInfo) => {
  try {
    if (record.provider.enabled) {
      await providerStore.disableProviderAction(record.provider.name)
      message.success('已禁用')
    } else {
      await providerStore.enableProviderAction(record.provider.name)
      message.success('已启用')
    }
  } catch (error: any) {
    message.error('操作失败: ' + (error.message || '未知错误'))
  }
}

// 删除 Provider
const handleDelete = (record: ProviderInfo) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除 Provider "${record.provider.name}" 吗？此操作不可撤销。`,
    okText: '确定',
    cancelText: '取消',
    onOk: async () => {
      try {
        await providerStore.removeProvider(record.provider.name)
        message.success('删除成功')
      } catch (error: any) {
        message.error('删除失败: ' + (error.message || '未知错误'))
      }
    }
  })
}

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入 Provider 名称', trigger: 'blur' },
    { pattern: /^[a-z0-9_-]+$/, message: '名称只能包含小写字母、数字、下划线和连字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择 Provider 类型', trigger: 'change' }
  ],
  base_url: [
    { required: true, message: '请输入 Base URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的 URL', trigger: 'blur' }
  ],
  timeout: [
    { required: true, message: '请输入超时时间', trigger: 'blur' },
    { type: 'number', min: 1, max: 300, message: '超时时间必须在 1-300 秒之间', trigger: 'blur' }
  ],
  fallback_models: [
    { required: true, message: '至少需要一个 Fallback 模型', trigger: 'change', type: 'array' }
  ]
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="providers-layout">
    <Sidebar />
    <Topbar />

    <div class="main-content">
      <div class="page-header">
        <h2 class="page-title">Provider 管理</h2>
        <a-button type="primary" @click="openCreateModal">
          <template #icon>
            <PlusOutlined />
          </template>
          新增provider
        </a-button>
      </div>

      <a-card class="table-card">
        <a-table
          :columns="columns"
          :data-source="providersForDisplay"
          :loading="providerStore.loading"
          :pagination="{ pageSize: 10 }"
          :scroll="{ x: 1200 }"
          row-key="name"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'enabled'">
              <a-tag :color="record.enabled ? 'success' : 'default'">
                {{ record.enabled ? '已启用' : '已禁用' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'is_running'">
              <a-tag :color="record.is_running ? 'success' : 'error'">
                {{ record.is_running ? '运行中' : '已停止' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'fallback_models'">
              <a-tag v-for="model in record.fallback_models" :key="model" color="blue">
                {{ model }}
              </a-tag>
              <span v-if="record.fallback_models.length === 0" class="text-gray">-</span>
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatTime(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-space>
                <a-button
                  type="link"
                  size="small"
                  @click="openEditModal(record)"
                >
                  <template #icon>
                    <EditOutlined />
                  </template>
                  编辑
                </a-button>
                <a-button
                  type="link"
                  size="small"
                  @click="handleToggleEnabled(record._providerInfo)"
                >
                  <template #icon>
                    <StopOutlined v-if="record.enabled" />
                    <CheckCircleOutlined v-else />
                  </template>
                  {{ record.enabled ? '禁用' : '启用' }}
                </a-button>
                <a-button
                  type="link"
                  size="small"
                  danger
                  @click="handleDelete(record._providerInfo)"
                >
                  <template #icon>
                    <DeleteOutlined />
                  </template>
                  删除
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-card>

      <!-- 创建/编辑表单模态框 -->
      <a-modal
        v-model:open="formModalVisible"
        :title="formMode === 'create' ? '创建 Provider' : '编辑 Provider'"
        width="600px"
        @ok="handleSubmit"
        @cancel="formModalVisible = false"
      >
        <a-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          layout="vertical"
          style="margin-top: 24px"
        >
          <a-form-item label="名称" name="name">
            <a-input
              v-model:value="formData.name"
              placeholder="例如: openai-primary"
              :disabled="formMode === 'edit'"
            />
          </a-form-item>

          <a-form-item label="类型" name="type">
            <a-select v-model:value="formData.type" placeholder="选择 Provider 类型">
              <a-select-option
                v-for="option in providerTypeOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="Base URL" name="base_url">
            <a-input
              v-model:value="formData.base_url"
              placeholder="https://api.openai.com/v1"
            />
          </a-form-item>

          <a-form-item label="超时时间 (秒)" name="timeout">
            <a-input-number
              v-model:value="formData.timeout"
              :min="1"
              :max="300"
              style="width: 100%"
            />
          </a-form-item>

          <a-form-item label="API Key" name="api_key">
            <a-input-password
              v-model:value="formData.api_key"
              placeholder="sk-..."
            />
          </a-form-item>

          <a-form-item label="Fallback 模型" name="fallback_models">
            <a-select
              v-model:value="formData.fallback_models"
              mode="tags"
              placeholder="输入模型名称并回车"
              :token-separators="[',']"
            >
            </a-select>
          </a-form-item>

          <a-form-item label="状态">
            <a-checkbox v-model:checked="formData.enabled">
              启用此 Provider
            </a-checkbox>
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </div>
</template>

<style scoped>
.providers-layout {
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

.text-gray {
  color: #9CA3AF;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #9CA3AF;
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
