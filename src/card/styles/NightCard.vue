<script setup lang="ts">
import type { CardTemplate } from "@/types";
import CardSignature from "@/card/components/CardSignature.vue";
import HighlightLayer from "@/card/components/HighlightLayer.vue";
import MainText from "@/card/components/MainText.vue";
import { useCardRender } from "@/card/composables/useCardRender";

const props = defineProps<{ template: CardTemplate; title: string; keyword: string }>();
const { layout, highlights, textX, textAnchor, particles } = useCardRender(props);
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

    <!-- 深夜星光装饰 -->
    <g>
      <circle cx="710" cy="220" r="116" :fill="template.accent" opacity="0.96" />
      <circle cx="755" cy="185" r="120" :fill="template.background" />
      <g
        v-for="(particle, index) in particles"
        :key="`star-${index}`"
        :transform="`translate(${particle.x} ${particle.y})`"
        :fill="index % 3 === 0 ? template.accent : template.foreground"
        :opacity="0.25 + (index % 4) * 0.14"
      >
        <circle v-if="index % 2" :r="particle.size / 4" />
        <path v-else d="M0 -11 L3 -3 L11 0 L3 3 L0 11 L-3 3 L-11 0 L-3 -3 Z" />
      </g>
      <path
        d="M135 980 C315 900 590 1070 770 955"
        fill="none"
        :stroke="template.muted"
        stroke-width="4"
      />
    </g>

    <!-- 高亮层（night 使用更高透明度） -->
    <HighlightLayer :template="template" :rects="highlights" :highlight-opacity="0.9" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />

    <!-- 签名 -->
    <CardSignature :template="template" />
  </svg>
</template>
