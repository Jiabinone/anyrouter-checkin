# 前端开发规范

开发 Vue 前端代码时，**必须**遵循以下规范。完整文档：`docs/frontend/技术规范.md`

## 技术栈

| 组件 | 技术 |
|------|------|
| 框架 | Vue 3 (Composition API) |
| 语言 | TypeScript |
| 构建 | Vite |
| UI 组件 | shadcn-vue / reka-ui |
| 样式 | TailwindCSS |
| 状态管理 | Pinia |
| 路由 | Vue Router 4 |
| HTTP | Axios（封装后使用） |

## 组件编写规范

```vue
<script setup lang="ts">
// 1. Vue 导入
import { ref, computed, onMounted } from 'vue'
// 2. 外部库导入
import { toast } from 'vue-sonner'
// 3. 内部模块导入
import { Button } from '@/components/ui/button'
import { getAccounts, type Account } from '@/api/account'

// 4. Props/Emits
const props = defineProps<{ title: string }>()
const emit = defineEmits<{ submit: [data: FormData] }>()

// 5. 响应式状态
const loading = ref(false)
const accounts = ref<Account[]>([])

// 6. 计算属性
const isEmpty = computed(() => accounts.value.length === 0)

// 7. 方法
async function loadData() {
  loading.value = true
  try {
    accounts.value = await getAccounts()
  } finally {
    loading.value = false
  }
}

// 8. 生命周期
onMounted(loadData)
</script>

<template>
  <!-- 模板 -->
</template>
```

## 类型安全

```typescript
// ❌ 禁止
const data: any = {}

// ✅ 正确
interface Account {
  id: number
  name: string
  status: number
}
const accounts = ref<Account[]>([])
const result = await getAccounts() as Account[]
```

## 样式规范

```vue
<!-- ❌ 禁止内联 style -->
<div style="margin-top: 16px; color: red;">

<!-- ✅ 使用 TailwindCSS -->
<div class="mt-4 text-destructive">

<!-- 条件类名使用 cn() -->
<div :class="cn('rounded-md', isActive && 'bg-primary')">
```

## API 调用

```typescript
// 从 @/api 导入，禁止直接使用 axios
import { getAccounts, createAccount } from '@/api/account'

// API 函数定义示例
export function getAccounts(): Promise<Account[]> {
  return request.get('/accounts')
}
```

## Toast 通知

```typescript
import { toast } from 'vue-sonner'

toast.success('操作成功')
toast.error('操作失败: ' + error.message)
toast.loading('处理中...')
```

## UI 组件导入

```typescript
// shadcn-vue 组件
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'

// 图标
import { Plus, Trash2, Settings } from 'lucide-vue-next'
```

## 禁止事项

1. 使用 `any` 类型
2. 使用 Options API（必须 Composition API）
3. 直接操作 DOM
4. 模板中使用复杂表达式（提取为 computed）
5. 使用内联 style
6. 硬编码 API 地址
7. 直接调用 axios（使用封装的 request）
8. 使用 `// @ts-ignore` 忽略类型错误
9. 使用 `var` 声明变量
10. 组件超过 300 行（应拆分）

## 开发命令

```bash
npm run dev        # 开发服务器
npm run build      # 生产构建
npm run type-check # 类型检查
npm run lint       # 代码检查
```
