<script setup lang="ts">
import CardSvg from "@/card/CardSvg.vue";
import type { CardTemplate } from "@/types";

const { template, title, keyword } = defineProps<{
  template: CardTemplate;
  title: string;
  keyword: string;
}>();

const emit = defineEmits<{
  export: [template: CardTemplate];
}>();
</script>

<template>
  <article
    class="overflow-hidden rounded-3xl border border-zinc-300 bg-white shadow-md transition-[transform,box-shadow] duration-150 hover:-translate-y-0.75 hover:shadow-[0_18px_50px_rgb(0_0_0/9%)]"
  >
    <div class="checkerboard-bg p-3.5">
      <CardSvg :template="template" :title="title" :keyword="keyword" />
    </div>

    <footer class="flex items-center justify-between gap-3.5 px-4.25 pt-3.75 pb-4.25">
      <div class="grid gap-0.75">
        <strong class="text-[15px]">{{ template.name }}</strong>
        <span class="text-xs text-zinc-400 uppercase">{{ template.kind }}</span>
      </div>

      <button
        type="button"
        class="shrink-0 cursor-pointer rounded-[10px] border border-zinc-300 bg-white px-3 py-2.25 text-[13px] font-bold text-zinc-800 hover:border-zinc-900 hover:bg-zinc-900 hover:text-white"
        @click="emit('export', template)"
      >
        导出 PNG
      </button>
    </footer>
  </article>
</template>

<style scoped>
.checkerboard-bg {
  background:
    linear-gradient(45deg, #f5f5f5 25%, transparent 25%) 0 0 / 18px 18px,
    linear-gradient(45deg, transparent 75%, #f5f5f5 75%) 0 0 / 18px 18px,
    linear-gradient(45deg, transparent 75%, #f5f5f5 75%) 9px -9px / 18px 18px,
    linear-gradient(45deg, #f5f5f5 25%, #fff 25%) 9px -9px / 18px 18px;
}
</style>
