<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getConfigs, updateConfigs, testTelegram } from '@/api/system'
import { Save, Send } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'

const recommendedTemplate = [
  '<b>AnyRouter 签到系统</b>',
  '你好，<code>{{.Name}}</code>',
  '状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}',
  '结果：',
  '<pre>{{.Result}}</pre>',
].join('\n')

const legacyTemplates = new Set([
  '签到结果: {{.Result}}',
  '【{{.Name}}】签到{{if .Success}}成功{{else}}失败{{end}}: {{.Result}}',
  '签到通知\n账号：{{.Name}}\n状态：{{if .Success}}成功{{else}}失败{{end}}\n结果：{{.Result}}',
  '<b>签到通知</b>\n账号：<code>{{.Name}}</code>\n状态：{{if .Success}}<b>成功</b>{{else}}<b>失败</b>{{end}}\n结果：\n<pre>{{.Result}}</pre>',
  '<b>签到提醒</b>\n你好，<code>{{.Name}}</code>\n状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}\n结果：\n<pre>{{.Result}}</pre>',
])

const form = ref({
  'telegram.enabled': 'false',
  'telegram.bot_token': '',
  'telegram.chat_id': '',
  'telegram.template': recommendedTemplate,
})

const telegramEnabled = computed({
  get: () => form.value['telegram.enabled'] === 'true',
  set: (value: boolean) => {
    form.value['telegram.enabled'] = value ? 'true' : 'false'
  },
})

const loading = ref(false)
const testing = ref(false)

async function loadConfigs() {
  const data = await getConfigs('telegram')
  Object.assign(form.value, data)

  const currentTemplate = form.value['telegram.template']
  if (!currentTemplate || legacyTemplates.has(currentTemplate)) {
    form.value['telegram.template'] = recommendedTemplate
  }
}

async function handleSave() {
  loading.value = true
  try {
    await updateConfigs('telegram', form.value)
    toast.success('保存成功')
  } finally {
    loading.value = false
  }
}

async function handleTest() {
  testing.value = true
  try {
    await testTelegram()
    toast.success('测试消息已发送')
  } catch (e) {
    toast.error('发送失败: ' + (e as Error).message)
  } finally {
    testing.value = false
  }
}

onMounted(loadConfigs)
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold">
      Telegram 推送配置
    </h1>

    <Card class="max-w-2xl">
      <CardHeader>
        <CardTitle>推送设置</CardTitle>
        <CardDescription>配置签到结果推送到 Telegram</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="space-y-6">
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>启用推送</Label>
              <p class="text-sm text-muted-foreground">
                开启 Telegram 通知
              </p>
            </div>
            <Switch v-model="telegramEnabled" />
          </div>

          <div class="space-y-2">
            <Label>Bot Token</Label>
            <Input
              v-model="form['telegram.bot_token']"
              class="font-mono text-sm"
              placeholder="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
            />
            <p class="text-xs text-muted-foreground">
              从 @BotFather 获取
            </p>
          </div>

          <div class="space-y-2">
            <Label>Chat ID</Label>
            <Input
              v-model="form['telegram.chat_id']"
              class="font-mono text-sm"
              placeholder="-1001234567890"
            />
            <p class="text-xs text-muted-foreground">
              群组/频道/个人 Chat ID
            </p>
          </div>

          <div class="space-y-2">
            <Label>消息模板（HTML）</Label>
            <Textarea
              v-model="form['telegram.template']"
              class="font-mono text-sm h-24"
              :placeholder="recommendedTemplate"
            />
            <p
              v-pre
              class="text-xs text-muted-foreground"
            >
              可用变量: {{.Name}} 账号名, {{.Success}} 是否成功, {{.Result}} 结果（支持 Telegram HTML）
            </p>
          </div>

          <div class="flex gap-2 pt-4">
            <Button
              :disabled="loading"
              @click="handleSave"
            >
              <Save class="w-4 h-4" />
              {{ loading ? '保存中...' : '保存配置' }}
            </Button>
            <Button
              variant="outline"
              :disabled="testing"
              @click="handleTest"
            >
              <Send class="w-4 h-4" />
              {{ testing ? '发送中...' : '测试推送' }}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
