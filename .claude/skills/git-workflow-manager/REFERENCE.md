# Git Workflow Manager - 完整命令参考

本文档包含所有 Git 命令和脚本的详细说明。

## 快速信息获取命令

### 今日提交状态速查

```bash
# 一键检查今日远程提交状态（推荐首选）
git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | head -5

# 返回说明:
# - 有输出 = 今日已有远程提交，需要用 --amend
# - 无输出 = 今日尚未提交到远程，可正常提交
```

```bash
# 完整状态报告（一条命令获取所有关键信息）
echo "=== 今日提交状态 ===" && \
echo "当前日期: $(date +%Y-%m-%d)" && \
echo "当前分支: $(git branch --show-current)" && \
echo "本地今日提交数: $(git log --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' ')" && \
echo "远程今日提交数: $(git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' ')" && \
echo "本地领先远程: $(git rev-list --count origin/$(git branch --show-current)..HEAD 2>/dev/null || echo 0) 个提交" && \
echo "本地落后远程: $(git rev-list --count HEAD..origin/$(git branch --show-current) 2>/dev/null || echo 0) 个提交"
```

### 提交详情速查

```bash
# 查看指定提交的精简信息（替换 COMMIT_ID）
git show --stat --format="提交: %H%n作者: %an <%ae>%n日期: %ci%n标题: %s%n%n%b" COMMIT_ID | head -30

# 仅查看提交影响的文件列表
git diff-tree --no-commit-id --name-status -r COMMIT_ID

# 查看提交的完整diff
git show COMMIT_ID --stat
```

```bash
# 查看最近N次提交的精简列表
git log --oneline -N --format="%h | %ci | %s"

# 查看最近提交的详细信息
git log -1 --format="哈希: %H%n短哈希: %h%n作者: %an%n邮箱: %ae%n日期: %ci%n标题: %s%n正文:%n%b"
```

### 变更文件速查

```bash
# 查看当前工作区变更概览
echo "=== 变更概览 ===" && \
echo "已暂存文件: $(git diff --cached --name-only | wc -l | tr -d ' ')" && \
echo "未暂存修改: $(git diff --name-only | wc -l | tr -d ' ')" && \
echo "未跟踪文件: $(git ls-files --others --exclude-standard | wc -l | tr -d ' ')"

# 列出所有变更文件（精简）
git status --porcelain

# 查看暂存区文件列表
git diff --cached --name-only

# 查看某次提交修改的文件及修改行数
git show --stat COMMIT_ID --format=""
```

### 分支状态速查

```bash
# 分支跟踪状态一览
git branch -vv --format="%(refname:short) -> %(upstream:short) [%(upstream:track)]"

# 检查是否有需要推送的提交
git log --oneline origin/$(git branch --show-current)..HEAD 2>/dev/null

# 检查是否需要拉取远程更新
git log --oneline HEAD..origin/$(git branch --show-current) 2>/dev/null
```

## 组合快捷命令

### 提交前完整状态检查

```bash
echo "========== Git 状态检查 ==========" && \
echo "📅 日期: $(date +%Y-%m-%d' '%H:%M:%S)" && \
echo "🌿 分支: $(git branch --show-current)" && \
echo "" && \
echo "📊 今日提交情况:" && \
REMOTE_TODAY=$(git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' ') && \
if [ "$REMOTE_TODAY" -gt 0 ]; then \
  echo "  ⚠️  远程已有 $REMOTE_TODAY 次提交，需要用 --amend"; \
  echo "  最近远程提交:"; \
  git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | head -3 | sed 's/^/    /'; \
else \
  echo "  ✅ 今日尚未提交到远程，可正常提交"; \
fi && \
echo "" && \
echo "📁 工作区状态:" && \
echo "  已暂存: $(git diff --cached --name-only | wc -l | tr -d ' ') 文件" && \
echo "  未暂存: $(git diff --name-only | wc -l | tr -d ' ') 文件" && \
echo "  未跟踪: $(git ls-files --others --exclude-standard | wc -l | tr -d ' ') 文件" && \
echo "" && \
echo "🔄 同步状态:" && \
echo "  本地领先: $(git rev-list --count origin/$(git branch --show-current)..HEAD 2>/dev/null || echo 0) 个提交" && \
echo "  本地落后: $(git rev-list --count HEAD..origin/$(git branch --show-current) 2>/dev/null || echo 0) 个提交" && \
echo "=================================="
```

### 快速决策：我应该怎么提交？

```bash
REMOTE_TODAY=$(git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' ') && \
LOCAL_AHEAD=$(git rev-list --count origin/$(git branch --show-current)..HEAD 2>/dev/null || echo 0) && \
if [ "$REMOTE_TODAY" -gt 0 ]; then \
  echo "📌 建议操作: git commit --amend + git push --force-with-lease"; \
  echo "   原因: 今日远程已有提交，需要合并"; \
elif [ "$LOCAL_AHEAD" -gt 1 ]; then \
  echo "📌 建议操作: git reset --soft HEAD~$LOCAL_AHEAD + git commit + git push"; \
  echo "   原因: 本地有多个未推送提交，建议合并后推送"; \
else \
  echo "📌 建议操作: git commit + git push"; \
  echo "   原因: 今日首次提交，正常流程即可"; \
fi
```

## 预置脚本工具

### git-quick.sh - 极简查询

```bash
# 查看是否需要 amend
.claude/skills/git-workflow-manager/scripts/git-quick.sh need-amend
# 输出: true 或 false

# 获取 JSON 格式完整状态
.claude/skills/git-workflow-manager/scripts/git-quick.sh json
# 输出: {"date":"2025-01-15","branch":"dev","remote_today":1,"need_amend":true,...}

# 更多命令
.claude/skills/git-workflow-manager/scripts/git-quick.sh remote-today  # 今日远程提交数
.claude/skills/git-workflow-manager/scripts/git-quick.sh ahead         # 本地领先提交数
.claude/skills/git-workflow-manager/scripts/git-quick.sh changes-count # 变更文件数统计
.claude/skills/git-workflow-manager/scripts/git-quick.sh last-subject  # 最近提交标题
```

### git-status-check.sh - 详细报告

```bash
# 今日提交状态（默认）
.claude/skills/git-workflow-manager/scripts/git-status-check.sh today

# 提交决策建议
.claude/skills/git-workflow-manager/scripts/git-status-check.sh decision

# 查看指定提交详情
.claude/skills/git-workflow-manager/scripts/git-status-check.sh commit abc123

# 完整状态报告
.claude/skills/git-workflow-manager/scripts/git-status-check.sh full
```

## 常用 Git 命令

### 基础操作

```bash
# 查看状态
git status
git diff
git log --oneline -10

# 查看今日提交
git log --since="today" --oneline --author="$(git config user.name)"

# 检查当日提交数量
git log --since="today" --oneline --author="$(git config user.name)" | wc -l
```

### 提交操作

```bash
# 首次提交
git commit -m "feat(auth): 实现用户登录功能"

# 修改最后一次提交
git commit --amend -m "新的提交信息"

# 合并提交（每日限制用）
git commit --amend -m "feat(auth): 完整实现认证系统

- 实现用户登录功能
- 添加token刷新机制
- 优化权限验证
- 修复认证bug"

# 强制推送
git push --force-with-lease origin main
```

### 本地开发管理

```bash
# 临时提交
git commit -m "wip: 开发进度保存"

# 合并本地提交
git reset --soft HEAD~N
git commit -m "feat(module): 当日完整开发总结"
```

### 分支管理

```bash
# 创建新分支
git checkout -b feature/user-management

# 合并分支
git merge feature/user-management

# 删除分支
git branch -d feature/user-management

# 查看分支跟踪状态
git branch -vv
```

## 分支推送策略

### 主要分支

- **`main`（生产）**：仅在准备上线时推送，执行"每日唯一提交"规则
- **`dev`（测试）**：允许推送到远程以触发测试部署，可一天多次推送

### 推荐工作流

```bash
# 1. 同步主干并更新 dev
git checkout main
git pull origin main
git checkout dev || git checkout -b dev origin/dev
git pull origin dev

# 2. 在 dev 分支开发，可多次提交
git add .
git commit -m "wip: 功能进度"

# 3. 推送到测试环境（触发 dev 部署）
git push origin dev

# 4. 功能验证通过后合并回 main（遵守每日唯一提交）
git checkout main
git merge --ff-only dev
git reset --soft HEAD~N  # 如需压缩提交
git commit -m "feat(module): 功能上线"
git push origin main
```

### 分支状态检查

```bash
git branch -vv              # 查看所有分支跟踪状态
git branch -vv | grep dev    # 确保 dev -> origin/dev

# 检查是否存在 main/dev 之外仍在跟踪远程的分支
git branch -vv | grep -v "main" | grep -v "dev" | grep -v "no upstream"
```

## 紧急情况处理

### 误操作导致多次提交

```bash
# 查看当日所有提交
git log --since="today" --oneline

# 合并最近N次提交
git reset --soft HEAD~N
git commit -m "feat(module): 当日开发完整总结"
git push --force-with-lease origin main
```

### 跨日提交混乱

```bash
# 确认当前日期
date

# 查看最近提交
git log --oneline -5

# 如果跨日，重新整理提交
git reset --soft HEAD~N
git commit -m "feat(module): 跨日整理提交

- 完成的功能1
- 完成的功能2
- 修复的问题"

git push --force-with-lease origin main
```

## 冲突解决

```bash
# 更新远程分支
git fetch origin
git rebase origin/main

# 解决冲突后继续
git add <resolved-files>
git rebase --continue

# 放弃rebase
git rebase --abort
```
