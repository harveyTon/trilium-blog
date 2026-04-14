<template>
  <header class="article-header">
    <div class="article-header-main">
      <h1 class="article-title">{{ title }}</h1>
      <div class="article-meta">
        <time>{{ formattedDate }}</time>
        <div v-if="$slots.actions" class="article-header-actions">
          <slot name="actions" />
        </div>
      </div>
    </div>
    <span class="artalk-pv-count" style="display: none"></span>
  </header>
</template>

<script>
export default {
  name: "ArticleHeader",
  props: {
    title: {
      type: String,
      required: true,
    },
    formattedDate: {
      type: String,
      required: true,
    },
  },
};
</script>

<style scoped>
.article-header {
  display: block;
}

.article-header-main {
  min-width: 0;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.article-header-actions {
  display: flex;
  align-items: center;
}

@media (max-width: 768px) {
  .article-header-actions {
    width: auto;
  }
}
</style>
