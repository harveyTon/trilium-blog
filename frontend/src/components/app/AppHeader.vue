<template>
  <el-header class="app-header">
    <div class="header-content">
      <router-link to="/" class="blog-brand">
        <img :src="logoSrc" :alt="site.name + ' - 返回首页'" class="blog-logo" />
        <span class="blog-title">{{ site.name }}</span>
      </router-link>
      <div class="header-actions">
        <GlobalSearchBox />
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
    </div>
  </el-header>
</template>

<script>
import { ElHeader, ElIcon } from "element-plus";
import { Moon, Sunny } from "@element-plus/icons-vue";
import logoSrc from "../../assets/logo.png";
import GlobalSearchBox from "./GlobalSearchBox.vue";

export default {
  name: "AppHeader",
  components: {
    ElHeader,
    ElIcon,
    Moon,
    Sunny,
    GlobalSearchBox,
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

<style scoped>
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .header-content {
    gap: 12px;
  }

  .header-actions {
    min-width: 0;
  }

  .blog-title {
    display: none;
  }
}
</style>
