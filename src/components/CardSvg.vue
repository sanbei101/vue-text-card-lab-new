<script setup lang="ts">
import { computed } from "vue";

import { buildTextLayout, findHighlightRects, hashText, seededValue } from "../card-engine";
import type { CardTemplate } from "../types";

const props = defineProps<{
  template: CardTemplate;
  title: string;
  keyword: string;
}>();

const layout = computed(() => buildTextLayout(props.title, props.template));

const highlights = computed(() => findHighlightRects(layout.value, props.keyword, props.template));

const textAnchor = computed(() => {
  if (props.template.textBox.align === "center") return "middle";

  if (props.template.textBox.align === "right") return "end";

  return "start";
});

const textX = computed(() => {
  const box = props.template.textBox;

  if (box.align === "center") return box.x + box.width / 2;

  if (box.align === "right") return box.x + box.width;

  return box.x;
});

const seed = computed(() => hashText(`${props.title}:${props.template.id}`));

const particles = computed(() =>
  Array.from({ length: 14 }, (_, index) => ({
    x: 60 + seededValue(seed.value, index) * 780,
    y: 90 + seededValue(seed.value, index + 20) * 1020,
    size: 8 + seededValue(seed.value, index + 40) * 22,
    rotation: seededValue(seed.value, index + 60) * 180,
  })),
);
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

      <linearGradient :id="`gradient-${template.id}`" x1="0" y1="0" x2="1" y2="1">
        <stop offset="0" :stop-color="template.background" />
        <stop offset="1" :stop-color="template.muted" stop-opacity="0.36" />
      </linearGradient>

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

    <rect
      x="18"
      y="18"
      width="864"
      height="1164"
      :rx="template.radius"
      :fill="template.kind === 'lavender' ? `url(#gradient-${template.id})` : template.background"
      :stroke="template.frame ?? template.muted"
      :stroke-width="template.kind === 'mono' ? 7 : 2"
    />

    <!-- 清透问号 -->
    <g
      v-if="template.kind === 'question'"
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

    <!-- 便签 -->
    <g v-if="template.kind === 'memo'">
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

    <!-- 杂志 -->
    <g v-if="template.kind === 'editorial'">
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

    <!-- 柔雾紫 -->
    <g v-if="template.kind === 'lavender'">
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

    <!-- 绿格手账 -->
    <g v-if="template.kind === 'notebook'">
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

    <!-- 红印宣言 -->
    <g v-if="template.kind === 'stamp'">
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

    <!-- 蓝图网格 -->
    <g v-if="template.kind === 'grid'">
      <rect
        x="18"
        y="18"
        width="864"
        height="1164"
        :rx="template.radius"
        :fill="`url(#grid-${template.id})`"
      />
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

    <!-- 深夜星光 -->
    <g v-if="template.kind === 'night'">
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

    <!-- 橙色爆炸 -->
    <g v-if="template.kind === 'burst'">
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

    <!-- 纸张摘录 -->
    <g v-if="template.kind === 'paper'">
      <path
        d="M70 110
           Q450 70 830 110
           L815 1080
           Q450 1125 85 1080 Z"
        fill="#f9f5ec"
        :stroke="template.muted"
        stroke-width="3"
        filter="url(#shadow-beige-paper)"
      />
      <text
        x="118"
        y="190"
        :fill="template.accent"
        font-size="86"
        font-family="Georgia, serif"
        font-weight="700"
      >
        “
      </text>
      <path d="M120 975 H440" :stroke="template.accent" stroke-width="5" />
      <text
        x="120"
        y="1032"
        :fill="template.foreground"
        font-size="24"
        font-family="serif"
        letter-spacing="4"
      >
        A SMALL PIECE OF THOUGHT
      </text>
    </g>

    <!-- 粉色气泡 -->
    <g v-if="template.kind === 'bubble'">
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

    <!-- 极简边框 -->
    <g v-if="template.kind === 'mono'">
      <text
        x="90"
        y="120"
        :fill="template.foreground"
        font-size="26"
        font-family="Arial, sans-serif"
        letter-spacing="6"
      >
        TEXT CARD / 001
      </text>
      <line x1="90" y1="150" x2="810" y2="150" :stroke="template.foreground" stroke-width="3" />
      <line x1="90" y1="1035" x2="810" y2="1035" :stroke="template.foreground" stroke-width="3" />
      <circle cx="110" cy="1090" r="12" :fill="template.foreground" />
      <circle
        cx="150"
        cy="1090"
        r="12"
        fill="none"
        :stroke="template.foreground"
        stroke-width="3"
      />
    </g>

    <!-- 高亮层 -->
    <g
      v-for="(rect, index) in highlights"
      :key="`highlight-${index}`"
      :transform="`rotate(${rect.rotation} ${rect.x + rect.width / 2} ${rect.y + rect.height / 2})`"
    >
      <rect
        :x="rect.x"
        :y="rect.y"
        :width="rect.width"
        :height="rect.height"
        :rx="Math.min(18, rect.height / 2)"
        :fill="template.accent"
        :opacity="template.kind === 'night' ? 0.9 : 0.78"
      />
    </g>

    <!-- 主文字 -->
    <g
      :fill="template.foreground"
      :font-family="template.fontFamily"
      :font-weight="template.fontWeight"
      :text-anchor="textAnchor"
    >
      <text
        v-for="(line, index) in layout.lines"
        :key="`${line}-${index}`"
        :x="textX"
        :y="template.textBox.y + layout.fontSize + index * layout.lineHeightPx"
        :font-size="layout.fontSize"
        dominant-baseline="alphabetic"
      >
        {{ line }}
      </text>
    </g>

    <!-- 小型签名 -->
    <g
      v-if="!['memo', 'editorial', 'paper', 'mono'].includes(template.kind)"
      :fill="template.foreground"
      opacity="0.62"
    >
      <circle cx="112" cy="1078" r="18" />
      <text x="145" y="1087" font-size="24" :font-family="template.fontFamily" font-weight="600">
        此刻想说
      </text>
    </g>
  </svg>
</template>
