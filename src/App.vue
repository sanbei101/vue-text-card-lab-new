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
  <main
    class="mx-auto min-h-screen min-w-[320px] w-[min(1480px,calc(100%-32px))]12 pb-20 font-sans text-zinc-900 max-[720px]:w-[min(calc(100%-20px),1480px)] max-[720px]:pt-3.5">
    <header
      class="flex items-end justify-between gap-8 rounded-[28px] border border-zinc-300 bg-[radial-gradient(circle_at_88%_20%,rgb(255_225_80/45%),transparent_24%),linear-gradient(135deg,#fff,#f7f7f8)] px-9.5 py-9 shadow-[0_16px_50px_rgb(0_0_0/6%)] max-[720px]:items-start max-[720px]:p-6.25">
      <div>
        <span class="text-xs font-extrabold tracking-[0.18em] text-zinc-500">Vue / SVG / Auto Layout</span>
        <h1 class="my-2.5 text-[clamp(36px,5vw,68px)] leading-none font-bold tracking-normal">
          文字卡片生成器
        </h1>
        <p class="m-0 max-w-180 text-[17px] leading-[1.7] text-zinc-500">
          输入任何标题，实时生成 12 套不同风格。
          每张卡片都会重新计算断行、字号、高亮位置和装饰布局。
        </p>
      </div>

      <div
        class="grid size-30.5 shrink-0 rotate-6 place-items-center rounded-full bg-[#ffe04a] text-5xl leading-[0.8] font-black text-zinc-950 max-[720px]:size-21.5 max-[720px]:text-[34px]">
        12
        <small class="text-sm tracking-[0.12em]">种模板</small>
      </div>
    </header>

    <section
      class="mt-6 grid grid-cols-[1.5fr_1fr_1fr] gap-4.5 rounded-3xl border border-zinc-300 bg-white p-5.5 max-[1080px]:grid-cols-2 max-[720px]:grid-cols-1">
      <label class="flex min-w-0 flex-col gap-2.5 max-[1080px]:col-span-full max-[720px]:col-auto">
        <span class="text-[13px] font-extrabold">标题内容</span>
        <textarea v-model="title"
          class="min-h-31.5 w-full resize-y rounded-[14px] border border-zinc-300 bg-zinc-50 px-4 py-3.5 text-lg leading-[1.65] text-zinc-900 outline-none transition-[border-color,box-shadow] duration-150 focus:border-zinc-900 focus:shadow-[0_0_0_3px_rgb(24_24_27/8%)]"
          maxlength="140" rows="4" placeholder="输入一段你想说的话" />
        <small class="text-zinc-400">{{ title.length }} / 140</small>
      </label>

      <div class="flex min-w-0 flex-col gap-2.5">
        <span class="text-[13px] font-extrabold">高亮关键词</span>
        <div class="flex gap-2.5">
          <input v-model="customKeyword"
            class="w-full rounded-[14px] border border-zinc-300 bg-zinc-50 px-3.5 py-3 text-zinc-900 outline-none transition-[border-color,box-shadow] duration-150 focus:border-zinc-900 focus:shadow-[0_0_0_3px_rgb(24_24_27/8%)]"
            :disabled="autoKeyword" placeholder="例如：实习">

          <label
            class="inline-flex shrink-0 items-center gap-1.5 whitespace-nowrap rounded-xl border border-zinc-300 bg-zinc-50 px-2.5 text-[13px]">
            <input v-model="autoKeyword" type="checkbox">
            <span>自动识别</span>
          </label>
        </div>

        <small class="text-zinc-400">
          当前高亮：
          <strong>{{ keyword || '无' }}</strong>
        </small>
      </div>

      <div class="flex min-w-0 flex-col gap-2.5">
        <span class="text-[13px] font-extrabold">快速示例</span>
        <div class="flex flex-wrap gap-2">
          <button v-for="example in examples" :key="example"
            class="max-w-full cursor-pointer overflow-hidden rounded-full border border-zinc-300 bg-zinc-50 px-2.75 py-2.25 text-zinc-600 text-ellipsis whitespace-nowrap hover:border-zinc-400 hover:bg-white hover:text-zinc-950"
            type="button" @click="fillExample(example)">
            {{ example }}
          </button>
        </div>
      </div>
    </section>

    <section
      class="mt-14 mb-5.5 flex items-end justify-between gap-10 max-[720px]:mt-10.5 max-[720px]:flex-col max-[720px]:items-start max-[720px]:gap-1.25">
      <div>
        <span class="text-xs font-extrabold tracking-[0.18em] text-zinc-500">GENERATED CARDS</span>
        <h2 class="my-2.5 text-[clamp(28px,4vw,46px)] font-bold tracking-normal">
          同一段文字，不同的设计语法
        </h2>
      </div>

      <p class="m-0 mb-1.25 max-w-127.5 leading-[1.7] text-zinc-500">
        不是简单换背景色：每个模板都有自己的文字区域、
        对齐方式、装饰系统、最大行数和字号范围。
      </p>
    </section>

    <section class="grid grid-cols-3 gap-5.5 max-[1080px]:grid-cols-2 max-[720px]:grid-cols-1">
      <CardTile v-for="template in templates" :key="template.id" :template="template" :title="title" :keyword="keyword"
        @export="exportTemplate" />
    </section>

    <div id="export-stage" class="pointer-events-none fixed top-[-10000px] left-[-10000px] -z-10 w-225"
      aria-hidden="true">
      <CardSvg :template="selectedTemplate" :title="title" :keyword="keyword" />
    </div>

    <div v-if="exporting"
      class="fixed right-6 bottom-6 z-100 rounded-[13px] bg-zinc-900 px-4.25 py-3.25 font-bold text-white shadow-[0_12px_30px_rgb(0_0_0/22%)]">
      正在生成高清 PNG…
    </div>
  </main>
</template>
