import { computed, onMounted, onUnmounted, ref, watch } from "vue";

const STORAGE_KEY = "trilium-blog-reading-mode";
const WIDTH_OPTIONS = new Set(["compact", "comfortable"]);
const DENSITY_OPTIONS = new Set(["comfortable", "relaxed"]);
const FONT_SIZE_OPTIONS = new Set(["compact", "comfortable", "large"]);

const DEFAULT_PREFERENCES = {
  enabled: false,
  readingTocCollapsed: true,
  width: "comfortable",
  density: "relaxed",
  fontSize: "comfortable",
};

function sanitizePreferences(candidate = {}) {
  return {
    enabled: Boolean(candidate.enabled),
    readingTocCollapsed:
      typeof candidate.readingTocCollapsed === "boolean"
        ? candidate.readingTocCollapsed
        : typeof candidate.tocCollapsed === "boolean"
          ? candidate.tocCollapsed
          : DEFAULT_PREFERENCES.readingTocCollapsed,
    width: WIDTH_OPTIONS.has(candidate.width)
      ? candidate.width
      : DEFAULT_PREFERENCES.width,
    density: DENSITY_OPTIONS.has(candidate.density)
      ? candidate.density
      : DEFAULT_PREFERENCES.density,
    fontSize: FONT_SIZE_OPTIONS.has(candidate.fontSize)
      ? candidate.fontSize
      : DEFAULT_PREFERENCES.fontSize,
  };
}

function readStoredPreferences() {
  if (typeof window === "undefined") {
    return DEFAULT_PREFERENCES;
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return DEFAULT_PREFERENCES;
    }
    return sanitizePreferences(JSON.parse(raw));
  } catch {
    return DEFAULT_PREFERENCES;
  }
}

export function useArticleReadingMode() {
  const stored = readStoredPreferences();
  const enabled = ref(stored.enabled);
  const readingTocCollapsed = ref(stored.readingTocCollapsed);
  const width = ref(stored.width);
  const density = ref(stored.density);
  const fontSize = ref(stored.fontSize);
  const hydrated = ref(false);

  const readingModeClass = computed(() => (enabled.value ? "reading-mode" : ""));

  function persist() {
    if (!hydrated.value || typeof window === "undefined") {
      return;
    }

    window.localStorage.setItem(
      STORAGE_KEY,
        JSON.stringify({
          enabled: enabled.value,
          readingTocCollapsed: readingTocCollapsed.value,
          width: width.value,
          density: density.value,
          fontSize: fontSize.value,
      })
    );
  }

  function enterReadingMode() {
    enabled.value = true;
  }

  function exitReadingMode() {
    enabled.value = false;
  }

  function toggleReadingMode() {
    enabled.value = !enabled.value;
  }

  function toggleReadingTocCollapsed(forceValue) {
    readingTocCollapsed.value =
      typeof forceValue === "boolean" ? forceValue : !readingTocCollapsed.value;
  }

  function cycleWidth() {
    width.value = width.value === "comfortable" ? "compact" : "comfortable";
  }

  function cycleDensity() {
    density.value = density.value === "relaxed" ? "comfortable" : "relaxed";
  }

  function setWidth(nextWidth) {
    if (WIDTH_OPTIONS.has(nextWidth)) {
      width.value = nextWidth;
    }
  }

  function setDensity(nextDensity) {
    if (DENSITY_OPTIONS.has(nextDensity)) {
      density.value = nextDensity;
    }
  }

  function setFontSize(nextFontSize) {
    if (FONT_SIZE_OPTIONS.has(nextFontSize)) {
      fontSize.value = nextFontSize;
    }
  }

  function onKeydown(event) {
    if (event.key === "Escape" && enabled.value) {
      exitReadingMode();
    }
  }

  onMounted(() => {
    hydrated.value = true;
    window.addEventListener("keydown", onKeydown);
  });

  onUnmounted(() => {
    window.removeEventListener("keydown", onKeydown);
  });

  watch([enabled, readingTocCollapsed, width, density, fontSize], persist);

  return {
    enabled,
    readingTocCollapsed,
    width,
    density,
    fontSize,
    readingModeClass,
    enterReadingMode,
    exitReadingMode,
    toggleReadingMode,
    toggleReadingTocCollapsed,
    cycleWidth,
    cycleDensity,
    setWidth,
    setDensity,
    setFontSize,
  };
}
