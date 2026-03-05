# UI Style Guide

## 设计风格

Modern SaaS Console

参考：

- Vercel Dashboard
- Stripe Console
- OpenAI Platform

---

## 颜色

### 主题色

| 名称 | Hex | 用途 |
|------|-----|------|
| Primary | `#10A37F` | 主按钮、链接、激活状态 |
| Primary Hover | `#0D8A6A` | 主按钮悬停 |
| Secondary | `#6366F1` | 次要按钮、标签 |

### 中性色

| 名称 | Hex | 用途 |
|------|-----|------|
| Background | `#F9FAFB` | 页面背景 |
| Surface | `#FFFFFF` | 卡片、模态框背景 |
| Border | `#E5E7EB` | 边框、分割线 |
| Text Primary | `#111827` | 主要文本 |
| Text Secondary | `#6B7280` | 次要文本 |
| Text Tertiary | `#9CA3AF` | 辅助文本 |

### 状态色

| 名称 | Hex | 用途 |
|------|-----|------|
| Success | `#10B981` | 成功状态、已启用 |
| Warning | `#F59E0B` | 警告状态 |
| Error | `#EF4444` | 错误状态、已禁用 |
| Info | `#3B82F6` | 信息提示 |

---

## 布局

### 整体布局

```
┌─────────────────────────────────────────────────────────────┐
│ Topbar (60px)                                               │
├──────────┬──────────────────────────────────────────────────┤
│          │                                                  │
│ Sidebar  │ Main Content                                     │
│ (240px)  │                                                  │
│          │                                                  │
│          │                                                  │
└──────────┴──────────────────────────────────────────────────┘
```

### Topbar

- **高度**: 60px
- **背景**: Surface (#FFFFFF)
- **边框**: 底部 1px Border
- **左侧**: Logo + 产品名称
- **右侧**: 用户菜单、通知

### Sidebar

- **宽度**: 240px
- **背景**: Surface (#FFFFFF)
- **边框**: 右侧 1px Border
- **导航项**: 图标 + 文字

### Main Content

- **背景**: Background (#F9FAFB)
- **内边距**: 24px
- **最大宽度**: 1400px
- **居中显示**

---

## 组件风格

### 按钮

| 类型 | 背景 | 文字 | 边框 | 高度 | 内边距 |
|------|------|------|------|------|--------|
| Primary | Primary | Surface | 无 | 36px | 0 16px |
| Secondary | Surface | Text Primary | Border | 36px | 0 16px |
| Ghost | 透明 | Text Primary | 无 | 36px | 0 16px |
| Danger | Error | Surface | 无 | 36px | 0 16px |

圆角: 6px

字体: 500 14px

### 输入框

- **高度**: 36px
- **边框**: 1px Border
- **圆角**: 6px
- **内边距**: 0 12px
- **字体**: 14px

Focus 状态:

- 边框色: Primary
- 阴影: 0 0 0 3px rgba(16, 163, 127, 0.1)

### 表格

- **边框**: 1px Border
- **圆角**: 8px
- **背景**: Surface
- **行高**: 48px

表头:

- 背景: Background
- 字体: 600 13px
- 文字色: Text Secondary
- 内边距: 0 16px

表格行:

- 字体: 14px
- 文字色: Text Primary
- 内边距: 0 16px
- 边框: 底部 1px Border

悬停:

- 背景: Background

### 卡片

- **背景**: Surface
- **边框**: 1px Border
- **圆角**: 12px
- **内边距**: 24px
- **阴影**: 0 1px 3px rgba(0,0,0,0.05)

### 标签

- **高度**: 24px
- **内边距**: 0 8px
- **圆角**: 4px
- **字体**: 500 12px

状态标签:

- Success: Success 背景 + Success 文字
- Error: Error 背景 + Error 文字
- Warning: Warning 背景 + Warning 文字

### 模态框

- **背景**: Surface
- **圆角**: 12px
- **阴影**: 0 20px 25px -5px rgba(0,0,0,0.1)
- **最大宽度**: 560px
- **内边距**: 24px

### 抽屉

- **宽度**: 480px
- **背景**: Surface
- **内边距**: 24px

---

## 响应式

### 断点

| 名称 | 宽度 | 说明 |
|------|------|------|
| Mobile | < 640px | 移动端 |
| Tablet | 640px - 1024px | 平板 |
| Desktop | > 1024px | 桌面端 |

### 移动端适配

- Sidebar 隐藏，使用汉堡菜单
- 表格横向滚动
- 卡片单列布局

---

## 字体

### 字体族

```
-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif
```

### 字号

| 用途 | 字号 | 字重 |
|------|------|------|
| H1 | 24px | 600 |
| H2 | 20px | 600 |
| H3 | 16px | 600 |
| Body | 14px | 400 |
| Small | 12px | 400 |
| Caption | 11px | 400 |

---

## 间距

### 间距系统

基于 4px 网格:

| 名称 | 值 |
|------|-----|
| xs | 4px |
| sm | 8px |
| md | 16px |
| lg | 24px |
| xl | 32px |
| 2xl | 48px |

---

## 图标

使用 Heroicons 图标库

- 尺寸: 20px (默认)
- 颜色: 继承文本颜色

---

## 动画

### 过渡时间

| 名称 | 时间 |
|------|------|
| Fast | 150ms |
| Base | 200ms |
| Slow | 300ms |

### 缓动函数

```
cubic-bezier(0.4, 0, 0.2, 1)
```

---

## 暗色模式

暂不支持，仅提供浅色主题。
