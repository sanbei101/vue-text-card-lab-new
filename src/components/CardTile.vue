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
  <article class="card-tile">
    <div class="card-tile__preview">
      <CardSvg
        :template="template"
        :title="title"
        :keyword="keyword"
      />
    </div>

    <footer class="card-tile__footer">
      <div>
        <strong>{{ template.name }}</strong>
        <span>{{ template.kind }}</span>
      </div>

      <button
        type="button"
        class="ghost-button"
        @click="emit('export', template)"
      >
        导出 PNG
      </button>
    </footer>
  </article>
</template>
