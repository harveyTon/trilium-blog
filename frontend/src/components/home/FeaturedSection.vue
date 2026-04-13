<template>
  <section v-if="items.length" class="featured-section" aria-label="精选文章">
    <div class="featured-header">
      <p class="featured-label">精选文章</p>
      <h2 class="featured-heading">优先阅读</h2>
    </div>
    <div class="featured-grid">
      <FeaturedCard v-for="item in items" :key="item.noteId" :post="item" />
    </div>
  </section>
</template>

<script>
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
};
</script>

<style scoped>
.featured-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
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
</style>
