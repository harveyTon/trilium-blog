<template>
  <div class="home">
    <div v-if="loading" class="article-list">
      <el-skeleton
        :rows="3"
        animated
        class="article-skeleton"
        v-for="i in 5"
        :key="i"
      />
    </div>
    <div v-else>
      <div class="article-list">
        <h3 class="article-latest">最新文章</h3>
        <div
          v-for="article in articles"
          :key="article.noteId"
          class="article-item"
        >
          <div class="article-beici">
            {{ article.title.charAt(0).toUpperCase() }}
          </div>
          <h2 class="article-title">
            <router-link
              :to="{ name: 'Article', params: { noteId: article.noteId } }"
            >
              {{
                article.title.charAt(0).toUpperCase() + article.title.slice(1)
              }}
            </router-link>
          </h2>
          <div class="article-info">
            <span>{{ formatDate(article.dateModified) }}</span>
          </div>

          <div class="article-image">
            <router-link
              :to="{ name: 'Article', params: { noteId: article.noteId } }"
            >
              <img
                :src="getRandomImage(article.noteId)"
                :alt="article.title"
                loading="lazy"
              />
            </router-link>
          </div>

          <div class="article-more">
            <router-link
              :to="{ name: 'Article', params: { noteId: article.noteId } }"
            >
              阅读更多
            </router-link>
          </div>
        </div>
      </div>
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="totalArticles"
        @current-change="handleCurrentChange"
        layout="prev, pager, next"
        background
        class="pagination"
      />
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { ElPagination, ElSkeleton } from "element-plus";
import { onMounted, ref, watch } from "vue";
import { useBlogStore } from "../store";

export default {
  name: "HomePage",
  components: {
    ElPagination,
    ElSkeleton,
  },
  setup() {
    const blogStore = useBlogStore();
    const articles = ref([]);
    const loading = ref(true);
    const currentPage = ref(1);
    const pageSize = ref(9);
    const totalArticles = ref(0);

    const fetchArticles = async () => {
      loading.value = true;
      try {
        const response = await axios.get(`/api/articles`, {
          params: {
            page: currentPage.value,
            pageSize: pageSize.value,
            t: new Date().getTime(),
          },
        });
        articles.value = response.data.articles;
        totalArticles.value = response.data.totalArticles;
      } catch (error) {
        console.error("Fetch Articles Error:", error);
      } finally {
        loading.value = false;
      }
    };

    const handleCurrentChange = (page) => {
      currentPage.value = page;
      fetchArticles();
    };

    onMounted(() => {
      fetchArticles();
      if (blogStore.blogInfo.blogTitle) {
        document.title = `${blogStore.blogInfo.title} - Powered By Trilium Blog`;
      }
    });

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    };

    const getRandomImage = (noteId) => {
      return `https://88900.net/api/imageproxy/https://picsum.photos/seed/${noteId}/840/420`;
    };

    watch(
      () => blogStore.blogInfo,
      (newBlogInfo) => {
        if (newBlogInfo.blogTitle) {
          document.title = `${newBlogInfo.title} - Powered By Trilium Blog`;
        }
      },
      { immediate: true }
    );

    return {
      articles,
      loading,
      currentPage,
      pageSize,
      totalArticles,
      handleCurrentChange,
      formatDate,
      getRandomImage,
    };
  },
};
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0px;
}

a {
  color: #333;
  text-decoration: none;
  background-color: transparent;
}

.article-latest {
  font-size: 1.5rem;
  letter-spacing: 3px;
  text-align: center;
  padding-bottom: 20px;
  position: relative;
  font-weight: 300;
  letter-spacing: 3px;
  margin-bottom: 30px;
}

.article-latest:after {
  content: "";
  position: absolute;
  width: 50px;
  height: 1px;
  background: rgba(51, 51, 51, 0.2);
  bottom: 0px;
  left: calc(50% - 25px);
}

.article-list {
  display: flex;
  flex-direction: column;
  box-shadow: 0px 0px 12px rgba(0, 0, 0, 0.12);
  background-color: #ffffff;
  padding: 100px 200px;
}

.article-item {
  display: flex;
  flex-direction: column;
  position: relative;
  margin-bottom: 180px;
}

.article-item:last-child {
  margin-bottom: 30px;
}

.article-beici {
  position: absolute;
  top: -70px;
  left: -100px;
  font-size: 13em;
  opacity: 0.1;
  font-weight: bold;
  z-index: 0;
}

.article-title {
  margin-bottom: 20px;
  font-size: 36px;
  line-height: 1.4;
  position: relative;
  font-weight: bold;
}

.article-summary {
  color: #5c6b7f;
  margin-bottom: 20px;
  line-height: 1.6;
  font-size: 16px;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
  padding-top: 15px;
}

.article-info {
  margin-bottom: 30px;
  color: #888888;
  position: relative;
}

.article-info span {
  font-weight: 300;
  position: relative;
  padding-right: 30px;
}

.article-date {
  margin-bottom: 30px;
  color: #888888;
  position: relative;
}

.read-more-link {
  color: #409eff;
  text-decoration: none;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s ease;
}

.read-more-link:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.article-more {
  position: relative;
  margin-top: 30px;
}

.article-more a {
  background: #2c3e50;
  color: #fff;
  padding: 13px 40px;
  margin-right: 15px;
  display: inline-block;
  border: 1px solid #2c3e50;
  font-size: 16px;
}

.article-more a i {
  margin-right: 5px;
}

.article-image {
  overflow: hidden;
}

.article-image img {
  width: 100%;
  height: 420px;
  object-fit: cover;
  transition: transform 0.3s ease-in-out, filter 0.3s ease-in-out;
}

.article-image img:hover {
  filter: brightness(80%);
}

.article-skeleton {
  margin-bottom: 30px;
}

@media (max-width: 768px) {
  .home {
    padding: 0;
    max-width: 100%;
  }

  .article-list {
    gap: 10px;
    padding: 10px 20px;
  }

  .article-item {
    border-radius: 0;
    box-shadow: none;
    margin-bottom: 60px;
  }

  .article-title {
    font-size: 24px;
    margin-bottom: 10px;
  }

  .article-summary {
    font-size: 14px;
    margin-bottom: 15px;
  }

  .article-footer {
    margin-top: 15px;
    padding-top: 10px;
  }

  .article-date,
  .read-more-link {
    font-size: 12px;
  }

  .article-content {
    padding: 15px 10px;
  }

  .article-image img {
    height: 150px;
  }

  .article-beici {
    display: none;
  }
}

html.dark .home {
  color: #e0e0e0;
}

html.dark a {
  color: #bbb;
}

html.dark .article-latest {
  color: #e0e0e0;
}

html.dark .article-latest:after {
  background: rgba(224, 224, 224, 0.2);
}

html.dark .article-list {
  background-color: #2e2e2e;
  box-shadow: 0px 0px 12px rgba(255, 255, 255, 0.1);
}

html.dark .article-title {
  color: #ffffff;
}

html.dark .article-info {
  color: #bbbbbb;
}

html.dark .article-date {
  color: #bbbbbb;
}

html.dark .article-more a {
  background-color: #dddddd;
  color: #1e1e1e;
  border-color: #bbbbbb;
}

html.dark .article-more a:hover {
  background-color: #bbbbbb;
  color: #1e1e1e;
}

html.dark .article-image img {
  filter: brightness(70%);
}

html.dark .article-skeleton {
  background-color: #444444;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

html.dark .pagination {
  --el-pagination-bg-color: #2e2e2e;
  --el-pagination-text-color: #e0e0e0;
  --el-pagination-button-color: #e0e0e0;
  --el-pagination-button-bg-color: #444444;
  --el-pagination-button-disabled-color: #888888;
  --el-pagination-button-disabled-bg-color: #333333;
  --el-pagination-hover-color: #ffffff;
}

:deep(html.dark .el-pagination) {
  --el-pagination-bg-color: #2e2e2e;
  --el-pagination-text-color: #e0e0e0;
  --el-pagination-button-color: #e0e0e0;
  --el-pagination-button-bg-color: #444444;
  --el-pagination-button-disabled-color: #888888;
  --el-pagination-button-disabled-bg-color: #333333;
  --el-pagination-hover-color: #ffffff;
}

:deep(html.dark .el-pagination .btn-prev),
:deep(html.dark .el-pagination .btn-next) {
  background-color: #444444;
  color: #e0e0e0;
}

:deep(html.dark .el-pagination .btn-prev:disabled),
:deep(html.dark .el-pagination .btn-next:disabled) {
  background-color: #333333;
  color: #888888;
}

:deep(html.dark .el-pagination .el-pager li) {
  background-color: #444444;
  color: #e0e0e0;
}

:deep(html.dark .el-pagination .el-pager li.active) {
  background-color: #409eff;
  color: #ffffff;
}

:deep(html.dark .el-pagination .el-pager li:hover) {
  color: #409eff;
}
</style>
