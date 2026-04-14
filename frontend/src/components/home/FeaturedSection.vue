<template>
  <section v-if="items.length" class="featured-section" :aria-label="t('featured.sectionLabel')">
    <div class="featured-header">
      <p class="featured-label">{{ t('featured.kicker') }}</p>
      <h2 class="featured-heading">{{ t('featured.heading') }}</h2>
    </div>

    <div class="featured-carousel">
      <button
        v-if="showDesktopControls"
        class="featured-nav featured-nav--prev"
        type="button"
        :aria-label="t('featured.ariaPrev')"
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
        v-if="showDesktopControls"
        class="featured-nav featured-nav--next"
        type="button"
        :aria-label="t('featured.ariaNext')"
        @click="goNext"
      >
        ›
      </button>
    </div>

    <div v-if="showControls" class="featured-dots">
      <button
        v-for="page in totalPages"
        :key="page"
        type="button"
        :aria-label="t('featured.ariaDot', { page })"
        :class="['featured-dot', { 'is-active': page - 1 === currentPage }]"
        @click="goToPage(page - 1)"
      ></button>
    </div>
  </section>
</template>

<script>
import { computed, ref, watch } from "vue";
import FeaturedCard from "./FeaturedCard.vue";
import { t } from "../../i18n";

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
    const currentPage = ref(0);
    const touchStartX = ref(0);

    const totalPages = computed(() => props.items.length);
    const showControls = computed(() => props.items.length > 1);
    const showDesktopControls = computed(() => props.items.length > 1);
    const maxPage = computed(() => Math.max(0, totalPages.value - 1));
    const trackStyle = computed(() => ({
      transform: `translateX(-${currentPage.value * 100}%)`,
    }));
    const itemStyle = computed(() => ({
      flex: "0 0 100%",
      maxWidth: "100%",
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

    const goToPage = (page) => {
      currentPage.value = Math.max(0, Math.min(page, maxPage.value));
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

    watch(() => props.items.length, clampCurrentPage, { immediate: true });

    return {
      t,
      currentPage,
      totalPages,
      showControls,
      showDesktopControls,
      trackStyle,
      itemStyle,
      goPrev,
      goNext,
      goToPage,
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
  gap: 18px;
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
  position: relative;
  width: 100%;
}

.featured-carousel::after {
  content: "";
  position: absolute;
  left: 18px;
  right: 18px;
  bottom: -2px;
  height: 30px;
  border-radius: 0 0 24px 24px;
  background: linear-gradient(
    180deg,
    color-mix(in srgb, var(--bg) 0%, transparent) 0%,
    color-mix(in srgb, var(--bg) 92%, var(--surface) 8%) 100%
  );
  pointer-events: none;
  z-index: 0;
}

.featured-viewport {
  padding: 6px 2px 8px;
  margin: -6px -2px -8px;
  overflow: hidden;
  min-width: 0;
  width: 100%;
  position: relative;
  z-index: 1;
}

.featured-track {
  display: flex;
  transition: transform 0.42s ease;
  will-change: transform;
  align-items: stretch;
}

.featured-item-shell {
  min-width: 100%;
  flex-shrink: 0;
  display: flex;
}

.featured-item-shell :deep(.featured-card) {
  height: 100%;
  width: 100%;
}

.featured-item-shell :deep(.featured-link) {
  height: 100%;
  box-sizing: border-box;
}

.featured-nav {
  position: absolute;
  top: 50%;
  z-index: 2;
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
  box-shadow: var(--shadow-sm);
  transition: transform 0.16s ease, border-color 0.16s ease, color 0.16s ease;
}

.featured-nav--prev {
  left: -21px;
  transform: translate(-100%, -50%);
}

.featured-nav--next {
  right: -21px;
  transform: translate(100%, -50%);
}

.featured-nav:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.featured-nav--prev:hover {
  transform: translate(-100%, calc(-50% - 1px));
}

.featured-nav--next:hover {
  transform: translate(100%, calc(-50% - 1px));
  border-color: var(--accent);
}

.featured-dots {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 4px;
  padding: 2px 0 0;
  position: relative;
  z-index: 1;
}

.featured-dot {
  width: 8px;
  height: 8px;
  padding: 0;
  border: 0;
  border-radius: 999px;
  background: var(--border);
  transition: transform 0.18s ease, background 0.18s ease;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--surface) 65%, transparent);
}

.featured-dot.is-active {
  background: var(--accent);
  transform: scale(1.25);
}

@media (max-width: 768px) {
  .featured-carousel::after {
    left: 10px;
    right: 10px;
    height: 22px;
    border-radius: 0 0 20px 20px;
  }

  .featured-viewport {
    padding: 6px 0 6px;
    margin: -6px 0 -6px;
  }

  .featured-nav {
    display: none;
  }

  .featured-heading {
    font-size: 24px;
  }
}
</style>
