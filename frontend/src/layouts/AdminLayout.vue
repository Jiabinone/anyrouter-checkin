<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { LayoutDashboard, Users, Clock, Bell, Settings, LogOut, KeyRound } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { useAuthStore } from '@/stores/auth'
import { changePassword, getProfile } from '@/api/auth'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import ModeToggle from '@/components/common/ModeToggle.vue'

const route = useRoute()
const auth = useAuthStore()

onMounted(async () => {
  if (!auth.username) {
    const profile = await getProfile()
    auth.setUser(profile.username)
  }
})

const userInitial = computed(() => {
  const name = auth.username || 'U'
  return name.charAt(0).toUpperCase()
})

const avatarColor = computed(() => {
  const name = auth.username || 'U'
  const code = name.charCodeAt(0)
  const colors = [
    'from-rose-500 to-pink-600',
    'from-violet-500 to-purple-600',
    'from-blue-500 to-cyan-600',
    'from-emerald-500 to-teal-600',
    'from-amber-500 to-orange-600',
  ]
  return colors[code % colors.length]
})

const menuItems = [
  { path: '/dashboard', name: '仪表盘', icon: LayoutDashboard },
  { path: '/account', name: '账号管理', icon: Users },
  { path: '/system/cron', name: '定时任务', icon: Clock },
  { path: '/system/push', name: '推送配置', icon: Bell },
  { path: '/system/config', name: '更多设置', icon: Settings },
]

function isActiveRoute(path: string) {
  return route.path === path || route.path.startsWith(path + '/')
}

const showPasswordDialog = ref(false)
const passwordForm = ref({ oldPassword: '', newPassword: '', confirmPassword: '' })
const submitting = ref(false)

async function handleChangePassword() {
  if (!passwordForm.value.oldPassword || !passwordForm.value.newPassword) {
    toast.error('请填写完整信息')
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    toast.error('新密码至少6位')
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    toast.error('两次密码不一致')
    return
  }
  submitting.value = true
  try {
    await changePassword(passwordForm.value.oldPassword, passwordForm.value.newPassword)
    toast.success('密码修改成功')
    showPasswordDialog.value = false
    passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    submitting.value = false
  }
}

function handleLogout() {
  auth.logout()
  window.location.href = '/login'
}
</script>

<template>
  <div class="flex h-screen">
    <aside class="w-64 bg-card border-r flex flex-col">
      <div class="p-6">
        <h1 class="text-xl font-bold">
          AnyRouter
        </h1>
        <p class="text-sm text-muted-foreground">
          管理后台
        </p>
      </div>

      <Separator />

      <nav class="p-4 space-y-1 flex-1">
        <RouterLink
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-3 py-2 rounded-md transition-colors"
          :class="isActiveRoute(item.path)
            ? 'bg-primary text-primary-foreground'
            : 'hover:bg-muted'"
        >
          <component
            :is="item.icon"
            class="w-5 h-5"
          />
          <span>{{ item.name }}</span>
        </RouterLink>
      </nav>
    </aside>

    <main class="flex-1 flex flex-col overflow-hidden">
      <header class="h-14 border-b bg-card flex items-center justify-end px-6 gap-3">
        <ModeToggle />
        <DropdownMenu>
          <DropdownMenuTrigger class="flex items-center gap-2 rounded-full py-1 pl-1 pr-3 transition-colors hover:bg-accent focus:outline-none focus-visible:ring-1 focus-visible:ring-ring">
            <div :class="['h-8 w-8 rounded-full bg-gradient-to-br flex items-center justify-center text-white text-sm font-semibold shadow-sm ring-2 ring-background', avatarColor]">
              {{ userInitial }}
            </div>
            <span class="text-sm font-medium">{{ auth.username }}</span>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            class="w-52"
          >
            <div class="flex items-center gap-3 px-3 py-3">
              <div :class="['h-9 w-9 shrink-0 rounded-full bg-gradient-to-br flex items-center justify-center text-white font-semibold shadow-sm', avatarColor]">
                {{ userInitial }}
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-sm font-medium truncate">{{ auth.username }}</span>
                <span class="text-xs text-muted-foreground">管理员</span>
              </div>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              class="cursor-pointer"
              @click="showPasswordDialog = true"
            >
              <KeyRound class="mr-2 h-4 w-4" />
              修改密码
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              class="text-destructive focus:text-destructive cursor-pointer"
              @click="handleLogout"
            >
              <LogOut class="mr-2 h-4 w-4" />
              退出登录
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </header>
      <div class="flex-1 overflow-auto bg-muted p-6">
        <RouterView />
      </div>
    </main>

    <Dialog v-model:open="showPasswordDialog">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>修改密码</DialogTitle>
          <DialogDescription>请输入原密码和新密码</DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-2">
          <div class="space-y-2">
            <Label>原密码</Label>
            <Input
              v-model="passwordForm.oldPassword"
              type="password"
              placeholder="请输入原密码"
            />
          </div>
          <div class="space-y-2">
            <Label>新密码</Label>
            <Input
              v-model="passwordForm.newPassword"
              type="password"
              placeholder="至少6位"
            />
          </div>
          <div class="space-y-2">
            <Label>确认密码</Label>
            <Input
              v-model="passwordForm.confirmPassword"
              type="password"
              placeholder="再次输入新密码"
            />
          </div>
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            @click="showPasswordDialog = false"
          >
            取消
          </Button>
          <Button
            :disabled="submitting"
            @click="handleChangePassword"
          >
            {{ submitting ? '提交中...' : '确认修改' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
