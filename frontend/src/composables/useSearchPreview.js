import { ref } from "vue";
import { fetchSearchPreview } from "../api/search";
import { useSearchStore } from "../store/search";
import { t } from "../i18n";

export function useSearchPreview() {
  const searchStore = useSearchStore();
  const debounceTimer = ref(null);
  let requestToken = 0;

  const closePreview = () => {
    searchStore.resetPreview();
  };

  const openPreview = () => {
    searchStore.setPreviewOpen(true);
  };

  const loadPreview = (query) => {
    searchStore.setQuery(query);

    if (debounceTimer.value) {
      window.clearTimeout(debounceTimer.value);
      debounceTimer.value = null;
    }

    if (!query.trim()) {
      closePreview();
      return;
    }

    openPreview();
    debounceTimer.value = window.setTimeout(async () => {
      const currentToken = ++requestToken;
      searchStore.setLoading(true);
      searchStore.setError("");
      try {
        const data = await fetchSearchPreview(query);
        if (currentToken !== requestToken) return;
        searchStore.setPreviewItems(data.items || []);
        searchStore.setPreviewOpen(true);
      } catch (err) {
        if (currentToken !== requestToken) return;
        console.error("Fetch Search Preview Error:", err);
        searchStore.setError(t('search.previewError'));
        searchStore.setPreviewItems([]);
      } finally {
        if (currentToken === requestToken) {
          searchStore.setLoading(false);
        }
      }
    }, 180);
  };

  return {
    searchStore,
    loadPreview,
    closePreview,
    openPreview,
  };
}
