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
        d="M92 170 H808 V970 L740 1038 H92 Z"
        fill="#fff9cf"
        :stroke="template.muted"
        stroke-width="3"
        filter="url(#shadow-yellow-memo)"
      />
      <path d="M740 970 H808 L740 1038 Z" fill="#f0d66d" />
      <rect
        x="180"
        y="110"
        width="210"
        height="78"
        rx="10"
        fill="#ffd4ba"
        opacity="0.85"
        transform="rotate(-5 180 110)"
      />
      <path
        d="M135 930 C230 865 315 995 420 918 S635 895 760 945"
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