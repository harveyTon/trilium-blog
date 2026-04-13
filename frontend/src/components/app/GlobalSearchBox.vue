<template>
  <div class="global-search" @focusin="handleFocusIn" @focusout="handleFocusOut">
    <input
      ref="inputRef"
      v-model="query"
      class="global-search-input"
      type="search"
      placeholder="搜索文章标题与内容"
      aria-label="搜索文章标题与内容"
      @input="handleInput"
      @keydown="handleKeydown"
    />

    <SearchPreviewPanel
      v-if="showPreview"
      :items="searchStore.previewItems"
      :loading="searchStore.loading"
      :error="searchStore.error"
      :active-index="activeIndex"
      @select-item="selectItem"
      @highlight-index="setActiveIndex"
      @submit-search="submitSearch"
    />
  </div>
</template>

<script>
import { computed, ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import SearchPreviewPanel from "./SearchPreviewPanel.vue";
import { useSearchPreview } from "../../composables/useSearchPreview";

export default {
  name: "GlobalSearchBox",
  components: {
    SearchPreviewPanel,
  },
  setup() {
    const router = useRouter();
    const route = useRoute();
    const inputRef = ref(null);
    const activeIndex = ref(-1);
    const blurTimer = ref(null);
    const { searchStore, loadPreview, closePreview } = useSearchPreview();
    const query = ref((route.query.q || "").toString());

    const showPreview = computed(() => {
      if (!searchStore.previewOpen) return false;
      if (!query.value.trim()) return false;
      return true;
    });

    const syncFromRoute = () => {
      query.value = (route.query.q || "").toString();
    };

    watch(() => route.fullPath, syncFromRoute);

    const setActiveIndex = (index) => {
      activeIndex.value = index;
    };

    const submitSearch = () => {
      const trimmed = query.value.trim();
      if (!trimmed) {
        closePreview();
        return;
      }
      router.push({ name: "Search", query: { q: trimmed } });
      closePreview();
    };

    const selectItem = (item) => {
      router.push({ name: "Article", params: { noteId: item.noteId } });
      closePreview();
    };

    const handleInput = () => {
      activeIndex.value = -1;
      loadPreview(query.value);
    };

    const handleFocusIn = () => {
      if (blurTimer.value) {
        window.clearTimeout(blurTimer.value);
        blurTimer.value = null;
      }
      if (query.value.trim()) {
        loadPreview(query.value);
      }
    };

    const handleFocusOut = () => {
      blurTimer.value = window.setTimeout(() => {
        closePreview();
      }, 120);
    };

    const handleKeydown = (event) => {
      const items = searchStore.previewItems;
      if (event.key === "ArrowDown" && items.length) {
        event.preventDefault();
        activeIndex.value = (activeIndex.value + 1 + items.length) % items.length;
        return;
      }
      if (event.key === "ArrowUp" && items.length) {
        event.preventDefault();
        activeIndex.value = (activeIndex.value - 1 + items.length) % items.length;
        return;
      }
      if (event.key === "Enter") {
        event.preventDefault();
        if (activeIndex.value >= 0 && items[activeIndex.value]) {
          selectItem(items[activeIndex.value]);
          return;
        }
        submitSearch();
        return;
      }
      if (event.key === "Escape") {
        closePreview();
      }
    };

    return {
      inputRef,
      query,
      activeIndex,
      searchStore,
      showPreview,
      setActiveIndex,
      submitSearch,
      selectItem,
      handleInput,
      handleFocusIn,
      handleFocusOut,
      handleKeydown,
    };
  },
};
</script>

<style scoped>
.global-search {
  position: relative;
  width: min(420px, 100%);
}

.global-search-input {
  width: 100%;
  border: 1px solid rgba(255, 255, 255, 0.18);
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-inverse);
  border-radius: 999px;
  padding: 10px 16px;
  font-size: 14px;
  box-sizing: border-box;
}

.global-search-input::placeholder {
  color: rgba(255, 255, 255, 0.72);
}

html.dark .global-search-input {
  background: rgba(255, 255, 255, 0.06);
}

@media (max-width: 768px) {
  .global-search {
    width: 100%;
  }
}
</style>
