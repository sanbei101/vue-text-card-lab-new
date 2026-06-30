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
    <defs>
      <pattern :id="`grid-${template.id}`" width="52" height="52" patternUnits="userSpaceOnUse">
        <path
          d="M 52 0 L 0 0 0 52"
          fill="none"
          :stroke="template.muted"
          stroke-width="2"
          opacity="0.48"
        />
      </pattern>
    </defs>

    <!-- 背景（网格图案） -->
    <rect
      x="18"
      y="18"
      width="864"
      height="1164"
      :rx="template.radius"
      :fill="`url(#grid-${template.id})`"
    />

    <!-- 蓝图网格装饰 -->
    <g>
      <circle
        cx="740"
        cy="190"
        r="115"
        fill="none"
        :stroke="template.accent"
        stroke-width="18"
        opacity="0.18"
      />
      <path d="M90 1020 H810" :stroke="template.accent" stroke-width="8" />
      <path d="M90 1052 H570" :stroke="template.foreground" stroke-width="3" opacity="0.7" />
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />

    <!-- 签名 -->
    <CardSignature :template="template" />
  </svg>
</template>
