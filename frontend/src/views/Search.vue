<template>
  <section class="search-page">
    <header class="search-header">
      <p class="search-label">{{ t('search.title') }}</p>
      <h1 class="search-title">{{ t('search.heading') }}</h1>
      <p class="search-description">
        <span v-if="query">{{ t('search.currentKeyword') }}{{ query }}</span>
        <span v-else>{{ t('search.placeholder') }}</span>
      </p>
    </header>

    <div class="search-state">
      <SearchEmptyState
        v-if="!hasQuery"
        :description="t('search.previewHint')"
      />
      <template v-else>
        <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" />
        <div v-else-if="loading" class="search-loading">
          <el-skeleton :rows="6" animated />
        </div>
        <div v-else class="search-result-shell">
          <SearchFilters />
          <SearchEmptyState v-if="!items.length" :description="t('search.noResults')" />
          <SearchResultList v-else :items="items" :query="query" :total="total" />
        </div>
      </template>
    </div>
  </section>
</template>

<script>
import { ElAlert } from "element-plus";
import SearchEmptyState from "../components/search/SearchEmptyState.vue";
import SearchFilters from "../components/search/SearchFilters.vue";
import SearchResultList from "../components/search/SearchResultList.vue";
import { useSearch } from "../composables/useSearch";
import { t } from "../i18n";

export default {
  name: "SearchPage",
  components: {
    ElAlert,
    SearchEmptyState,
    SearchFilters,
    SearchResultList,
  },
  setup() {
    const { query, hasQuery, loading, error, items, total } = useSearch();
    return {
      t,
      query,
      hasQuery,
      loading,
      error,
      items,
      total,
    };
  },
};
</script>

<style scoped>
.search-page {
  max-width: var(--content-w);
  margin: 0 auto;
  padding: 24px 16px 48px;
}

.search-header {
  margin-bottom: 24px;
}

.search-label {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-faint);
}

.search-title {
  margin: 0;
  font-size: 32px;
  line-height: 1.2;
  color: var(--text);
}

.search-description {
  margin: 12px 0 0;
  color: var(--text-soft);
  line-height: 1.7;
}

.search-state {
  min-height: 240px;
}

.search-loading {
  padding-top: 8px;
}

.search-result-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
