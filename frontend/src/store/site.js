import { defineStore } from "pinia";
import { fetchSite } from "../api/blog";
import { setLocale } from "../i18n";

export const useSiteStore = defineStore("site", {
  state: () => ({
    site: {
      title: "",
      subtitle: "",
      domain: "",
      locale: "zh-CN",
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
        if (data.locale) {
          setLocale(data.locale);
        }
        document.title = [data.title, data.subtitle].filter(Boolean).join(" | ");
        this.loaded = true;
      } catch (error) {
        console.error("Fetch Site Info Error:", error);
      }
    },
  },
});
