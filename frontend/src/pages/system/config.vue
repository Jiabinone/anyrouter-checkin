<script setup lang="ts">
import { Settings } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold">
      更多设置
    </h1>

    <Card class="max-w-2xl">
      <CardHeader>
        <CardTitle>系统配置</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="flex flex-col items-center justify-center py-12 text-muted-foreground">
          <Settings class="w-16 h-16 mb-4" />
          <p>更多配置项开发中...</p>
          <p class="text-sm mt-2">
            可扩展: 代理设置、日志级别、数据备份等
          </p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
