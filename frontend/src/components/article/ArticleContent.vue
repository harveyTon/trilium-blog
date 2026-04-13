<template>
  <div ref="root" class="article-content" v-html="contentHtml"></div>
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
  methods: {
    getRootElement() {
      return this.$refs.root ?? null;
    },
  },
};
</script>
