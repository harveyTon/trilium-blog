<template>
  <p v-if="content" class="summary-preview" :data-summary-type="type || 'fallback'">{{ content }}</p>
</template>

<script>
export default {
  name: "SummaryPreview",
  props: {
    content: {
      type: String,
      default: "",
    },
    type: {
      type: String,
      default: "",
    },
  },
};
</script>

<style scoped>
.summary-preview {
  margin: 10px 0 0;
  font-size: 16px;
  line-height: 1.75;
  color: var(--text-soft);
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

@media (max-width: 768px) {
  .summary-preview {
    margin-top: 10px;
    font-size: 15px;
    line-height: 1.7;
    -webkit-line-clamp: 3;
  }
}
</style>
