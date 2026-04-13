<template>
  <div
    v-if="items?.length >= 3"
    :class="[
      'article-toc-wrapper',
      { 'is-collapsed': collapsed, 'is-reading-mode': readingMode },
    ]"
  >
    <aside class="article-toc">
      <button
        v-if="readingMode && collapsed"
        type="button"
        class="toc-collapsed-trigger"
        @click="$emit('toggle-collapse')"
      >
        目录
      </button>
      <div v-else class="toc-panel">
        <div class="toc-header">
          <div class="toc-title">目录</div>
          <button type="button" class="toc-toggle" @click="$emit('toggle-collapse')">
            {{ collapsed ? "展开" : "收起" }}
          </button>
        </div>
        <div v-if="!collapsed" class="toc-links">
          <a
            v-for="item in items"
            :key="item.id"
            :href="'#' + item.id"
            :class="['toc-link', 'toc-level-' + item.level, { 'is-active': activeHeading === item.id }]"
            @click.prevent="$emit('scroll-to-heading', item.id)"
          >
            {{ item.title }}
          </a>
        </div>
      </div>
    </aside>
  </div>
</template>

<script>
export default {
  name: "ArticleTOC",
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    activeHeading: {
      type: String,
      default: "",
    },
    collapsed: {
      type: Boolean,
      default: false,
    },
    readingMode: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["scroll-to-heading", "toggle-collapse"],
};
</script>

<style scoped>
.article-toc-wrapper.is-reading-mode {
  width: auto;
}

.article-toc-wrapper.is-reading-mode .article-toc {
  overflow: visible;
}

.toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.toc-collapsed-trigger,
.toc-toggle {
  border: none;
  background: transparent;
  color: var(--accent);
  cursor: pointer;
  font-size: 12px;
  padding: 0 6px;
}

.toc-collapsed-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 52px;
  min-height: 34px;
  border: 1px solid color-mix(in srgb, var(--border-soft) 88%, transparent 12%);
  border-radius: 999px;
  padding: 0 12px;
  background: color-mix(in srgb, var(--surface) 90%, transparent 10%);
  box-shadow: var(--shadow-sm);
  color: var(--text-soft);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.toc-links {
  display: flex;
  flex-direction: column;
}

@media (max-width: 1024px) {
  .article-toc-wrapper {
    display: block;
    width: 100%;
  }
}
</style>
