---
name: collaborating-with-gemini
description: Delegates coding tasks to Gemini CLI for prototyping, debugging, and code review. Use when needing algorithm implementation, bug analysis, or code quality feedback. Supports multi-turn sessions via SESSION_ID.
---

## How to Invoke (CRITICAL)

**MUST use Bash tool with `run_in_background: true`:**

```bash
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
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
| `--PROMPT` | Yes | - | Task instruction to send to Gemini |
| `--cd` | Yes | - | Workspace root directory (use absolute path) |
| `--sandbox` | No | False | Run in sandbox mode |
| `--SESSION_ID` | No | - | Resume a previous session |
| `--return-all-messages` | No | False | Include full reasoning trace |

## Multi-turn Sessions

**Always capture `SESSION_ID`** from the first response for follow-up:

```bash
# Initial task
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Analyze Vue component in frontend/src/views/OrdersView.vue"

# Continue with SESSION_ID
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --SESSION_ID "uuid-from-response" \
  --PROMPT "Improve the styling"
```

## Common Patterns

**Frontend/UI Review (request unified diff):**
```bash
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Review frontend/src/views/PlayersView.vue and provide unified diff for UI improvements. OUTPUT: Unified Diff Patch ONLY."
```

**Debug with full trace:**
```bash
python /Users/jiabinfeng/Documents/work/go/q9jy-backend/.claude/skills/collaborating-with-gemini/scripts/gemini_bridge.py \
  --cd "/Users/jiabinfeng/Documents/work/go/q9jy-backend" \
  --PROMPT "Debug this error: ..." \
  --return-all-messages
```

## Note

Gemini excels at **frontend/UI/CSS** tasks. For backend logic, prefer using `collaborating-with-codex`.
