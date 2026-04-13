<template>
  <div
    ref="root"
    :class="['article-content', { 'is-reading-mode': readingMode }]"
    v-html="contentHtml"
  ></div>
</template>

<script>
export default {
  name: "ArticleContent",
  props: {
    contentHtml: {
      type: String,
      default: "",
    },
    readingMode: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    getRootElement() {
      return this.$refs.root ?? null;
    },
  },
};
</script>
