<template>
  <div class="search-preview-panel">
    <div v-if="loading" class="search-preview-state">{{ t('search.searching') }}</div>
    <div v-else-if="error" class="search-preview-state">{{ error }}</div>
    <div v-else-if="!items.length" class="search-preview-state">{{ t('search.noResults') }}</div>
    <ul v-else class="search-preview-list">
      <li v-for="(item, index) in items" :key="item.noteId">
        <button
          type="button"
          :class="['search-preview-item', { 'is-active': index === activeIndex }]"
          @mousedown.prevent="$emit('select-item', item)"
          @mouseenter="$emit('highlight-index', index)"
        >
          <span class="search-preview-title">{{ item.title }}</span>
          <span v-if="item.match?.snippet" class="search-preview-snippet">
            {{ item.match.snippet }}
          </span>
        </button>
      </li>
    </ul>
    <button
      v-if="items.length"
      type="button"
      class="search-preview-more"
      @mousedown.prevent="$emit('submit-search')"
    >
      {{ t('search.viewAll') }}
    </button>
  </div>
</template>

<script>
import { t } from "../../i18n";

export default {
  name: "SearchPreviewPanel",
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    loading: {
      type: Boolean,
      default: false,
    },
    error: {
      type: String,
      default: "",
    },
    activeIndex: {
      type: Number,
      default: -1,
    },
  },
  emits: ["select-item", "highlight-index", "submit-search"],
};
</script>

<style scoped>
.search-preview-panel {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  padding: 8px;
  z-index: 1200;
}

.search-preview-state {
  padding: 12px;
  color: var(--text-soft);
  font-size: 14px;
}

.search-preview-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.search-preview-item {
  width: 100%;
  border: none;
  background: transparent;
  text-align: left;
  padding: 10px 12px;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.search-preview-item.is-active,
.search-preview-item:hover {
  background: var(--surface-muted);
}

.search-preview-title {
  color: var(--text);
  font-weight: 600;
}

.search-preview-snippet {
  color: var(--text-soft);
  font-size: 13px;
  line-height: 1.5;
}

.search-preview-more {
  width: 100%;
  margin-top: 8px;
  border: none;
  background: var(--surface-muted);
  color: var(--accent);
  border-radius: 10px;
  padding: 10px 12px;
  cursor: pointer;
  font-weight: 600;
}
</style>
