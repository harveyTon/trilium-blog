<template>
  <section v-if="items.length" class="featured-section" aria-label="精选文章">
    <div class="featured-header">
      <p class="featured-label">优先阅读</p>
      <h2 class="featured-heading">精选文章</h2>
    </div>

    <div class="featured-carousel">
      <button
        v-if="showControls"
        class="featured-nav featured-nav--prev"
        type="button"
        aria-label="查看上一组精选文章"
        @click="goPrev"
      >
        ‹
      </button>

      <div
        class="featured-viewport"
        @touchstart.passive="handleTouchStart"
        @touchend.passive="handleTouchEnd"
      >
        <div class="featured-track" :style="trackStyle">
          <div
            v-for="item in items"
            :key="item.noteId"
            class="featured-item-shell"
            :style="itemStyle"
          >
            <FeaturedCard :post="item" />
          </div>
        </div>
      </div>

      <button
        v-if="showControls"
        class="featured-nav featured-nav--next"
        type="button"
        aria-label="查看下一组精选文章"
        @click="goNext"
      >
        ›
      </button>
    </div>

    <div v-if="showControls" class="featured-dots" aria-hidden="true">
      <span
        v-for="page in totalPages"
        :key="page"
        :class="['featured-dot', { 'is-active': page - 1 === currentPage }]"
      ></span>
    </div>
  </section>
</template>

<script>
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import FeaturedCard from "./FeaturedCard.vue";

export default {
  name: "FeaturedSection",
  components: {
    FeaturedCard,
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const visibleCount = ref(3);
    const currentPage = ref(0);
    const touchStartX = ref(0);

    const syncVisibleCount = () => {
      visibleCount.value = window.innerWidth <= 768 ? 1 : 3;
    };

    const totalPages = computed(() => Math.max(1, props.items.length - visibleCount.value + 1));
    const showControls = computed(() => props.items.length > visibleCount.value);
    const maxPage = computed(() => Math.max(0, totalPages.value - 1));
    const trackStyle = computed(() => ({
      transform: `translateX(calc(-1 * ${currentPage.value} * (100% + var(--featured-gap)) / ${visibleCount.value}))`,
    }));
    const itemStyle = computed(() => ({
      flex: `0 0 calc((100% - var(--featured-gap) * ${visibleCount.value - 1}) / ${visibleCount.value})`,
      maxWidth: `calc((100% - var(--featured-gap) * ${visibleCount.value - 1}) / ${visibleCount.value})`,
    }));

    const clampCurrentPage = () => {
      currentPage.value = Math.min(currentPage.value, maxPage.value);
    };

    const goPrev = () => {
      currentPage.value = currentPage.value <= 0 ? maxPage.value : currentPage.value - 1;
    };

    const goNext = () => {
      currentPage.value = currentPage.value >= maxPage.value ? 0 : currentPage.value + 1;
    };

    const handleTouchStart = (event) => {
      touchStartX.value = event.changedTouches[0]?.clientX || 0;
    };

    const handleTouchEnd = (event) => {
      const touchEndX = event.changedTouches[0]?.clientX || 0;
      const deltaX = touchEndX - touchStartX.value;
      if (Math.abs(deltaX) < 32 || !showControls.value) {
        return;
      }
      if (deltaX > 0) {
        goPrev();
      } else {
        goNext();
      }
    };

    watch([() => props.items.length, visibleCount], clampCurrentPage, { immediate: true });

    onMounted(() => {
      syncVisibleCount();
      window.addEventListener("resize", syncVisibleCount);
    });

    onUnmounted(() => {
      window.removeEventListener("resize", syncVisibleCount);
    });

    return {
      currentPage,
      totalPages,
      showControls,
      trackStyle,
      itemStyle,
      goPrev,
      goNext,
      handleTouchStart,
      handleTouchEnd,
    };
  },
};
</script>

<style scoped>
.featured-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 12px;
}

.featured-header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.featured-label {
  margin: 0;
  font-size: 12px;
  color: var(--text-faint);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.featured-heading {
  margin: 0;
  color: var(--text);
  font-size: 28px;
}

.featured-carousel {
  --featured-gap: 16px;
  position: relative;
  display: flex;
  align-items: center;
  gap: 14px;
}

.featured-viewport {
  overflow: hidden;
  flex: 1;
  min-width: 0;
}

.featured-track {
  display: flex;
  gap: var(--featured-gap);
  transition: transform 0.42s ease;
  will-change: transform;
}

.featured-item-shell {
  min-width: 0;
}

.featured-item-shell :deep(.featured-card) {
  height: 100%;
}

.featured-item-shell :deep(.featured-link) {
  height: 100%;
  box-sizing: border-box;
}

.featured-nav {
  width: 42px;
  height: 42px;
  border: 1px solid var(--border-soft);
  border-radius: 999px;
  background: color-mix(in srgb, var(--surface) 90%, white 10%);
  color: var(--text);
  font-size: 28px;
  line-height: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  box-shadow: var(--shadow-sm);
  transition: transform 0.16s ease, border-color 0.16s ease, color 0.16s ease;
}

.featured-nav:hover {
  transform: translateY(-1px);
  border-color: var(--accent);
  color: var(--accent);
}

.featured-dots {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.featured-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--border);
  transition: transform 0.18s ease, background 0.18s ease;
}

.featured-dot.is-active {
  background: var(--accent);
  transform: scale(1.25);
}

@media (max-width: 768px) {
  .featured-carousel {
    --featured-gap: 12px;
    gap: 10px;
  }

  .featured-nav {
    width: 36px;
    height: 36px;
    font-size: 24px;
  }

  .featured-heading {
    font-size: 24px;
  }
}
</style>
