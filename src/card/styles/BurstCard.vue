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

    <!-- 橙色爆炸装饰 -->
    <g>
      <g :stroke="template.accent" stroke-width="18" stroke-linecap="round" opacity="0.78">
        <line x1="450" y1="90" x2="450" y2="180" />
        <line x1="720" y1="160" x2="660" y2="230" />
        <line x1="815" y1="430" x2="720" y2="430" />
        <line x1="720" y1="980" x2="655" y2="900" />
        <line x1="180" y1="980" x2="245" y2="900" />
        <line x1="80" y1="430" x2="175" y2="430" />
        <line x1="180" y1="160" x2="245" y2="230" />
      </g>
      <path
        d="M450 155
           L525 270
           L660 245
           L650 385
           L785 450
           L670 540
           L720 675
           L575 690
           L520 825
           L405 740
           L285 830
           L240 690
           L95 675
           L145 535
           L30 450
           L165 385
           L150 245
           L285 270 Z"
        fill="#ffb778"
        opacity="0.34"
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