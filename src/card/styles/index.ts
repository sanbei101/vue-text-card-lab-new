import type { Component } from "vue";

import type { TemplateKind } from "../../types";
import BubbleCard from "./BubbleCard.vue";
import BurstCard from "./BurstCard.vue";
import EditorialCard from "./EditorialCard.vue";
import GridCard from "./GridCard.vue";
import LavenderCard from "./LavenderCard.vue";
import MemoCard from "./MemoCard.vue";
import MonoCard from "./MonoCard.vue";
import NightCard from "./NightCard.vue";
import NotebookCard from "./NotebookCard.vue";
import PaperCard from "./PaperCard.vue";
import QuestionCard from "./QuestionCard.vue";
import StampCard from "./StampCard.vue";

export const kindToComponent: Record<TemplateKind, Component> = {
  question: QuestionCard,
  memo: MemoCard,
  editorial: EditorialCard,
  lavender: LavenderCard,
  notebook: NotebookCard,
  stamp: StampCard,
  grid: GridCard,
  night: NightCard,
  burst: BurstCard,
  paper: PaperCard,
  bubble: BubbleCard,
  mono: MonoCard,
};
