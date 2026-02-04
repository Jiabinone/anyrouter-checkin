#!/bin/bash
# Git 快速查询脚本 - 极简版
# 用途: 单行输出，适合 AI 快速解析
# 用法: ./git-quick.sh <命令>

case "$1" in
    # 今日远程提交数
    remote-today)
        git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' '
        ;;

    # 今日本地提交数
    local-today)
        git log --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' '
        ;;

    # 本地领先远程的提交数
    ahead)
        git rev-list --count origin/$(git branch --show-current)..HEAD 2>/dev/null || echo 0
        ;;

    # 本地落后远程的提交数
    behind)
        git rev-list --count HEAD..origin/$(git branch --show-current) 2>/dev/null || echo 0
        ;;

    # 是否需要 amend (true/false)
    need-amend)
        COUNT=$(git log origin/$(git branch --show-current) --since="$(date +%Y-%m-%d) 00:00:00" --oneline --author="$(git config user.name)" 2>/dev/null | wc -l | tr -d ' ')
        [ "$COUNT" -gt 0 ] && echo "true" || echo "false"
        ;;

    # 当前分支名
    branch)
        git branch --show-current
        ;;

    # 最近一次提交的短哈希
    last-hash)
        git rev-parse --short HEAD
        ;;

    # 最近一次提交的完整哈希
    last-hash-full)
        git rev-parse HEAD
        ;;

    # 最近一次提交的标题
    last-subject)
        git log -1 --format="%s"
        ;;

    # 变更文件数统计 (格式: staged,unstaged,untracked)
    changes-count)
        STAGED=$(git diff --cached --name-only | wc -l | tr -d ' ')
        UNSTAGED=$(git diff --name-only | wc -l | tr -d ' ')
        UNTRACKED=$(git ls-files --others --exclude-standard | wc -l | tr -d ' ')
        echo "$STAGED,$UNSTAGED,$UNTRACKED"
        ;;

    # 是否有未提交的变更 (true/false)
    has-changes)
        [ -n "$(git status --porcelain)" ] && echo "true" || echo "false"
        ;;

    # 工作区是否干净 (true/false)
    is-clean)
        [ -z "$(git status --porcelain)" ] && echo "true" || echo "false"
        ;;

    # 当前用户名
    user)
        git config user.name
        ;;

    # 今日日期
    date)
        date +%Y-%m-%d
        ;;

    # JSON 格式的完整状态
    json)
        BRANCH=$(git branch --show-current)
        USER=$(git config user.name)
        TODAY=$(date +%Y-%m-%d)
        REMOTE_TODAY=$(git log origin/$BRANCH --since="$TODAY 00:00:00" --oneline --author="$USER" 2>/dev/null | wc -l | tr -d ' ')
        LOCAL_TODAY=$(git log --since="$TODAY 00:00:00" --oneline --author="$USER" 2>/dev/null | wc -l | tr -d ' ')
        AHEAD=$(git rev-list --count origin/$BRANCH..HEAD 2>/dev/null || echo 0)
        BEHIND=$(git rev-list --count HEAD..origin/$BRANCH 2>/dev/null || echo 0)
        STAGED=$(git diff --cached --name-only | wc -l | tr -d ' ')
        UNSTAGED=$(git diff --name-only | wc -l | tr -d ' ')
        UNTRACKED=$(git ls-files --others --exclude-standard | wc -l | tr -d ' ')
        NEED_AMEND="false"
        [ "$REMOTE_TODAY" -gt 0 ] && NEED_AMEND="true"

        echo "{\"date\":\"$TODAY\",\"branch\":\"$BRANCH\",\"user\":\"$USER\",\"remote_today\":$REMOTE_TODAY,\"local_today\":$LOCAL_TODAY,\"ahead\":$AHEAD,\"behind\":$BEHIND,\"staged\":$STAGED,\"unstaged\":$UNSTAGED,\"untracked\":$UNTRACKED,\"need_amend\":$NEED_AMEND}"
        ;;

    # 帮助
    *)
        echo "用法: $0 <命令>"
        echo ""
        echo "可用命令:"
        echo "  remote-today   - 今日远程提交数"
        echo "  local-today    - 今日本地提交数"
        echo "  ahead          - 本地领先远程提交数"
        echo "  behind         - 本地落后远程提交数"
        echo "  need-amend     - 是否需要 amend (true/false)"
        echo "  branch         - 当前分支名"
        echo "  last-hash      - 最近提交短哈希"
        echo "  last-hash-full - 最近提交完整哈希"
        echo "  last-subject   - 最近提交标题"
        echo "  changes-count  - 变更文件数 (staged,unstaged,untracked)"
        echo "  has-changes    - 是否有未提交变更 (true/false)"
        echo "  is-clean       - 工作区是否干净 (true/false)"
        echo "  user           - 当前用户名"
        echo "  date           - 今日日期"
        echo "  json           - JSON 格式完整状态"
        ;;
esac
