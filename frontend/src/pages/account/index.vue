<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAccounts, createAccount, deleteAccount, checkinAccount, verifySession, type Account } from '@/api/account'
import { Plus, Trash2, Play, RefreshCw } from 'lucide-vue-next'
import { formatTime } from '@/utils/time'
import { toast } from 'vue-sonner'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const accounts = ref<Account[]>([])
const showModal = ref(false)
const loading = ref(false)

const form = ref({ name: '', session: '' })
const sessionInfo = ref<{ user_id: number; username: string; role: number } | null>(null)

async function loadAccounts() {
  accounts.value = await getAccounts()
}

async function handleVerify() {
  if (!form.value.session) return
  try {
    sessionInfo.value = await verifySession(form.value.session)
    if (sessionInfo.value && !form.value.name) {
      form.value.name = `${sessionInfo.value.username} (${sessionInfo.value.user_id})`
    }
    toast.success('Session 验证成功')
  } catch (e) {
    toast.error('Session 验证失败: ' + (e as Error).message)
    sessionInfo.value = null
  }
}

async function handleCreate() {
  if (!form.value.name || !form.value.session) {
    toast.error('请填写完整信息')
    return
  }
  try {
    await createAccount(form.value)
    showModal.value = false
    form.value = { name: '', session: '' }
    sessionInfo.value = null
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

onMounted(loadAccounts)
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold">
        账号管理
      </h1>
      <Button @click="showModal = true">
        <Plus class="w-4 h-4" />
        添加账号
      </Button>
    </div>

    <div class="bg-card rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>名称</TableHead>
            <TableHead>AnyRouter用户</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>最后签到</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="account in accounts"
            :key="account.id"
          >
            <TableCell>{{ account.id }}</TableCell>
            <TableCell>{{ account.name }}</TableCell>
            <TableCell>{{ account.username }} (ID: {{ account.user_id }})</TableCell>
            <TableCell>
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
            <TableCell class="text-sm">
              {{ formatTime(account.last_checkin) }}
            </TableCell>
            <TableCell>
              <div class="flex gap-2">
                <Button
                  variant="ghost"
                  size="icon-sm"
                  :disabled="loading"
                  title="立即签到"
                  @click="handleCheckin(account.id)"
                >
                  <Play class="w-4 h-4" />
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
                        确定要删除账号 "{{ account.name }}" 吗？此操作不可撤销。
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
          <DialogTitle>添加账号</DialogTitle>
        </DialogHeader>

        <div class="space-y-4">
          <div class="space-y-2">
            <Label>账号名称</Label>
            <Input
              v-model="form.name"
              placeholder="备注名称"
              autocomplete="off"
            />
          </div>

          <div class="space-y-2">
            <Label>Session Cookie</Label>
            <Textarea
              v-model="form.session"
              class="h-24"
              placeholder="从浏览器复制的 session cookie"
              autocomplete="off"
            />
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
            @click="showModal = false"
          >
            取消
          </Button>
          <Button @click="handleCreate">
            添加
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
