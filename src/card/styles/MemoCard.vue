<script setup lang="ts">
import type { CardTemplate } from "@/types";

import HighlightLayer from "@/card/components/HighlightLayer.vue";
import MainText from "@/card/components/MainText.vue";
import { useCardRender } from "@/card/composables/useCardRender";

const props = defineProps<{ template: CardTemplate; title: string; keyword: string }>();
const { layout, highlights, textX, textAnchor } = useCardRender(props);
</script>

<template>
  <svg
    class="block h-auto w-full overflow-visible"
    viewBox="0 0 900 1200"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    :aria-label="template.name"
  >
    <defs>
      <filter id="shadow-yellow-memo" x="-20%" y="-20%" width="140%" height="140%">
        <feDropShadow dx="0" dy="12" stdDeviation="14" flood-opacity="0.15" />
      </filter>
    </defs>

    <!-- 背景 -->
    <rect
      x="18"
      y="18"
      width="864"
      height="1164"
      :rx="template.radius"
      :fill="template.background"
    />

    <!-- 便签装饰 -->
    <g>
      <path
        d="M40 40 H860 V1100 L802 1100 H40 Z"
        fill="#fff9cf"
        :stroke="template.muted"
        stroke-width="3"
        filter="url(#shadow-yellow-memo)"
      />
      <path d="M802 1100 H860 L802 1124 Z" fill="#f0d66d" />
      <rect
        x="120"
        y="90"
        width="210"
        height="78"
        rx="10"
        fill="#ffd4ba"
        opacity="0.85"
        transform="rotate(-5 225 90)"
      />
      <path
        d="M145 920 C240 855 325 985 430 908 S640 885 760 935"
        fill="none"
        :stroke="template.accent"
        stroke-width="12"
        stroke-linecap="round"
        opacity="0.65"
      />
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />
  </svg>
</template>
