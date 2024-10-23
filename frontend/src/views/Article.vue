<template>
  <div class="article">
    <el-skeleton :loading="loading" animated>
      <template #template>
        <el-skeleton-item
          variant="p"
          style="width: 100%; height: 32px; margin-bottom: 16px"
        />
        <el-skeleton-item
          variant="text"
          style="width: 30%; margin-bottom: 16px"
        />
        <el-skeleton-item
          variant="p"
          style="width: 100%; height: 16px; margin-bottom: 12px"
          v-for="i in 10"
          :key="i"
        />
      </template>
      <template #default>
        <div class="article-container">
          <div class="article-layout">
            <el-card v-if="article" class="article-card">
              <template #header>
                <div class="card-header">
                  <div class="article-fword">
                    {{ article.title.charAt(0).toUpperCase() }}
                  </div>
                  <h1 class="article-title">
                    {{
                      article.title.charAt(0).toUpperCase() +
                      article.title.slice(1)
                    }}
                  </h1>
                  <div class="article-date">
                    {{ formatDate(article.dateModified) }}
                  </div>
                  <span class="artalk-pv-count" style="display: none"></span>
                </div>
              </template>

              <div class="article-content" v-html="processedContent"></div>

              <div v-if="article.pageUrl" class="article-source">
                剪贴自：<a
                  :href="article.pageUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  >{{ article.pageUrl }}</a
                >
              </div>
              <div class="article-comments">
                <h2>评论</h2>
                <div ref="artalkContainer"></div>
              </div>
            </el-card>
            <el-empty v-else description="文章未找到"></el-empty>
          </div>

          <el-affix
            :offset="60"
            class="article-anchor-wrapper"
            v-if="anchors.length >= 3"
          >
            <el-popover
              placement="right"
              :width="200"
              trigger="click"
              :visible="anchorVisible"
              @hide="anchorVisible = false"
            >
              <template #reference>
                <el-button
                  class="anchor-toggle"
                  @click="anchorVisible = !anchorVisible"
                >
                  <el-icon><Menu /></el-icon>
                </el-button>
              </template>
              <el-scrollbar max-height="calc(100vh - 120px)">
                <el-anchor
                  :bounds="0"
                  :offset="70"
                  @change="handleAnchorChange"
                >
                  <el-anchor-link
                    v-for="(anchor, index) in anchors"
                    :key="index"
                    :href="'#' + anchor.id"
                    :title="anchor.title"
                  />
                </el-anchor>
              </el-scrollbar>
            </el-popover>
          </el-affix>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<script>
import { Menu } from "@element-plus/icons-vue";
import { Fancybox } from "@fancyapps/ui";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import Artalk from "artalk";
import "artalk/dist/Artalk.css";
import axios from "axios";
import { storeToRefs } from "pinia";
import Prism from "prismjs";
import { nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useBlogStore } from "../store";

export default {
  name: "ArticlePage",
  components: { Menu },
  setup() {
    const route = useRoute();
    const blogStore = useBlogStore();
    const { blogInfo } = storeToRefs(blogStore);
    const article = ref(null);
    const processedContent = ref("");
    const loading = ref(true);
    const anchors = ref([]);
    const activeAnchor = ref("");
    const observer = ref(null);
    const anchorVisible = ref(false);
    const artalkContainer = ref(null);
    let artalkInstance = null;

    const isDarkMode = () =>
      document.documentElement.classList.contains("dark");

    const initArtalk = () => {
      if (artalkContainer.value && article.value) {
        artalkInstance = Artalk.init({
          el: artalkContainer.value,
          pageKey: "/post/" + route.params.noteId,
          pageTitle: article.value.title,
          pvEl: ".artalk-pv-count",
          server: "https://comments.uto.to",
          site: blogInfo.value.blogName,
          darkMode: isDarkMode(),
        });
      }
    };

    const handleAnchorChange = (link) => {
      activeAnchor.value = link.slice(1);
    };

    const scrollToAnchor = (id) => {
      const element = document.getElementById(id);
      if (element) {
        const yOffset = -80;
        const y =
          element.getBoundingClientRect().top + window.pageYOffset + yOffset;
        window.scrollTo({ top: y, behavior: "smooth" });
      }
    };

    const setupIntersectionObserver = () => {
      observer.value = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              activeAnchor.value = entry.target.id;
            }
          });
        },
        { threshold: 0.5 }
      );

      anchors.value.forEach((anchor) => {
        const element = document.getElementById(anchor.id);
        if (element) observer.value.observe(element);
      });
    };

    onMounted(() => {
      fetchArticle().then(() => {
        nextTick(() => {
          initArtalk();
        });
      });
    });

    onUnmounted(() => {
      if (observer.value) {
        observer.value.disconnect();
      }
      if (artalkInstance) {
        artalkInstance.destroy();
      }
    });
    const observeDarkMode = () => {
      const targetNode = document.documentElement;
      const config = { attributes: true, attributeFilter: ["class"] };

      const callback = (mutationsList) => {
        for (let mutation of mutationsList) {
          if (
            mutation.type === "attributes" &&
            mutation.attributeName === "class"
          ) {
            const isDark = isDarkMode();
            if (artalkInstance) {
              artalkInstance.setDarkMode(isDark);
            }
          }
        }
      };

      const observer = new MutationObserver(callback);
      observer.observe(targetNode, config);

      return observer;
    };

    let darkModeObserver;

    onMounted(() => {
      darkModeObserver = observeDarkMode();
    });

    onUnmounted(() => {
      Fancybox.destroy();
      if (darkModeObserver) {
        darkModeObserver.disconnect();
      }
    });

    const fetchArticle = async () => {
      loading.value = true;
      try {
        const noteId = route.params.noteId;
        const response = await axios.get(`/api/articles/${noteId}`, {
          params: {
            t: new Date().getTime(),
          },
        });
        article.value = response.data;
        await processContent();
        updateTitle();
      } catch (error) {
        console.error("Fetch Article Error:", error);
      } finally {
        loading.value = false;
        nextTick(() => {
          window.scrollTo({ top: 0, behavior: "smooth" });
        });
      }
    };

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    };

    const processContent = async () => {
      if (!article.value) return;

      const tempDiv = document.createElement("div");
      tempDiv.innerHTML = article.value.content;

      tempDiv.querySelectorAll("img").forEach((img) => {
        const parent = img.parentElement;
        if (parent.tagName !== "A") {
          const wrapper = document.createElement("a");
          wrapper.href = img.src;
          wrapper.target = "_blank";
          wrapper.dataset.fancybox = "gallery";
          parent.replaceChild(wrapper, img);
          wrapper.appendChild(img);
        } else if (parent.href !== img.src) {
          parent.href = img.src;
          parent.target = "_blank";
          parent.dataset.fancybox = "gallery";
        }
        console.log("Process Image:", img.src);
      });

      anchors.value = [];
      tempDiv
        .querySelectorAll("h1, h2, h3, h4, h5, h6")
        .forEach((heading, index) => {
          const id = `heading-${index}`;
          heading.id = id;
          anchors.value.push({
            id: id,
            title: heading.textContent,
          });
        });

      tempDiv.querySelectorAll("pre code").forEach((block) => {
        Prism.highlightElement(block);
      });

      processedContent.value = tempDiv.innerHTML;

      nextTick(() => {
        document
          .querySelectorAll(".article-content pre code")
          .forEach((block) => {
            Prism.highlightElement(block);
          });

        Fancybox.bind("[data-fancybox]", {});
        setupIntersectionObserver();
      });
    };

    const updateTitle = () => {
      if (article.value && blogInfo.value.title) {
        document.title = `${article.value.title} - ${blogInfo.value.title} - Powered By Trilium Blog`;
      }
    };

    watch([article, blogInfo], updateTitle, { immediate: true });

    return {
      blogInfo,
      article,
      processedContent,
      loading,
      anchors,
      activeAnchor,
      scrollToAnchor,
      formatDate,
      handleAnchorChange,
      anchorVisible,
      artalkContainer,
    };
  },
};
</script>

<style>
@import "artalk/dist/Artalk.css";

.article-container {
  display: flex;
  justify-content: center;
  position: relative;
}

.article-layout {
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
}

.article-card {
  width: 100%;
  box-sizing: border-box;
  padding: 100px 200px;
}

.article-anchor-wrapper {
  position: fixed;
  left: calc(50% - 450px);
  top: 335px;
  z-index: 1000;
}

.anchor-toggle {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
}

.article-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
  font-size: 36px;
  line-height: 1.4;
  position: relative;
  font-weight: bold;
}
.article-fword {
  position: absolute;
  top: 60px;
  left: 100px;
  font-size: 14em;
  opacity: 0.1;
  font-weight: bold;
  z-index: 0;
}
.article-date {
  font-size: 14px;
  color: #909399;
}

.article-content {
  font-size: 16px;
  line-height: 1.6;
  color: #303133;
  overflow-wrap: break-word;
  word-wrap: break-word;
}
.article-content a {
  color: #909399;
  text-decoration: underline gray;
  font-style: italic;
}
.article-source a::after {
  content: url("../assets/icons/external-link.svg");
  margin-left: 4px;
  vertical-align: middle;
}
.article-content h1,
.article-content h2,
.article-content h3,
.article-content h4,
.article-content h5,
.article-content h6 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  scroll-margin-top: 80px;
}

.article-content p {
  margin-bottom: 16px;
}

.article-content img {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 0 auto;
}

.article-content pre {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 16px;
  overflow: auto;
  max-width: 100%;
}

.article-content code {
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier,
    monospace;
  font-size: 14px;
  max-width: 100%;
  overflow-x: auto;
}

.article-source {
  padding: 10px 0 10px 0;
  margin: 20px 0 20px 0;
  margin-bottom: 0.5rem;
  color: #666;
  border: 1px solid #ddd;
  border-left: 0;
  border-right: 0;
  font-style: italic;
}

.article-source a {
  color: #409eff;
  text-decoration: none;
}

.article-source a:hover {
  text-decoration: underline;
}

@media (max-width: 1200px) {
  .article-anchor-wrapper {
    left: 20px;
  }
}

@media (max-width: 768px) {
  .article-layout {
    max-width: 100%;
  }

  .article-card {
    padding: 0;
    border-radius: 0;
    box-shadow: none;
  }

  :deep(.el-card__body) {
    padding: 10px;
  }

  .article-title {
    font-size: 20px;
    margin-bottom: 5px;
  }

  .article-date {
    font-size: 12px;
  }

  .article-content {
    font-size: 14px;
  }

  .article-content h1 {
    font-size: 20px;
  }
  .article-content h2 {
    font-size: 18px;
  }
  .article-content h3 {
    font-size: 16px;
  }
  .article-content h4,
  h5,
  h6 {
    font-size: 14px;
  }

  .article-content pre {
    padding: 10px;
  }

  .article-content code {
    font-size: 12px;
  }

  .article-source {
    font-size: 12px;
  }

  .article-comments h2 {
    font-size: 18px;
    margin-bottom: 15px;
  }

  .article-anchor-wrapper {
    display: none;
  }

  .article-fword {
    display: none;
  }
}

.article-comments {
  margin-top: 40px;
}

.article-comments h2 {
  font-size: 20px;
  margin-bottom: 20px;
}

.atk-copyright {
  display: none !important;
}

html.dark .article-card {
  background-color: #2e2e2e;
  color: #e0e0e0;
  border-color: #444;
}

html.dark .article-title {
  color: #ffffff;
}

html.dark .article-date {
  color: #bbbbbb;
}

html.dark .article-content {
  color: #dddddd;
}

html.dark .article-content a {
  color: #80a0ff;
}

html.dark .article-content a:hover {
  color: #ffffff;
}

html.dark .article-source {
  color: #bbbbbb;
  border-color: #444;
}

html.dark .article-source a {
  color: #80a0ff;
}

html.dark .article-source a:hover {
  color: #ffffff;
}

html.dark .article-content h1,
html.dark .article-content h2,
html.dark .article-content h3,
html.dark .article-content h4,
html.dark .article-content h5,
html.dark .article-content h6 {
  color: #ffffff;
}

html.dark .article-content pre {
  background-color: #444444;
  color: #e0e0e0;
}

html.dark .article-content code {
  color: #80a0ff;
}

html.dark .article-comments h2 {
  color: #ffffff;
}

html.dark .anchor-toggle {
  background-color: #444444;
  color: #e0e0e0;
}

.fancybox__container {
  --fancybox-bg: rgba(24, 24, 27, 0.98);
}

html.dark .fancybox__container {
  --fancybox-bg: rgba(0, 0, 0, 0.98);
}
</style>
