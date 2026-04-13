import { createRouter, createWebHistory } from "vue-router";
import HomePage from "../views/Home.vue";
import NotFound from "../views/NotFound.vue";

const routes = [
  { path: "/", name: "HomePage", component: HomePage },
  {
    path: "/search",
    name: "Search",
    component: () => import("../views/Search.vue"),
  },
  {
    path: "/post/:noteId",
    name: "Article",
    component: () => import("../views/Article.vue"),
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: NotFound,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    }

    if (to.fullPath !== from.fullPath) {
      return { top: 0 };
    }

    return undefined;
  },
});

export default router;
