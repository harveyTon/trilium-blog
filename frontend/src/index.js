import { createRouter, createWebHistory } from "vue-router";
import HomePage from "../views/Home.vue";
import NotFound from "../views/NotFound.vue";

const routes = [
  { path: "/", name: "HomePage", component: HomePage },
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
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
