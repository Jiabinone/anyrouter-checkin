# 贡献指南

感谢你对 AnyRouter 签到系统的关注与贡献。

## 开发环境

- Go 1.23+
- Node.js 20+
- SQLite 3

## 项目结构

- 后端：`backend/`
- 前端：`frontend/`
- 文档：`docs/`

## 本地开发

后端：

```bash
cd backend
make dev
```

前端：

```bash
cd frontend
npm ci
npm run dev
```

## 代码规范

- 后端遵循 `Handler → Service → Repository → Model` 分层
- 时间处理必须使用 `dromara/carbon`
- 前端使用 `<script setup lang="ts">`，禁止 `any` 类型
- 样式使用 TailwindCSS，禁止内联 style
- API 调用统一从 `@/api` 导入

详细规范参考：
- `docs/backend/技术规范.md`
- `docs/frontend/技术规范.md`

## 提交流程

1. 确保代码可编译、无警告错误
2. 运行检查命令

后端：

```bash
cd backend
make lint
make build
```

前端：

```bash
cd frontend
npm run lint
npm run type-check
npm run build
```

## 提交规范（建议）

- 建议使用：`type(scope): description`
- 示例：`feat(push): 支持 HTML 模板`

## 提交 PR

- 说明变更背景与原因
- 关联相关 issue（如有）
- 附上测试结果或截图（前端改动）
