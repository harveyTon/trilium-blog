<template>
  <div v-if="items?.length >= 3" :class="['article-toc-wrapper', { 'is-collapsed': collapsed }]">
    <aside class="article-toc">
      <div class="toc-panel">
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
  },
  emits: ["scroll-to-heading", "toggle-collapse"],
};
</script>

<style scoped>
.toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.toc-toggle {
  border: none;
  background: transparent;
  color: var(--accent);
  cursor: pointer;
  font-size: 12px;
  padding: 0 6px;
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
