<template>
  <section v-if="items.length" class="featured-section" aria-label="精选文章">
    <div class="featured-header">
      <p class="featured-label">优先阅读</p>
      <h2 class="featured-heading">精选文章</h2>
    </div>

    <div v-if="isScrollable" class="featured-window">
      <div class="featured-track" :style="trackStyle">
        <div v-for="(item, index) in scrollingItems" :key="`${item.noteId}-${index}`" class="featured-item-shell">
          <FeaturedCard :post="item" />
        </div>
      </div>
    </div>

    <div v-else class="featured-grid">
      <FeaturedCard v-for="item in items" :key="item.noteId" :post="item" />
    </div>
  </section>
</template>

<script>
import { computed } from "vue";
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
    const isScrollable = computed(() => props.items.length > 3);
    const scrollingItems = computed(() => (isScrollable.value ? [...props.items, ...props.items] : props.items));
    const trackStyle = computed(() => ({
      "--featured-count": String(props.items.length),
      "--featured-duration": `${Math.max(props.items.length * 8, 28)}s`,
    }));

    return {
      isScrollable,
      scrollingItems,
      trackStyle,
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

.featured-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 16px;
}

.featured-window {
  --featured-gap: 16px;
  --featured-card-width: minmax(280px, 1fr);
  position: relative;
  overflow: hidden;
  padding-bottom: 4px;
  mask-image: linear-gradient(to right, transparent 0, black 5%, black 95%, transparent 100%);
}

.featured-window::after {
  content: "";
  display: block;
  clear: both;
}

.featured-window {
  overflow: hidden;
}

.featured-track {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: minmax(280px, 1fr);
  gap: var(--featured-gap);
  animation: featured-marquee var(--featured-duration) linear infinite;
  will-change: transform;
}

.featured-window:hover .featured-track {
  animation-play-state: paused;
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

@keyframes featured-marquee {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(calc(-1 * (280px + var(--featured-gap)) * var(--featured-count)));
  }
}

@media (max-width: 768px) {
  .featured-window {
    --featured-gap: 12px;
  }

  .featured-track {
    grid-auto-columns: minmax(240px, 82vw);
  }

  @keyframes featured-marquee {
    from {
      transform: translateX(0);
    }
    to {
      transform: translateX(calc(-1 * (240px + var(--featured-gap)) * var(--featured-count)));
    }
  }
}
</style>
