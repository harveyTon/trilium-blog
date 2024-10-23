import axios from "axios";
import { defineStore } from "pinia";
export const useBlogStore = defineStore("blog", {
  state: () => ({
    blogInfo: {
      blogName: "",
      blogTitle: "",
    },
  }),
  actions: {
    async fetchBlogInfo() {
      try {
        const response = await axios.get(`/api/info`, {
          params: {
            t: new Date().getTime(),
          },
        });
        this.blogInfo = response.data;
        document.title = response.data.blogTitle;
      } catch (error) {
        console.error("Fetch Blog Info Error:", error);
      }
    },
  },
});
