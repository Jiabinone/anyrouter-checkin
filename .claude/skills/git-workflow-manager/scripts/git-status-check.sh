#!/bin/bash
# Git 状态快速检查脚本
# 用途: 为 AI 提供精简的 Git 状态信息，减少上下文窗口占用
# 用法: ./git-status-check.sh [选项]
#   选项:
#     today     - 检查今日提交状态（默认）
#     commit    - 查看指定提交详情 (需要第二个参数: commit_id)
#     changes   - 查看当前变更
#     branch    - 查看分支状态
#     decision  - 获取提交决策建议
#     full      - 完整状态报告

set -e

# 颜色定义（可选，某些终端可能不支持）
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取当前日期（格式: YYYY-MM-DD）
TODAY=$(date +%Y-%m-%d)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
GIT_USER=$(git config user.name 2>/dev/null || echo "unknown")

# 函数: 检查今日提交状态
check_today() {
    echo "=== 今日提交状态 ==="
    echo "日期: $TODAY"
    echo "分支: $CURRENT_BRANCH"
    echo "用户: $GIT_USER"
    echo ""

    # 本地今日提交
    LOCAL_COUNT=$(git log --since="$TODAY 00:00:00" --oneline --author="$GIT_USER" 2>/dev/null | wc -l | tr -d ' ')
    echo "本地今日提交: $LOCAL_COUNT 次"

    # 远程今日提交
    REMOTE_COUNT=$(git log origin/$CURRENT_BRANCH --since="$TODAY 00:00:00" --oneline --author="$GIT_USER" 2>/dev/null | wc -l | tr -d ' ')
    echo "远程今日提交: $REMOTE_COUNT 次"

    echo ""
    if [ "$REMOTE_COUNT" -gt 0 ]; then
        echo "⚠️  状态: 今日已有远程提交"
        echo "📌 操作: 需要使用 --amend 合并提交"
        echo ""
        echo "最近远程提交:"
        git log origin/$CURRENT_BRANCH --since="$TODAY 00:00:00" --oneline --author="$GIT_USER" 2>/dev/null | head -3 | sed 's/^/  /'
    else
        echo "✅ 状态: 今日尚未提交到远程"
        echo "📌 操作: 可正常提交"
    fi
}

# 函数: 查看提交详情
show_commit() {
    local COMMIT_ID=$1
    if [ -z "$COMMIT_ID" ]; then
        echo "错误: 请提供提交ID"
        echo "用法: $0 commit <commit_id>"
        exit 1
    fi

    echo "=== 提交详情: $COMMIT_ID ==="
    git show --stat --format="哈希: %H%n短哈希: %h%n作者: %an <%ae>%n日期: %ci%n标题: %s%n%n正文:%n%b" "$COMMIT_ID" 2>/dev/null | head -40

    echo ""
    echo "=== 修改文件列表 ==="
    git diff-tree --no-commit-id --name-status -r "$COMMIT_ID" 2>/dev/null
}

# 函数: 查看当前变更
show_changes() {
    echo "=== 工作区变更 ==="

    STAGED=$(git diff --cached --name-only | wc -l | tr -d ' ')
    UNSTAGED=$(git diff --name-only | wc -l | tr -d ' ')
    UNTRACKED=$(git ls-files --others --exclude-standard | wc -l | tr -d ' ')

    echo "已暂存: $STAGED 文件"
    echo "未暂存: $UNSTAGED 文件"
    echo "未跟踪: $UNTRACKED 文件"
    echo ""

    if [ "$STAGED" -gt 0 ]; then
        echo "--- 已暂存文件 ---"
        git diff --cached --name-status | head -20
        echo ""
    fi

    if [ "$UNSTAGED" -gt 0 ]; then
        echo "--- 未暂存修改 ---"
        git diff --name-status | head -20
        echo ""
    fi

    if [ "$UNTRACKED" -gt 0 ]; then
        echo "--- 未跟踪文件 ---"
        git ls-files --others --exclude-standard | head -20
    fi
}

# 函数: 查看分支状态
show_branch() {
    echo "=== 分支状态 ==="
    echo "当前分支: $CURRENT_BRANCH"
    echo ""

    # 本地领先/落后远程
    AHEAD=$(git rev-list --count origin/$CURRENT_BRANCH..HEAD 2>/dev/null || echo 0)
    BEHIND=$(git rev-list --count HEAD..origin/$CURRENT_BRANCH 2>/dev/null || echo 0)

    echo "本地领先远程: $AHEAD 个提交"
    echo "本地落后远程: $BEHIND 个提交"
    echo ""

    if [ "$AHEAD" -gt 0 ]; then
        echo "--- 待推送提交 ---"
        git log --oneline origin/$CURRENT_BRANCH..HEAD 2>/dev/null | head -10
    fi

    if [ "$BEHIND" -gt 0 ]; then
        echo ""
        echo "--- 待拉取提交 ---"
        git log --oneline HEAD..origin/$CURRENT_BRANCH 2>/dev/null | head -10
    fi
}

# 函数: 提交决策建议
show_decision() {
    echo "=== 提交决策建议 ==="

    REMOTE_TODAY=$(git log origin/$CURRENT_BRANCH --since="$TODAY 00:00:00" --oneline --author="$GIT_USER" 2>/dev/null | wc -l | tr -d ' ')
    LOCAL_AHEAD=$(git rev-list --count origin/$CURRENT_BRANCH..HEAD 2>/dev/null || echo 0)

    if [ "$REMOTE_TODAY" -gt 0 ]; then
        echo "📌 建议: git commit --amend + git push --force-with-lease"
        echo "   原因: 今日远程已有 $REMOTE_TODAY 次提交，需要合并"
        echo ""
        echo "   命令示例:"
        echo '   git add .'
        echo '   git commit --amend -m "feat(scope): 更新后的提交信息"'
        echo '   git push --force-with-lease origin '$CURRENT_BRANCH
    elif [ "$LOCAL_AHEAD" -gt 1 ]; then
        echo "📌 建议: git reset --soft HEAD~$LOCAL_AHEAD + git commit + git push"
        echo "   原因: 本地有 $LOCAL_AHEAD 个未推送提交，建议合并后推送"
        echo ""
        echo "   命令示例:"
        echo "   git reset --soft HEAD~$LOCAL_AHEAD"
        echo '   git commit -m "feat(scope): 合并后的提交信息"'
        echo '   git push origin '$CURRENT_BRANCH
    else
        echo "📌 建议: git commit + git push"
        echo "   原因: 今日首次提交，正常流程即可"
        echo ""
        echo "   命令示例:"
        echo '   git add .'
        echo '   git commit -m "feat(scope): 提交信息"'
        echo '   git push origin '$CURRENT_BRANCH
    fi
}

# 函数: 完整状态报告
show_full() {
    echo "=========================================="
    echo "         Git 完整状态报告"
    echo "=========================================="
    echo ""
    check_today
    echo ""
    echo "----------------------------------------"
    show_changes
    echo ""
    echo "----------------------------------------"
    show_branch
    echo ""
    echo "----------------------------------------"
    show_decision
    echo ""
    echo "=========================================="
}

# 主逻辑
case "${1:-today}" in
    today)
        check_today
        ;;
    commit)
        show_commit "$2"
        ;;
    changes)
        show_changes
        ;;
    branch)
        show_branch
        ;;
    decision)
        show_decision
        ;;
    full)
        show_full
        ;;
    *)
        echo "Git 状态快速检查脚本"
        echo ""
        echo "用法: $0 [选项] [参数]"
        echo ""
        echo "选项:"
        echo "  today     - 检查今日提交状态（默认）"
        echo "  commit    - 查看指定提交详情 (需要: commit_id)"
        echo "  changes   - 查看当前变更"
        echo "  branch    - 查看分支状态"
        echo "  decision  - 获取提交决策建议"
        echo "  full      - 完整状态报告"
        echo ""
        echo "示例:"
        echo "  $0 today"
        echo "  $0 commit abc123"
        echo "  $0 full"
        ;;
esac
