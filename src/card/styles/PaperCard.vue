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
      <filter id="shadow-beige-paper" x="-20%" y="-20%" width="140%" height="140%">
        <feDropShadow dx="0" dy="12" stdDeviation="14" flood-opacity="0.12" />
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

    <!-- 纸张摘录装饰 -->
    <g>
      <path
        d="M70 110
           Q450 70 830 110
           L815 1080
           Q450 1125 85 1080 Z"
        fill="#f9f5ec"
        :stroke="template.muted"
        stroke-width="3"
        filter="url(#shadow-beige-paper)"
      />
      <text
        x="118"
        y="190"
        :fill="template.accent"
        font-size="86"
        font-family="Georgia, serif"
        font-weight="700"
      >
        "
      </text>
      <path d="M120 975 H440" :stroke="template.accent" stroke-width="5" />
      <text
        x="120"
        y="1032"
        :fill="template.foreground"
        font-size="24"
        font-family="serif"
        letter-spacing="4"
      >
        A SMALL PIECE OF THOUGHT
      </text>
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />
  </svg>
</template>
