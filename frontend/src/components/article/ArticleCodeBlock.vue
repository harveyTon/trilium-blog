<template>
  <div
    ref="hostEl"
    class="article-code-block"
    :class="{
      'is-dark': isDark,
      'shows-line-numbers': showDesktopLineNumbers,
      'is-scrollable': isScrollable,
    }"
  >
    <div class="code-toolbar">
      <span class="code-language">{{ languageLabel || "Code" }}</span>
      <button
        type="button"
        class="code-copy-button"
        :aria-label="copyLabel"
        @click="copyCode"
      >
        {{ copyLabel }}
      </button>
    </div>

    <div ref="scroller" class="code-block-scroller" @scroll="updateScrollState">
      <div class="code-block-rendered" v-html="renderedHtml"></div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { highlightCodeBlock } from "../../composables/useCodeHighlighter";

const props = defineProps({
  code: {
    type: String,
    default: "",
  },
  languageId: {
    type: String,
    default: "plaintext",
  },
  languageLabel: {
    type: String,
    default: "Code",
  },
  showLineNumbers: {
    type: Boolean,
    default: false,
  },
});

const renderedHtml = ref('<pre class="shiki shiki-loading"><code></code></pre>');
const copyLabel = ref("Copy");
const isDark = ref(false);
const isScrollable = ref(false);
const scroller = ref(null);
const hostEl = ref(null);
let resetCopyTimer = null;
let resizeObserver = null;
let visibilityObserver = null;
let hasRendered = false;

const showDesktopLineNumbers = computed(() => props.showLineNumbers);

function syncTheme() {
  isDark.value = document.documentElement.classList.contains("dark");
}

async function renderCode() {
  const result = await highlightCodeBlock({
    code: props.code,
    language: props.languageId,
  });
  renderedHtml.value = result.html;
  await nextTick();
  updateScrollState();
}

function updateScrollState() {
  if (!scroller.value) return;
  isScrollable.value = scroller.value.scrollWidth > scroller.value.clientWidth + 4;
}

async function copyCode() {
  try {
    await navigator.clipboard.writeText(props.code || "");
    copyLabel.value = "Copied";
  } catch (error) {
    console.error("Failed to copy code block:", error);
    copyLabel.value = "Copy";
  }

  window.clearTimeout(resetCopyTimer);
  resetCopyTimer = window.setTimeout(() => {
    copyLabel.value = "Copy";
  }, 1400);
}

function ensureRendered() {
  if (hasRendered) return;
  hasRendered = true;
  renderCode();
}

onMounted(() => {
  syncTheme();

  const themeObserver = new MutationObserver(() => {
    syncTheme();
  });
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
  });

  onBeforeUnmount(() => {
    themeObserver.disconnect();
  });

  visibilityObserver = new IntersectionObserver(
    ([entry]) => {
      if (entry.isIntersecting) {
        ensureRendered();
        visibilityObserver?.disconnect();
        visibilityObserver = null;
      }
    },
    { rootMargin: "200px" }
  );
  if (hostEl.value) {
    visibilityObserver.observe(hostEl.value);
  }

  resizeObserver = new ResizeObserver(() => {
    updateScrollState();
  });
  if (scroller.value) {
    resizeObserver.observe(scroller.value);
  }
});

watch(
  () => [props.code, props.languageId],
  async () => {
    if (hasRendered) {
      await renderCode();
    }
  }
);

onBeforeUnmount(() => {
  window.clearTimeout(resetCopyTimer);
  visibilityObserver?.disconnect();
  resizeObserver?.disconnect();
});
</script>

<style scoped>
.article-code-block {
  --code-font-size: 15px;
  --code-line-gap-ratio: 0.8;
  --code-line-height: calc(var(--code-font-size) * (1 + var(--code-line-gap-ratio)));
  --code-line-hover: color-mix(in srgb, var(--surface-muted) 92%, var(--accent) 8%);
  --code-surface: color-mix(in srgb, var(--surface) 88%, white 12%);
  --code-surface-dark: color-mix(in srgb, var(--surface) 92%, black 8%);
  --code-border: color-mix(in srgb, var(--border-soft) 78%, transparent 22%);
  --code-toolbar-text: color-mix(in srgb, var(--text-soft) 90%, var(--text-faint) 10%);
  --code-shadow: 0 14px 28px rgba(15, 23, 42, 0.04);
  position: relative;
  margin: 28px 0 36px;
  border: 1px solid var(--code-border);
  border-radius: 10px;
  background: var(--code-surface);
  box-shadow: var(--code-shadow);
  overflow: hidden;
}

:global(.article-shell.reading-mode) .article-code-block {
  --code-font-size: 14px;
  --code-line-gap-ratio: 0.75;
  border-radius: 12px;
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.05);
}

:global(.article-shell.reading-mode) .article-code-block .code-toolbar {
  min-height: 32px;
  padding: 8px 14px 7px;
}

:global(.article-shell.reading-mode) .article-code-block .code-language,
:global(.article-shell.reading-mode) .article-code-block .code-copy-button {
  font-size: 11px;
  letter-spacing: 0.04em;
}

.article-code-block.is-dark {
  --code-line-hover: color-mix(in srgb, var(--surface-muted) 74%, transparent 26%);
  background: var(--code-surface-dark);
  box-shadow: 0 16px 32px rgba(0, 0, 0, 0.16);
}

.code-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 34px;
  padding: 8px 14px 6px;
  border-bottom: 1px solid color-mix(in srgb, var(--code-border) 88%, transparent 12%);
}

.code-language {
  color: var(--code-toolbar-text);
  font-size: 12px;
  line-height: 1;
  font-weight: 600;
}

.code-copy-button {
  border: none;
  background: transparent;
  color: var(--code-toolbar-text);
  cursor: pointer;
  font-size: 12px;
  line-height: 1;
  padding: 0;
  transition: color 160ms ease, opacity 160ms ease;
}

.code-copy-button:hover {
  color: var(--text);
}

.code-block-scroller {
  position: relative;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: thin;
  scrollbar-color: color-mix(in srgb, var(--text-faint) 25%, transparent 75%) transparent;
}

.code-block-scroller::-webkit-scrollbar {
  height: 9px;
}

.code-block-scroller::-webkit-scrollbar-track {
  background: transparent;
}

.code-block-scroller::-webkit-scrollbar-thumb {
  background: color-mix(in srgb, var(--text-faint) 22%, transparent 78%);
  border-radius: 999px;
}

.article-code-block.is-scrollable .code-block-scroller::after {
  content: "";
  position: absolute;
  inset: 0 0 0 auto;
  width: 18px;
  pointer-events: none;
  background: linear-gradient(to left, color-mix(in srgb, var(--code-surface) 100%, transparent 0%), transparent);
}

.article-code-block.is-dark.is-scrollable .code-block-scroller::after {
  background: linear-gradient(to left, color-mix(in srgb, var(--code-surface-dark) 100%, transparent 0%), transparent);
}

.code-block-rendered :deep(pre.shiki) {
  margin: 0;
  padding: 18px 0;
  background: transparent !important;
  min-width: 100%;
  width: max-content;
  font-family: var(--mono) !important;
  font-size: var(--code-font-size) !important;
  line-height: var(--code-line-height) !important;
}

.code-block-rendered :deep(.shiki),
.code-block-rendered :deep(.shiki span) {
  color: var(--shiki-light);
}

.article-code-block.is-dark .code-block-rendered :deep(.shiki),
.article-code-block.is-dark .code-block-rendered :deep(.shiki span) {
  color: var(--shiki-dark);
}

.code-block-rendered :deep(code) {
  display: block;
  font-family: var(--mono) !important;
  font-size: var(--code-font-size) !important;
  line-height: var(--code-line-height) !important;
}

.code-block-rendered :deep(.line) {
  display: inline-block;
  min-width: 100%;
  padding: 0 24px;
  border-radius: 6px;
  box-decoration-break: clone;
  -webkit-box-decoration-break: clone;
  font-family: var(--mono) !important;
  font-size: var(--code-font-size) !important;
  line-height: var(--code-line-height) !important;
  transition: background-color 160ms ease, color 160ms ease;
}

.code-block-rendered :deep(.line span) {
  font-family: inherit !important;
  font-size: inherit !important;
  line-height: inherit !important;
  white-space: pre;
}

.article-code-block.shows-line-numbers .code-block-rendered :deep(pre.shiki) {
  counter-reset: line;
}

.article-code-block.shows-line-numbers .code-block-rendered :deep(.line) {
  padding-left: 64px;
  position: relative;
}

.article-code-block.shows-line-numbers .code-block-rendered :deep(.line::before) {
  counter-increment: line;
  content: counter(line);
  position: absolute;
  left: 22px;
  width: 26px;
  top: 50%;
  transform: translateY(-50%);
  text-align: right;
  color: color-mix(in srgb, var(--text-faint) 72%, transparent 28%);
  font-size: 12px;
  line-height: 1;
}

@media (hover: hover) and (pointer: fine) {
  .code-block-rendered :deep(.line:hover) {
    background-color: var(--code-line-hover);
    box-shadow: inset 0 0 0 999px var(--code-line-hover);
  }

  .code-copy-button {
    opacity: 0.78;
  }

  .article-code-block:hover .code-copy-button {
    opacity: 1;
  }
}

@media (max-width: 768px) {
  .article-code-block {
    margin: 24px 0 32px;
  }

  .code-toolbar {
    padding-inline: 12px;
  }

  .code-block-rendered :deep(.line) {
    padding-inline: 18px;
  }

  .article-code-block.shows-line-numbers .code-block-rendered :deep(.line) {
    padding-left: 18px;
  }

  .article-code-block.shows-line-numbers .code-block-rendered :deep(.line::before) {
    content: none;
  }
}
</style>
