<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Trash2, Play, Pause, PlayCircle } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { getCronTasks, createCronTask, updateCronTask, deleteCronTask, triggerCronTask, type CronTask } from '@/api/cron'
import { getAccounts, type Account } from '@/api/account'
import { formatTime } from '@/utils/time'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const tasks = ref<CronTask[]>([])
const accounts = ref<Account[]>([])
const showModal = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  cron_expr: '0 8 * * *',
  account_ids: [] as number[],
  status: 1,
})

async function loadData() {
  tasks.value = await getCronTasks()
  accounts.value = await getAccounts()
}

function openCreate() {
  editingId.value = null
  form.value = { name: '', cron_expr: '0 8 * * *', account_ids: [], status: 1 }
  showModal.value = true
}

function openEdit(task: CronTask) {
  editingId.value = task.id
  const parsedIds = parseAccountIds(task.account_ids)
  form.value = {
    name: task.name,
    cron_expr: task.cron_expr,
    account_ids: parsedIds,
    status: task.status,
  }
  showModal.value = true
}

async function handleSave() {
  if (!form.value.name || !form.value.cron_expr) {
    toast.error('请填写完整信息')
    return
  }

  const data = {
    ...form.value,
    account_ids: JSON.stringify(form.value.account_ids),
  }

  if (editingId.value) {
    await updateCronTask(editingId.value, data)
    toast.success('任务已更新')
  } else {
    await createCronTask(data)
    toast.success('任务已创建')
  }

  showModal.value = false
  await loadData()
}

async function handleToggle(task: CronTask) {
  await updateCronTask(task.id, { ...task, status: task.status === 1 ? 0 : 1 })
  toast.success(task.status === 1 ? '任务已暂停' : '任务已启用')
  await loadData()
}

async function handleDelete(id: number) {
  await deleteCronTask(id)
  toast.success('任务已删除')
  await loadData()
}

async function handleTrigger(id: number) {
  await triggerCronTask(id)
  toast.success('任务已触发')
}

function parseAccountIds(value: string) {
  try {
    const parsed = JSON.parse(value || '[]') as unknown
    if (!Array.isArray(parsed)) return []
    return parsed
      .map((id) => Number(id))
      .filter((id) => Number.isFinite(id))
  } catch {
    return []
  }
}

function getAccountCount(value: string) {
  return parseAccountIds(value).length
}

function setAccountChecked(id: number, checked: boolean | 'indeterminate') {
  const idx = form.value.account_ids.indexOf(id)
  if (checked === true) {
    if (idx === -1) form.value.account_ids.push(id)
    return
  }
  if (idx > -1) form.value.account_ids.splice(idx, 1)
}

onMounted(loadData)
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">
        定时任务
      </h1>
      <Button @click="openCreate">
        <Plus class="w-4 h-4" />
        创建任务
      </Button>
    </div>

    <div class="bg-card rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>名称</TableHead>
            <TableHead>Cron表达式</TableHead>
            <TableHead>账号数</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>下次执行</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="task in tasks"
            :key="task.id"
            class="cursor-pointer hover:bg-muted/50"
            @click="openEdit(task)"
          >
            <TableCell>{{ task.name }}</TableCell>
            <TableCell class="font-mono text-sm">
              {{ task.cron_expr }}
            </TableCell>
            <TableCell>{{ getAccountCount(task.account_ids) }}</TableCell>
            <TableCell>
              <Badge
                v-if="task.status === 1"
                class="bg-green-500"
              >
                运行中
              </Badge>
              <Badge
                v-else
                variant="secondary"
              >
                已暂停
              </Badge>
            </TableCell>
            <TableCell class="text-sm">
              {{ formatTime(task.next_run) }}
            </TableCell>
            <TableCell @click.stop>
              <div class="flex gap-2">
                <Button
                  variant="ghost"
                  size="icon-sm"
                  title="立即执行"
                  @click="handleTrigger(task.id)"
                >
                  <PlayCircle class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  :title="task.status === 1 ? '暂停' : '启用'"
                  @click="handleToggle(task)"
                >
                  <Pause
                    v-if="task.status === 1"
                    class="w-4 h-4"
                  />
                  <Play
                    v-else
                    class="w-4 h-4"
                  />
                </Button>
                <AlertDialog>
                  <AlertDialogTrigger as-child>
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      title="删除"
                    >
                      <Trash2 class="w-4 h-4 text-destructive" />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>确认删除</AlertDialogTitle>
                      <AlertDialogDescription>
                        确定要删除任务 "{{ task.name }}" 吗？此操作不可撤销。
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>取消</AlertDialogCancel>
                      <AlertDialogAction @click="handleDelete(task.id)">
                        删除
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <p
        v-if="tasks.length === 0"
        class="text-center text-muted-foreground py-8"
      >
        暂无任务
      </p>
    </div>

    <Dialog v-model:open="showModal">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ editingId ? '编辑任务' : '创建任务' }}</DialogTitle>
        </DialogHeader>

        <div class="space-y-4">
          <div class="space-y-2">
            <Label>任务名称</Label>
            <Input
              v-model="form.name"
              placeholder="每日签到"
            />
          </div>

          <div class="space-y-2">
            <Label>Cron 表达式</Label>
            <Input
              v-model="form.cron_expr"
              class="font-mono"
              placeholder="0 8 * * *"
            />
            <p class="text-xs text-muted-foreground">
              格式: 分 时 日 月 周 (例: 0 8 * * * = 每天8点)
            </p>
          </div>

          <div class="space-y-2">
            <Label>关联账号</Label>
            <div class="space-y-2 max-h-40 overflow-auto">
              <div
                v-for="acc in accounts"
                :key="acc.id"
                class="flex items-center gap-2"
              >
                <Checkbox
                  :id="`acc-${acc.id}`"
                  :model-value="form.account_ids.includes(acc.id)"
                  @update:model-value="setAccountChecked(acc.id, $event)"
                />
                <Label
                  :for="`acc-${acc.id}`"
                  class="cursor-pointer font-normal"
                >
                  {{ acc.name }} ({{ acc.username }})
                </Label>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            @click="showModal = false"
          >
            取消
          </Button>
          <Button @click="handleSave">
            保存
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
