<template>
  <div class="home">
    <div v-if="loading" class="post-list" aria-busy="true">
      <div v-for="i in 5" :key="i" class="post-card post-card--skeleton">
        <div class="skeleton-meta"></div>
        <div class="skeleton-title"></div>
        <div class="skeleton-title skeleton-title--short"></div>
      </div>
    </div>

    <div v-else>
      <el-empty v-if="!posts.length && !fetchError" description="暂无已发布文章" />
      <div v-if="fetchError" class="fetch-error">
        <p>{{ fetchError }}</p>
        <el-button type="primary" @click="loadPosts">重试</el-button>
      </div>
      <nav class="post-list" aria-label="文章列表">
        <router-link
          v-for="post in posts"
          :key="post.noteId"
          :to="{ name: 'Article', params: { noteId: post.noteId } }"
          class="post-card"
          :aria-label="post.title + '，' + formatDate(post.dateModified)"
        >
          <div class="post-meta">
            <time :datetime="post.dateModified">{{ formatDate(post.dateModified) }}</time>
          </div>
          <h2 class="post-title">{{ post.title }}</h2>
          <p v-if="post.summary" class="post-summary">{{ sanitizeSummary(post.summary) }}</p>
          <div class="post-read-cue" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
        </router-link>
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

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
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
      formatDate,
      sanitizeSummary,
      loadPosts,
    };
  },
};
</script>

<style scoped>
/* ── Page shell ── */
.home {
  max-width: var(--list-w);
  margin: 0 auto;
  padding: 0 16px;
}

/* ── List ── */
.post-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.post-list:focus-within {
  outline: none;
}

/* ── Card ── */
.post-card {
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 28px 0;
  border-bottom: 1px solid var(--border-soft);
  text-decoration: none;
  transition: border-color 160ms ease;
}

.post-card:first-child {
  border-top: 1px solid var(--border-soft);
}

.post-card:hover {
  border-bottom-color: var(--border);
}

.post-card:hover .post-title {
  color: var(--link);
}

.post-card:hover .post-read-cue {
  opacity: 1;
  transform: translateX(0);
}

/* ── Meta ── */
.post-meta {
  margin-bottom: 10px;
}

.post-meta time {
  color: var(--text-faint);
  font-size: 13px;
  font-variant-numeric: tabular-nums;
  letter-spacing: 0.01em;
}

/* ── Title ── */
.post-title {
  margin: 0;
  color: var(--text);
  font-size: clamp(18px, 3vw, 22px);
  line-height: 1.35;
  font-weight: 600;
  letter-spacing: -0.005em;
  padding-right: 40px;
  transition: color 160ms ease;
}

/* ── Read cue (right arrow) ── */
.post-read-cue {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateX(-4px) translateY(-50%);
  color: var(--accent);
  opacity: 0;
  transition: opacity 160ms ease, transform 160ms ease;
  pointer-events: none;
}

/* ── Summary ── */
.post-summary {
  margin: 8px 0 0;
  color: var(--text-soft);
  font-size: 14px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  padding-right: 40px;
}

/* ── Skeleton ── */
.post-card--skeleton {
  pointer-events: none;
}

.skeleton-meta {
  width: 120px;
  height: 12px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 14px;
  animation: skeleton-pulse 1.6s ease-in-out infinite;
}

.skeleton-title {
  width: 70%;
  height: 20px;
  background: var(--border-soft);
  border-radius: 4px;
  margin-bottom: 8px;
  animation: skeleton-pulse 1.6s ease-in-out infinite;
}

.skeleton-title--short {
  width: 45%;
  animation-delay: 0.15s;
}

@keyframes skeleton-pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

/* ── Pagination wrapper ── */
.pagination-wrapper {
  margin-top: 40px;
  padding-bottom: 24px;
  display: flex;
  justify-content: center;
}

/* Override Element Plus pagination defaults */
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

/* ── Responsive ── */
@media (max-width: 768px) {
  .home {
    padding: 0 10px;
  }

  .post-title {
    font-size: 17px;
    padding-right: 32px;
  }
}
</style>
