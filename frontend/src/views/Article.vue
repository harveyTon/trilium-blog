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
        <div v-if="post" :class="['article-shell', post.toc && post.toc.length >= 3 ? 'has-toc' : '']">
          <aside v-if="post.toc && post.toc.length >= 3" class="article-toc">
            <div class="toc-panel">
              <div class="toc-title">目录</div>
              <a
                v-for="item in post.toc"
                :key="item.id"
                :href="'#' + item.id"
                :class="['toc-link', 'toc-level-' + item.level, { 'is-active': activeHeading === item.id }]"
                @click.prevent="scrollToHeading(item.id)"
              >
                {{ item.title }}
              </a>
            </div>
          </aside>

          <main class="article-main">
            <header class="article-header">
              <h1 class="article-title">{{ post.title }}</h1>
              <div class="article-meta">
                <time>{{ formatDate(post.dateModified) }}</time>
              </div>
              <span class="artalk-pv-count" style="display: none"></span>
            </header>

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
          </main>
        </div>

        <el-empty v-else-if="!loadError" description="文章未找到"></el-empty>
        <div v-else class="load-error">
          <p>加载失败，请检查网络后重试</p>
          <el-button type="primary" @click="loadPost">重试</el-button>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<script>
import { ElButton } from "element-plus";
import { Fancybox } from "@fancyapps/ui";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import Artalk from "artalk";
import "artalk/dist/Artalk.css";
import hljs from "highlight.js/lib/core";
import javascript from "highlight.js/lib/languages/javascript";
import bash from "highlight.js/lib/languages/bash";
import "highlight.js/styles/atom-one-dark.css";
import "highlight.js/styles/atom-one-light.css";

hljs.registerLanguage("javascript", javascript);
hljs.registerLanguage("bash", bash);
import { storeToRefs } from "pinia";
import { nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { fetchPost } from "../api/blog";
import { useSiteStore } from "../store";

export default {
  name: "ArticlePage",
  components: { ElButton },
  setup() {
    const route = useRoute();
    const siteStore = useSiteStore();
    const { site } = storeToRefs(siteStore);
    const post = ref(null);
    const loading = ref(true);
    const loadError = ref(false);
    const activeHeading = ref("");
    const artalkContainer = ref(null);
    let artalkInstance = null;
    let darkModeObserver = null;
    let headingObserver = null;
    let tocScrollHandler = null;

    const isDarkMode = () =>
      document.documentElement.classList.contains("dark");

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    };

    const applyHighlightTheme = () => {
      const dark = isDarkMode();
      for (let i = document.styleSheets.length - 1; i >= 0; i--) {
        const sheet = document.styleSheets[i];
        try {
          const href = sheet.href || "";
          if (href.includes("atom-one-dark.css")) {
            sheet.disabled = !dark;
          } else if (href.includes("atom-one-light.css")) {
            sheet.disabled = dark;
          }
        } catch {}
      }
    };

    const highlightCode = () => {
      applyHighlightTheme();
      document.querySelectorAll("pre code").forEach((el) => {
        const code = el.textContent ?? "";
        const languageMatch = el.className.match(/language-(\S+)/);
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
        img.loading = "lazy";
        const parent = img.parentElement;
        if (!parent) return;
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
      if (!site.value.comments.enabled || !artalkContainer.value || !post.value) return;
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

    const setupHeadingObserver = () => {
      if (headingObserver) {
        headingObserver.disconnect();
        headingObserver = null;
      }
      const headings = document.querySelectorAll(".article-content h1, .article-content h2, .article-content h3");
      if (!headings.length) return;

      headingObserver = new IntersectionObserver(
        (entries) => {
          for (let i = entries.length - 1; i >= 0; i--) {
            if (entries[i].isIntersecting) {
              activeHeading.value = entries[i].target.id;
              return;
            }
          }
        },
        { rootMargin: "-80px 0px -70% 0px", threshold: 0 }
      );
      headings.forEach((h) => headingObserver.observe(h));
    };

    const enhanceContent = async () => {
      await nextTick();
      highlightCode();
      setupGallery();
      initComments();
      setupHeadingObserver();
    };

    const setupTocTracking = () => {
      if (!post.value?.toc?.length || post.value.toc.length < 3) return;

      const tocEl = document.querySelector('.article-toc');
      const shellEl = document.querySelector('.article-shell');
      if (!tocEl || !shellEl) return;

      const headerH = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--header-h')) || 64;

      const updateTocPosition = () => {
        const shellRect = shellEl.getBoundingClientRect();
        const tocWidth = 220;
        const gap = 32;
        const stickyTop = headerH + 24;
        const tocHeight = tocEl.offsetHeight;
        const shellBottom = shellRect.bottom;

        if (shellRect.top <= stickyTop && shellBottom > tocHeight) {
          tocEl.style.position = 'fixed';
          tocEl.style.top = stickyTop + 'px';
          tocEl.style.left = (shellRect.left - tocWidth - gap) + 'px';
          tocEl.style.bottom = 'auto';
          tocEl.style.width = tocWidth + 'px';
        } else if (shellBottom <= tocHeight) {
          tocEl.style.position = 'absolute';
          tocEl.style.top = 'auto';
          tocEl.style.bottom = '0';
          tocEl.style.left = '0';
          tocEl.style.width = tocWidth + 'px';
        } else {
          tocEl.style.position = '';
          tocEl.style.top = '';
          tocEl.style.left = '';
          tocEl.style.bottom = '';
          tocEl.style.width = '';
        }
      };

      tocScrollHandler = updateTocPosition;
      window.addEventListener('scroll', tocScrollHandler, { passive: true });
      updateTocPosition();
    };

    const syncTitle = () => {
      if (post.value && site.value.title) {
        document.title = `${post.value.title} - ${site.value.title}`;
      }
    };

    const scrollToHeading = (id) => {
      const el = document.getElementById(id);
      if (el) {
        const y = el.getBoundingClientRect().top + window.scrollY - 90;
        window.scrollTo({ top: y, behavior: "smooth" });
      }
    };

    const loadPost = async () => {
      loading.value = true;
      loadError.value = false;
      if (headingObserver) {
        headingObserver.disconnect();
        headingObserver = null;
      }
      if (tocScrollHandler) {
        window.removeEventListener('scroll', tocScrollHandler);
        tocScrollHandler = null;
      }
      activeHeading.value = "";
      try {
        post.value = await fetchPost(route.params.noteId);
        await enhanceContent();
        setupTocTracking();
        syncTitle();
      } catch {
        loadError.value = true;
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
        applyHighlightTheme();
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
      if (darkModeObserver) darkModeObserver.disconnect();
      if (headingObserver) headingObserver.disconnect();
      if (tocScrollHandler) window.removeEventListener('scroll', tocScrollHandler);
    });

    watch(() => route.params.noteId, loadPost);
    watch([post, site], syncTitle, { immediate: true });

    return {
      site,
      post,
      loading,
      loadError,
      activeHeading,
      artalkContainer,
      formatDate,
      loadPost,
      scrollToHeading,
    };
  },
};
</script>

<style>
@import "artalk/dist/Artalk.css";

/* ── Shell: dual-mode layout, no dead columns ── */
.article-shell {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 24px 80px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 32px;
}

.article-shell:not(.has-toc) {
  gap: 0;
}

/* ── TOC sidebar ── */
.article-toc {
  width: 220px;
  flex-shrink: 0;
  align-self: start;
  max-height: calc(100vh - var(--header-h) - 48px);
  overflow-y: auto;
  overflow-x: hidden;
  overscroll-behavior: contain;
}

.toc-panel {
  background: var(--surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius-md);
  padding: 14px 12px;
  box-shadow: var(--shadow-sm);
}

.toc-title {
  margin: 0 0 10px;
  padding: 0 8px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-faint);
}

.toc-link {
  display: block;
  padding: 7px 8px 7px 14px;
  color: var(--text-soft);
  font-size: 14px;
  line-height: 1.45;
  border-left: 2px solid transparent;
  border-radius: 6px;
  text-decoration: none;
  transition: color 140ms ease, background 140ms ease, border-color 140ms ease;
}

.toc-link:hover {
  color: var(--text);
  background: var(--surface-muted);
}

.toc-link.is-active {
  color: var(--text);
  background: var(--surface-muted);
  border-left-color: var(--accent);
  font-weight: 600;
}

.toc-link.toc-level-2 { padding-left: 24px; }
.toc-link.toc-level-3 { padding-left: 36px; }

/* ── Article main content ── */
.article-main {
  flex: 1;
  min-width: 0;
  max-width: var(--content-w);
}

.article-header {
  padding-bottom: 24px;
  margin-bottom: 32px;
  border-bottom: 1px solid var(--border-soft);
}

.article-title {
  margin: 0;
  color: var(--text);
  font-size: clamp(28px, 4vw, 44px);
  line-height: 1.18;
  font-weight: 650;
  letter-spacing: -0.01em;
  word-break: break-word;
}

.article-meta {
  margin-top: 14px;
  color: var(--text-faint);
  font-size: 15px;
}

.article-source {
  margin-top: 28px;
  padding-top: 16px;
  border-top: 1px solid var(--border-soft);
  color: var(--text-faint);
  font-size: 14px;
}

.article-source a {
  color: var(--link);
  text-decoration: none;
}

.article-source a:hover {
  text-decoration: underline;
}

/* ── Prose content ── */
.article-content {
  color: var(--text);
  font-size: 17px;
  line-height: 1.85;
  word-break: break-word;
}

.article-content > * + * {
  margin-top: 0;
}

.article-content p {
  margin: 0 0 1.2em;
}

.article-content h2 {
  margin: 2.4em 0 0.85em;
  font-size: 28px;
  line-height: 1.28;
  font-weight: 620;
  color: var(--text);
  scroll-margin-top: 90px;
}

.article-content h3 {
  margin: 2em 0 0.75em;
  font-size: 22px;
  line-height: 1.35;
  font-weight: 600;
  color: var(--text);
  scroll-margin-top: 90px;
}

.article-content ul,
.article-content ol {
  margin: 0 0 1.2em 1.4em;
  padding: 0;
}

.article-content li {
  margin: 0.4em 0;
}

.article-content strong {
  font-weight: 650;
}

.article-content a {
  color: var(--link);
  text-decoration: none;
}

.article-content a:hover {
  color: var(--link-hover);
  text-decoration: underline;
}

.article-content img {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-sm);
  display: block;
  margin: 8px 0;
}

.article-content blockquote {
  margin: 1.5em 0;
  padding: 12px 20px;
  border-left: 3px solid var(--accent);
  background: var(--surface-muted);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  color: var(--text-soft);
}

.article-content blockquote p:last-child {
  margin-bottom: 0;
}

/* ── Table ── */
.article-content table {
  width: 100%;
  border-collapse: collapse;
  margin: 24px 0 32px;
  font-size: 15px;
  line-height: 1.6;
  display: block;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.article-content th,
.article-content td {
  border-bottom: 1px solid var(--border-soft);
  padding: 12px 14px;
  text-align: left;
  vertical-align: top;
  white-space: nowrap;
}

.article-content th {
  color: var(--text);
  font-weight: 600;
  background: var(--surface-muted);
  position: sticky;
  top: 0;
}

.article-content td {
  white-space: normal;
}

/* ── Code ── */
.article-content pre {
  margin: 24px 0 32px;
  padding: 20px 24px;
  border-radius: var(--radius-md);
  background: #0f1722;
  color: #e6edf3;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.03);
  /* Let pre be as wide as its container allows, not more */
  box-sizing: border-box;
}

.article-content pre code {
  font-family: var(--mono);
  font-size: 14px;
  line-height: 1.65;
  background: none;
  padding: 0;
  border: none;
  border-radius: 0;
  color: inherit;
}

.article-content :not(pre) > code {
  font-family: var(--mono);
  font-size: 0.9em;
  padding: 0.18em 0.42em;
  border-radius: 6px;
  background: #f2f5f8;
  color: #203040;
  border: 1px solid #e1e7ee;
}

html.dark .article-content :not(pre) > code {
  background: #1f2836;
  color: #c0ccda;
  border-color: #2e3a4a;
}

.article-content hr {
  border: none;
  border-top: 1px solid var(--border-soft);
  margin: 2.5em 0;
}

/* ── Comments ── */
.article-comments {
  margin-top: 48px;
  padding-top: 32px;
  border-top: 1px solid var(--border-soft);
}

.article-comments h2 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text);
}

/* ── Responsive ── */
@media (max-width: 1024px) {
  .article-shell.has-toc {
    flex-direction: column;
    padding: 24px 16px 64px;
  }

  .article-shell:not(.has-toc) {
    padding: 24px 16px 64px;
  }

  .article-toc {
    display: none;
  }

  .article-main {
    max-width: 100%;
    flex: none;
  }
}

@media (max-width: 768px) {
  .article-title {
    font-size: 26px;
    line-height: 1.22;
  }

  .article-content h2 {
    font-size: 22px;
  }

  .article-content h3 {
    font-size: 18px;
  }
}

.load-error {
  text-align: center;
  color: var(--text-faint);
  padding: 40px 0;
  font-size: 0.95rem;
}

.load-error .el-button {
  margin-top: 16px;
}
</style>
