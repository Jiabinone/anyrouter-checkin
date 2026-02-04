---
globs: ["frontend/**/*.vue", "frontend/**/*.ts", "frontend/**/*.tsx"]
---

# 前端开发规范

开发 Vue 前端代码时，**必须**遵循 `docs/frontend/技术规范.md` 中定义的所有规范。

## 核心要点速查

### 技术栈
- Vue 3 (Composition API + `<script setup>`)
- TypeScript（严格类型检查）
- Vite / shadcn-vue / TailwindCSS
- Pinia / Vue Router / Axios

### 组件编写顺序
```vue
<script setup lang="ts">
// 1. Vue 导入
// 2. 外部库导入
// 3. 内部模块导入
// 4. Props/Emits 定义
// 5. 响应式状态
// 6. 计算属性
// 7. 方法
// 8. 生命周期
</script>

<template>
  <!-- 模板 -->
</template>
```

### 类型安全
```typescript
// ❌ 禁止
const data: any = {}

// ✅ 正确
interface Account {
  id: number
  name: string
}
const accounts = ref<Account[]>([])
```

### 样式规范
```vue
<!-- ❌ 禁止内联 style -->
<div style="margin-top: 16px;">

<!-- ✅ 使用 TailwindCSS -->
<div class="mt-4">

<!-- 条件类名使用 cn() -->
<div :class="cn('rounded-md', isActive && 'bg-primary')">
```

### API 调用
```typescript
// 从 @/api 导入，禁止直接使用 axios
import { getAccounts } from '@/api/account'
const accounts = await getAccounts()
```

### Toast 通知
```typescript
import { toast } from 'vue-sonner'
toast.success('操作成功')
toast.error('操作失败')
```

### 禁止事项
1. 使用 `any` 类型
2. 使用 Options API
3. 直接操作 DOM
4. 模板中使用复杂表达式
5. 使用内联 style
6. 硬编码 API 地址
7. 直接调用 axios
8. 使用 `// @ts-ignore`
9. 使用 `var` 声明
10. 组件超过 300 行
