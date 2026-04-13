<template>
  <header class="article-header">
    <div class="article-header-main">
      <h1 class="article-title">{{ title }}</h1>
      <div class="article-meta">
        <time>{{ formattedDate }}</time>
      </div>
    </div>
    <div v-if="$slots.actions" class="article-header-actions">
      <slot name="actions" />
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
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
}

.article-header-main {
  min-width: 0;
  flex: 1;
}

.article-header-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .article-header {
    flex-direction: column;
    gap: 16px;
  }

  .article-header-actions {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
