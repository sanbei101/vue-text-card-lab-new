<script setup lang="ts">
import { computed } from "vue";

import { kindToComponent } from "@/card/styles";
import type { CardTemplate } from "@/types";

const props = defineProps<{ template: CardTemplate; title: string; keyword: string }>();

const StyleComponent = computed(() => kindToComponent[props.template.kind]);
</script>

<template>
  <component :is="StyleComponent" :template="template" :title="title" :keyword="keyword" />
</template>
