<template>
  <header class="article-header">
    <h1 class="article-title">{{ title }}</h1>
    <div class="article-meta">
      <time>{{ formattedDate }}</time>
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
