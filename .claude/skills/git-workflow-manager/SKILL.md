---
name: git-workflow-manager
description: Manage Git workflows with Q9JY project commit standards and daily commit limits. Use when committing code, checking commit status, managing branches, or when user mentions git, commit, push, or pull.
allowed-tools: Bash
---

# Git Workflow Manager

专业的 Git 工作流管理助手，严格遵循 Q9JY 项目的提交规范和每日提交限制规则。

## 核心功能

- **提交规范化**：自动生成符合 `type(scope): description` 格式的提交信息
- **每日提交限制**：强制执行每日唯一提交规则（每人每天最终只能有一次远程提交）
- **质量检查**：提交前自动检查代码格式、测试和构建
- **分支管理**：规范的分支策略和推送管理

## 快速开始

### 检查今日提交状态

使用快速脚本检查是否需要 amend：

```bash
.claude/skills/git-workflow-manager/scripts/git-quick.sh need-amend
```

输出 `true` 表示今日已有提交，需要使用 `--amend`；输出 `false` 表示可以正常提交。

### 获取完整状态

```bash
.claude/skills/git-workflow-manager/scripts/git-status-check.sh today
```

### 提交决策建议

```bash
.claude/skills/git-workflow-manager/scripts/git-status-check.sh decision
```

## 提交规范

### 标准格式

```
<类型>(<范围>): <简短描述>

<详细描述（可选）>
```

### 类型说明

- **feat**: 新功能
- **fix**: 修复 bug
- **docs**: 文档更新
- **style**: 代码格式调整
- **refactor**: 代码重构
- **perf**: 性能优化
- **test**: 测试相关
- **chore**: 构建/工具变动

### 范围说明

常用范围：`auth`, `user`, `payment`, `api`, `db`, `config`, `ui`, `admin`

## 每日提交限制规则

### 强制性规则

- **每日唯一提交**：每个开发者每天最终只能有一次代码提交到远程仓库
- **多次开发合并**：如当日已有提交，必须使用 `git commit --amend` 合并

### 操作流程

#### 情况 A：今日首次提交

```bash
git add .
git commit -m "feat(user): 实现用户管理功能"
git push origin main
```

#### 情况 B：今日已有提交（必须合并）

```bash
git add .
git commit --amend -m "feat(user): 实现用户管理功能和权限控制

- 完成用户CRUD操作
- 添加权限验证中间件
- 优化用户状态管理
- 修复用户搜索bug"

git push --force-with-lease origin main
```

## 提交前检查清单

### 通用检查

- [ ] 检查当日提交状态（使用 `git-quick.sh need-amend`）
- [ ] 提交信息格式正确
- [ ] 代码已测试通过
- [ ] 相关文档已更新

### 后端检查

```bash
cd backend
make fmt      # 代码格式化
make check    # 编译检查
make gen-docs # 更新 Swagger 文档（如修改了模型）
```

### 前端检查

```bash
cd frontend
npm run lint       # 代码检查
npm run type-check # 类型检查
```

## 常用脚本

### git-quick.sh - 快速查询

```bash
# 检查是否需要 amend
.claude/skills/git-workflow-manager/scripts/git-quick.sh need-amend

# 获取 JSON 格式状态
.claude/skills/git-workflow-manager/scripts/git-quick.sh json

# 今日远程提交数
.claude/skills/git-workflow-manager/scripts/git-quick.sh remote-today

# 本地领先提交数
.claude/skills/git-workflow-manager/scripts/git-quick.sh ahead
```

### git-status-check.sh - 详细报告

```bash
# 今日提交状态
.claude/skills/git-workflow-manager/scripts/git-status-check.sh today

# 提交决策建议
.claude/skills/git-workflow-manager/scripts/git-status-check.sh decision

# 完整状态报告
.claude/skills/git-workflow-manager/scripts/git-status-check.sh full
```

## 分支管理策略

### 主要分支

- **main**：生产环境分支，执行每日唯一提交规则
- **dev**：开发分支，可多次推送用于测试部署

### 推荐工作流

```bash
# 1. 在 dev 分支开发
git checkout dev
git add .
git commit -m "wip: 功能进度"
git push origin dev  # 触发测试部署

# 2. 功能验证通过后合并到 main
git checkout main
git merge --ff-only dev
git reset --soft HEAD~N  # 压缩提交
git commit -m "feat(module): 功能上线"
git push origin main
```

## 详细文档

更多详细信息请参考：

- [完整命令参考](REFERENCE.md) - 所有 Git 命令和脚本的详细说明
- [最佳实践](BEST_PRACTICES.md) - 团队协作和质量保证指南
- [故障排除](TROUBLESHOOTING.md) - 常见问题和解决方案

## 提交信息示例

### ✅ 好的示例

```bash
feat(auth): 实现token刷新轮换机制

- 添加token自动轮换功能，提升安全性
- 实现黑名单管理，及时撤销失效token
- 优化refresh token映射逻辑
- 完善认证中间件错误处理

关联issue: #123
```

### ❌ 不好的示例

```bash
# 过于简单
feat: 更新代码

# 没有类型标识
修复了一些bug

# 中英文混杂
fix: 修复login的bug
```

---
