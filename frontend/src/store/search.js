import { defineStore } from "pinia";

export const useSearchStore = defineStore("search", {
  state: () => ({
    query: "",
    previewOpen: false,
    previewItems: [],
    loading: false,
    error: "",
  }),
  actions: {
    setQuery(query) {
      this.query = query;
    },
    setPreviewOpen(open) {
      this.previewOpen = open;
    },
    setPreviewItems(items) {
      this.previewItems = items;
    },
    setLoading(loading) {
      this.loading = loading;
    },
    setError(error) {
      this.error = error;
    },
    resetPreview() {
      this.previewOpen = false;
      this.previewItems = [];
      this.loading = false;
      this.error = "";
    },
  },
});
