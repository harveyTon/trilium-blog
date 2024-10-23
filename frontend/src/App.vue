<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <router-link to="/" class="blog-brand" @click="refreshHome">
          <img src="./assets/logo.png" alt="Logo" class="blog-logo" />
          <span class="blog-title">{{ blogStore.blogInfo.blogName }}</span>
        </router-link>
        <el-icon class="theme-icon" @click="toggleDark()">
          <component :is="isDark ? 'Sunny' : 'Moon'" />
        </el-icon>
      </div>
    </el-header>
    <el-main class="app-main">
      <div class="main-content">
        <router-view></router-view>
      </div>
    </el-main>
    <el-footer class="app-footer">
      <div class="footer-content">
        <div class="footer-text">
          <p>
            © {{ new Date().getFullYear() }} {{ blogStore.blogInfo.blogName }}.
            保留所有权利。
          </p>
          <p>
            由
            <a href="https://github.com/harveyTon/trilium-blog" target="_blank"
              >Trilium Blog</a
            >
            &
            <a href="https://github.com/TriliumNext/Notes" target="_blank"
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
import { useBlogStore } from "./store";

export default {
  name: "App",
  components: {
    Moon,
    Sunny,
  },
  setup() {
    const blogStore = useBlogStore();
    const isDark = useDark({
      storageKey: "vueuse-color-scheme",
      valueDark: "dark",
      valueLight: "light",
    });
    const toggleDark = useToggle(isDark);

    const applyTheme = (dark) => {
      document.documentElement.classList.toggle("dark", dark);
    };

    const refreshHome = () => {
      // TODO: PAGE-1
    };

    watch(
      isDark,
      (newValue) => {
        applyTheme(newValue);
      },
      { immediate: true }
    );

    onMounted(() => {
      blogStore.fetchBlogInfo();
      applyTheme(isDark.value);
    });

    return {
      blogStore,
      isDark,
      toggleDark,
      Moon,
      Sunny,
      refreshHome,
    };
  },
};
</script>

<style>
@import url("https://ik.imagekit.io/tigerton/LXGWWenKai-Bold/result.css");
@import "element-plus/theme-chalk/dark/css-vars.css";

body {
  margin: 0;
  padding: 0;
  background-color: #f5f7fa;
  color: #2c3e50;
}

#app {
  font-family: "LXGW WenKai", "Helvetica Neue", Helvetica, "PingFang SC",
    "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
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
  background-color: #2c3e50;
  color: #ffffff;
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
  color: #ffffff;
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
  color: #ffffff;
  text-decoration: none;
}

.footer-text a:hover {
  text-decoration: underline;
}

.el-backtop {
  background-color: #409eff;
  color: #ffffff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.el-backtop:hover {
  background-color: #66b1ff;
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

html.dark {
  background-color: #1a1a1a;
  color: #ffffff;
}

html.dark .app-header,
html.dark .app-footer {
  background-color: #333333;
}

html.dark .el-main {
  background-color: #1a1a1a;
}
html.dark .el-backtop {
  color: #909399;
  background-color: #333333;
}

html.dark .el-backtop:hover {
  color: aliceblue;
  background-color: #1a1a1a;
}
.theme-switch {
  margin-left: 20px;
}

@media (max-width: 768px) {
  .header-content {
    justify-content: space-between;
  }

  .theme-switch {
    margin-left: 10px;
  }
}

.theme-icon {
  font-size: 24px;
  cursor: pointer;
  color: #ffffff;
  transition: color 0.3s ease;
}

.theme-icon:hover {
  color: #409eff;
}

@media (max-width: 768px) {
  .header-content {
    justify-content: space-between;
  }

  .theme-icon {
    font-size: 20px;
  }
}
</style>
