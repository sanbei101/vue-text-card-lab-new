<script setup lang="ts">
import type { CardTemplate, TextLayout } from "@/types";

defineProps<{
  template: CardTemplate;
  layout: TextLayout;
  x: number;
  anchor: "start" | "middle" | "end";
}>();
</script>

<template>
  <g
    :fill="template.foreground"
    :font-family="template.fontFamily"
    :font-weight="template.fontWeight"
    :text-anchor="anchor"
  >
    <text
      v-for="(line, index) in layout.lines"
      :key="`${line}-${index}`"
      :x="x"
      :y="template.textBox.y + layout.fontSize + index * layout.lineHeightPx"
      :font-size="layout.fontSize"
      dominant-baseline="alphabetic"
    >
      {{ line }}
    </text>
  </g>
</template>
