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
          <div class="search-skeleton-filters">
            <div class="search-skeleton-filter-label"></div>
            <div class="search-skeleton-filter-chip"></div>
          </div>
          <div v-for="i in 4" :key="i" class="search-skeleton-card">
            <div class="search-skeleton-card-title"></div>
            <div class="search-skeleton-card-date"></div>
            <div class="search-skeleton-card-snippet"></div>
            <div class="search-skeleton-card-snippet search-skeleton-card-snippet--short"></div>
          </div>
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
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.search-skeleton-filters {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-skeleton-filter-label {
  width: 40px;
  height: 13px;
  background: var(--border-soft);
  border-radius: 4px;
  animation: skeleton-pulse-search 1.6s ease-in-out infinite;
}

.search-skeleton-filter-chip {
  width: 56px;
  height: 30px;
  border: 1px solid var(--border-soft);
  border-radius: 999px;
  background: var(--surface);
  animation: skeleton-pulse-search 1.6s ease-in-out infinite;
  animation-delay: 0.05s;
}

.search-skeleton-card {
  border: 1px solid var(--border-soft);
  background: var(--surface);
  border-radius: var(--radius-md);
  padding: 18px 20px;
  animation: skeleton-pulse-search 1.6s ease-in-out infinite;
  animation-delay: 0.1s;
}

.search-skeleton-card-title {
  width: 60%;
  height: 20px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 10px;
}

.search-skeleton-card-date {
  width: 100px;
  height: 13px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 12px;
}

.search-skeleton-card-snippet {
  width: 85%;
  height: 14px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 8px;
}

.search-skeleton-card-snippet--short {
  width: 50%;
  margin-bottom: 0;
}

@keyframes skeleton-pulse-search {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

.search-result-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
