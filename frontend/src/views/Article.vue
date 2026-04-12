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
          v-for="i in 10"
          :key="i"
          variant="p"
          style="width: 100%; height: 16px; margin-bottom: 12px"
        />
      </template>
      <template #default>
        <div class="article-container">
          <div class="article-layout">
            <el-card v-if="post" class="article-card">
              <template #header>
                <div class="card-header">
                  <div class="article-fword">
                    {{ post.title.charAt(0).toUpperCase() }}
                  </div>
                  <h1 class="article-title">{{ post.title }}</h1>
                  <div class="article-date">
                    {{ formatDate(post.dateModified) }}
                  </div>
                  <span class="artalk-pv-count" style="display: none"></span>
                </div>
              </template>

              <div class="article-content" v-html="post.contentHtml"></div>

              <div v-if="post.pageUrl" class="article-source">
                剪贴自：
                <a
                  :href="post.pageUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {{ post.pageUrl }}
                </a>
              </div>

              <div v-if="site.comments.enabled" class="article-comments">
                <h2>评论</h2>
                <div ref="artalkContainer"></div>
              </div>
            </el-card>
            <el-empty v-else description="文章未找到"></el-empty>
          </div>

          <el-affix
            v-if="post && post.toc && post.toc.length >= 3"
            :offset="60"
            class="article-anchor-wrapper"
          >
            <el-popover
              placement="right"
              :width="220"
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
                <el-anchor :bounds="0" :offset="70">
                  <el-anchor-link
                    v-for="item in post.toc"
                    :key="item.id"
                    :href="'#' + item.id"
                    :title="item.title"
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
import hljs from "highlight.js";
import "highlight.js/styles/atom-one-dark.css";
import { storeToRefs } from "pinia";
import { nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { fetchPost } from "../api/blog";
import { useSiteStore } from "../store";

export default {
  name: "ArticlePage",
  components: { Menu },
  setup() {
    const route = useRoute();
    const siteStore = useSiteStore();
    const { site } = storeToRefs(siteStore);
    const post = ref(null);
    const loading = ref(true);
    const anchorVisible = ref(false);
    const artalkContainer = ref(null);
    let artalkInstance = null;
    let darkModeObserver = null;

    const isDarkMode = () =>
      document.documentElement.classList.contains("dark");

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    };

    const highlightCode = () => {
      document.querySelectorAll("pre code").forEach((el) => {
        const code = el.textContent ?? "";
        const classList = el.className;
        const languageMatch = classList.match(/language-(\S+)/);
        if (languageMatch) {
          el.innerHTML = hljs.highlight(code, {
            language: languageMatch[1],
          }).value;
        } else {
          el.innerHTML = hljs.highlightAuto(code).value;
        }
        el.classList.add("hljs");
      });
    };

    const setupGallery = () => {
      document.querySelectorAll(".article-content img").forEach((img) => {
        const parent = img.parentElement;
        if (!parent) {
          return;
        }
        if (parent.tagName !== "A") {
          const wrapper = document.createElement("a");
          wrapper.href = img.src;
          wrapper.target = "_blank";
          wrapper.dataset.fancybox = "gallery";
          parent.replaceChild(wrapper, img);
          wrapper.appendChild(img);
          return;
        }
        parent.href = img.src;
        parent.target = "_blank";
        parent.dataset.fancybox = "gallery";
      });
      Fancybox.bind("[data-fancybox]", {});
    };

    const destroyComments = () => {
      if (artalkInstance) {
        artalkInstance.destroy();
        artalkInstance = null;
      }
    };

    const initComments = () => {
      destroyComments();
      if (!site.value.comments.enabled || !artalkContainer.value || !post.value) {
        return;
      }
      artalkInstance = Artalk.init({
        el: artalkContainer.value,
        pageKey: `/posts/${route.params.noteId}`,
        pageTitle: post.value.title,
        pvEl: ".artalk-pv-count",
        server: site.value.comments.server,
        site: site.value.comments.site || site.value.name,
        darkMode: isDarkMode(),
      });
    };

    const enhanceContent = async () => {
      await nextTick();
      highlightCode();
      setupGallery();
      initComments();
    };

    const syncTitle = () => {
      if (post.value && site.value.title) {
        document.title = `${post.value.title} - ${site.value.title}`;
      }
    };

    const loadPost = async () => {
      loading.value = true;
      try {
        post.value = await fetchPost(route.params.noteId);
        await enhanceContent();
        syncTitle();
      } finally {
        loading.value = false;
        if (typeof window.scrollTo === "function") {
          window.scrollTo({ top: 0, behavior: "smooth" });
        }
      }
    };

    const observeDarkMode = () => {
      const observer = new MutationObserver(() => {
        if (artalkInstance) {
          artalkInstance.setDarkMode(isDarkMode());
        }
      });
      observer.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ["class"],
      });
      return observer;
    };

    onMounted(async () => {
      await loadPost();
      darkModeObserver = observeDarkMode();
    });

    onUnmounted(() => {
      Fancybox.destroy();
      destroyComments();
      if (darkModeObserver) {
        darkModeObserver.disconnect();
      }
    });

    watch(() => route.params.noteId, loadPost);
    watch([post, site], syncTitle, { immediate: true });

    return {
      site,
      post,
      loading,
      anchorVisible,
      artalkContainer,
      formatDate,
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
  width: min(860px, 100%);
}

.article-card {
  border-radius: 8px;
}

.card-header {
  position: relative;
}

.article-fword {
  position: absolute;
  top: -28px;
  left: -12px;
  font-size: 6rem;
  opacity: 0.08;
  font-weight: 700;
  line-height: 1;
}

.article-title {
  margin: 0 0 16px;
  line-height: 1.3;
  position: relative;
}

.article-date,
.article-source {
  color: #7a7a7a;
}

.article-content {
  line-height: 1.8;
  overflow-wrap: anywhere;
}

.article-content img {
  max-width: 100%;
  height: auto;
  border-radius: 6px;
}

.article-content pre {
  overflow-x: auto;
}

.article-comments {
  margin-top: 40px;
}

.article-anchor-wrapper {
  margin-left: 20px;
}

.anchor-toggle {
  border-radius: 6px;
}

@media (max-width: 960px) {
  .article-container {
    display: block;
  }

  .article-anchor-wrapper {
    display: none;
  }
}
</style>
