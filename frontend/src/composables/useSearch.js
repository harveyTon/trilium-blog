import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { fetchSearchResults } from "../api/search";
import { t } from "../i18n";

export function useSearch() {
  const route = useRoute();
  const loading = ref(false);
  const error = ref("");
  const result = ref(null);

  const query = computed(() => (route.query.q || "").toString().trim());
  const hasQuery = computed(() => query.value.length > 0);
  const items = computed(() => result.value?.items || []);
  const total = computed(() => result.value?.total || 0);

  const load = async () => {
    if (!hasQuery.value) {
      result.value = null;
      error.value = "";
      loading.value = false;
      return;
    }

    loading.value = true;
    error.value = "";
    try {
      result.value = await fetchSearchResults(query.value);
    } catch (err) {
      console.error("Fetch Search Results Error:", err);
      error.value = t('search.fetchError');
      result.value = null;
    } finally {
      loading.value = false;
    }
  };

  watch(query, () => {
    load();
  }, { immediate: true });

  return {
    query,
    hasQuery,
    loading,
    error,
    result,
    items,
    total,
    load,
  };
}
