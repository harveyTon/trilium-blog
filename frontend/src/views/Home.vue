<template>
  <div class="home">
    <div class="home-sections">
      <FeaturedSection :items="featuredPosts" />
    </div>

    <section class="home-feed-section">
      <div class="home-section-header">
        <p class="home-section-kicker">{{ t('home.sectionKicker') }}</p>
        <h1 class="home-section-title">{{ t('home.sectionTitle') }}</h1>
      </div>

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

      <div v-else-if="fetchError" class="fetch-error">
        <p>{{ fetchError }}</p>
        <el-button type="primary" @click="loadPosts">{{ t('home.retry') }}</el-button>
      </div>

      <el-empty v-else-if="!posts.length" :description="t('home.noPosts')" />
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
        :prev-text="'←'"
        :next-text="'→'"
        layout="prev, pager, next"
        :ellipsis-item="'...'"
      />
    </div>
  </div>
</template>

<script>
import { ElButton, ElEmpty, ElPagination } from "element-plus";
import { onActivated, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import FeaturedSection from "../components/home/FeaturedSection.vue";
import PostFeed from "../components/home/PostFeed.vue";
import { fetchFeaturedPosts, fetchPosts } from "../api/blog";
import { t, locale } from "../i18n";
import { useSiteStore } from "../store";

const parsePageQuery = (value) => {
  const raw = Array.isArray(value) ? value[0] : value;
  const page = Number.parseInt(raw || "1", 10);
  return Number.isFinite(page) && page > 0 ? page : 1;
};

export default {
  name: "HomePage",
  components: {
    ElButton,
    ElEmpty,
    ElPagination,
    FeaturedSection,
    PostFeed,
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const siteStore = useSiteStore();
    const posts = ref([]);
    const featuredPosts = ref([]);
    const loading = ref(true);
    const fetchError = ref(null);
    const currentPage = ref(parsePageQuery(route.query.page));
    const featuredLoaded = ref(false);
    const loadedPage = ref(0);
    const pagination = ref({
      page: 1,
      pageSize: 9,
      total: 0,
      totalPages: 0,
    });

    const loadFeaturedPosts = async () => {
      if (featuredLoaded.value) {
        return;
      }
      try {
        const featuredResponse = await fetchFeaturedPosts();
        featuredPosts.value = featuredResponse.items || [];
        featuredLoaded.value = true;
      } catch {
        featuredPosts.value = [];
      }
    };

    const loadPosts = async (page = currentPage.value) => {
      loading.value = true;
      fetchError.value = null;
      try {
        const postResponse = await fetchPosts(page);
        posts.value = postResponse.items;
        currentPage.value = postResponse.page;
        loadedPage.value = postResponse.page;
        pagination.value = {
          page: postResponse.page,
          pageSize: postResponse.pageSize,
          total: postResponse.total,
          totalPages: postResponse.totalPages,
        };
      } catch {
        fetchError.value = t('home.fetchError');
        posts.value = [];
      } finally {
        loading.value = false;
      }
    };

    const handleCurrentChange = async (page) => {
      const query = page > 1 ? { page: String(page) } : {};
      await router.push({ name: "HomePage", query });
      window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const pad = (n) => String(n).padStart(2, "0");

    const getDay = (dateString) => {
      return pad(new Date(dateString).getDate());
    };

    const getMonth = (dateString) => {
      const d = new Date(dateString);
      return `${d.getFullYear()}${t('date.yearSuffix')}-${pad(d.getMonth() + 1)}${t('date.monthSuffix')}`;
    };

    const formatFullDate = (dateString) => {
      const d = new Date(dateString);
      const h = d.getHours();
      const m = pad(d.getMinutes());
      const period = h < 12 ? t('date.am') : t('date.pm');
      const displayHour = h > 12 ? h - 12 : h === 0 ? 12 : h;
      return `${d.getFullYear()}${t('date.yearSuffix')}${d.getMonth() + 1}${t('date.monthSuffix')}${d.getDate()}${t('date.daySuffix')} ${period}${displayHour}:${m}`;
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

    const buildSiteTitle = () => {
      const baseTitle = [siteStore.site.title, siteStore.site.subtitle].filter(Boolean).join(" | ");
      if (!baseTitle) {
        return "";
      }
      return currentPage.value > 1 ? `${baseTitle} - ${t('home.pageSuffix', { page: currentPage.value })}` : baseTitle;
    };

    const syncTitle = () => {
      const nextTitle = buildSiteTitle();
      if (nextTitle) {
        document.title = nextTitle;
      }
    };

    onMounted(async () => {
      await loadFeaturedPosts();
      await loadPosts(currentPage.value);
      syncTitle();
    });

    onActivated(() => {
      syncTitle();
    });

    watch(
      () => route.query.page,
      async (nextPage) => {
        const parsedPage = parsePageQuery(nextPage);
        if (parsedPage === loadedPage.value && posts.value.length) {
          return;
        }
        currentPage.value = parsedPage;
        await loadPosts(parsedPage);
        syncTitle();
      }
    );

    watch(
      () => [siteStore.site.title, siteStore.site.subtitle, currentPage.value],
      syncTitle,
      { immediate: true }
    );

    return {
      t,
      locale,
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
  padding: 10px 16px 32px;
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
  margin-top: 44px;
}

.home-section-header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.home-section-kicker {
  margin: 0;
  font-size: 12px;
  color: var(--text-faint);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.home-section-title {
  margin: 0;
  color: var(--text);
  font-size: clamp(30px, 4vw, 42px);
  line-height: 1.16;
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
  margin: 20px 0 8px;
  padding-bottom: 12px;
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

  .home-section-title {
    font-size: 28px;
  }

  .home-feed-section {
    margin-top: 34px;
  }

  .post-item--skeleton {
    display: block;
  }

  .skeleton-date {
    display: none;
  }
}
</style>
