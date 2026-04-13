<template>
  <article class="featured-card">
    <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }" class="featured-link">
      <p class="featured-kicker">精选</p>
      <h2 class="featured-title">{{ post.title }}</h2>
      <p class="featured-date">{{ post.dateModified }}</p>
      <SummaryPreview :content="resolvedSummary" :type="summaryType" />
    </router-link>
  </article>
</template>

<script>
import { computed, toRef } from "vue";
import SummaryPreview from "./SummaryPreview.vue";
import { useSummaryStatus } from "../../composables/useSummaryStatus";

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
  setup(props) {
    const { preferredSummary } = useSummaryStatus(toRef(props, "post"));
    const resolvedSummary = computed(() => preferredSummary.value.text || props.post.summary || "");
    const summaryType = computed(() => preferredSummary.value.type);

    return {
      resolvedSummary,
      summaryType,
    };
  },
};
</script>

<style scoped>
.featured-card {
  border: 1px solid var(--border-soft);
  background: var(--surface);
  border-radius: var(--radius-lg);
  height: 100%;
}

.featured-link {
  display: flex;
  flex-direction: column;
  padding: 22px 24px;
  height: 100%;
  box-sizing: border-box;
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

.featured-link :deep(.summary-preview-wrap) {
  margin-top: 14px;
}
</style>
