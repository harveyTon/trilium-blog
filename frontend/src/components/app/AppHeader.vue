<template>
  <el-header class="app-header">
    <div class="header-content">
      <router-link to="/" class="blog-brand">
        <img :src="logoSrc" :alt="site.name + ' - 返回首页'" class="blog-logo" />
        <span class="blog-title">{{ site.name }}</span>
      </router-link>
      <button
        class="theme-toggle"
        :aria-label="isDark ? '切换到亮色模式' : '切换到暗色模式'"
        @click="$emit('toggle-theme')"
      >
        <el-icon>
          <component :is="isDark ? 'Sunny' : 'Moon'" />
        </el-icon>
      </button>
    </div>
  </el-header>
</template>

<script>
import { ElHeader, ElIcon } from "element-plus";
import { Moon, Sunny } from "@element-plus/icons-vue";
import logoSrc from "../../assets/logo.png";

export default {
  name: "AppHeader",
  components: {
    ElHeader,
    ElIcon,
    Moon,
    Sunny,
  },
  props: {
    isDark: {
      type: Boolean,
      required: true,
    },
    site: {
      type: Object,
      required: true,
    },
  },
  emits: ["toggle-theme"],
  setup() {
    return {
      logoSrc,
    };
  },
};
</script>
