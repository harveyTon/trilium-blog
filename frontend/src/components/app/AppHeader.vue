<template>
  <el-header class="app-header">
    <div class="header-content">
      <router-link :to="{ name: 'HomePage' }" class="blog-brand">
        <img :src="logoSrc" :alt="site.name + ' - 返回首页'" class="blog-logo" />
        <span class="blog-title">{{ site.name }}</span>
      </router-link>
      <div class="header-search">
        <GlobalSearchBox />
      </div>
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
.header-content {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  min-height: var(--header-h);
}

.blog-brand {
  min-width: 0;
  position: relative;
  z-index: 2;
}

.header-search {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: min(560px, calc(100% - 320px));
  z-index: 1;
}

.header-search :deep(.global-search) {
  width: 100%;
}

.theme-toggle {
  position: relative;
  z-index: 2;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .header-content {
    gap: 10px;
  }

  .header-search {
    position: static;
    transform: none;
    width: 100%;
    order: 3;
    margin-top: 8px;
  }

  .blog-title {
    display: none;
  }

  .header-content {
    flex-wrap: wrap;
  }
}
</style>
