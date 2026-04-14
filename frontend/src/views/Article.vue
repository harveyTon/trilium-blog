<template>
  <div class="article">
    <ReadingProgressBar
      :progress="readingProgress"
      :top-offset="isReadingMode ? readingProgressOffset : null"
    />
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
        <div
          v-if="post"
          :class="[
            'article-shell',
            post.toc && post.toc.length >= 3 ? 'has-toc' : '',
            readingModeClass,
            `reading-width-${readingWidth}`,
            `reading-density-${readingDensity}`,
            `reading-font-size-${readingFontSize}`,
          ]"
        >
          <div
            v-if="isReadingMode"
            ref="readingTopbarRef"
            :class="['reading-topbar', { 'is-hidden': !readingTopbarVisible }]"
          >
            <div class="reading-topbar-inner">
              <button type="button" class="reading-topbar-link" @click="goHome">{{ t('article.home') }}</button>
              <div class="reading-topbar-actions">
                <button
                  v-if="post.toc && post.toc.length >= 3"
                  type="button"
                  :class="['reading-topbar-link', 'reading-topbar-state', { 'is-active': !readingTocCollapsed }]"
                  @click="toggleReadingTocCollapsed(!readingTocCollapsed)"
                >
                  <span>{{ t('article.toc') }}</span>
                  <span class="reading-topbar-state-text">{{ readingTocCollapsed ? t('article.tocCollapsed') : t('article.tocExpanded') }}</span>
                </button>
                <div ref="readingSettingsRef" class="reading-settings">
                  <button
                    type="button"
                    :class="['reading-topbar-link', { 'is-active': readingSettingsOpen }]"
                    @click="toggleReadingSettings"
                  >
                    Aa
                  </button>
                  <div v-if="readingSettingsOpen" class="reading-settings-popover">
                    <div class="reading-settings-group">
                      <div class="reading-settings-label">{{ t('article.width') }}</div>
                      <div class="reading-settings-options">
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingWidth === 'comfortable' }]"
                          @click="setReadingWidth('comfortable')"
                        >
                          {{ t('article.comfortable') }}
                        </button>
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingWidth === 'compact' }]"
                          @click="setReadingWidth('compact')"
                        >
                          {{ t('article.compact') }}
                        </button>
                      </div>
                    </div>
                    <div class="reading-settings-group">
                      <div class="reading-settings-label">{{ t('article.density') }}</div>
                      <div class="reading-settings-options">
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingDensity === 'relaxed' }]"
                          @click="setReadingDensity('relaxed')"
                        >
                          {{ t('article.comfortable') }}
                        </button>
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingDensity === 'comfortable' }]"
                          @click="setReadingDensity('comfortable')"
                        >
                          {{ t('article.compact') }}
                        </button>
                      </div>
                    </div>
                    <div class="reading-settings-group">
                      <div class="reading-settings-label">{{ t('article.fontSize') }}</div>
                      <div class="reading-settings-options is-triple">
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingFontSize === 'compact' }]"
                          @click="setReadingFontSize('compact')"
                        >
                          {{ t('article.small') }}
                        </button>
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingFontSize === 'comfortable' }]"
                          @click="setReadingFontSize('comfortable')"
                        >
                          {{ t('article.medium') }}
                        </button>
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': readingFontSize === 'large' }]"
                          @click="setReadingFontSize('large')"
                        >
                          {{ t('article.large') }}
                        </button>
                      </div>
                    </div>
                    <div class="reading-settings-group">
                      <div class="reading-settings-label">{{ t('article.readingTheme') }}</div>
                      <div class="reading-settings-options">
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': !isDarkTheme }]"
                          @click="setThemeMode(false)"
                        >
                          {{ t('article.lightTheme') }}
                        </button>
                        <button
                          type="button"
                          :class="['reading-settings-option', { 'is-active': isDarkTheme }]"
                          @click="setThemeMode(true)"
                        >
                          {{ t('article.darkTheme') }}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
                <button type="button" class="reading-topbar-link is-exit" @click="exitReadingMode">
                  {{ t('article.exitReading') }}
                </button>
              </div>
            </div>
          </div>
          <ArticleTOC
            :items="post.toc"
            :active-heading="activeHeading"
            :collapsed="currentTocCollapsed"
            :reading-mode="isReadingMode"
            @scroll-to-heading="scrollToHeading"
            @toggle-collapse="isReadingMode ? toggleReadingTocCollapsed() : toggleStandardTocCollapsed()"
          />

          <main class="article-main">
            <ArticleHeader :title="post.title" :formatted-date="formatDate(post.dateModified)">
              <template v-if="!isReadingMode" #actions>
                <button type="button" class="reading-mode-trigger" @click="enterReadingMode">
                  {{ t('article.enterReading') }}
                </button>
              </template>
            </ArticleHeader>
            <ArticleSummaryBlock
              v-if="!isReadingMode"
              :summary="summaryState.ai"
              :enabled="summaryState.aiEnabled"
            />
            <ArticleContent
              ref="articleContentRef"
              :content-html="post.contentHtml"
              :reading-mode="isReadingMode"
            />
            <SourceLinkBlock v-if="!isReadingMode" :page-url="post.pageUrl" />
          </main>
        </div>

        <el-empty v-else-if="!loadError" :description="t('article.notFound')"></el-empty>
        <div v-else class="load-error">
          <p>{{ t('article.fetchError') }}</p>
          <el-button type="primary" @click="loadPost">{{ t('article.retry') }}</el-button>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<script>
import { ElButton } from "element-plus";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import { storeToRefs } from "pinia";
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useDark } from "@vueuse/core";
import { fetchPost } from "../api/blog";
import { fetchPostSummary, normalizeSummaryPayload } from "../api/summary";
import ReadingProgressBar from "../components/app/ReadingProgressBar.vue";
import ArticleContent from "../components/article/ArticleContent.vue";
import ArticleHeader from "../components/article/ArticleHeader.vue";
import ArticleSummaryBlock from "../components/article/ArticleSummaryBlock.vue";
import ArticleTOC from "../components/article/ArticleTOC.vue";
import SourceLinkBlock from "../components/article/SourceLinkBlock.vue";
import { useArticleEnhancements } from "../composables/useArticleEnhancements";
import { useArticleReadingMode } from "../composables/useArticleReadingMode";
import { useReadingProgress } from "../composables/useReadingProgress";
import { t } from "../i18n";
import { useSiteStore } from "../store";

export default {
  name: "ArticlePage",
  components: {
    ElButton,
    ArticleContent,
    ArticleHeader,
    ArticleSummaryBlock,
    ArticleTOC,
    ReadingProgressBar,
    SourceLinkBlock,
  },
  setup() {
    const route = useRoute();
    const router = useRouter();
    const siteStore = useSiteStore();
    const { site } = storeToRefs(siteStore);
    const post = ref(null);
    const summarySource = ref(null);
    const loading = ref(true);
    const loadError = ref(false);
    const activeHeading = ref("");
    const articleContentRef = ref(null);
    const standardTocCollapsed = ref(window.innerWidth <= 1024);
    const readingTopbarRef = ref(null);
    const readingSettingsRef = ref(null);
    const readingSettingsOpen = ref(false);
    const readingTopbarVisible = ref(true);
    const readingTopbarHeight = ref(44);
    const summaryState = computed(() => normalizeSummaryPayload(summarySource.value || post.value));
    const readingProgressOffset = computed(() =>
      isReadingMode.value ? (readingTopbarVisible.value ? readingTopbarHeight.value : 0) : null
    );
    const { progress: readingProgress, updateReadingProgress } = useReadingProgress();
    const {
      enabled: isReadingMode,
      readingTocCollapsed,
      width: readingWidth,
      density: readingDensity,
      fontSize: readingFontSize,
      readingModeClass,
      enterReadingMode,
      exitReadingMode,
      toggleReadingTocCollapsed,
      setWidth,
      setDensity,
      setFontSize,
    } = useArticleReadingMode();
    const currentTocCollapsed = computed(() =>
      isReadingMode.value ? readingTocCollapsed.value : standardTocCollapsed.value
    );
    let headingObserver = null;
    let summaryPollTimer = null;
    let lastScrollY = 0;

    const syncReadingTopbarHeight = () => {
      const nextHeight = readingTopbarRef.value?.getBoundingClientRect?.().height;
      if (nextHeight) {
        readingTopbarHeight.value = Math.round(nextHeight);
      }
    };

    const isDark = useDark();
    const isDarkTheme = computed(() => isDark.value);

    const toggleTheme = () => {
      isDark.value = !isDark.value;
    };

    const setThemeMode = (dark) => {
      isDark.value = dark;
    };

    const formatDate = (dateString) => {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    };

    const { enhanceArticleContent, cleanupEnhancements } = useArticleEnhancements();

    const setupHeadingObserver = () => {
      if (headingObserver) {
        headingObserver.disconnect();
        headingObserver = null;
      }
      const root = articleContentRef.value?.getRootElement?.();
      if (!root) return;

      const headings = root.querySelectorAll("h1, h2, h3");
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
      const root = articleContentRef.value?.getRootElement?.();
      if (!root || !post.value) return;
      await enhanceArticleContent({
        root,
        codeBlocks: post.value.codeBlocks || [],
      });
      setupHeadingObserver();
      updateReadingProgress();
    };

    const applyReadingModeDocumentState = (enabled) => {
      document.body.classList.toggle("article-reading-mode", enabled);
    };

    const goHome = async () => {
      await router.push({ path: "/" });
      window.scrollTo(0, 0);
    };

    const toggleReadingSettings = () => {
      readingSettingsOpen.value = !readingSettingsOpen.value;
    };

    const closeReadingSettings = () => {
      readingSettingsOpen.value = false;
    };

    const setReadingWidth = (value) => {
      setWidth(value);
    };

    const setReadingDensity = (value) => {
      setDensity(value);
    };

    const setReadingFontSize = (value) => {
      setFontSize(value);
    };

    const toggleStandardTocCollapsed = (forceValue) => {
      standardTocCollapsed.value =
        typeof forceValue === "boolean" ? forceValue : !standardTocCollapsed.value;
    };

    const handlePointerDown = (event) => {
      if (!readingSettingsOpen.value) {
        return;
      }
      if (readingSettingsRef.value?.contains(event.target)) {
        return;
      }
      closeReadingSettings();
    };

    const handleReadingScroll = () => {
      if (!isReadingMode.value) {
        readingTopbarVisible.value = true;
        lastScrollY = window.scrollY;
        return;
      }

      const currentScrollY = window.scrollY;
      const delta = currentScrollY - lastScrollY;

      if (currentScrollY <= 24) {
        readingTopbarVisible.value = true;
      } else if (delta > 10) {
        readingTopbarVisible.value = false;
        closeReadingSettings();
      } else if (delta < -8) {
        readingTopbarVisible.value = true;
      }

      lastScrollY = currentScrollY;
    };

    const syncTitle = () => {
      if (post.value && site.value.title) {
        const siteTitle = [site.value.title, site.value.subtitle].filter(Boolean).join(" | ");
        document.title = siteTitle ? `${post.value.title} - ${siteTitle}` : post.value.title;
      }
    };

    const stopSummaryPolling = () => {
      if (summaryPollTimer) {
        window.clearTimeout(summaryPollTimer);
        summaryPollTimer = null;
      }
    };

    const hasReadySummary = (summaries) => {
      if (!summaries) {
        return false;
      }
      if (summaries.aiEnabled && summaries.ai?.status === "ready" && summaries.ai?.text) {
        return true;
      }
      if (summaries.code?.status === "ready" && summaries.code?.text) {
        return true;
      }
      return false;
    };

    const shouldFetchSummaryImmediately = (source) => {
      const normalized = normalizeSummaryPayload(source);
      if (!normalized.aiEnabled) {
        return false;
      }
      if (!source?.summaries) {
        return true;
      }
      if (normalized.ai?.status === "ready" && normalized.ai?.text) {
        return false;
      }
      return !hasReadySummary(source.summaries) || ["pending", "processing"].includes(normalized.ai?.status);
    };

    const shouldPollSummary = () => {
      if (!summaryState.value.aiEnabled) {
        return false;
      }
      if (!summaryState.value.ai) {
        return true;
      }
      return ["pending", "processing"].includes(summaryState.value.ai.status);
    };

    const syncSummary = (summaries) => {
      if (!post.value) return;
      const nextSummaryText = preferredSummaryText(summaries, post.value.summary);
      post.value = {
        ...post.value,
        summaries,
        summary: nextSummaryText,
      };
      summarySource.value = {
        ...post.value,
        summaries,
        summary: nextSummaryText,
      };
    };

    const refreshSummaryStatus = async (noteId) => {
      const summaries = await fetchPostSummary(noteId);
      if (route.params.noteId !== noteId || !post.value) {
        return null;
      }
      syncSummary(summaries);
      return summaries;
    };

    const pollSummaryStatus = (noteId) => {
      stopSummaryPolling();
      if (!shouldPollSummary()) {
        return;
      }

      summaryPollTimer = window.setTimeout(async () => {
        try {
          await refreshSummaryStatus(noteId);
        } catch (error) {
          console.error("Failed to poll post summary:", error);
        } finally {
          if (route.params.noteId === noteId && shouldPollSummary()) {
            pollSummaryStatus(noteId);
          }
        }
      }, 2500);
    };

    const scrollToHeading = (id) => {
      const el = document.getElementById(id);
      if (el) {
        const y = el.getBoundingClientRect().top + window.scrollY - 90;
        window.scrollTo({ top: y, behavior: "smooth" });
        if (isReadingMode.value) {
          toggleReadingTocCollapsed(true);
        }
      }
    };

    const loadPost = async () => {
      loading.value = true;
      loadError.value = false;
      if (headingObserver) {
        headingObserver.disconnect();
        headingObserver = null;
      }
      activeHeading.value = "";
      if (!isReadingMode.value) {
        standardTocCollapsed.value = window.innerWidth <= 1024;
      }
      stopSummaryPolling();
      cleanupEnhancements();
      try {
        const fetchedPost = await fetchPost(route.params.noteId);
        post.value = fetchedPost;
        summarySource.value = fetchedPost;
        if (shouldFetchSummaryImmediately(fetchedPost)) {
          try {
            await refreshSummaryStatus(route.params.noteId);
          } catch (error) {
            console.error("Failed to fetch initial post summary:", error);
          }
        }
        loading.value = false;
        await enhanceContent();
        syncTitle();
        pollSummaryStatus(route.params.noteId);
      } catch {
        loadError.value = true;
        post.value = null;
        summarySource.value = null;
      } finally {
        loading.value = false;
        if (typeof window.scrollTo === "function") {
          window.scrollTo({ top: 0, behavior: "smooth" });
        }
      }
    };

    onMounted(async () => {
      await loadPost();
      document.addEventListener("pointerdown", handlePointerDown);
      window.addEventListener("scroll", handleReadingScroll, { passive: true });
      window.addEventListener("resize", syncReadingTopbarHeight, { passive: true });
      lastScrollY = window.scrollY;
      await nextTick();
      syncReadingTopbarHeight();
    });

    onUnmounted(() => {
      stopSummaryPolling();
      cleanupEnhancements();
      applyReadingModeDocumentState(false);
      if (headingObserver) headingObserver.disconnect();
      document.removeEventListener("pointerdown", handlePointerDown);
      window.removeEventListener("scroll", handleReadingScroll);
      window.removeEventListener("resize", syncReadingTopbarHeight);
    });

    watch(() => route.params.noteId, loadPost);
    watch([post, site], syncTitle, { immediate: true });
    watch(isReadingMode, async (enabled) => {
      applyReadingModeDocumentState(enabled);
      readingTopbarVisible.value = true;
      if (!enabled) {
        closeReadingSettings();
      }
      lastScrollY = window.scrollY;
      await nextTick();
      syncReadingTopbarHeight();
      updateReadingProgress();
    }, { immediate: true });

    return {
      t,
      site,
      post,
      summaryState,
      loading,
      loadError,
      activeHeading,
      currentTocCollapsed,
      readingTocCollapsed,
      standardTocCollapsed,
      isReadingMode,
      readingWidth,
      readingDensity,
      readingFontSize,
      readingModeClass,
      readingProgress,
      articleContentRef,
      readingTopbarRef,
      readingSettingsRef,
      readingSettingsOpen,
      readingTopbarVisible,
      readingProgressOffset,
      isDarkTheme,
      formatDate,
      loadPost,
      scrollToHeading,
      enterReadingMode,
      exitReadingMode,
      toggleStandardTocCollapsed,
      toggleReadingTocCollapsed,
      toggleReadingSettings,
      setReadingWidth,
      setReadingDensity,
      setReadingFontSize,
      setThemeMode,
      goHome,
    };
  },
};

function preferredSummaryText(summaries, fallback) {
  if (summaries?.aiEnabled && summaries?.ai?.status === "ready" && summaries.ai?.text) {
    return summaries.ai.text;
  }
  if (summaries?.code?.status === "ready" && summaries.code?.text) {
    return summaries.code.text;
  }
  return fallback || "";
}
</script>

<style>
/* ── Shell: dual-mode layout, no dead columns ── */
.article-shell {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 24px 80px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 32px;
  overflow-x: clip;
}

body.article-reading-mode .app-header,
body.article-reading-mode .app-footer,
body.article-reading-mode .el-backtop {
  display: none !important;
}

body.article-reading-mode .app-main {
  padding-top: 0;
}

body.article-reading-mode {
  background: color-mix(in srgb, var(--bg) 94%, white 6%);
}

html.dark body.article-reading-mode {
  background: color-mix(in srgb, var(--bg) 96%, black 4%);
}

.article-shell:not(.has-toc) {
  gap: 0;
}

.article-shell.reading-mode {
  --article-reading-width: 760px;
  --article-reading-font-size: 18px;
  --article-reading-line-height: 1.92;
  --article-reading-block-gap: 1.28em;
  --article-reading-mobile-padding: 16px;
  position: relative;
  max-width: 100%;
  min-height: 100vh;
  padding: 76px 32px 96px;
  gap: 24px;
}

.article-shell.reading-mode.reading-width-compact {
  --article-reading-width: 700px;
}

.article-shell.reading-mode.reading-font-size-compact {
  --article-reading-font-size: 16px;
  --article-reading-line-height: 1.82;
}

.article-shell.reading-mode.reading-font-size-large {
  --article-reading-font-size: 20px;
  --article-reading-line-height: 2;
}

.article-shell.reading-mode.reading-density-comfortable {
  --article-reading-line-height: 1.82;
  --article-reading-block-gap: 1.18em;
}

.reading-topbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1200;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  min-height: 44px;
  padding: 0 20px;
  border-bottom: 1px solid color-mix(in srgb, var(--border-soft) 88%, transparent 12%);
  background: color-mix(in srgb, var(--bg) 82%, transparent 18%);
  backdrop-filter: blur(16px);
  transition: transform 220ms ease, opacity 220ms ease;
}

.reading-topbar.is-hidden {
  opacity: 0;
  transform: translateY(calc(-100% - 6px));
  pointer-events: none;
}

.reading-topbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  width: 100%;
}

.reading-topbar-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.reading-topbar-link,
.reading-mode-trigger {
  border: 1px solid color-mix(in srgb, var(--border-soft) 88%, transparent 12%);
  background: color-mix(in srgb, var(--surface) 92%, transparent 8%);
  color: var(--text-soft);
  border-radius: 999px;
  min-height: 32px;
  padding: 0 12px;
  font: inherit;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.04em;
  cursor: pointer;
  white-space: nowrap;
  transition: transform 160ms ease, border-color 160ms ease, color 160ms ease, background 160ms ease;
}

@media (hover: hover) and (pointer: fine) {
  .reading-topbar-link:hover,
  .reading-mode-trigger:hover {
    color: var(--text);
    border-color: color-mix(in srgb, var(--accent) 28%, var(--border-soft) 72%);
    transform: translateY(-1px);
  }

  .reading-settings-option:hover {
    color: var(--text);
    transform: translateY(-1px);
  }
}

.reading-topbar-link.is-exit,
.reading-mode-trigger {
  color: var(--text);
}

.reading-topbar-link.is-active {
  color: var(--text);
  border-color: color-mix(in srgb, var(--accent) 28%, var(--border-soft) 72%);
  background: color-mix(in srgb, var(--surface) 86%, var(--accent) 14%);
}

.reading-topbar-state {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.reading-topbar-state-text {
  color: var(--text-faint);
  font-size: 11px;
  letter-spacing: 0;
  font-weight: 600;
}

.reading-settings {
  position: relative;
}

.reading-settings-popover {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 196px;
  padding: 14px;
  border: 1px solid color-mix(in srgb, var(--border-soft) 90%, transparent 10%);
  border-radius: 16px;
  background: color-mix(in srgb, var(--surface) 94%, transparent 6%);
  box-shadow: 0 18px 44px rgba(15, 23, 42, 0.14);
  backdrop-filter: blur(18px);
}

.reading-settings-group + .reading-settings-group {
  margin-top: 12px;
}

.reading-settings-label {
  margin-bottom: 8px;
  color: var(--text-faint);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.reading-settings-options {
  display: flex;
  gap: 8px;
}

.reading-settings-options.is-triple .reading-settings-option {
  min-width: 0;
}

.reading-settings-option {
  flex: 1;
  min-height: 34px;
  border: 1px solid color-mix(in srgb, var(--border-soft) 88%, transparent 12%);
  border-radius: 12px;
  background: color-mix(in srgb, var(--surface-muted) 90%, transparent 10%);
  color: var(--text-soft);
  font: inherit;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: border-color 160ms ease, color 160ms ease, background 160ms ease, transform 160ms ease;
}

.reading-settings-option.is-active {
  color: var(--text);
  border-color: color-mix(in srgb, var(--accent) 28%, var(--border-soft) 72%);
  background: color-mix(in srgb, var(--surface) 84%, var(--accent) 16%);
}

/* ── TOC sidebar ── */
.article-toc-wrapper {
  position: sticky;
  top: calc(var(--header-h) + 24px);
  align-self: start;
  flex-shrink: 0;
  width: 220px;
}

.article-toc {
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

.article-shell.reading-mode .article-main {
  max-width: var(--article-reading-width);
  margin: 0 auto;
}

.article-header {
  padding-bottom: 24px;
  margin-bottom: 32px;
  border-bottom: 1px solid var(--border-soft);
}

.article-shell.reading-mode .article-header {
  padding-bottom: 28px;
  margin-bottom: 36px;
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
  max-width: 100%;
}

.article-source a {
  color: var(--link);
  text-decoration: none;
  word-break: break-word;
  overflow-wrap: anywhere;
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
  overflow-x: clip;
}

.article-content.is-reading-mode {
  font-size: var(--article-reading-font-size);
  line-height: var(--article-reading-line-height);
}

.article-content > * + * {
  margin-top: 0;
}

.article-content p {
  margin: 0 0 1.2em;
}

.article-content.is-reading-mode p {
  margin: 0 0 var(--article-reading-block-gap);
}

.article-content h2 {
  margin: 2.4em 0 0.85em;
  font-size: 28px;
  line-height: 1.28;
  font-weight: 620;
  color: var(--text);
  scroll-margin-top: 90px;
}

.article-content.is-reading-mode h2 {
  margin-top: 2.65em;
  margin-bottom: 0.95em;
}

.article-content h3 {
  margin: 2em 0 0.75em;
  font-size: 22px;
  line-height: 1.35;
  font-weight: 600;
  color: var(--text);
  scroll-margin-top: 90px;
}

.article-content.is-reading-mode h3 {
  margin-top: 2.15em;
  margin-bottom: 0.82em;
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

.article-content.is-reading-mode img {
  margin: 0;
}

.article-content .image-gallery-group {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  align-items: start;
}

.article-content .image-gallery-group img {
  width: 100%;
  margin: 0;
}

.article-content.is-reading-mode figure,
.article-shell.reading-mode .table-scroll-wrapper,
.article-content.is-reading-mode .image-gallery-group {
  margin: 32px 0 40px;
}

.article-content.is-reading-mode figure img,
.article-content.is-reading-mode > p > a[data-fancybox] > img,
.article-content.is-reading-mode > p > img {
  width: auto;
  max-width: min(100%, 1000px);
  margin-inline: auto;
}

.article-content blockquote {
  margin: 1.6em 0;
  padding: 12px 16px;
  border-left: none;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 6px;
  color: var(--text-soft);
}

.article-content blockquote p {
  margin: 0;
  line-height: 1.7;
}

html.dark .article-content blockquote {
  background: rgba(255, 255, 255, 0.04);
}

.article-content figure,
.table-scroll-wrapper {
  width: 100%;
  max-width: 100%;
  margin: 24px 0 32px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

/* ── Table ── */

.article-content table {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
  margin: 0;
  font-size: 15px;
  line-height: 1.6;
  table-layout: auto;
}

.article-content.is-reading-mode table {
  font-size: 14px;
}

.article-content th,
.article-content td {
  border-bottom: 1px solid var(--border-soft);
  padding: 12px 14px;
  text-align: left;
  vertical-align: top;
  white-space: nowrap;
  word-break: normal;
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
.article-code-block-host {
  display: block;
}

.article-content pre:not(.shiki) {
  margin: 24px 0 32px;
  padding: 18px 22px;
  border-radius: 10px;
  background: color-mix(in srgb, var(--surface) 90%, white 10%);
  color: var(--text);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border: 1px solid var(--border-soft);
  box-sizing: border-box;
}

.article-content pre:not(.shiki) code {
  font-family: var(--mono);
  font-size: 15px;
  line-height: 1.8;
  background: none;
  padding: 0;
  border: none;
  border-radius: 0;
  color: inherit;
}

.article-content :not(pre) > code {
  font-family: var(--mono);
  font-size: 0.88em;
  padding: 0.16em 0.4em;
  border-radius: 6px;
  background: color-mix(in srgb, var(--surface-muted) 82%, white 18%);
  color: color-mix(in srgb, var(--text) 86%, var(--text-soft) 14%);
  border: 1px solid color-mix(in srgb, var(--border-soft) 80%, transparent 20%);
}

html.dark .article-content :not(pre) > code {
  background: color-mix(in srgb, var(--surface-muted) 68%, black 32%);
  color: color-mix(in srgb, var(--text) 88%, white 12%);
  border-color: color-mix(in srgb, var(--border-soft) 72%, transparent 28%);
}

.article-content hr {
  border: none;
  border-top: 1px solid var(--border-soft);
  margin: 2.5em 0;
}

.article-shell.reading-mode .article-source,
.article-shell.reading-mode .article-summary-block {
  display: none;
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

  .article-toc-wrapper {
    display: none;
  }

  .article-shell.reading-mode .article-toc-wrapper.is-reading-mode {
    display: block;
  }

  .article-main {
    max-width: 100%;
    min-width: 0;
    overflow-x: clip;
  }
}

@media (max-width: 768px) {
  .reading-topbar {
    align-items: center;
    padding: 6px 12px;
  }

  .reading-topbar-inner {
    gap: 8px;
  }

  .reading-topbar-actions {
    gap: 6px;
    overflow-x: auto;
    overflow-y: visible;
    scrollbar-width: none;
  }

  .reading-topbar-link {
    min-height: 30px;
    padding-inline: 10px;
    font-size: 11px;
  }

  .reading-topbar-state-text {
    display: none;
  }

  .reading-settings-popover {
    position: fixed;
    top: 52px;
    right: 12px;
    width: min(220px, calc(100vw - 24px));
    z-index: 1250;
  }

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

  .article-shell.reading-mode {
    padding: 68px var(--article-reading-mobile-padding) 72px;
  }

  .article-shell.reading-mode.reading-width-compact {
    --article-reading-mobile-padding: 10px;
  }

  .article-shell.reading-mode .article-main {
    max-width: 100%;
  }

  .article-content figure,
  .table-scroll-wrapper {
    margin-inline: -4px;
    padding-inline: 4px;
  }

  .article-content table {
    min-width: 560px;
    font-size: 14px;
  }

  .article-content th,
  .article-content td {
    min-width: 110px;
    padding: 10px 12px;
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
