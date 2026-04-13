<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <router-link to="/" class="blog-brand">
          <img src="./assets/logo.png" :alt="site.name + ' - 返回首页'" class="blog-logo" />
          <span class="blog-title">{{ site.name }}</span>
        </router-link>
        <button
          class="theme-toggle"
          :aria-label="isDark ? '切换到亮色模式' : '切换到暗色模式'"
          @click="toggleDark()"
        >
          <el-icon>
            <component :is="isDark ? 'Sunny' : 'Moon'" />
          </el-icon>
        </button>
      </div>
    </el-header>
    <el-main class="app-main">
      <div class="main-content">
        <router-view v-slot="{ Component }">
          <keep-alive include="HomePage">
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </div>
    </el-main>
    <el-footer class="app-footer">
      <div class="footer-content">
        <div class="footer-text">
          <p>
            © {{ new Date().getFullYear() }} {{ site.name }}.
            保留所有权利。
          </p>
          <p>
            由
            <a href="https://github.com/harveyTon/trilium-blog" target="_blank"
              >Trilium Blog</a
            >
            &
            <a href="https://github.com/TriliumNext/Trilium" target="_blank"
              >TriliumNext Notes</a
            >
            强力驱动
          </p>
        </div>
      </div>
    </el-footer>
    <el-backtop></el-backtop>
  </el-container>
</template>

<script>
import { Moon, Sunny } from "@element-plus/icons-vue";
import { useDark, useToggle } from "@vueuse/core";
import { onMounted, watch } from "vue";
import { storeToRefs } from "pinia";
import { useSiteStore } from "./store";

export default {
  name: "App",
  components: {
    Moon,
    Sunny,
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
      Moon,
      Sunny,
    };
  },
};
</script>

<style>
@import url("https://ik.imagekit.io/tigerton/LXGWWenKai-Bold/result.css");
@import "element-plus/theme-chalk/dark/css-vars.css";

:root {
  --bg-page: #f5f7fa;
  --bg-surface: #ffffff;
  --bg-header-footer: #2c3e50;

  --text-primary: #2c3e50;
  --text-secondary: #636363;
  --text-muted: #6b6b6b;
  --text-inverse: #ffffff;

  --accent: #409eff;
  --accent-hover: #66b1ff;

  --border-color: rgba(0, 0, 0, 0.12);

  --font-family: "LXGW WenKai", "Helvetica Neue", Helvetica, "PingFang SC",
    "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
}

html.dark {
  --bg-page: #1a1a1a;
  --bg-surface: #1f1f1f;
  --bg-header-footer: #333333;

  --text-primary: #f5f5f5;
  --text-secondary: #bdbdbd;
  --text-muted: #a0a0a0;

  --border-color: rgba(255, 255, 255, 0.12);

  --accent: #d0d0d0;
  --accent-hover: #bdbdbd;
  --text-on-accent: #1a1a1a;
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--bg-page);
  color: var(--text-primary);
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
  background-color: var(--bg-header-footer);
  color: var(--text-inverse);
  padding: 20px 0;
}

.app-header {
  position: fixed;
  width: 100%;
  z-index: 1000;
  height: 60px;
  display: flex;
  align-items: center;
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
  padding-top: 80px;
  padding-bottom: 40px;
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
  background-color: var(--accent);
  color: var(--text-inverse);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.3s, color 0.3s;
}

.el-backtop:hover {
  background-color: var(--accent-hover);
}

@media (max-width: 768px) {
  .header-content,
  .footer-content,
  .main-content {
    padding: 0 10px;
  }

  .app-main {
    padding: 0;
    padding-top: 60px;
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
  background-color: var(--accent);
  color: var(--text-on-accent);
}

html.dark .el-backtop:hover {
  background-color: var(--accent-hover);
  color: var(--text-on-accent);
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
  transition: color 0.3s ease;
  border-radius: 8px;
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
