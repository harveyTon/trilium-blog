<template>
  <div v-if="pageUrl" class="article-source">
    剪贴自：
    <a :href="pageUrl" target="_blank" rel="noopener noreferrer">
      {{ pageUrl }}
    </a>
  </div>
</template>

<script>
export default {
  name: "SourceLinkBlock",
  props: {
    pageUrl: {
      type: String,
      default: "",
    },
  },
};
</script>
