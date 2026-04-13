import { defineStore } from "pinia";
import { fetchSite } from "../api/blog";

export const useSiteStore = defineStore("site", {
  state: () => ({
    site: {
      name: "",
      title: "",
      domain: "",
      comments: {
        enabled: false,
        server: "",
        site: "",
      },
      imageProxy: {
        enabled: false,
        baseUrl: "",
      },
    },
    loaded: false,
  }),
  actions: {
    async fetchSiteInfo() {
      try {
        const data = await fetchSite();
        this.site = data;
        document.title = data.title || "";
        this.loaded = true;
      } catch (error) {
        console.error("Fetch Site Info Error:", error);
      }
    },
  },
});
