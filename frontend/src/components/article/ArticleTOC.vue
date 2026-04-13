<template>
  <div v-if="items?.length >= 3" class="article-toc-wrapper">
    <aside class="article-toc">
      <div class="toc-panel">
        <div class="toc-title">目录</div>
        <a
          v-for="item in items"
          :key="item.id"
          :href="'#' + item.id"
          :class="['toc-link', 'toc-level-' + item.level, { 'is-active': activeHeading === item.id }]"
          @click.prevent="$emit('scroll-to-heading', item.id)"
        >
          {{ item.title }}
        </a>
      </div>
    </aside>
  </div>
</template>

<script>
export default {
  name: "ArticleTOC",
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    activeHeading: {
      type: String,
      default: "",
    },
  },
  emits: ["scroll-to-heading"],
};
</script>
