<template>
  <section v-if="shouldRender" class="article-summary-block" :data-summary-status="summary?.status || 'idle'">
    <p class="article-summary-label">AI Summary</p>
    <p v-if="summary?.status === 'ready' && summary?.text" class="article-summary-text">{{ summary.text }}</p>
    <div v-else-if="showPending" class="article-summary-pending">
      <span class="article-summary-spinner" aria-hidden="true"></span>
      <p class="article-summary-hint">正在生成 AI 摘要...</p>
    </div>
    <p v-else-if="summary?.status === 'failed'" class="article-summary-hint">
      AI 摘要暂时生成失败，稍后会自动重试。
    </p>
  </section>
</template>

<script>
import { computed } from "vue";

export default {
  name: "ArticleSummaryBlock",
  props: {
    summary: {
      type: Object,
      default: null,
    },
    enabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const isLoading = computed(() => ["pending", "processing"].includes(props.summary?.status));
    const showPending = computed(() => props.enabled && (!props.summary || isLoading.value));
    const shouldRender = computed(
      () =>
        Boolean(
          props.enabled &&
            (showPending.value || props.summary?.status === "ready" || props.summary?.status === "failed")
        )
    );

    return {
      isLoading,
      showPending,
      shouldRender,
    };
  },
};
</script>

<style scoped>
.article-summary-block {
  margin-bottom: 24px;
  padding: 18px 20px;
  border-radius: 16px;
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-lighter);
}

.article-summary-label {
  margin: 0 0 10px;
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--el-text-color-secondary);
}

.article-summary-text,
.article-summary-hint {
  margin: 0;
  line-height: 1.75;
  color: var(--el-text-color-primary);
}

.article-summary-pending {
  display: flex;
  align-items: center;
  gap: 10px;
}

.article-summary-spinner {
  width: 16px;
  height: 16px;
  border-radius: 999px;
  border: 2px solid var(--el-border-color);
  border-top-color: var(--el-color-primary);
  animation: article-summary-spin 0.9s linear infinite;
}

@keyframes article-summary-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
