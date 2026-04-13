<template>
  <div class="search-filters">
    <span class="search-filter-label">搜索范围</span>
    <button type="button" class="search-filter-chip is-active">全部</button>
  </div>
</template>

<script>
export default {
  name: "SearchFilters",
};
</script>

<style scoped>
.search-filters {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.search-filter-label {
  color: var(--text-faint);
  font-size: 13px;
}

.search-filter-chip {
  border: 1px solid var(--border-soft);
  background: var(--surface);
  color: var(--text-soft);
  border-radius: 999px;
  padding: 6px 12px;
  cursor: default;
}

.search-filter-chip.is-active {
  border-color: var(--accent);
  color: var(--accent);
}
</style>
