<template>
  <section class="search-page">
    <header class="search-header">
      <p class="search-label">站内搜索</p>
      <h1 class="search-title">搜索结果</h1>
      <p class="search-description">
        <span v-if="query">当前关键词：{{ query }}</span>
        <span v-else>请输入关键词开始搜索文章内容。</span>
      </p>
    </header>

    <div class="search-state">
      <SearchEmptyState
        v-if="!hasQuery"
        description="输入关键词后，可在这里查看完整搜索结果"
      />
      <template v-else>
        <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" />
        <div v-else-if="loading" class="search-loading">
          <el-skeleton :rows="6" animated />
        </div>
        <div v-else class="search-result-shell">
          <SearchFilters />
          <SearchEmptyState v-if="!items.length" description="没有找到相关内容" />
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
