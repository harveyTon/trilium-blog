<template>
  <el-header class="app-header">
    <div class="header-content">
      <button
        class="search-toggle"
        type="button"
        aria-label="打开搜索"
        @click="mobileSearchOpen = true"
      >
        <el-icon><Search /></el-icon>
      </button>

      <router-link :to="{ name: 'HomePage' }" class="blog-brand" @click="handleBrandClick">
        <img :src="logoSrc" :alt="site.title + ' - 返回首页'" class="blog-logo" />
        <span class="blog-brand-copy">
          <span class="blog-title">{{ site.title }}</span>
          <span v-if="site.subtitle" class="blog-subtitle">{{ site.subtitle }}</span>
        </span>
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

    <transition name="mobile-search-fade">
      <div v-if="mobileSearchOpen" class="mobile-search-modal" @click.self="mobileSearchOpen = false">
        <div class="mobile-search-panel">
          <div class="mobile-search-header">
            <p class="mobile-search-title">搜索文章</p>
            <button
              class="mobile-search-close"
              type="button"
              aria-label="关闭搜索"
              @click="mobileSearchOpen = false"
            >
              ×
            </button>
          </div>
          <GlobalSearchBox />
        </div>
      </div>
    </transition>
  </el-header>
</template>

<script>
import { ElHeader, ElIcon } from "element-plus";
import { Moon, Search, Sunny } from "@element-plus/icons-vue";
import { ref, watch } from "vue";
import { useRoute } from "vue-router";
import logoSrc from "../../assets/logo.png";
import GlobalSearchBox from "./GlobalSearchBox.vue";

export default {
  name: "AppHeader",
  components: {
    ElHeader,
    ElIcon,
    Moon,
    Search,
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
    const route = useRoute();
    const mobileSearchOpen = ref(false);

    const handleBrandClick = () => {
      mobileSearchOpen.value = false;

      if (route.name === "HomePage" && !route.query.page) {
        window.scrollTo({ top: 0, behavior: "smooth" });
        return;
      }

      requestAnimationFrame(() => {
        window.scrollTo({ top: 0, behavior: "auto" });
      });
    };

    watch(
      () => route.fullPath,
      () => {
        mobileSearchOpen.value = false;
      }
    );

    return {
      logoSrc,
      handleBrandClick,
      mobileSearchOpen,
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

.search-toggle {
  display: none;
}

.blog-brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  position: relative;
  z-index: 2;
}

.blog-brand-copy {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.blog-subtitle {
  color: var(--text-faint);
  font-size: 12px;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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
    display: grid;
    grid-template-columns: 40px minmax(0, 1fr) 40px;
    grid-template-areas: "search brand theme";
    gap: 10px;
    align-items: center;
    min-height: 44px;
  }

  .search-toggle,
  .theme-toggle {
    width: 40px;
    height: 40px;
    border: 1px solid rgba(255, 255, 255, 0.14);
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.06);
    color: var(--text-inverse);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0;
  }

  .search-toggle {
    grid-area: search;
  }

  .blog-brand {
    grid-area: brand;
    justify-self: center;
    align-self: center;
    gap: 10px;
  }

  .header-search {
    display: none;
  }

  .theme-toggle {
    grid-area: theme;
    justify-self: end;
  }

  .blog-title {
    display: inline;
    font-size: 1.05rem;
    white-space: nowrap;
  }

  .blog-subtitle {
    display: none;
  }

  .mobile-search-modal {
    position: fixed;
    inset: 0;
    z-index: 1100;
    background: rgba(12, 18, 26, 0.52);
    backdrop-filter: blur(10px);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding: calc(var(--header-h) + 12px) 12px 12px;
    box-sizing: border-box;
  }

  .mobile-search-panel {
    width: min(100%, 540px);
    border-radius: 20px;
    background: var(--brand);
    border: 1px solid rgba(255, 255, 255, 0.12);
    box-shadow: 0 20px 60px rgba(8, 12, 18, 0.28);
    padding: 14px;
  }

  .mobile-search-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 12px;
  }

  .mobile-search-title {
    margin: 0;
    color: var(--text-inverse);
    font-size: 0.95rem;
    font-weight: 700;
  }

  .mobile-search-close {
    width: 34px;
    height: 34px;
    border-radius: 999px;
    border: 1px solid rgba(255, 255, 255, 0.12);
    background: rgba(255, 255, 255, 0.08);
    color: var(--text-inverse);
    font-size: 24px;
    line-height: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0;
  }

  .mobile-search-panel :deep(.global-search) {
    width: 100%;
  }

  .mobile-search-panel :deep(.global-search-input) {
    border-color: rgba(255, 255, 255, 0.14);
    background: rgba(255, 255, 255, 0.08);
  }

  .mobile-search-panel :deep(.search-preview-panel) {
    position: static;
    margin-top: 10px;
  }

  .mobile-search-fade-enter-active,
  .mobile-search-fade-leave-active {
    transition: opacity 0.18s ease;
  }

  .mobile-search-fade-enter-from,
  .mobile-search-fade-leave-to {
    opacity: 0;
  }
}
</style>
