import { computed } from "vue";

import { buildTextLayout, findHighlightRects, hashText, seededValue } from "../../card-engine";
import type { CardTemplate } from "../../types";

export interface CardRenderProps {
  template: CardTemplate;
  title: string;
  keyword: string;
}

export function useCardRender(props: CardRenderProps) {
  const layout = computed(() => buildTextLayout(props.title, props.template));

  const highlights = computed(() =>
    findHighlightRects(layout.value, props.keyword, props.template),
  );

  const textAnchor = computed<"start" | "middle" | "end">(() => {
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

  return {
    layout,
    highlights,
    textAnchor,
    textX,
    seed,
    particles,
  };
}
