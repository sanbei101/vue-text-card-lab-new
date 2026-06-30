<script setup lang="ts">
import type { CardTemplate } from "../../types";
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
    />

    <!-- 杂志装饰 -->
    <g>
      <rect x="62" y="78" width="776" height="54" :fill="template.foreground" />
      <text
        x="78"
        y="118"
        fill="#fff"
        font-size="27"
        font-family="Arial, sans-serif"
        font-weight="700"
        letter-spacing="7"
      >
        DAILY NOTE / 2026
      </text>
      <rect x="62" y="1020" width="250" height="18" :fill="template.accent" />
      <text
        x="62"
        y="1102"
        :fill="template.foreground"
        font-size="28"
        font-family="Arial, sans-serif"
        letter-spacing="3"
      >
        WORDS MATTER.
      </text>
      <circle
        cx="760"
        cy="1045"
        r="78"
        fill="none"
        :stroke="template.foreground"
        stroke-width="4"
      />
      <text
        x="760"
        y="1060"
        text-anchor="middle"
        :fill="template.foreground"
        font-size="44"
        font-family="Arial Black, sans-serif"
      >
        12
      </text>
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />
  </svg>
</template>
