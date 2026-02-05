<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Save, Send, FileText } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { getConfigs, updateConfigs, testTelegram } from '@/api/system'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'

const recommendedTemplate = [
  '<b>AnyRouter 签到系统</b>',
  '用户名：<code>{{.Username}}</code>',
  '状态：{{if .Success}}<b>成功 ✅</b>{{else}}<b>失败 ❌</b>{{end}}',
  '结果：',
  '<pre>{{.Result}}</pre>',
].join('\n')

const form = ref({
  'telegram.enabled': 'false',
  'telegram.api_base': 'https://api.telegram.org',
  'telegram.bot_token': '',
  'telegram.chat_id': '',
  'telegram.proxy_url': '',
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
}

async function handleSave() {
  loading.value = true
  try {
    await updateConfigs('telegram', form.value)
    toast.success('保存成功')
  } catch (e) {
    const message = e instanceof Error ? e.message : '保存失败'
    toast.error('保存失败: ' + message)
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

function applyTemplate() {
  form.value['telegram.template'] = recommendedTemplate
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
            <Label>API 地址 / 反代地址</Label>
            <Input
              v-model="form['telegram.api_base']"
              class="font-mono text-sm"
              placeholder="https://api.telegram.org"
            />
            <p class="text-xs text-muted-foreground">
              国内无法直连可填写反代地址
            </p>
          </div>

          <div class="space-y-2">
            <Label>代理地址</Label>
            <Input
              v-model="form['telegram.proxy_url']"
              class="font-mono text-sm"
              placeholder="http://127.0.0.1:7890"
            />
            <p class="text-xs text-muted-foreground">
              支持 HTTP/HTTPS 代理（可选）
            </p>
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
              可用变量: {{.Username}} (用户名), {{.Success}} 是否成功, {{.Result}} 结果（支持 Telegram HTML）
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
            <Button
              variant="outline"
              @click="applyTemplate"
            >
              <FileText class="w-4 h-4" />
              预设模板
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
