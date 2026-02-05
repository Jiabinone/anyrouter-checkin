<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Trash2, Play, RefreshCw, Power } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { getAccounts, createAccount, updateAccount, updateAccountStatus, deleteAccount, checkinAccount, verifySession, refreshAccount, type Account } from '@/api/account'
import { formatTime } from '@/utils/time'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const accounts = ref<Account[]>([])
const showModal = ref(false)
const loading = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = ref({ session: '' })
const sessionInfo = ref<{ user_id: number; username: string; role: number } | null>(null)

function formatBalance(balance: Account['balance']): string {
  if (balance === null || balance === undefined) {
    return '-'
  }
  if (typeof balance === 'number') {
    return balance.toFixed(2)
  }
  const trimmed = balance.trim()
  if (!trimmed) {
    return '-'
  }
  const parsed = Number(trimmed)
  if (Number.isNaN(parsed)) {
    return trimmed
  }
  return parsed.toFixed(2)
}

function resetForm() {
  form.value = { session: '' }
  sessionInfo.value = null
}

function openCreate() {
  resetForm()
  isEditing.value = false
  editingId.value = null
  showModal.value = true
}

function openEdit(account: Account) {
  resetForm()
  isEditing.value = true
  editingId.value = account.id
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  resetForm()
  isEditing.value = false
  editingId.value = null
}

async function loadAccounts() {
  accounts.value = await getAccounts()
}

async function handleVerify() {
  if (!form.value.session) return
  try {
    sessionInfo.value = await verifySession(form.value.session)
    toast.success('Session 验证成功')
  } catch (e) {
    toast.error('Session 验证失败: ' + (e as Error).message)
    sessionInfo.value = null
  }
}

async function handleCreate() {
  if (isEditing.value) {
    if (!editingId.value) {
      toast.error('编辑账号失败：缺少账号ID')
      return
    }
    const payload = form.value.session ? { session: form.value.session } : {}
    try {
      await updateAccount(editingId.value, payload)
      closeModal()
      toast.success('账号更新成功')
      await loadAccounts()
    } catch (e) {
      toast.error('更新失败: ' + (e as Error).message)
    }
    return
  }

  if (!form.value.session) {
    toast.error('请填写 Session')
    return
  }
  try {
    await createAccount(form.value)
    closeModal()
    toast.success('账号添加成功')
    await loadAccounts()
  } catch (e) {
    toast.error('创建失败: ' + (e as Error).message)
  }
}

async function handleDelete(id: number) {
  await deleteAccount(id)
  toast.success('账号已删除')
  await loadAccounts()
}

async function handleToggleStatus(account: Account) {
  loading.value = true
  const targetStatus = account.status === 1 ? 0 : 1
  try {
    await updateAccountStatus(account.id, targetStatus)
    toast.success(targetStatus === 1 ? '账号已启用' : '账号已禁用')
    await loadAccounts()
  } catch (e) {
    toast.error('状态更新失败: ' + (e as Error).message)
  } finally {
    loading.value = false
  }
}

async function handleCheckin(id: number) {
  loading.value = true
  try {
    const result = await checkinAccount(id)
    if (result.success) {
      toast.success('签到成功')
    } else {
      toast.error('签到失败: ' + result.result)
    }
    await loadAccounts()
  } finally {
    loading.value = false
  }
}

async function handleRefresh(id: number) {
  loading.value = true
  try {
    await refreshAccount(id)
    toast.success('账号信息已刷新')
    await loadAccounts()
  } catch (e) {
    toast.error('刷新失败: ' + (e as Error).message)
  } finally {
    loading.value = false
  }
}

onMounted(loadAccounts)
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">
        账号管理
      </h1>
      <Button @click="openCreate">
        <Plus class="w-4 h-4" />
        添加账号
      </Button>
    </div>

    <div class="bg-card rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[50px]">
              ID
            </TableHead>
            <TableHead class="w-[200px]">
              AnyRouter用户
            </TableHead>
            <TableHead class="w-[100px]">
              状态
            </TableHead>
            <TableHead class="text-right w-[100px]">
              余额($)
            </TableHead>
            <TableHead class="w-[160px]">
              最后签到
            </TableHead>
            <TableHead class="w-[140px]">
              操作
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="account in accounts"
            :key="account.id"
            class="cursor-pointer"
            @click="openEdit(account)"
          >
            <TableCell class="w-[50px]">
              {{ account.id }}
            </TableCell>
            <TableCell class="w-[200px]">
              {{ account.username }} (ID: {{ account.user_id }})
            </TableCell>
            <TableCell class="w-[100px]">
              <Badge
                v-if="account.status === 1"
                class="bg-green-500"
              >
                启用
              </Badge>
              <Badge
                v-else
                variant="secondary"
              >
                禁用
              </Badge>
            </TableCell>
            <TableCell class="w-[100px] text-right tabular-nums">
              {{ formatBalance(account.balance) }}
            </TableCell>
            <TableCell class="w-[160px] text-sm">
              {{ formatTime(account.last_checkin) }}
            </TableCell>
            <TableCell
              class="w-[140px]"
              @click.stop
            >
              <div class="flex gap-2">
                <Button
                  variant="ghost"
                  size="icon-sm"
                  :disabled="loading"
                  :title="account.status === 1 ? '禁用账号' : '启用账号'"
                  @click="handleToggleStatus(account)"
                >
                  <Power
                    class="w-4 h-4"
                    :class="account.status === 1 ? 'text-green-500' : 'text-muted-foreground'"
                  />
                </Button>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  :disabled="loading || account.status !== 1"
                  :title="account.status === 1 ? '立即签到' : '账号已禁用'"
                  @click="handleCheckin(account.id)"
                >
                  <Play class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  :disabled="loading || account.status !== 1"
                  :title="account.status === 1 ? '刷新账号信息' : '账号已禁用'"
                  @click="handleRefresh(account.id)"
                >
                  <RefreshCw class="w-4 h-4" />
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
                        确定要删除账号 "{{ account.username }} (ID: {{ account.user_id }})" 吗？此操作不可撤销。
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>取消</AlertDialogCancel>
                      <AlertDialogAction @click="handleDelete(account.id)">
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
        v-if="accounts.length === 0"
        class="text-center text-muted-foreground py-8"
      >
        暂无账号
      </p>
    </div>

    <Dialog v-model:open="showModal">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ isEditing ? '编辑账号' : '添加账号' }}</DialogTitle>
        </DialogHeader>

        <div class="space-y-4">
          <div class="space-y-2">
            <Label>Session Cookie</Label>
            <Textarea
              v-model="form.session"
              class="h-24"
              :placeholder="isEditing ? '留空表示不修改' : '从浏览器复制的 session cookie'"
              autocomplete="off"
            />
            <p class="text-xs text-muted-foreground">
              {{ isEditing ? '出于安全考虑不回显 Session，留空表示不修改' : 'Session 建议先验证后再保存' }}
            </p>
            <Button
              variant="link"
              size="sm"
              class="p-0 h-auto"
              @click="handleVerify"
            >
              <RefreshCw class="w-4 h-4" />
              验证 Session
            </Button>
          </div>

          <div
            v-if="sessionInfo"
            class="p-3 bg-muted rounded text-sm"
          >
            <p>用户名: {{ sessionInfo.username }}</p>
            <p>用户ID: {{ sessionInfo.user_id }}</p>
            <p>角色: {{ sessionInfo.role }}</p>
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            @click="closeModal"
          >
            取消
          </Button>
          <Button @click="handleCreate">
            {{ isEditing ? '保存' : '添加' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
