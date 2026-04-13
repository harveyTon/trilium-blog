<template>
  <div class="home">
    <section class="home-intro">
      <p class="home-kicker">内容入口</p>
      <h1 class="home-title">搜索、精选与最新文章</h1>
      <p class="home-description">
        先搜索，再浏览精选，最后顺着时间线继续阅读。
      </p>
      <div class="home-search-entry">
        <GlobalSearchBox />
      </div>
    </section>

    <div v-if="loading" class="post-list" aria-busy="true">
      <div v-for="i in 5" :key="i" class="post-item post-item--skeleton">
        <div class="skeleton-date"></div>
        <div class="skeleton-right">
          <div class="skeleton-title"></div>
          <div class="skeleton-line"></div>
          <div class="skeleton-meta"></div>
          <div class="skeleton-summary"></div>
        </div>
      </div>
    </div>

    <div v-else class="home-sections">
      <div v-if="fetchError" class="fetch-error">
        <p>{{ fetchError }}</p>
        <el-button type="primary" @click="loadPosts">重试</el-button>
      </div>

      <FeaturedSection :items="featuredPosts" />

      <section class="home-feed-section">
        <div class="home-section-header">
          <p class="home-section-kicker">最新文章</p>
          <h2 class="home-section-title">按更新时间浏览</h2>
        </div>

        <el-empty v-if="!posts.length && !fetchError" description="暂无已发布文章" />
        <PostFeed
          v-else
          :items="posts"
          :get-day="getDay"
          :get-month="getMonth"
          :format-full-date="formatFullDate"
          :sanitize-summary="sanitizeSummary"
        />
      </section>

      <div v-if="pagination.totalPages > 1" class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          @current-change="handleCurrentChange"
          :prev-text="'←'" :next-text="'→'"
          layout="prev, pager, next"
          :ellipsis-item="'...'"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ElButton, ElEmpty, ElPagination } from "element-plus";
import { onMounted, ref, watch } from "vue";
import GlobalSearchBox from "../components/app/GlobalSearchBox.vue";
import FeaturedSection from "../components/home/FeaturedSection.vue";
import PostFeed from "../components/home/PostFeed.vue";
import { fetchFeaturedPosts, fetchPosts } from "../api/blog";
import { useSiteStore } from "../store";

export default {
  name: "HomePage",
  components: {
    ElButton,
    ElEmpty,
    ElPagination,
    GlobalSearchBox,
    FeaturedSection,
    PostFeed,
  },
  setup() {
    const siteStore = useSiteStore();
    const posts = ref([]);
    const featuredPosts = ref([]);
    const loading = ref(true);
    const fetchError = ref(null);
    const currentPage = ref(1);
    const pagination = ref({
      page: 1,
      pageSize: 9,
      total: 0,
      totalPages: 0,
    });

    const loadPosts = async () => {
      loading.value = true;
      fetchError.value = null;
      try {
        const [postResponse, featuredResponse] = await Promise.all([
          fetchPosts(currentPage.value),
          fetchFeaturedPosts(),
        ]);
        posts.value = postResponse.items;
        featuredPosts.value = featuredResponse.items || [];
        pagination.value = {
          page: postResponse.page,
          pageSize: postResponse.pageSize,
          total: postResponse.total,
          totalPages: postResponse.totalPages,
        };
      } catch {
        fetchError.value = "加载失败，请检查网络后重试";
        posts.value = [];
        featuredPosts.value = [];
      } finally {
        loading.value = false;
      }
    };

    const handleCurrentChange = async (page) => {
      currentPage.value = page;
      await loadPosts();
      window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const pad = (n) => String(n).padStart(2, "0");

    const getDay = (dateString) => {
      return pad(new Date(dateString).getDate());
    };

    const getMonth = (dateString) => {
      const d = new Date(dateString);
      return `${d.getFullYear()}年-${pad(d.getMonth() + 1)}月`;
    };

    const formatFullDate = (dateString) => {
      const d = new Date(dateString);
      const h = d.getHours();
      const m = pad(d.getMinutes());
      const period = h < 12 ? "上午" : "下午";
      const displayHour = h > 12 ? h - 12 : h === 0 ? 12 : h;
      return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${period}${displayHour}:${m}`;
    };

    const sanitizeSummary = (text) => {
      if (!text) return "";
      return text
        .replace(/\u200b/g, "")
        .replace(/\u200c/g, "")
        .replace(/\u200d/g, "")
        .replace(/\ufeff/g, "")
        .replace(/\u00ad/g, "")
        .replace(/\ufffd/g, "")
        .replace(/[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]/g, "");
    };

    const syncTitle = () => {
      if (siteStore.site.title) {
        document.title = siteStore.site.title;
      }
    };

    onMounted(async () => {
      await loadPosts();
      syncTitle();
    });

    watch(() => siteStore.site.title, syncTitle, { immediate: true });

    return {
      posts,
      featuredPosts,
      loading,
      fetchError,
      currentPage,
      pagination,
      handleCurrentChange,
      getDay,
      getMonth,
      formatFullDate,
      sanitizeSummary,
      loadPosts,
    };
  },
};
</script>

<style scoped>
.home {
  max-width: var(--list-w);
  margin: 0 auto;
  padding: 0 16px 32px;
}

.home-intro {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px 0 28px;
}

.home-kicker,
.home-section-kicker {
  margin: 0;
  font-size: 12px;
  color: var(--text-faint);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.home-title,
.home-section-title {
  margin: 0;
  color: var(--text);
}

.home-title {
  font-size: clamp(30px, 4vw, 42px);
  line-height: 1.16;
}

.home-description {
  margin: 0;
  color: var(--text-soft);
  line-height: 1.7;
}

.home-search-entry {
  margin-top: 8px;
}

.home-search-entry :deep(.global-search) {
  width: min(560px, 100%);
}

.home-search-entry :deep(.global-search-input) {
  border-color: var(--border-soft);
  background: var(--surface);
  color: var(--text);
  box-shadow: var(--shadow-sm);
}

.home-search-entry :deep(.global-search-input::placeholder) {
  color: var(--text-faint);
}

.home-sections {
  display: flex;
  flex-direction: column;
  gap: 36px;
}

.home-feed-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.home-section-header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.post-list {
  display: flex;
  flex-direction: column;
}

.post-item--skeleton {
  display: flex;
  align-items: flex-start;
  gap: 28px;
  padding: 18px 0 22px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  pointer-events: none;
}

.skeleton-date {
  width: 92px;
  height: 92px;
  min-width: 92px;
  border: 1px solid var(--border-soft);
  background: var(--surface-muted);
  animation: skeleton-pulse 1.6s ease-in-out infinite;
}

.skeleton-right {
  flex: 1;
  min-width: 0;
}

.skeleton-title {
  width: 70%;
  height: 24px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 14px;
  animation: skeleton-pulse 1.6s ease-in-out infinite;
}

.skeleton-line {
  width: 120px;
  height: 1px;
  background: var(--border-soft);
  margin-bottom: 12px;
  animation: skeleton-pulse 1.6s ease-in-out infinite;
  animation-delay: 0.05s;
}

.skeleton-meta {
  width: 40%;
  height: 12px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 12px;
  animation: skeleton-pulse 1.6s ease-in-out infinite;
  animation-delay: 0.1s;
}

.skeleton-summary {
  width: 90%;
  height: 16px;
  background: var(--border-soft);
  border-radius: 4px;
  animation: skeleton-pulse 1.6s ease-in-out infinite;
  animation-delay: 0.15s;
}

@keyframes skeleton-pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

.pagination-wrapper {
  margin-top: 4px;
  padding-bottom: 24px;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  font-weight: 500;
  gap: 4px;
}

:deep(.el-pagination button) {
  min-width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  background: var(--surface);
  border: 1px solid var(--border-soft);
  color: var(--text-soft);
  transition: border-color 160ms, color 160ms, background 160ms;
}

:deep(.el-pagination button:hover) {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--surface);
}

:deep(.el-pagination .el-pager li) {
  min-width: 36px;
  height: 36px;
  line-height: 36px;
  border-radius: var(--radius-sm);
  background: var(--surface);
  border: 1px solid var(--border-soft);
  color: var(--text-soft);
  font-weight: 500;
  transition: border-color 160ms, color 160ms, background 160ms;
}

:deep(.el-pager li:hover) {
  border-color: var(--accent);
  color: var(--accent);
}

:deep(.el-pager li.is-active) {
  background: var(--brand);
  border-color: var(--brand);
  color: var(--text-inverse);
}

:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next) {
  font-family: inherit;
}

html.dark :deep(.el-pagination button),
html.dark :deep(.el-pagination .el-pager li) {
  background: var(--surface);
  border-color: var(--border);
  color: var(--text-soft);
}

html.dark :deep(.el-pager li.is-active) {
  background: var(--brand);
  border-color: var(--brand);
  color: var(--text-inverse);
}

.fetch-error {
  text-align: center;
  color: var(--text-faint);
  padding: 40px 0;
  font-size: 0.95rem;
}

.fetch-error .el-button {
  margin-top: 16px;
}

@media (max-width: 768px) {
  .home {
    padding: 0 10px 32px;
  }

  .home-title {
    font-size: 28px;
  }

  .post-item--skeleton {
    display: block;
  }

  .skeleton-date {
    display: none;
  }
}
</style>
