<template>
  <p v-if="content" class="summary-preview" :data-summary-type="type || 'fallback'">
    <span v-if="type === 'ai'" class="summary-badge">AI</span>
    {{ content }}
  </p>
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
.summary-badge {
  display: inline-block;
  margin-right: 6px;
  padding: 1px 6px 2px;
  vertical-align: 1px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--accent) 12%, transparent);
  border: 1px solid color-mix(in srgb, var(--accent) 28%, transparent);
  color: white;
  color: var(--accent);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  line-height: 1.1;
}

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
  .summary-badge {
    font-size: 10px;
    margin-right: 5px;
  }

  .summary-preview {
    font-size: 15px;
    line-height: 1.7;
    -webkit-line-clamp: 3;
  }
}
</style>
