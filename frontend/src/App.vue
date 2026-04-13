<template>
  <el-container class="app-container">
    <AppHeader :is-dark="isDark" :site="site" @toggle-theme="toggleDark()" />
    <el-main class="app-main">
      <div class="main-content">
        <router-view v-slot="{ Component }">
          <keep-alive include="HomePage">
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </div>
    </el-main>
    <AppFooter :site="site" />
    <el-backtop></el-backtop>
  </el-container>
</template>

<script>
import { useDark, useToggle } from "@vueuse/core";
import { onMounted, watch } from "vue";
import { storeToRefs } from "pinia";
import { useSiteStore } from "./store";
import AppFooter from "./components/app/AppFooter.vue";
import AppHeader from "./components/app/AppHeader.vue";

export default {
  name: "App",
  components: {
    AppHeader,
    AppFooter,
  },
  setup() {
    const siteStore = useSiteStore();
    const { site } = storeToRefs(siteStore);
    const isDark = useDark({
      storageKey: "vueuse-color-scheme",
      valueDark: "dark",
      valueLight: "light",
    });
    const toggleDark = useToggle(isDark);

    const applyTheme = (dark) => {
      document.documentElement.classList.toggle("dark", dark);
    };

    watch(
      isDark,
      (newValue) => {
        applyTheme(newValue);
      },
      { immediate: true }
    );

    onMounted(() => {
      siteStore.fetchSiteInfo();
      applyTheme(isDark.value);
    });

    return {
      site,
      isDark,
      toggleDark,
    };
  },
};
</script>

<style>
@import url("https://ik.imagekit.io/tigerton/LXGWWenKai-Bold/result.css");
@import "element-plus/theme-chalk/dark/css-vars.css";

:root {
  --bg: #eef1f4;
  --surface: #ffffff;
  --surface-muted: #f7f8fa;
  --border: #d9dee5;
  --border-soft: #e7ebf0;

  --text: #243447;
  --text-soft: #5c6b7a;
  --text-faint: #8a96a3;
  --text-inverse: #ffffff;

  --brand: #2d4359;
  --brand-strong: #223447;
  --accent: #f08c2e;

  --link: #2f6db3;
  --link-hover: #1f5a9b;

  --shadow-sm: 0 2px 8px rgba(20, 30, 40, 0.04);
  --shadow-md: 0 8px 24px rgba(20, 30, 40, 0.08);

  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;

  --content-w: 760px;
  --list-w: 860px;
  --header-h: 64px;

  --font-family: "LXGW WenKai", "Helvetica Neue", Helvetica, "PingFang SC",
    "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  --mono: "SFMono-Regular", "JetBrains Mono", "Cascadia Code", ui-monospace,
    monospace;
}

html.dark {
  --bg: #141a22;
  --surface: #1b222d;
  --surface-muted: #1f2836;
  --border: #2e3a4a;
  --border-soft: #253040;

  --text: #dce4ed;
  --text-soft: #94a3b4;
  --text-faint: #5f7082;
  --text-inverse: #141a22;

  --brand: #3d5570;
  --brand-strong: #4d6a85;
  --accent: #e89430;

  --link: #6aadec;
  --link-hover: #8dc4f5;
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--bg);
  color: var(--text);
}

:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

#app {
  font-family: var(--font-family);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header,
.app-footer {
  background-color: var(--brand);
  color: var(--text-inverse);
  padding: 20px 0;
}

.app-header {
  position: fixed;
  width: 100%;
  z-index: 1000;
  height: var(--header-h);
  display: flex;
  align-items: center;
  box-shadow: var(--shadow-sm);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.blog-brand {
  display: flex;
  align-items: center;
  text-decoration: none;
}

.blog-logo {
  height: 40px;
  margin-right: 10px;
}

.blog-title {
  font-size: 1.5em;
  color: var(--text-inverse);
  line-height: 1;
}

.app-main {
  flex: 1;
  padding-top: calc(var(--header-h) + 16px);
  padding-bottom: 40px;
  overflow: visible;
}

.main-content {
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: 0px;
  box-sizing: border-box;
}

.app-footer {
  height: auto;
  margin-top: auto;
}

.footer-content {
  text-align: center;
}

.footer-text {
  font-size: 0.9em;
  line-height: 1.5;
}

.footer-text a {
  color: var(--text-inverse);
  text-decoration: none;
}

.footer-text a:hover {
  text-decoration: underline;
}

.el-backtop {
  background-color: var(--brand);
  color: var(--text-inverse);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.3s, color 0.3s;
}

.el-backtop:hover {
  background-color: var(--brand-strong);
}

@media (max-width: 768px) {
  .header-content,
  .footer-content,
  .main-content {
    padding: 0 10px;
  }

  .app-main {
    padding: 0;
    padding-top: calc(var(--header-h) + 8px);
    padding-bottom: 20px;
  }

  .footer-text {
    font-size: 0.8em;
  }

  .app-header {
    height: 50px;
    padding: 0;
  }

  .blog-logo {
    height: 30px;
  }

  .blog-title {
    font-size: 1.2em;
  }

  .el-backtop {
    right: 15px !important;
    bottom: 15px !important;
  }
}

html.dark .el-backtop {
  background-color: var(--brand);
  color: var(--text-inverse);
}

html.dark .el-backtop:hover {
  background-color: var(--brand-strong);
}

.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--text-inverse);
  font-size: 24px;
  transition: color 0.2s ease;
  border-radius: var(--radius-sm);
}

.theme-toggle:hover {
  color: var(--accent);
}

@media (max-width: 768px) {
  .header-content {
    justify-content: space-between;
  }

  .theme-toggle {
    font-size: 20px;
    width: 40px;
    height: 40px;
  }
}
</style>
