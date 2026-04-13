<template>
  <section v-if="shouldRender" class="article-summary-block" :data-summary-status="summary?.status || 'idle'">
    <button class="article-summary-toggle" type="button" @click="collapsed = !collapsed">
      <div class="article-summary-heading">
        <span class="article-summary-chip">AI</span>
        <div>
          <p class="article-summary-label">AI Summary</p>
          <p class="article-summary-subtitle">{{ subtitle }}</p>
        </div>
      </div>
      <span class="article-summary-chevron" :class="{ 'is-collapsed': collapsed }" aria-hidden="true">⌃</span>
    </button>

    <transition name="summary-expand">
      <div v-if="!collapsed" class="article-summary-body">
        <p v-if="summary?.status === 'ready' && summary?.text" class="article-summary-text">{{ summary.text }}</p>
        <div v-else-if="showPending" class="article-summary-pending">
          <span class="article-summary-spinner" aria-hidden="true"></span>
          <p class="article-summary-hint">正在生成更自然的 AI 摘要，生成完成后会自动更新。</p>
        </div>
        <p v-else-if="summary?.status === 'failed'" class="article-summary-hint">
          AI 摘要暂时生成失败，稍后会自动重试。
        </p>
      </div>
    </transition>
  </section>
</template>

<script>
import { computed, ref, watch } from "vue";

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
    const collapsed = ref(false);
    const isLoading = computed(() => ["pending", "processing"].includes(props.summary?.status));
    const showPending = computed(() => props.enabled && (!props.summary || isLoading.value));
    const subtitle = computed(() => {
      if (props.summary?.status === "ready" && props.summary?.text) {
        return "已生成，可快速浏览文章核心内容";
      }
      if (showPending.value) {
        return "正在异步生成中";
      }
      return "暂时不可用";
    });
    const shouldRender = computed(
      () =>
        Boolean(
          props.enabled &&
            (showPending.value || props.summary?.status === "ready" || props.summary?.status === "failed")
        )
    );

    watch(
      () => props.summary?.status,
      (status) => {
        if (status === "ready") {
          collapsed.value = false;
        }
      },
      { immediate: true }
    );

    return {
      collapsed,
      isLoading,
      showPending,
      subtitle,
      shouldRender,
    };
  },
};
</script>

<style scoped>
.article-summary-block {
  margin-bottom: 24px;
  border-radius: 20px;
  background:
    radial-gradient(circle at top right, color-mix(in srgb, var(--accent) 18%, transparent), transparent 40%),
    linear-gradient(135deg, color-mix(in srgb, var(--surface) 88%, white 12%), var(--el-fill-color-light));
  border: 1px solid color-mix(in srgb, var(--accent) 12%, var(--el-border-color-lighter));
  box-shadow: 0 22px 44px color-mix(in srgb, var(--accent) 8%, transparent);
  overflow: hidden;
}

.article-summary-toggle {
  width: 100%;
  border: 0;
  padding: 18px 20px 16px;
  background: transparent;
  color: inherit;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  text-align: left;
  cursor: pointer;
}

.article-summary-heading {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.article-summary-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 14px;
  background: linear-gradient(135deg, color-mix(in srgb, var(--accent) 82%, white 18%), var(--accent));
  color: white;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.12em;
  box-shadow: 0 14px 26px color-mix(in srgb, var(--accent) 25%, transparent);
  flex-shrink: 0;
}

.article-summary-label {
  margin: 0;
  font-size: 0.85rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--el-text-color-secondary);
}

.article-summary-subtitle {
  margin: 4px 0 0;
  color: var(--text-soft);
  font-size: 13px;
  line-height: 1.5;
}

.article-summary-chevron {
  font-size: 18px;
  color: var(--text-faint);
  transition: transform 0.22s ease;
  flex-shrink: 0;
}

.article-summary-chevron.is-collapsed {
  transform: rotate(180deg);
}

.article-summary-body {
  padding: 0 20px 20px;
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

@media (max-width: 768px) {
  .article-summary-toggle {
    padding: 16px 16px 14px;
  }

  .article-summary-body {
    padding: 0 16px 16px;
  }

  .article-summary-chip {
    width: 36px;
    height: 36px;
    border-radius: 12px;
    font-size: 11px;
  }
}

@keyframes article-summary-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes summary-expand-in {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes summary-expand-out {
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(-4px);
  }
}

.summary-expand-enter-active {
  animation: summary-expand-in 0.2s ease;
}

.summary-expand-leave-active {
  animation: summary-expand-out 0.18s ease;
}
</style>
