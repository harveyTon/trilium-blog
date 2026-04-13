import { computed, onMounted, onUnmounted, ref, watch } from "vue";

const STORAGE_KEY = "trilium-blog-reading-mode";
const WIDTH_OPTIONS = new Set(["compact", "comfortable"]);
const DENSITY_OPTIONS = new Set(["comfortable", "relaxed"]);

const DEFAULT_PREFERENCES = {
  enabled: false,
  tocCollapsed: true,
  width: "comfortable",
  density: "relaxed",
};

function sanitizePreferences(candidate = {}) {
  return {
    enabled: Boolean(candidate.enabled),
    tocCollapsed:
      typeof candidate.tocCollapsed === "boolean"
        ? candidate.tocCollapsed
        : DEFAULT_PREFERENCES.tocCollapsed,
    width: WIDTH_OPTIONS.has(candidate.width)
      ? candidate.width
      : DEFAULT_PREFERENCES.width,
    density: DENSITY_OPTIONS.has(candidate.density)
      ? candidate.density
      : DEFAULT_PREFERENCES.density,
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
  const tocCollapsed = ref(stored.tocCollapsed);
  const width = ref(stored.width);
  const density = ref(stored.density);
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
        tocCollapsed: tocCollapsed.value,
        width: width.value,
        density: density.value,
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

  function toggleTocCollapsed(forceValue) {
    tocCollapsed.value =
      typeof forceValue === "boolean" ? forceValue : !tocCollapsed.value;
  }

  function cycleWidth() {
    width.value = width.value === "comfortable" ? "compact" : "comfortable";
  }

  function cycleDensity() {
    density.value = density.value === "relaxed" ? "comfortable" : "relaxed";
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

  watch([enabled, tocCollapsed, width, density], persist);

  return {
    enabled,
    tocCollapsed,
    width,
    density,
    readingModeClass,
    enterReadingMode,
    exitReadingMode,
    toggleReadingMode,
    toggleTocCollapsed,
    cycleWidth,
    cycleDensity,
  };
}
