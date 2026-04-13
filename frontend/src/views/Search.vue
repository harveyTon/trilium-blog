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
      <el-empty
        v-if="!query"
        description="输入关键词后，可在这里查看完整搜索结果"
      />
      <el-skeleton v-else :rows="6" animated />
    </div>
  </section>
</template>

<script>
import { computed } from "vue";
import { useRoute } from "vue-router";

export default {
  name: "SearchPage",
  setup() {
    const route = useRoute();
    const query = computed(() => (route.query.q || "").toString().trim());

    return {
      query,
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
</style>
