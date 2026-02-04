# AnyRouter 项目规范

## 项目概述

AnyRouter 签到中台管理系统，采用前后端分离架构：
- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue 3 + TypeScript + shadcn-vue + TailwindCSS

## 目录结构

```
project-root/
├── backend/              # Go 后端服务
│   ├── cmd/server/       # 程序入口
│   ├── internal/         # 内部业务逻辑
│   │   ├── handler/      # HTTP 处理器
│   │   ├── service/      # 业务逻辑层
│   │   ├── repository/   # 数据访问层
│   │   └── model/        # 数据模型
│   └── pkg/              # 可复用包
├── frontend/             # Vue 前端
│   └── src/
│       ├── api/          # API 请求
│       ├── components/   # 组件
│       ├── pages/        # 页面
│       └── stores/       # Pinia 状态
└── docs/                 # 技术文档
    ├── backend/技术规范.md
    └── frontend/技术规范.md
```

## 开发命令

**后端**:
```bash
cd backend
make dev          # 开发模式（热重载）
make build        # 编译
make gen-docs     # 生成 Swagger 文档
```

**前端**:
```bash
cd frontend
npm run dev       # 开发服务器
npm run build     # 生产构建
npm run type-check
```

## 核心规范

### 后端规范
- 分层架构：`Handler → Service → Repository → Model`
- 时间处理：**必须**使用 `dromara/carbon`，禁止 `time.Now()`
- 响应格式：`response.Success/Error/Unauthorized`
- 所有 Handler 必须添加 Swagger 注解
- 详细规范见：`docs/backend/技术规范.md`

### 前端规范
- 组件写法：Composition API + `<script setup lang="ts">`
- 类型安全：禁止使用 `any` 类型
- 样式：使用 TailwindCSS，禁止内联 style
- API 调用：从 `@/api` 导入
- 详细规范见：`docs/frontend/技术规范.md`

## 禁止事项

1. 硬编码敏感信息（密钥、密码）
2. 跳过分层架构直接操作数据库
3. 忽略错误处理
4. 提交未经测试的代码
5. 修改与需求无关的代码

## 工作流

### 上下文检索（编码前必做）

在修改任何代码前，**必须**先理解相关上下文。使用 `rg`（ripgrep）进行代码检索：

```bash
# 查找函数/结构体定义
rg "func \w+Account" backend/
rg "interface Account" frontend/src/

# 查找引用和调用
rg "useAuthStore" frontend/src/
rg "repository\.DB" backend/internal/

# 查找文件
find . -name "*.go" -path "*/handler/*"
find . -name "*.vue" -path "*/pages/*"
```

**原则**：禁止基于假设修改代码。先检索、再理解、再动手。

### 调用 Gemini 协作

项目内置 Gemini 桥接脚本，**前端/UI/CSS 任务强烈推荐使用**。

**调用方式**：
```bash
python .claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "$(pwd)" \
  --PROMPT "你的任务描述"
```

**前端 UI 审查**（获取 Unified Diff）：
```bash
python .claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "$(pwd)" \
  --PROMPT "Review frontend/src/pages/account/index.vue and suggest UI improvements. OUTPUT: Unified Diff Patch ONLY."
```

**多轮对话**（使用 SESSION_ID 继续）：
```bash
# 首次调用，从返回结果中获取 SESSION_ID
python .claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "$(pwd)" \
  --PROMPT "Analyze the dark mode theme in frontend/src/assets/main.css"

# 后续追问，携带 SESSION_ID
python .claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "$(pwd)" \
  --SESSION_ID "<上次返回的ID>" \
  --PROMPT "Improve the color contrast for cards"
```

**使用场景**：

| 场景 | 是否调用 Gemini |
|------|----------------|
| 前端 UI/CSS/样式调整 | **推荐** |
| Vue 组件结构设计 | **推荐** |
| 后端逻辑/算法 | 不推荐 |
| 简单 bug 修复 | 不需要 |

**关键约束**：
- Gemini 返回的代码仅作为原型参考，**必须重构后再使用**
- Prompt 中**必须**要求返回 `Unified Diff Patch`，禁止让 Gemini 直接修改文件
- Gemini 对后端逻辑的理解有缺陷，忽略其后端建议
