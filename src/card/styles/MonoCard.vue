<script setup lang="ts">
import { useCardRender } from "../composables/useCardRender";
import HighlightLayer from "../components/HighlightLayer.vue";
import MainText from "../components/MainText.vue";
import type { CardTemplate } from "../../types";

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
    <!-- 背景 -->
    <rect
      x="18"
      y="18"
      width="864"
      height="1164"
      :rx="template.radius"
      :fill="template.background"
      :stroke="template.frame ?? template.muted"
      stroke-width="7"
    />

    <!-- 极简边框装饰 -->
    <g>
      <text
        x="90"
        y="120"
        :fill="template.foreground"
        font-size="26"
        font-family="Arial, sans-serif"
        letter-spacing="6"
      >
        TEXT CARD / 001
      </text>
      <line x1="90" y1="150" x2="810" y2="150" :stroke="template.foreground" stroke-width="3" />
      <line x1="90" y1="1035" x2="810" y2="1035" :stroke="template.foreground" stroke-width="3" />
      <circle cx="110" cy="1090" r="12" :fill="template.foreground" />
      <circle
        cx="150"
        cy="1090"
        r="12"
        fill="none"
        :stroke="template.foreground"
        stroke-width="3"
      />
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />
  </svg>
</template>