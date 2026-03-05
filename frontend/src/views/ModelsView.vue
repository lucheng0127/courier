<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, CopyOutlined } from '@ant-design/icons-vue'
import { useProviderStore } from '@/stores/providers'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'

const providerStore = useProviderStore()

// 按 Provider 分组的模型数据
interface ProviderModels {
  name: string
  type: string
  enabled: boolean
  is_running: boolean
  models: string[]
}

const providerModelsList = computed<ProviderModels[]>(() => {
  return providerStore.providers.map(p => ({
    name: p.provider.name,
    type: p.provider.type,
    enabled: p.provider.enabled,
    is_running: p.is_running,
    models: p.provider.fallback_models
  }))
})

const loading = ref(false)

// 加载所有 Provider 的模型列表
const loadAllModels = async () => {
  loading.value = true
  try {
    await providerStore.fetchProviders()
  } finally {
    loading.value = false
  }
}

// 复制模型标识
const copyModelId = (providerName: string, modelName: string) => {
  const fullId = `${providerName}/${modelName}`
  navigator.clipboard.writeText(fullId)
  message.success('已复制: ' + fullId)
}

// 展开/折叠状态
const activeKeys = ref<string[]>([])

onMounted(() => {
  loadAllModels()
})
</script>

<template>
  <div class="models-layout">
    <Sidebar />
    <Topbar />

    <div class="main-content">
      <div class="page-header">
        <h2 class="page-title">模型列表</h2>
        <a-button @click="loadAllModels">
          <template #icon>
            <ReloadOutlined />
          </template>
          刷新
        </a-button>
      </div>

      <a-spin :spinning="loading">
        <!-- 空状态 -->
        <div v-if="providerModelsList.length === 0 && !loading" class="empty-state">
          <a-empty description="暂无可用模型">
            <p class="empty-hint">请联系管理员添加 Provider</p>
          </a-empty>
        </div>

        <!-- Provider 分组折叠面板 -->
        <a-collapse
          v-else
          v-model:active-key="activeKeys"
          class="models-collapse"
          ghost
        >
          <a-collapse-panel
            v-for="item in providerModelsList"
            :key="item.name"
            class="provider-panel"
          >
            <template #header>
              <div class="panel-header">
                <div class="provider-info">
                  <span class="provider-name">{{ item.name }}</span>
                  <a-tag size="small" color="blue">{{ item.type }}</a-tag>
                  <a-tag
                    size="small"
                    :color="item.enabled ? 'success' : 'default'"
                  >
                    {{ item.enabled ? '已启用' : '已禁用' }}
                  </a-tag>
                  <a-tag
                    size="small"
                    :color="item.is_running ? 'success' : 'error'"
                  >
                    {{ item.is_running ? '运行中' : '已停止' }}
                  </a-tag>
                </div>
                <span class="model-count">{{ item.models.length }} 个模型</span>
              </div>
            </template>

            <!-- Provider 下的模型为空 -->
            <div v-if="item.models.length === 0" class="panel-empty">
              暂无模型
            </div>

            <!-- 模型列表 -->
            <div v-else class="models-grid">
              <div
                v-for="model in item.models"
                :key="model"
                class="model-card"
              >
                <div class="model-header">
                  <div class="model-name">{{ model }}</div>
                </div>
                <div class="model-footer">
                  <div class="full-id">
                    {{ item.name }}/{{ model }}
                  </div>
                  <a-button
                    type="text"
                    size="small"
                    @click="copyModelId(item.name, model)"
                  >
                    <template #icon>
                      <CopyOutlined />
                    </template>
                    复制
                  </a-button>
                </div>
              </div>
            </div>
          </a-collapse-panel>
        </a-collapse>
      </a-spin>
    </div>
  </div>
</template>

<style scoped>
.models-layout {
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

.models-collapse {
  background: transparent;
}

.provider-panel {
  background: #FFFFFF;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  margin-bottom: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.provider-panel :deep(.ant-collapse-header) {
  padding: 16px 24px;
  border-radius: 12px !important;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding-right: 24px;
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.provider-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.model-count {
  font-size: 14px;
  color: #6B7280;
}

.models-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  padding: 16px 24px 24px;
}

.model-card {
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  padding: 16px;
  transition: all 0.2s;
}

.model-card:hover {
  border-color: #10A37F;
  box-shadow: 0 2px 8px rgba(16, 163, 127, 0.1);
}

.model-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.model-name {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.model-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #E5E7EB;
}

.full-id {
  font-size: 12px;
  color: #10A37F;
  font-family: monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin-right: 8px;
}

.panel-empty {
  padding: 24px;
  text-align: center;
  color: #9CA3AF;
}

.empty-state {
  background: #FFFFFF;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  padding: 60px 24px;
  text-align: center;
}

.empty-hint {
  color: #6B7280;
  margin-top: 16px;
  margin-bottom: 0;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .main-content {
    padding: 16px;
    margin-left: 0;
  }

  .models-grid {
    grid-template-columns: 1fr;
    padding: 16px;
  }

  .panel-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .model-card {
    padding: 12px;
  }

  .model-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .full-id {
    margin-right: 0;
    margin-bottom: 4px;
  }
}
</style>
