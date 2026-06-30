<script setup lang="ts">
import CardSvg from './CardSvg.vue'
import type { CardTemplate } from '../types'

const props = defineProps<{
  template: CardTemplate
  title: string
  keyword: string
}>()

const emit = defineEmits<{
  export: [template: CardTemplate]
}>()
</script>

<template>
  <article
    class="overflow-hidden rounded-3xl border border-zinc-300 bg-white shadow-[0_10px_35px_rgb(0_0_0_/_4%)] transition-[transform,box-shadow] duration-150 hover:-translate-y-[3px] hover:shadow-[0_18px_50px_rgb(0_0_0_/_9%)]">
    <div
      class="bg-[linear-gradient(45deg,#f5f5f5_25%,transparent_25%)_0_0_/_18px_18px,linear-gradient(45deg,transparent_75%,#f5f5f5_75%)_0_0_/_18px_18px,linear-gradient(45deg,transparent_75%,#f5f5f5_75%)_9px_-9px_/_18px_18px,linear-gradient(45deg,#f5f5f5_25%,#fff_25%)_9px_-9px_/_18px_18px] p-3.5">
      <CardSvg :template="template" :title="title" :keyword="keyword" />
    </div>

    <footer class="flex items-center justify-between gap-3.5 px-[17px] pt-[15px] pb-[17px]">
      <div class="grid gap-[3px]">
        <strong class="text-[15px]">{{ template.name }}</strong>
        <span class="text-xs text-zinc-400 uppercase">{{ template.kind }}</span>
      </div>

      <button type="button"
        class="shrink-0 cursor-pointer rounded-[10px] border border-zinc-300 bg-white px-3 py-2.25 text-[13px] font-bold text-zinc-800 hover:border-zinc-900 hover:bg-zinc-900 hover:text-white"
        @click="emit('export', template)">
        导出 PNG
      </button>
    </footer>
  </article>
</template>
