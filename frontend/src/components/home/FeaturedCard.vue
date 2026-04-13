<template>
  <article class="featured-card">
    <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }" class="featured-link">
      <p class="featured-kicker">精选</p>
      <h2 class="featured-title">{{ post.title }}</h2>
      <p class="featured-date">{{ post.dateModified }}</p>
      <SummaryPreview :content="post.summary" />
    </router-link>
  </article>
</template>

<script>
import SummaryPreview from "./SummaryPreview.vue";

export default {
  name: "FeaturedCard",
  components: {
    SummaryPreview,
  },
  props: {
    post: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style scoped>
.featured-card {
  border: 1px solid var(--border-soft);
  background: var(--surface);
  border-radius: var(--radius-lg);
}

.featured-link {
  display: block;
  padding: 22px 24px;
  text-decoration: none;
}

.featured-kicker {
  margin: 0 0 10px;
  font-size: 12px;
  color: var(--accent);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.featured-title {
  margin: 0;
  color: var(--text);
  line-height: 1.3;
}

.featured-date {
  margin: 12px 0 0;
  color: var(--text-faint);
  font-size: 13px;
}
</style>
