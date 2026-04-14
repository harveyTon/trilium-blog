<template>
  <section class="search-result-list-shell">
    <p class="search-result-count">{{ t('search.resultCount', { total }) }}</p>
    <div class="search-result-list">
      <SearchResultCard
        v-for="item in items"
        :key="item.noteId"
        :item="item"
        :query="query"
      />
    </div>
  </section>
</template>

<script>
import SearchResultCard from "./SearchResultCard.vue";
import { t } from "../../i18n";

export default {
  name: "SearchResultList",
  methods: { t },
  components: {
    SearchResultCard,
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    total: {
      type: Number,
      default: 0,
    },
    query: {
      type: String,
      default: "",
    },
  },
};
</script>

<style scoped>
.search-result-list-shell {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.search-result-count {
  margin: 0;
  color: var(--text-soft);
}

.search-result-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
