<script setup lang="ts">
import type { CardTemplate } from "@/types";
import CardSignature from "@/card/components/CardSignature.vue";
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
    <!-- 背景 -->
    <rect
      x="18"
      y="18"
      width="864"
      height="1164"
      :rx="template.radius"
      :fill="template.background"
    />

    <!-- 粉色气泡装饰 -->
    <g>
      <circle cx="155" cy="210" r="92" fill="#fff" opacity="0.6" />
      <circle cx="750" cy="310" r="130" :fill="template.accent" opacity="0.12" />
      <circle cx="680" cy="1030" r="170" fill="#fff" opacity="0.5" />
      <path
        d="M115 180 C210 90 355 165 350 280 C345 390 195 410 125 330 C72 270 75 220 115 180 Z"
        fill="#fff"
        opacity="0.7"
      />
      <path
        d="M300 910 C430 825 575 870 640 965"
        fill="none"
        :stroke="template.accent"
        stroke-width="12"
        stroke-linecap="round"
        opacity="0.4"
      />
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />

    <!-- 签名 -->
    <CardSignature :template="template" />
  </svg>
</template>
