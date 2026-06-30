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

    <!-- 红印宣言装饰 -->
    <g>
      <rect
        x="55"
        y="55"
        width="790"
        height="1090"
        rx="8"
        fill="none"
        :stroke="template.accent"
        stroke-width="6"
      />
      <rect
        x="80"
        y="80"
        width="740"
        height="1040"
        rx="4"
        fill="none"
        :stroke="template.accent"
        stroke-width="2"
        stroke-dasharray="18 12"
        opacity="0.65"
      />
      <g transform="translate(690 945) rotate(-12)" opacity="0.8">
        <rect width="125" height="125" fill="none" :stroke="template.accent" stroke-width="10" />
        <text
          x="62"
          y="56"
          text-anchor="middle"
          :fill="template.accent"
          font-size="30"
          font-family="serif"
          font-weight="900"
        >
          今日
        </text>
        <text
          x="62"
          y="96"
          text-anchor="middle"
          :fill="template.accent"
          font-size="30"
          font-family="serif"
          font-weight="900"
        >
          要说
        </text>
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
