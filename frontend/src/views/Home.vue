<template>
  <div class="home">
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

    <div v-else>
      <el-empty v-if="!posts.length && !fetchError" description="暂无已发布文章" />
      <div v-if="fetchError" class="fetch-error">
        <p>{{ fetchError }}</p>
        <el-button type="primary" @click="loadPosts">重试</el-button>
      </div>
      <nav class="post-list" aria-label="文章列表">
        <article
          v-for="post in posts"
          :key="post.noteId"
          class="post-item"
        >
          <div class="post-date">
            <div class="post-day">{{ getDay(post.dateModified) }}</div>
            <div class="post-month">{{ getMonth(post.dateModified) }}</div>
          </div>
          <div class="post-body">
            <h2 class="post-title">
              <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }">{{ post.title }}</router-link>
            </h2>
            <div class="post-meta">{{ formatFullDate(post.dateModified) }}</div>
            <p v-if="post.summary" class="post-summary">{{ sanitizeSummary(post.summary) }}</p>
          </div>
        </article>
      </nav>
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
import { fetchPosts } from "../api/blog";
import { useSiteStore } from "../store";

export default {
  name: "HomePage",
  components: {
    ElButton,
    ElEmpty,
    ElPagination,
  },
  setup() {
    const siteStore = useSiteStore();
    const posts = ref([]);
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
        const response = await fetchPosts(currentPage.value);
        posts.value = response.items;
        pagination.value = {
          page: response.page,
          pageSize: response.pageSize,
          total: response.total,
          totalPages: response.totalPages,
        };
      } catch {
        fetchError.value = "加载失败，请检查网络后重试";
        posts.value = [];
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
  padding: 0 16px;
}

/* ── List ── */
.post-list {
  display: flex;
  flex-direction: column;
}

.post-list:focus-within {
  outline: none;
}

/* ── Item ── */
.post-item {
  display: flex;
  align-items: flex-start;
  gap: 28px;
  padding: 18px 0 22px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
}

.post-item:first-child {
  border-top: none;
}

/* ── Date block (left) ── */
.post-date {
  width: 92px;
  min-width: 92px;
  height: 92px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-soft);
  background: var(--surface-muted);
  flex-shrink: 0;
}

.post-day {
  font-size: 52px;
  line-height: 1;
  font-weight: 600;
  color: var(--text);
  font-variant-numeric: tabular-nums;
}

.post-month {
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-faint);
  line-height: 1.2;
}

/* ── Body (right) ── */
.post-body {
  flex: 1;
  min-width: 0;
}

/* ── Title ── */
.post-title {
  margin: 0;
  font-size: 20px;
  line-height: 1.4;
  font-weight: 700;
}

.post-title a {
  color: var(--text);
  text-decoration: none;
  position: relative;
  display: inline-block;
}

/* ── Short gradient divider line (key visual) ── */
.post-title a::after {
  content: "";
  display: block;
  width: 120px;
  height: 1px;
  margin-top: 14px;
  background: linear-gradient(
    to right,
    var(--border) 0%,
    var(--border) 35%,
    var(--border) 65%,
    transparent 100%
  );
  transition: background 0.2s ease;
}

.post-item:hover .post-title a::after {
  background: linear-gradient(
    to right,
    #FC9E42 0%,
    #FC9E42 35%,
    #FC9E42 65%,
    transparent 100%
  );
}

html.dark .post-item:hover .post-title a::after {
  background: linear-gradient(
    to right,
    #FC9E42 0%,
    #FC9E42 35%,
    #FC9E42 65%,
    transparent 100%
  );
}

/* ── Meta (hidden on desktop, shown on mobile) ── */
.post-meta {
  display: none;
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-faint);
  line-height: 1.4;
}

/* ── Summary ── */
.post-summary {
  margin: 10px 0 0;
  font-size: 16px;
  line-height: 1.75;
  color: var(--text-soft);
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

/* ── Skeleton ── */
.post-item--skeleton {
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

/* ── Pagination ── */
.pagination-wrapper {
  margin-top: 40px;
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

/* ── Fetch error ── */
.fetch-error {
  text-align: center;
  color: var(--text-faint);
  padding: 40px 0;
  font-size: 0.95rem;
}

.fetch-error .el-button {
  margin-top: 16px;
}

/* ── Dark mode ── */
html.dark .post-item {
  border-color: rgba(255, 255, 255, 0.04);
}

html.dark .post-date {
  border-color: var(--border);
  background: var(--surface-muted);
}

/* ── Responsive ── */
@media (max-width: 768px) {
  .home {
    padding: 0 10px;
  }

  .post-item {
    display: block;
    padding: 18px 0;
  }

  .post-date {
    display: none;
  }

  .post-title {
    font-size: 18px;
    line-height: 1.45;
  }

  .post-title a::after {
    width: 96px;
    margin-top: 12px;
  }

  .post-meta {
    display: block;
    margin-top: 10px;
    font-size: 12px;
  }

  .post-summary {
    margin-top: 10px;
    font-size: 15px;
    line-height: 1.7;
    -webkit-line-clamp: 3;
  }
}
</style>
