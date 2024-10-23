import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import { createPinia } from "pinia";
import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

import "prismjs";
import "prismjs/components/prism-bash";
import "prismjs/components/prism-cshtml";
import "prismjs/components/prism-css";
import "prismjs/components/prism-docker";
import "prismjs/components/prism-javascript";
import "prismjs/components/prism-markup";
import "prismjs/components/prism-php";
import "prismjs/themes/prism.css";

const app = createApp(App);
const pinia = createPinia();
app.use(ElementPlus);
app.use(router);
app.use(pinia);
app.mount("#app");
