<script setup lang="ts">
import type { CardTemplate } from "@/types";

defineProps<{
  template: CardTemplate;
}>();
</script>

<template>
  <g :fill="template.foreground" opacity="0.62">
    <circle cx="112" cy="1078" r="18" />
    <text x="145" y="1087" font-size="24" :font-family="template.fontFamily" font-weight="600">
      此刻想说
    </text>
  </g>
</template>
