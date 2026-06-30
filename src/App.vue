<script setup lang="ts">
import {
  computed,
  nextTick,
  ref,
  watch,
} from 'vue'
import {
  inferKeyword,
} from './card-engine'
import CardSvg from './components/CardSvg.vue'
import CardTile from './components/CardTile.vue'
import { templates } from './templates'
import type { CardTemplate } from './types'

const title = ref(
  '7月10号才能到岗，还能找到实习吗',
)

const customKeyword = ref('实习')
const autoKeyword = ref(true)
const selectedTemplateId = ref(templates[0].id)
const exporting = ref(false)

const keyword = computed(() => {
  if (!autoKeyword.value)
    return customKeyword.value

  return inferKeyword(title.value)
})

const selectedTemplate = computed(() =>
  templates.find(
    item => item.id === selectedTemplateId.value,
  ) ?? templates[0],
)

watch(
  title,
  () => {
    if (autoKeyword.value)
      customKeyword.value = inferKeyword(title.value)
  },
  {
    immediate: true,
  },
)

async function svgToPng(
  svg: SVGSVGElement,
  filename: string,
) {
  exporting.value = true

  try {
    await nextTick()

    if ('fonts' in document)
      await document.fonts.ready

    const cloned = svg.cloneNode(true) as SVGSVGElement

    cloned.setAttribute('width', '1800')
    cloned.setAttribute('height', '2400')

    const serialized =
      new XMLSerializer().serializeToString(cloned)

    const blob = new Blob(
      [
        `<?xml version="1.0" encoding="UTF-8"?>${serialized}`,
      ],
      {
        type: 'image/svg+xml;charset=utf-8',
      },
    )

    const url = URL.createObjectURL(blob)
    const image = new Image()

    await new Promise<void>((resolve, reject) => {
      image.onload = () => resolve()
      image.onerror = reject
      image.src = url
    })

    const canvas = document.createElement('canvas')
    canvas.width = 1800
    canvas.height = 2400

    const context = canvas.getContext('2d')

    if (!context)
      throw new Error('浏览器不支持 Canvas')

    context.drawImage(image, 0, 0, 1800, 2400)
    URL.revokeObjectURL(url)

    const pngUrl = canvas.toDataURL('image/png')
    const link = document.createElement('a')

    link.download = filename
    link.href = pngUrl
    link.click()
  }
  finally {
    exporting.value = false
  }
}

async function exportTemplate(template: CardTemplate) {
  selectedTemplateId.value = template.id
  await nextTick()

  const svg = document.querySelector(
    '#export-stage svg',
  ) as SVGSVGElement | null

  if (!svg)
    return

  await svgToPng(
    svg,
    `text-card-${template.id}.png`,
  )
}

function fillExample(value: string) {
  title.value = value
}

const examples = [
  '7月10号才能到岗，还能找到实习吗',
  '不是所有的坚持，都要立刻看到答案',
  '今天也要认真生活，慢慢变好',
  '成年人的崩溃，往往只差一句“你还好吗”',
]
</script>

<template>
  <main class="app-shell">
    <header class="hero">
      <div>
        <span class="eyebrow">Vue / SVG / Auto Layout</span>
        <h1>文字卡片生成器</h1>
        <p>
          输入任何标题，实时生成 12 套不同风格。
          每张卡片都会重新计算断行、字号、高亮位置和装饰布局。
        </p>
      </div>

      <div class="hero__badge">
        12
        <small>种模板</small>
      </div>
    </header>

    <section class="editor-panel">
      <label class="field field--wide">
        <span>标题内容</span>
        <textarea
          v-model="title"
          maxlength="140"
          rows="4"
          placeholder="输入一段你想说的话"
        />
        <small>{{ title.length }} / 140</small>
      </label>

      <div class="field">
        <span>高亮关键词</span>
        <div class="inline-control">
          <input
            v-model="customKeyword"
            :disabled="autoKeyword"
            placeholder="例如：实习"
          >

          <label class="toggle">
            <input
              v-model="autoKeyword"
              type="checkbox"
            >
            <span>自动识别</span>
          </label>
        </div>

        <small>
          当前高亮：
          <strong>{{ keyword || '无' }}</strong>
        </small>
      </div>

      <div class="field">
        <span>快速示例</span>
        <div class="example-list">
          <button
            v-for="example in examples"
            :key="example"
            type="button"
            @click="fillExample(example)"
          >
            {{ example }}
          </button>
        </div>
      </div>
    </section>

    <section class="section-heading">
      <div>
        <span class="eyebrow">GENERATED CARDS</span>
        <h2>同一段文字，不同的设计语法</h2>
      </div>

      <p>
        不是简单换背景色：每个模板都有自己的文字区域、
        对齐方式、装饰系统、最大行数和字号范围。
      </p>
    </section>

    <section class="card-grid">
      <CardTile
        v-for="template in templates"
        :key="template.id"
        :template="template"
        :title="title"
        :keyword="keyword"
        @export="exportTemplate"
      />
    </section>

    <div
      id="export-stage"
      aria-hidden="true"
    >
      <CardSvg
        :template="selectedTemplate"
        :title="title"
        :keyword="keyword"
      />
    </div>

    <div
      v-if="exporting"
      class="export-toast"
    >
      正在生成高清 PNG…
    </div>
  </main>
</template>
