# Git Workflow Manager - 故障排除

本文档包含常见问题和解决方案。

## 常见问题

### 问题1：误操作导致多次提交

**症状**：今日已经推送了多次提交到远程仓库

**解决方案**：

```bash
# 1. 查看当日所有提交
git log --since="today" --oneline

# 2. 合并最近N次提交（N为今日提交数）
git reset --soft HEAD~N

# 3. 重新创建一个规范的提交
git commit -m "feat(module): 当日开发完整总结

- 完成的功能1
- 完成的功能2
- 修复的问题"

# 4. 强制推送（使用 --force-with-lease 更安全）
git push --force-with-lease origin main
```

### 问题2：跨日提交混乱

**症状**：提交时间跨越了多天，不确定如何处理

**解决方案**：

```bash
# 1. 确认当前日期
date

# 2. 查看最近提交
git log --oneline -5

# 3. 如果跨日，重新整理提交
git reset --soft HEAD~N
git commit -m "feat(module): 跨日整理提交

- 完成的功能1
- 完成的功能2
- 修复的问题"

# 4. 推送
git push --force-with-lease origin main
```

### 问题3：不确定是否需要使用 --amend

**症状**：不知道今日是否已有提交

**解决方案**：

```bash
# 使用快速脚本检查
.claude/skills/git-workflow-manager/scripts/git-quick.sh need-amend

# 输出 true = 需要 --amend
# 输出 false = 可以正常提交
```

或者手动检查：

```bash
# 查看今日远程提交
git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)"

# 有输出 = 需要 --amend
# 无输出 = 可以正常提交
```

### 问题4：推送被拒绝（non-fast-forward）

**症状**：执行 `git push` 时提示 non-fast-forward 错误

**原因**：远程分支有新的提交，本地分支落后

**解决方案**：

```bash
# 1. 拉取远程更新
git fetch origin

# 2. 查看差异
git log HEAD..origin/$(git branch --show-current)

# 3. 选择合并策略

# 方案A：rebase（推荐，保持线性历史）
git rebase origin/$(git branch --show-current)

# 方案B：merge（保留分支历史）
git merge origin/$(git branch --show-current)

# 4. 解决冲突（如有）
git add <resolved-files>
git rebase --continue  # 或 git commit（如果是merge）

# 5. 推送
git push origin $(git branch --show-current)
```

### 问题5：提交信息格式不正确

**症状**：提交信息不符合 `type(scope): description` 格式

**解决方案**：

```bash
# 修改最后一次提交信息
git commit --amend -m "feat(user): 实现用户管理功能"

# 如果已经推送，需要强制推送
git push --force-with-lease origin main
```

### 问题6：忘记添加文件到提交

**症状**：提交后发现有文件忘记添加

**解决方案**：

```bash
# 1. 添加遗漏的文件
git add <forgotten-files>

# 2. 修改最后一次提交
git commit --amend --no-edit

# 3. 如果已经推送，需要强制推送
git push --force-with-lease origin main
```

### 问题7：提交了不应该提交的文件

**症状**：提交中包含了敏感文件或临时文件

**解决方案**：

```bash
# 1. 从暂存区移除文件
git reset HEAD <file>

# 2. 如果已经提交，修改最后一次提交
git reset --soft HEAD~1
git reset HEAD <file>
git commit -m "原提交信息"

# 3. 如果已经推送，需要强制推送
git push --force-with-lease origin main

# 4. 将文件添加到 .gitignore
echo "<file-pattern>" >> .gitignore
git add .gitignore
git commit -m "chore: 更新 .gitignore"
```

### 问题8：合并冲突

**症状**：执行 merge 或 rebase 时出现冲突

**解决方案**：

```bash
# 1. 查看冲突文件
git status

# 2. 编辑冲突文件，解决冲突标记
# <<<<<<< HEAD
# 当前分支的内容
# =======
# 合并分支的内容
# >>>>>>> branch-name

# 3. 标记冲突已解决
git add <resolved-files>

# 4. 继续合并或rebase
git rebase --continue  # 如果是rebase
git commit             # 如果是merge

# 5. 如果想放弃
git rebase --abort     # 放弃rebase
git merge --abort      # 放弃merge
```

### 问题9：分支跟踪错误

**症状**：分支没有正确跟踪远程分支

**解决方案**：

```bash
# 1. 查看当前分支跟踪状态
git branch -vv

# 2. 设置分支跟踪
git branch --set-upstream-to=origin/<branch> <local-branch>

# 或者在推送时设置
git push -u origin <branch>
```

### 问题10：本地有多个未推送的提交

**症状**：本地有多个提交还没推送到远程

**解决方案**：

```bash
# 1. 查看未推送的提交
git log origin/$(git branch --show-current)..HEAD

# 2. 合并这些提交
git reset --soft origin/$(git branch --show-current)

# 3. 创建一个新的提交
git commit -m "feat(module): 合并多个提交

- 功能1
- 功能2
- 修复问题"

# 4. 推送
git push origin $(git branch --show-current)
```

## 紧急恢复

### 恢复误删的提交

```bash
# 1. 查看所有操作历史
git reflog

# 2. 找到要恢复的提交ID
# 3. 恢复到该提交
git reset --hard <commit-id>
```

### 恢复误删的分支

```bash
# 1. 查看所有操作历史
git reflog

# 2. 找到分支最后的提交ID
# 3. 重新创建分支
git checkout -b <branch-name> <commit-id>
```

### 撤销最近的推送

```bash
# ⚠️ 警告：这会改变远程历史，谨慎使用

# 1. 本地回退
git reset --hard HEAD~1

# 2. 强制推送
git push --force-with-lease origin <branch>
```

## 预防措施

### 1. 使用 --force-with-lease 而非 --force

```bash
# ✅ 推荐：更安全，会检查远程是否有新提交
git push --force-with-lease origin main

# ❌ 不推荐：可能覆盖他人的提交
git push --force origin main
```

### 2. 提交前使用脚本检查

```bash
# 检查今日提交状态
.claude/skills/git-workflow-manager/scripts/git-quick.sh need-amend

# 获取提交建议
.claude/skills/git-workflow-manager/scripts/git-status-check.sh decision
```

### 3. 定期备份重要分支

```bash
# 创建备份分支
git branch backup-$(date +%Y%m%d) <branch>

# 推送备份到远程
git push origin backup-$(date +%Y%m%d)
```

### 4. 使用 .gitignore 避免提交不必要的文件

```bash
# 常见的 .gitignore 模式
.env
*.log
node_modules/
dist/
.DS_Store
```

## 获取帮助

### 查看 Git 命令帮助

```bash
# 查看命令帮助
git help <command>

# 例如
git help commit
git help rebase
```

### 查看脚本帮助

```bash
# git-quick.sh 帮助
.claude/skills/git-workflow-manager/scripts/git-quick.sh --help

# git-status-check.sh 帮助
.claude/skills/git-workflow-manager/scripts/git-status-check.sh --help
```

## 联系支持

如果遇到无法解决的问题：

1. 查看项目文档：[问题排查](../../../docs/common/问题排查.md)
2. 查看 Git 官方文档：https://git-scm.com/docs
3. 联系团队技术负责人
