<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAccounts, type Account } from '@/api/account'
import { getCronTasks, type CronTask } from '@/api/cron'
import { getLogs, type CheckinLog } from '@/api/system'
import { Users, Clock, CheckCircle } from 'lucide-vue-next'
import { formatTime } from '@/utils/time'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'

const accounts = ref<Account[]>([])
const tasks = ref<CronTask[]>([])
const logs = ref<CheckinLog[]>([])
const todayCheckinAccountCount = ref(0)

onMounted(async () => {
  accounts.value = await getAccounts()
  tasks.value = await getCronTasks()
  const logSummary = await getLogs()
  logs.value = logSummary.logs
  todayCheckinAccountCount.value = logSummary.today_checkin_account_count
})
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold">
      仪表盘
    </h1>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card>
        <CardContent class="pt-6">
          <div class="flex items-center gap-4">
            <Users class="w-10 h-10 text-primary" />
            <div>
              <p class="text-sm text-muted-foreground">
                账号数量
              </p>
              <p class="text-2xl font-bold">
                {{ accounts.length }}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="pt-6">
          <div class="flex items-center gap-4">
            <Clock class="w-10 h-10 text-primary" />
            <div>
              <p class="text-sm text-muted-foreground">
                定时任务
              </p>
              <p class="text-2xl font-bold">
                {{ tasks.filter(t => t.status === 1).length }}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="pt-6">
          <div class="flex items-center gap-4">
            <CheckCircle class="w-10 h-10 text-green-500" />
            <div>
              <p class="text-sm text-muted-foreground">
                今日签到
              </p>
              <p class="text-2xl font-bold">
                {{ todayCheckinAccountCount }}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>最近签到记录</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>账号ID</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>消息</TableHead>
              <TableHead>时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="log in logs.slice(0, 10)"
              :key="log.id"
            >
              <TableCell>{{ log.account_id }}</TableCell>
              <TableCell>
                <Badge
                  v-if="log.success"
                  variant="default"
                  class="bg-green-500"
                >
                  成功
                </Badge>
                <Badge
                  v-else
                  variant="destructive"
                >
                  失败
                </Badge>
              </TableCell>
              <TableCell class="text-sm">
                {{ log.message.slice(0, 50) }}
              </TableCell>
              <TableCell class="text-sm text-muted-foreground">
                {{ formatTime(log.created_at) }}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <p
          v-if="logs.length === 0"
          class="text-center text-muted-foreground py-4"
        >
          暂无记录
        </p>
      </CardContent>
    </Card>
  </div>
</template>
