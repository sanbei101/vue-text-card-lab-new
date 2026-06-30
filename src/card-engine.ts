import {
  clearCache as clearPretextCache,
  layoutWithLines,
  measureLineStats,
  measureNaturalWidth,
  prepareWithSegments,
  setLocale as setPretextLocale,
} from "@chenglou/pretext";
import type { PreparedTextWithSegments } from "@chenglou/pretext";

import type { CardTemplate, HighlightRect, TextLayout } from "@/types";

const DEFAULT_EMPTY_TEXT = "输入一个标题";
const PREPARED_CACHE_LIMIT = 320;

const graphemeSegmenter = new Intl.Segmenter(undefined, {
  granularity: "grapheme",
});

const BREAK_AFTER = new Set([
  "，",
  "。",
  "！",
  "？",
  "、",
  "；",
  "：",
  ",",
  ".",
  "!",
  "?",
  ";",
  ":",
  " ",
]);

const BAD_LINE_START = new Set([
  "，",
  "。",
  "！",
  "？",
  "、",
  "；",
  "：",
  ",",
  ".",
  "!",
  "?",
  ";",
  ":",
  ")",
  "）",
  "]",
  "】",
  "》",
  "」",
  "』",
]);

const BAD_LINE_END = new Set(["(", "（", "[", "【", "《", "「", "『"]);

/**
 * 越靠前越倾向使用完整文字框；
 * 越靠后越容易得到更均衡、更“海报化”的短行。
 */
const WIDTH_RATIOS = [1, 0.97, 0.94, 0.91, 0.88, 0.85, 0.82, 0.79, 0.76, 0.73] as const;

interface LayoutLineValue {
  text: string;
  width: number;
}

interface LayoutCandidate {
  lines: LayoutLineValue[];
  wrapWidth: number;
  score: number;
}

interface TextBoxWithOptionalLetterSpacing {
  letterSpacing?: number;
}

const preparedCache = new Map<string, PreparedTextWithSegments>();

function normalizeText(value: string): string {
  return value
    .replace(/\r\n?/g, "\n")
    .replace(/[ \t]+/g, " ")
    .trim();
}

function getGraphemes(value: string): string[] {
  return Array.from(graphemeSegmenter.segment(value), (segment) => segment.segment);
}

function graphemeLength(value: string): number {
  let count = 0;

  for (const _segment of graphemeSegmenter.segment(value)) {
    count += 1;
  }

  return count;
}

function firstGrapheme(value: string): string {
  return Array.from(graphemeSegmenter.segment(value))[0]?.segment ?? "";
}

function lastGrapheme(value: string): string {
  let result = "";

  for (const segment of graphemeSegmenter.segment(value)) {
    result = segment.segment;
  }

  return result;
}

function clamp(value: number, minimum: number, maximum: number): number {
  return Math.min(maximum, Math.max(minimum, value));
}

function getLetterSpacing(template: CardTemplate): number {
  const textBox = template.textBox as typeof template.textBox & TextBoxWithOptionalLetterSpacing;

  return textBox.letterSpacing ?? 0;
}

function createFont(template: CardTemplate, fontSize: number): string {
  return [template.fontWeight, `${fontSize}px`, template.fontFamily].join(" ");
}

function createPreparedCacheKey(text: string, font: string, letterSpacing: number): string {
  return [text, font, letterSpacing, "pre-wrap", "normal"].join("\u0000");
}

function setPreparedCache(key: string, value: PreparedTextWithSegments): void {
  preparedCache.set(key, value);

  if (preparedCache.size <= PREPARED_CACHE_LIMIT) {
    return;
  }

  const oldestKey = preparedCache.keys().next().value;

  if (typeof oldestKey === "string") {
    preparedCache.delete(oldestKey);
  }
}

function getPrepared(
  rawText: string,
  fontSize: number,
  template: CardTemplate,
): PreparedTextWithSegments {
  const text = normalizeText(rawText) || DEFAULT_EMPTY_TEXT;
  const font = createFont(template, fontSize);
  const letterSpacing = getLetterSpacing(template);

  const key = createPreparedCacheKey(text, font, letterSpacing);

  const cached = preparedCache.get(key);

  if (cached) {
    // Map 同时充当简单 LRU：命中后移动到最后。
    preparedCache.delete(key);
    preparedCache.set(key, cached);

    return cached;
  }

  const prepared = prepareWithSegments(text, font, {
    whiteSpace: "pre-wrap",
    wordBreak: "normal",
    letterSpacing,
  });

  setPreparedCache(key, prepared);

  return prepared;
}

/**
 * 清理本地缓存以及 Pretext 自身的共享缓存。
 *
 * 建议在以下情况调用：
 * - 动态字体加载完毕
 * - 切换语言/locale
 * - 大量模板被卸载
 */
export function clearCardEngineCache(): void {
  preparedCache.clear();
  clearPretextCache();
}

/**
 * 切换 Pretext 的分段 locale。
 * 已经创建的 prepared 对象不会被修改，因此这里同步清理本地缓存。
 */
export function setCardEngineLocale(locale?: string): void {
  preparedCache.clear();
  setPretextLocale(locale);
}

/**
 * 单行自然宽度测量。
 *
 * 仍保留和旧 card-engine 相同的导出 API，
 * 但底层已经改为 Pretext。
 */
export function measureText(value: string, fontSize: number, template: CardTemplate): number {
  if (!value) {
    return 0;
  }

  const prepared = getPrepared(value, fontSize, template);

  return measureNaturalWidth(prepared);
}

function createCandidateWidths(
  prepared: PreparedTextWithSegments,
  template: CardTemplate,
): number[] {
  const boxWidth = template.textBox.width;
  const minimumWidth = boxWidth * 0.68;
  const widths = new Set<number>();

  for (const ratio of WIDTH_RATIOS) {
    widths.add(Math.round(boxWidth * ratio * 10) / 10);
  }

  /**
   * 根据自然宽度和目标行数补充候选宽度。
   * 这能比固定 ratios 更容易找到视觉上均衡的断行点。
   */
  const naturalWidth = measureNaturalWidth(prepared);
  const maximumTargetLines = Math.min(template.textBox.maxLines, 6);

  for (let targetLines = 2; targetLines <= maximumTargetLines; targetLines += 1) {
    for (const factor of [0.94, 1, 1.06]) {
      const width = clamp((naturalWidth / targetLines) * factor, minimumWidth, boxWidth);

      widths.add(Math.round(width * 10) / 10);
    }
  }

  return [...widths].sort((a, b) => b - a);
}

function getIdealLineCount(text: string, maximumLines: number): number {
  const length = graphemeLength(text.replace(/\s/g, ""));

  if (length <= 9) {
    return 2;
  }

  if (length <= 18) {
    return 3;
  }

  if (length <= 32) {
    return 4;
  }

  return Math.min(5, maximumLines);
}

function scoreLinePunctuation(line: string, isLast: boolean): number {
  const first = firstGrapheme(line);
  const last = lastGrapheme(line);
  let score = 0;

  if (BAD_LINE_START.has(first)) {
    score += 1800;
  }

  if (BAD_LINE_END.has(last)) {
    score += 1100;
  }

  if (!isLast && BREAK_AFTER.has(last)) {
    score -= 95;
  }

  return score;
}

function scoreLayout(
  lines: LayoutLineValue[],
  wrapWidth: number,
  boxWidth: number,
  maximumLines: number,
  text: string,
): number {
  if (!lines.length) {
    return Number.POSITIVE_INFINITY;
  }

  const fillRatios = lines.map((line) => clamp(line.width / wrapWidth, 0, 1.5));

  const averageFill = fillRatios.reduce((sum, value) => sum + value, 0) / fillRatios.length;

  const fillVariance =
    fillRatios.reduce((sum, value) => sum + (value - averageFill) ** 2, 0) / fillRatios.length;

  let score = 0;

  // 行宽越均衡越好。
  score += fillVariance * 1450;

  // 避免整体过于松散。
  score += Math.max(0, 0.67 - averageFill) ** 2 * 1200;

  // 使用过窄的排版宽度时略加惩罚。
  score += (1 - wrapWidth / boxWidth) ** 2 * 360;

  const idealLineCount = getIdealLineCount(text, maximumLines);

  score += (lines.length - idealLineCount) ** 2 * 34;

  lines.forEach((line, index) => {
    const isLast = index === lines.length - 1;
    const length = graphemeLength(line.text);
    const fill = fillRatios[index];

    score += scoreLinePunctuation(line.text, isLast);

    // 非末行太短通常不好看。
    if (!isLast && fill < 0.54) {
      score += (0.54 - fill) ** 2 * 1700;
    }

    // 卡片标题出现 1～2 个字的孤行通常很突兀。
    if (length <= 2) {
      score += isLast ? 780 : 1250;
    } else if (length === 3 && isLast) {
      score += 160;
    }
  });

  const lastFill = fillRatios[fillRatios.length - 1];

  if (lastFill < 0.38) {
    score += (0.38 - lastFill) ** 2 * 2500;
  }

  // 最后一行明显比倒数第二行长时，视觉重心容易下坠。
  if (lines.length >= 2) {
    const previousFill = fillRatios[fillRatios.length - 2];

    if (lastFill > previousFill + 0.28) {
      score += (lastFill - previousFill) ** 2 * 520;
    }
  }

  return score;
}

function findBestCandidateAtSize(
  rawText: string,
  fontSize: number,
  template: CardTemplate,
): LayoutCandidate | null {
  const text = normalizeText(rawText) || DEFAULT_EMPTY_TEXT;

  const box = template.textBox;
  const lineHeightPx = fontSize * box.lineHeight;

  const prepared = getPrepared(text, fontSize, template);

  const candidateWidths = createCandidateWidths(prepared, template);

  let best: LayoutCandidate | null = null;

  for (const wrapWidth of candidateWidths) {
    const stats = measureLineStats(prepared, wrapWidth);

    if (
      stats.lineCount === 0 ||
      stats.lineCount > box.maxLines ||
      stats.lineCount * lineHeightPx > box.height
    ) {
      continue;
    }

    const result = layoutWithLines(prepared, wrapWidth, lineHeightPx);

    const lines = result.lines.map((line) => ({
      text: line.text,
      width: line.width,
    }));

    const score = scoreLayout(lines, wrapWidth, box.width, box.maxLines, text);

    if (!best || score < best.score) {
      best = {
        lines,
        wrapWidth,
        score,
      };
    }
  }

  return best;
}

function fitWithEllipsis(
  value: string,
  maximumWidth: number,
  fontSize: number,
  template: CardTemplate,
): string {
  const ellipsis = "…";
  const graphemes = getGraphemes(value);

  if (measureText(`${value}${ellipsis}`, fontSize, template) <= maximumWidth) {
    return `${value}${ellipsis}`;
  }

  let low = 0;
  let high = graphemes.length;
  let best = ellipsis;

  while (low <= high) {
    const middle = Math.floor((low + high) / 2);
    const candidate = `${graphemes.slice(0, middle).join("")}${ellipsis}`;

    if (measureText(candidate, fontSize, template) <= maximumWidth) {
      best = candidate;
      low = middle + 1;
    } else {
      high = middle - 1;
    }
  }

  return best;
}

/**
 * 输入极长、即使最小字号也放不下时的兜底。
 * 保证始终返回最多 maxLines 行，而不是让 SVG 完全溢出。
 */
function createFallbackCandidate(
  rawText: string,
  fontSize: number,
  template: CardTemplate,
): LayoutCandidate {
  const text = normalizeText(rawText) || DEFAULT_EMPTY_TEXT;

  const box = template.textBox;
  const lineHeightPx = fontSize * box.lineHeight;

  const prepared = getPrepared(text, fontSize, template);

  const result = layoutWithLines(prepared, box.width, lineHeightPx);

  const visibleLines = result.lines.slice(0, box.maxLines).map((line) => ({
    text: line.text,
    width: line.width,
  }));

  if (!visibleLines.length) {
    return {
      lines: [
        {
          text: DEFAULT_EMPTY_TEXT,
          width: measureText(DEFAULT_EMPTY_TEXT, fontSize, template),
        },
      ],
      wrapWidth: box.width,
      score: Number.POSITIVE_INFINITY,
    };
  }

  if (result.lines.length > box.maxLines) {
    const lastIndex = visibleLines.length - 1;
    const lastText = visibleLines[lastIndex].text;

    const truncated = fitWithEllipsis(lastText, box.width, fontSize, template);

    visibleLines[lastIndex] = {
      text: truncated,
      width: measureText(truncated, fontSize, template),
    };
  }

  return {
    lines: visibleLines,
    wrapWidth: box.width,
    score: Number.POSITIVE_INFINITY,
  };
}

/**
 * 构建标题卡片的最终文字布局。
 *
 * 流程：
 * 1. 二分寻找“能放下”的最大字号
 * 2. 在最大字号附近再比较几个稍小字号
 * 3. 用视觉评分选出更均衡的结果
 */
export function buildTextLayout(rawText: string, template: CardTemplate): TextLayout {
  const text = normalizeText(rawText) || DEFAULT_EMPTY_TEXT;

  const box = template.textBox;

  let low = box.minFontSize;
  let high = box.maxFontSize;
  let maximumFittingSize = box.minFontSize;

  let maximumCandidate = findBestCandidateAtSize(text, box.minFontSize, template);

  while (high - low > 0.8) {
    const fontSize = (low + high) / 2;

    const candidate = findBestCandidateAtSize(text, fontSize, template);

    if (candidate) {
      maximumFittingSize = fontSize;
      maximumCandidate = candidate;
      low = fontSize;
    } else {
      high = fontSize;
    }
  }

  /**
   * 最大字号不一定是最好看的字号。
   * 继续比较最大字号附近的几个值，
   * 同时对缩小字号施加成本，避免文字无意义地变小。
   */
  const nearbySizes = new Set<number>([
    maximumFittingSize,
    maximumFittingSize * 0.97,
    maximumFittingSize * 0.94,
    maximumFittingSize * 0.91,
    maximumFittingSize * 0.88,
  ]);

  let selectedSize = maximumFittingSize;
  let selectedCandidate = maximumCandidate;
  let selectedScore = maximumCandidate?.score ?? Number.POSITIVE_INFINITY;

  for (const rawSize of nearbySizes) {
    const fontSize = clamp(rawSize, box.minFontSize, box.maxFontSize);

    const candidate = findBestCandidateAtSize(text, fontSize, template);

    if (!candidate) {
      continue;
    }

    const shrinkRatio =
      maximumFittingSize > 0 ? (maximumFittingSize - fontSize) / maximumFittingSize : 0;

    const adjustedScore = candidate.score + shrinkRatio * 620;

    if (!selectedCandidate || adjustedScore < selectedScore) {
      selectedSize = fontSize;
      selectedCandidate = candidate;
      selectedScore = adjustedScore;
    }
  }

  if (!selectedCandidate) {
    selectedSize = box.minFontSize;
    selectedCandidate = createFallbackCandidate(text, selectedSize, template);
  }

  const fontSize = Math.floor(selectedSize * 10) / 10;

  const lineHeightPx = fontSize * box.lineHeight;

  return {
    lines: selectedCandidate.lines.map((line) => line.text),
    fontSize,
    lineHeightPx,
    totalHeight: selectedCandidate.lines.length * lineHeightPx,
  };
}

/**
 * 根据最终行布局计算关键词高亮区域。
 *
 * Pretext 负责宽度测量；这里保留原有的 SVG 高亮矩形输出格式。
 */
export function findHighlightRects(
  layout: TextLayout,
  keyword: string,
  template: CardTemplate,
): HighlightRect[] {
  const normalized = keyword.trim();

  if (!normalized) {
    return [];
  }

  const rectangles: HighlightRect[] = [];
  const box = template.textBox;

  layout.lines.forEach((line, lineIndex) => {
    let from = 0;

    while (from < line.length) {
      const index = line.indexOf(normalized, from);

      if (index < 0) {
        break;
      }

      const before = line.slice(0, index);
      const matched = line.slice(index, index + normalized.length);

      const beforeWidth = measureText(before, layout.fontSize, template);

      const matchedWidth = measureText(matched, layout.fontSize, template);

      const lineWidth = measureText(line, layout.fontSize, template);

      let lineStart = box.x;

      if (box.align === "center") {
        lineStart = box.x + (box.width - lineWidth) / 2;
      } else if (box.align === "right") {
        lineStart = box.x + box.width - lineWidth;
      }

      rectangles.push({
        x: lineStart + beforeWidth - 8,
        y: box.y + lineIndex * layout.lineHeightPx + layout.fontSize * 0.55,
        width: matchedWidth + 16,
        height: layout.fontSize * 0.5,
        rotation: (((lineIndex + index) % 3) - 1) * 1.8,
      });

      from = index + normalized.length;
    }
  });

  return rectangles;
}

/**
 * 简单的前端关键词推断。
 * 生产环境可以替换成分词、TF-IDF 或后端 NLP。
 */
export function inferKeyword(text: string): string {
  const trimmed = normalizeText(text);

  const quoted = trimmed.match(/[“「『《](.{2,8}?)[”」』》]/);

  if (quoted?.[1]) {
    return quoted[1];
  }

  const fragments = trimmed
    .split(/[，。！？、；：,\s!?]+/)
    .map((item) => item.trim())
    .filter((item) => graphemeLength(item) >= 2);

  if (!fragments.length) {
    return "";
  }

  const last = fragments[fragments.length - 1];

  if (graphemeLength(last) <= 6) {
    return last;
  }

  return getGraphemes(last).slice(-4).join("");
}

export function hashText(value: string): number {
  let hash = 2166136261;

  for (let index = 0; index < value.length; index += 1) {
    hash ^= value.charCodeAt(index);
    hash = Math.imul(hash, 16777619);
  }

  return hash >>> 0;
}

export function seededValue(seed: number, index: number): number {
  const value = Math.sin(seed * 12.9898 + index * 78.233) * 43758.5453;

  return value - Math.floor(value);
}
