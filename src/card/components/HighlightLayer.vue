<script setup lang="ts">
import type { CardTemplate, HighlightRect } from "@/types";

defineProps<{
  template: CardTemplate;
  rects: HighlightRect[];
  highlightOpacity?: number;
}>();
</script>

<template>
  <g
    v-for="(rect, index) in rects"
    :key="`highlight-${index}`"
    :transform="`rotate(${rect.rotation} ${rect.x + rect.width / 2} ${rect.y + rect.height / 2})`"
  >
    <rect
      :x="rect.x"
      :y="rect.y"
      :width="rect.width"
      :height="rect.height"
      :rx="Math.min(18, rect.height / 2)"
      :fill="template.accent"
      :opacity="highlightOpacity ?? 0.78"
    />
  </g>
</template>
