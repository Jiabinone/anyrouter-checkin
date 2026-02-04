# Git Workflow Manager Skill

专业的 Git 工作流管理技能，严格遵循 Q9JY 项目的提交规范和每日提交限制规则。

## 📁 文件结构

```
git-workflow-manager/
├── SKILL.md              # 主技能文件（Claude 自动读取）
├── README.md             # 本文件（技能说明）
├── REFERENCE.md          # 完整命令参考（按需加载）
├── BEST_PRACTICES.md     # 最佳实践指南（按需加载）
├── TROUBLESHOOTING.md    # 故障排除文档（按需加载）
└── scripts/              # 辅助脚本
    ├── git-quick.sh      # 快速查询脚本
    └── git-status-check.sh  # 详细状态检查脚本
```

## 🎯 技能特性

### 1. 渐进式披露

- **SKILL.md**：包含核心功能和快速开始指南
- **REFERENCE.md**：详细的命令参考，仅在需要时加载
- **BEST_PRACTICES.md**：最佳实践，仅在需要时加载
- **TROUBLESHOOTING.md**：故障排除，仅在需要时加载

### 2. 脚本自动化

- **git-quick.sh**：提供快速查询功能，输出简洁，适合 AI 解析
- **git-status-check.sh**：提供详细的状态报告和决策建议

### 3. 工具权限限制

使用 `allowed-tools: Bash` 限制技能只能使用 Bash 工具，确保安全性。

## 🚀 快速开始

### 检查今日提交状态

```bash
.claude/skills/git-workflow-manager/scripts/git-quick.sh need-amend
```

### 获取提交建议

```bash
.claude/skills/git-workflow-manager/scripts/git-status-check.sh decision
```

## 📖 文档说明

### SKILL.md

主技能文件，包含：
- 核心功能介绍
- 快速开始指南
- 提交规范说明
- 每日提交限制规则
- 常用脚本使用方法
- 提交信息示例

### REFERENCE.md

完整命令参考，包含：
- 所有快速信息获取命令
- 组合快捷命令
- 预置脚本工具详解
- 常用 Git 命令
- 分支推送策略
- 紧急情况处理

### BEST_PRACTICES.md

最佳实践指南，包含：
- 开发流程优化
- 团队协作建议
- 代码Review流程
- 质量指标
- 提交信息最佳实践
- 分支管理最佳实践
- 常见场景处理

### TROUBLESHOOTING.md

故障排除文档，包含：
- 常见问题及解决方案
- 紧急恢复方法
- 预防措施
- 获取帮助的方式

## 🔧 脚本说明

### git-quick.sh

快速查询脚本，提供以下功能：

```bash
# 检查是否需要 amend
git-quick.sh need-amend

# 获取 JSON 格式状态
git-quick.sh json

# 今日远程提交数
git-quick.sh remote-today

# 本地领先提交数
git-quick.sh ahead

# 变更文件数统计
git-quick.sh changes-count

# 最近提交标题
git-quick.sh last-subject
```

### git-status-check.sh

详细状态检查脚本，提供以下功能：

```bash
# 今日提交状态
git-status-check.sh today

# 提交决策建议
git-status-check.sh decision

# 查看指定提交详情
git-status-check.sh commit <commit-id>

# 完整状态报告
git-status-check.sh full
```

## 📋 使用场景

### 场景1：准备提交代码

```
用户："我要提交代码"

Claude 会：
1. 运行 git-quick.sh need-amend 检查今日提交状态
2. 根据结果建议使用正常提交或 --amend
3. 检查代码质量（格式、测试、构建）
4. 生成规范的提交信息
5. 执行提交操作
```

### 场景2：查看提交状态

```
用户："今天我提交过代码吗？"

Claude 会：
1. 运行 git-status-check.sh today
2. 显示今日提交状态
3. 提供下一步操作建议
```

### 场景3：解决提交问题

```
用户："我不小心提交了多次，怎么办？"

Claude 会：
1. 查看 TROUBLESHOOTING.md
2. 提供合并提交的解决方案
3. 指导执行修复步骤
```

## 🎓 技能优化

本技能按照 Claude Code 官方最佳实践优化：

1. **清晰的 description**：包含功能说明和触发关键词
2. **渐进式披露**：主文件简洁，详细内容分离到独立文件
3. **脚本自动化**：使用脚本代替复杂的命令组合
4. **工具权限限制**：使用 allowed-tools 限制工具访问
5. **模块化文档**：按主题分离文档，便于维护和查找

## 📚 相关资源

- [Claude Code Skills 官方文档](https://code.claude.com/docs/zh-CN/skills)
- [Git 官方文档](https://git-scm.com/docs)
- [项目 Git 提交规范](../../../docs/common/git提交规范.md)

---

*💡 提示：本技能会在用户提到 git、commit、push、pull 等关键词时自动激活。*
