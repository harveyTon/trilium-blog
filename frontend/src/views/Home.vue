<template>
  <div class="home">
    <div v-if="loading" class="article-list">
      <el-skeleton
        v-for="i in 5"
        :key="i"
        :rows="3"
        animated
        class="article-skeleton"
      />
    </div>
    <div v-else class="article-list">
      <h3 class="article-latest">最新文章</h3>
      <el-empty v-if="!posts.length" description="暂无已发布文章" />
      <div v-if="fetchError" class="fetch-error">{{ fetchError }}</div>
      <article
        v-for="post in posts"
        :key="post.noteId"
        class="article-item"
      >
        <div class="article-beici">
          {{ post.title ? post.title.charAt(0).toUpperCase() : '' }}
        </div>
        <h2 class="article-title">
          <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }">
            {{ post.title }}
          </router-link>
        </h2>
        <div class="article-info">
          <span>{{ formatDate(post.dateModified) }}</span>
        </div>
        <div class="article-more">
          <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }">
            阅读全文
          </router-link>
        </div>
      </article>
      <el-pagination
        v-if="pagination.totalPages > 1"
        v-model:current-page="currentPage"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        @current-change="handleCurrentChange"
        layout="prev, pager, next"
        background
        class="pagination"
      />
    </div>
  </div>
</template>

<script>
import { ElEmpty, ElPagination, ElSkeleton } from "element-plus";
import { onMounted, ref, watch } from "vue";
import { fetchPosts } from "../api/blog";
import { useSiteStore } from "../store";

export default {
  name: "HomePage",
  components: {
    ElEmpty,
    ElPagination,
    ElSkeleton,
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
    };

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
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
    };
  },
};
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
}

a {
  color: var(--text-primary);
  text-decoration: none;
}

.article-latest {
  font-size: 1.5rem;
  text-align: center;
  padding-bottom: 20px;
  position: relative;
  font-weight: 300;
  letter-spacing: 2px;
  margin-bottom: 30px;
}

.article-latest::after {
  content: "";
  position: absolute;
  width: 50px;
  height: 1px;
  background: var(--border-color);
  bottom: 0;
  left: calc(50% - 25px);
}

.article-list {
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 12px var(--border-color);
  background-color: var(--bg-surface);
  padding: clamp(40px, 8vw, 100px) clamp(20px, 15vw, 200px);
  min-height: 400px;
}

.article-item {
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
  margin-bottom: clamp(40px, 8vw, 100px);
}

.article-item:last-of-type {
  margin-bottom: 40px;
}

.article-beici {
  position: absolute;
  top: -20px;
  left: 0;
  font-size: clamp(3rem, 10vw, 8rem);
  opacity: 0.08;
  font-weight: bold;
  z-index: 0;
  pointer-events: none;
  line-height: 1;
  user-select: none;
}

.article-title {
  margin-bottom: 20px;
  font-size: clamp(1.3rem, 4vw, 2rem);
  line-height: 1.4;
  position: relative;
  font-weight: bold;
}

.article-info {
  margin-bottom: 24px;
  color: var(--text-secondary);
}

.article-more a {
  background: var(--bg-header-footer);
  color: var(--text-inverse);
  padding: 13px 40px;
  display: inline-block;
  border: 1px solid var(--bg-header-footer);
  font-size: 16px;
  border-radius: 6px;
}

.article-skeleton {
  margin-bottom: 30px;
}

.pagination {
  align-self: center;
}

.fetch-error {
  text-align: center;
  color: var(--text-muted);
  padding: 40px 0;
  font-size: 0.95rem;
}
</style>
