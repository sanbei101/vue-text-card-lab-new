import type { CardTemplate, HighlightRect, TextLayout } from "./types";

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
  ")",
  "）",
  "]",
  "】",
  "》",
]);

const BAD_LINE_END = new Set(["(", "（", "[", "【", "《"]);

const canvas = document.createElement("canvas");
const context = canvas.getContext("2d");

function setFont(fontSize: number, fontWeight: number, fontFamily: string) {
  if (!context) return;

  context.font = `${fontWeight} ${fontSize}px ${fontFamily}`;
}

export function measureText(value: string, fontSize: number, template: CardTemplate) {
  setFont(fontSize, template.fontWeight, template.fontFamily);

  return context?.measureText(value).width ?? [...value].length * fontSize;
}

function candidateBreaks(text: string) {
  const points = new Set<number>([0, text.length]);

  for (let i = 1; i < text.length; i += 1) {
    const previous = text[i - 1];
    const current = text[i];

    if (BREAK_AFTER.has(previous) || current === "\n" || previous === "\n") {
      points.add(i);
    }

    // 中文标题允许逐字断行，但标点附近会获得更高分。
    points.add(i);
  }

  return [...points].sort((a, b) => a - b);
}

function cleanLine(value: string) {
  return value.replace(/\n/g, "").trim();
}

function linePenalty(line: string, lineWidth: number, maxWidth: number, isLast: boolean) {
  if (!line) return Number.POSITIVE_INFINITY;

  const first = line[0];
  const last = line[line.length - 1];

  let penalty = 0;
  const fill = lineWidth / maxWidth;

  // 非末行更强调视觉均衡，末行允许短一些。
  penalty += isLast ? Math.pow(Math.max(0, 0.72 - fill), 2) * 520 : Math.pow(1 - fill, 2) * 820;

  if (BAD_LINE_START.has(first)) penalty += 1600;

  if (BAD_LINE_END.has(last)) penalty += 900;

  if (BREAK_AFTER.has(last)) penalty -= 100;

  if ([...line].length <= 2) penalty += isLast ? 680 : 1100;

  return penalty;
}

function layoutAtSize(rawText: string, fontSize: number, template: CardTemplate): string[] | null {
  const text = rawText
    .replace(/\r\n/g, "\n")
    .replace(/[ \t]+/g, " ")
    .trim();

  if (!text) return ["输入一个标题"];

  const points = candidateBreaks(text);
  const maxLines = template.textBox.maxLines;
  const maxWidth = template.textBox.width;

  const costs = new Map<string, number>();
  const paths = new Map<string, number[]>();

  costs.set("0|0", 0);
  paths.set("0|0", [0]);

  for (let lineIndex = 0; lineIndex < maxLines; lineIndex += 1) {
    for (let startIndex = 0; startIndex < points.length; startIndex += 1) {
      const start = points[startIndex];
      const key = `${lineIndex}|${start}`;
      const baseCost = costs.get(key);

      if (baseCost == null) continue;

      for (let endIndex = startIndex + 1; endIndex < points.length; endIndex += 1) {
        const end = points[endIndex];
        const rawLine = text.slice(start, end);

        if (rawLine.includes("\n") && !rawLine.endsWith("\n")) {
          break;
        }

        const line = cleanLine(rawLine);
        const width = measureText(line, fontSize, template);

        if (width > maxWidth) break;

        const isLast = end === text.length;
        const nextLineIndex = lineIndex + 1;

        const score = baseCost + linePenalty(line, width, maxWidth, isLast);

        const nextKey = `${nextLineIndex}|${end}`;
        const oldScore = costs.get(nextKey);

        if (oldScore == null || score < oldScore) {
          costs.set(nextKey, score);
          paths.set(nextKey, [...(paths.get(key) ?? [start]), end]);
        }

        if (rawLine.endsWith("\n")) break;
      }
    }
  }

  let best:
    | {
        cost: number;
        path: number[];
      }
    | undefined;

  for (let lineCount = 1; lineCount <= maxLines; lineCount += 1) {
    const key = `${lineCount}|${text.length}`;
    const cost = costs.get(key);
    const path = paths.get(key);

    if (cost == null || !path) continue;

    const lineCountPenalty = Math.pow(lineCount - Math.min(4, maxLines), 2) * 18;

    const totalCost = cost + lineCountPenalty;

    if (!best || totalCost < best.cost) {
      best = {
        cost: totalCost,
        path,
      };
    }
  }

  if (!best) return null;

  const lines: string[] = [];

  for (let i = 1; i < best.path.length; i += 1) {
    const line = cleanLine(text.slice(best.path[i - 1], best.path[i]));

    if (line) lines.push(line);
  }

  return lines;
}

export function buildTextLayout(text: string, template: CardTemplate): TextLayout {
  let low = template.textBox.minFontSize;
  let high = template.textBox.maxFontSize;
  let bestSize = low;
  let bestLines = layoutAtSize(text, low, template) ?? [text];

  while (high - low > 0.8) {
    const size = (low + high) / 2;
    const lines = layoutAtSize(text, size, template);

    const lineHeightPx = size * template.textBox.lineHeight;

    const fits =
      lines &&
      lines.length <= template.textBox.maxLines &&
      lines.length * lineHeightPx <= template.textBox.height;

    if (fits && lines) {
      bestSize = size;
      bestLines = lines;
      low = size;
    } else {
      high = size;
    }
  }

  const lineHeightPx = bestSize * template.textBox.lineHeight;

  return {
    lines: bestLines,
    fontSize: Math.floor(bestSize * 10) / 10,
    lineHeightPx,
    totalHeight: bestLines.length * lineHeightPx,
  };
}

export function findHighlightRects(
  layout: TextLayout,
  keyword: string,
  template: CardTemplate,
): HighlightRect[] {
  const normalized = keyword.trim();

  if (!normalized) return [];

  const rectangles: HighlightRect[] = [];
  const box = template.textBox;

  layout.lines.forEach((line, lineIndex) => {
    let from = 0;

    while (from < line.length) {
      const index = line.indexOf(normalized, from);

      if (index < 0) break;

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

export function inferKeyword(text: string) {
  const trimmed = text.trim();

  const quoted = trimmed.match(/[“「『《](.{2,8}?)[”」』》]/);

  if (quoted?.[1]) return quoted[1];

  const fragments = trimmed
    .split(/[，。！？、；：,\s!?]+/)
    .map((item) => item.trim())
    .filter((item) => item.length >= 2);

  if (!fragments.length) return "";

  const last = fragments[fragments.length - 1];

  if (last.length <= 6) return last;

  return last.slice(-4);
}

export function hashText(value: string) {
  let hash = 2166136261;

  for (let i = 0; i < value.length; i += 1) {
    hash ^= value.charCodeAt(i);
    hash = Math.imul(hash, 16777619);
  }

  return hash >>> 0;
}

export function seededValue(seed: number, index: number) {
  const x = Math.sin(seed * 12.9898 + index * 78.233) * 43758.5453;

  return x - Math.floor(x);
}
