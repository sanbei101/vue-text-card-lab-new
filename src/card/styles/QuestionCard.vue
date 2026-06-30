<script setup lang="ts">
import type { CardTemplate } from "../../types";
import CardSignature from "../components/CardSignature.vue";
import HighlightLayer from "../components/HighlightLayer.vue";
import MainText from "../components/MainText.vue";
import { useCardRender } from "../composables/useCardRender";

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
      <filter :id="`shadow-${template.id}`" x="-20%" y="-20%" width="140%" height="140%">
        <feDropShadow dx="0" dy="16" stdDeviation="18" flood-opacity="0.12" />
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
      :stroke="template.muted"
      stroke-width="2"
    />

    <!-- 清透问号装饰 -->
    <g
      :fill="template.muted"
      opacity="0.18"
      font-family="Arial Black, sans-serif"
      font-weight="900"
    >
      <text x="-80" y="260" font-size="310" transform="rotate(-18 -80 260)">?</text>
      <text x="300" y="390" font-size="430" transform="rotate(32 300 390)">?</text>
      <text x="560" y="1040" font-size="330" transform="rotate(-20 560 1040)">?</text>
      <circle cx="735" cy="490" r="72" />
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />

    <!-- 签名 -->
    <CardSignature :template="template" />
  </svg>
</template>
