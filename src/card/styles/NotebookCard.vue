<script setup lang="ts">
import { useCardRender } from "../composables/useCardRender";
import HighlightLayer from "../components/HighlightLayer.vue";
import MainText from "../components/MainText.vue";
import CardSignature from "../components/CardSignature.vue";
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
    />

    <!-- 绿格手账装饰 -->
    <g>
      <rect x="52" y="52" width="796" height="1096" rx="28" fill="#f7fbeF" opacity="0.74" />
      <g :stroke="template.muted" stroke-width="2" opacity="0.45">
        <line
          v-for="y in 14"
          :key="`line-${y}`"
          x1="80"
          x2="820"
          :y1="150 + y * 62"
          :y2="150 + y * 62"
        />
      </g>
      <line
        x1="145"
        y1="90"
        x2="145"
        y2="1110"
        :stroke="template.accent"
        stroke-width="4"
        opacity="0.6"
      />
      <circle v-for="y in 8" :key="`hole-${y}`" cx="74" :cy="150 + y * 120" r="11" fill="#d7e4cb" />
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />

    <!-- 签名 -->
    <CardSignature :template="template" />
  </svg>
</template>