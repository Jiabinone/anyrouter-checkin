---
name: collaborating-with-codex
description: Delegates coding tasks to Codex CLI for prototyping, debugging, and code review. Use when needing algorithm implementation, bug analysis, or code quality feedback. Supports multi-turn sessions via SESSION_ID.
---

## How to Invoke (CRITICAL)

**MUST use Bash tool with `run_in_background: true`:**

```bash
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-codex/scripts/codex_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Your task description here"
```

**After launching, use `TaskOutput` to get results:**
```
TaskOutput: task_id=<returned_task_id>, block=true
```

**Output:** JSON with `success`, `SESSION_ID`, `agent_messages`, and optional `error`.

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `--PROMPT` | Yes | - | Task instruction to send to Codex |
| `--cd` | Yes | - | Workspace root directory (use absolute path) |
| `--sandbox` | No | `read-only` | `read-only` / `workspace-write` / `danger-full-access` |
| `--SESSION_ID` | No | - | Resume a previous session |
| `--return-all-messages` | No | False | Include full reasoning trace |
| `--image` | No | - | Attach image files (comma-separated) |

## Multi-turn Sessions

**Always capture `SESSION_ID`** from the first response for follow-up:

```bash
# Initial task
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-codex/scripts/codex_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Analyze auth in backend/internal/handler/auth.go"

# Continue with SESSION_ID
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-codex/scripts/codex_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --SESSION_ID "uuid-from-response" \
  --PROMPT "Write unit tests for that"
```

## Common Patterns

**Code Review (request unified diff):**
```bash
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-codex/scripts/codex_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Review backend/internal/service/order.go and provide unified diff for improvements. OUTPUT: Unified Diff Patch ONLY."
```

**Debug with full trace:**
```bash
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-codex/scripts/codex_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Debug this error: ..." \
  --return-all-messages
```
