<template>
  <section v-if="summary?.status === 'ready' && summary?.text" class="article-summary-block">
    <p class="article-summary-label">AI Summary</p>
    <p class="article-summary-text">{{ summary.text }}</p>
  </section>
</template>

<script>
export default {
  name: "ArticleSummaryBlock",
  props: {
    summary: {
      type: Object,
      default: null,
    },
  },
};
</script>
