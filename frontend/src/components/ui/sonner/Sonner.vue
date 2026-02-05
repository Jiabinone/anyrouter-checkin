<script lang="ts" setup>
import type { ToasterProps } from "vue-sonner"
import { reactiveOmit } from "@vueuse/core"
import { CheckCircle2, Info, Loader2, XOctagon, AlertTriangle, X } from "lucide-vue-next"
import { Toaster as Sonner } from "vue-sonner"

const props = defineProps<ToasterProps>()
const delegatedProps = reactiveOmit(props, "toastOptions")
</script>

<template>
  <Sonner
    class="toaster group z-[9999]"
    :toast-options="{
      classes: {
        toast: 'group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:shadow-lg',
        description: 'group-[.toast]:text-muted-foreground',
        actionButton:
          'group-[.toast]:bg-primary group-[.toast]:text-primary-foreground',
        cancelButton:
          'group-[.toast]:bg-muted group-[.toast]:text-muted-foreground',
      },
    }"
    v-bind="delegatedProps"
  >
    <template #success-icon>
      <CheckCircle2 class="size-4" />
    </template>
    <template #info-icon>
      <Info class="size-4" />
    </template>
    <template #warning-icon>
      <AlertTriangle class="size-4" />
    </template>
    <template #error-icon>
      <XOctagon class="size-4" />
    </template>
    <template #loading-icon>
      <div>
        <Loader2 class="size-4 animate-spin" />
      </div>
    </template>
    <template #close-icon>
      <X class="size-4" />
    </template>
  </Sonner>
</template>
