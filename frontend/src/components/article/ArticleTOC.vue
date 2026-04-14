<template>
  <div
    v-if="items?.length >= 3 && (!readingMode || !collapsed)"
    :class="[
      'article-toc-wrapper',
      { 'is-collapsed': collapsed, 'is-reading-mode': readingMode },
    ]"
  >
    <aside class="article-toc">
      <div :class="['toc-panel', { 'is-reading-drawer': readingMode }]">
        <div class="toc-header">
          <div class="toc-title">{{ t('article.toc') }}</div>
          <button type="button" class="toc-toggle" @click="$emit('toggle-collapse')">
            {{ readingMode ? t('article.tocClose') : collapsed ? t('article.tocExpand') : t('article.tocCollapse') }}
          </button>
        </div>
        <div v-if="readingMode || !collapsed" class="toc-links">
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
import { t } from "../../i18n";

export default {
  name: "ArticleTOC",
  methods: { t },
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
  position: fixed;
  top: 60px;
  right: 24px;
  z-index: 1190;
}

.article-toc-wrapper.is-reading-mode .article-toc {
  overflow: visible;
  max-height: none;
}

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

.toc-panel.is-reading-drawer {
  width: min(320px, calc(100vw - 48px));
  max-height: calc(100vh - 92px);
  padding: 16px 14px 14px;
  border-radius: 18px;
  background: color-mix(in srgb, var(--surface) 95%, transparent 5%);
  box-shadow: 0 24px 64px rgba(15, 23, 42, 0.16);
  backdrop-filter: blur(18px);
}

.toc-panel.is-reading-drawer .toc-links {
  max-height: calc(100vh - 160px);
  overflow-y: auto;
  padding-right: 2px;
}

@media (max-width: 1024px) {
  .article-toc-wrapper {
    display: none;
    width: 100%;
  }

  .article-toc-wrapper.is-reading-mode {
    display: block;
    top: 56px;
    right: 16px;
  }

  .toc-panel.is-reading-drawer {
    width: min(320px, calc(100vw - 32px));
  }
}
</style>
