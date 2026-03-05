<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, Row, Col, Statistic, Table, Tag, Button, Spin } from 'ant-design-vue'
import {
  ReloadOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined
} from '@ant-design/icons-vue'
import { useDashboardStore } from '@/stores/dashboard'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'

const dashboardStore = useDashboardStore()

const loading = ref(false)

const columns = [
  {
    title: '时间',
    dataIndex: 'timestamp',
    key: 'timestamp',
    width: 180,
    customRender: ({ text }: { text: string }) => {
      return new Date(text).toLocaleString('zh-CN')
    }
  },
  {
    title: '模型',
    dataIndex: 'model',
    key: 'model',
    width: 200
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: 'Token 数',
    dataIndex: 'total_tokens',
    key: 'total_tokens',
    width: 120,
    customRender: ({ text }: { text: number }) => {
      return text.toLocaleString()
    }
  }
]

const loadData = async () => {
  loading.value = true
  try {
    await dashboardStore.fetchDashboardData()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})

const formatNumber = (num: number) => {
  return num.toLocaleString()
}

const formatPercent = (num: number) => {
  return num.toFixed(1) + '%'
}

const formatLatency = (ms: number) => {
  return ms.toFixed(0) + 'ms'
}
</script>

<template>
  <div class="dashboard-layout">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Topbar -->
    <Topbar />

    <!-- Main Content -->
    <div class="main-content">
      <Spin :spinning="loading">
        <!-- Statistics Cards -->
        <Row :gutter="[16, 16]" class="stats-row">
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic
                title="总请求数"
                :value="formatNumber(dashboardStore.stats.totalRequests)"
                :value-style="{ color: '#10A37F' }"
              />
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic
                title="总 Token 数"
                :value="formatNumber(dashboardStore.stats.totalTokens)"
                :value-style="{ color: '#6366F1' }"
              />
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic
                title="成功率"
                :value="formatPercent(dashboardStore.stats.successRate)"
                suffix="%"
                :value-style="{ color: '#10B981' }"
              >
                <template #prefix>
                  <ArrowUpOutlined v-if="dashboardStore.stats.successRate > 80" />
                  <ArrowDownOutlined v-else />
                </template>
              </Statistic>
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic
                title="平均延迟"
                :value="formatLatency(dashboardStore.stats.avgLatency)"
                :value-style="{ color: '#F59E0B' }"
              />
            </Card>
          </Col>
        </Row>

        <!-- Provider Status Card -->
        <Card title="Provider 状态" class="provider-card">
          <div v-if="dashboardStore.providers.length === 0" class="empty-state">
            暂无 Provider
          </div>
          <div v-else class="provider-list">
            <div
              v-for="provider in dashboardStore.providers"
              :key="provider.id"
              class="provider-item"
            >
              <div class="provider-info">
                <span class="provider-name">{{ provider.name }}</span>
                <span class="provider-type">({{ provider.type }})</span>
              </div>
              <Tag :color="provider.is_running ? 'success' : 'error'">
                {{ provider.is_running ? '运行中' : '已停止' }}
              </Tag>
            </div>
          </div>
        </Card>

        <!-- Recent Activity Table -->
        <Card title="最近活动" class="activity-card">
          <template #extra>
            <Button
              type="link"
              :icon="ReloadOutlined"
              @click="loadData"
            >
              刷新
            </Button>
          </template>

          <div v-if="dashboardStore.recentActivity.length === 0" class="empty-state">
            暂无活动记录
          </div>
          <Table
            v-else
            :columns="columns"
            :data-source="dashboardStore.recentActivity"
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <Tag :color="record.status === 'success' ? 'success' : 'error'">
                  {{ record.status === 'success' ? '成功' : '失败' }}
                </Tag>
              </template>
            </template>
          </Table>
        </Card>
      </Spin>
    </div>
  </div>
</template>

<style scoped>
.dashboard-layout {
  min-height: 100vh;
  background-color: #F9FAFB;
}

.main-content {
  padding: 24px;
  max-width: 1400px;
  margin-left: 240px;
  margin-top: 60px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.provider-card,
.activity-card {
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  margin-bottom: 16px;
}

.provider-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.provider-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 8px;
}

.provider-info {
  display: flex;
  gap: 8px;
  align-items: center;
}

.provider-name {
  font-weight: 500;
  color: #111827;
}

.provider-type {
  color: #6B7280;
  font-size: 12px;
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

  .topbar {
    left: 0;
  }
}
</style>
