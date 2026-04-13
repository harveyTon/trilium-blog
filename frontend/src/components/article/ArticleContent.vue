<template>
  <div class="article-content" v-html="contentHtml"></div>
</template>

<script>
export default {
  name: "ArticleContent",
  props: {
    contentHtml: {
      type: String,
      default: "",
    },
  },
};
</script>
