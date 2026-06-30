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
    <defs>
      <linearGradient :id="`gradient-${template.id}`" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" :stop-color="template.background" />
        <stop offset="1" :stop-color="template.muted" stop-opacity="0.36" />
      </linearGradient>
    </defs>

    <!-- 背景（渐变） -->
    <rect
      x="18"
      y="18"
      width="864"
      height="1164"
      :rx="template.radius"
      :fill="`url(#gradient-${template.id})`"
      :stroke="template.muted"
      stroke-width="2"
    />

    <!-- 柔雾紫装饰 -->
    <g>
      <circle cx="160" cy="180" r="128" fill="#ffffff" opacity="0.42" />
      <circle cx="760" cy="250" r="170" :fill="template.accent" opacity="0.09" />
      <path
        d="M145 950 C285 850 410 1040 575 920 S775 850 830 960"
        fill="none"
        :stroke="template.accent"
        stroke-width="7"
        stroke-linecap="round"
        opacity="0.4"
      />
      <g
        v-for="(particle, index) in particles.slice(0, 7)"
        :key="index"
        :transform="`translate(${particle.x} ${particle.y}) rotate(${particle.rotation})`"
      >
        <path
          d="M0 -14 L4 -4 L14 0 L4 4 L0 14 L-4 4 L-14 0 L-4 -4 Z"
          :fill="template.accent"
          opacity="0.3"
        />
      </g>
    </g>

    <!-- 高亮层 -->
    <HighlightLayer :template="template" :rects="highlights" />

    <!-- 主文字 -->
    <MainText :template="template" :layout="layout" :x="textX" :anchor="textAnchor" />

    <!-- 签名 -->
    <CardSignature :template="template" />
  </svg>
</template>
