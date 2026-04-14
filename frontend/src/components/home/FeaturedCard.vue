<template>
  <article class="featured-card">
    <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }" class="featured-link">
      <div class="featured-meta">
        <p class="featured-kicker">{{ t('featured.badge') }}</p>
        <p class="featured-date">{{ formattedDate }}</p>
      </div>
      <div class="featured-title-block">
        <h2 class="featured-title">{{ post.title }}</h2>
      </div>
      <div class="featured-summary-block">
        <SummaryPreview :content="resolvedSummary" :type="summaryType" :lines="4" :flush="true" />
      </div>
    </router-link>
  </article>
</template>

<script>
import { computed, toRef } from "vue";
import SummaryPreview from "./SummaryPreview.vue";
import { useSummaryStatus } from "../../composables/useSummaryStatus";
import { t, locale } from "../../i18n";

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
    const formattedDate = computed(() => {
      const date = new Date(props.post.dateModified);
      if (Number.isNaN(date.getTime())) {
        return props.post.dateModified;
      }

      return new Intl.DateTimeFormat(locale(), {
        year: "numeric",
        month: "short",
        day: "numeric",
      }).format(date);
    });

    return {
      t,
      locale,
      formattedDate,
      resolvedSummary,
      summaryType,
    };
  },
};
</script>

<style scoped>
.featured-card {
  --fc-meta-to-title: 16px;
  --fc-title-to-summary: 14px;
  --fc-px: 28px;
  --fc-py: 24px;
  --fc-title-lh: 1.3;
  --fc-summary-lh: 1.65;
  border: 1px solid var(--border-soft);
  border-radius: 16px;
  background: var(--surface);
  height: 100%;
  transition: transform 0.18s ease, border-color 0.18s ease;
  cursor: pointer;
}

.featured-card:hover {
  transform: translateY(-2px);
  border-color: color-mix(in srgb, var(--border-soft) 50%, var(--accent) 50%);
}

.featured-link {
  display: flex;
  flex-direction: column;
  padding: var(--fc-py) var(--fc-px);
  height: 100%;
  box-sizing: border-box;
  text-decoration: none;
  min-height: 280px;
  max-height: 360px;
  cursor: pointer;
  border-radius: inherit;
  overflow: hidden;
}

.featured-link:focus-visible {
  outline: none;
  border-color: color-mix(in srgb, var(--accent) 50%, var(--border-soft) 50%);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent) 14%, transparent);
}

.featured-card:hover .featured-title,
.featured-link:hover .featured-title {
  color: var(--link-hover);
}

.featured-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  height: 20px;
}

.featured-kicker {
  margin: 0;
  font-size: 11px;
  color: var(--text-faint);
  letter-spacing: 0.06em;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.featured-kicker::before {
  content: "";
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--accent) 82%, transparent);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent) 10%, transparent);
}

.featured-date {
  margin: 0;
  color: var(--text-faint);
  font-size: 11px;
  letter-spacing: 0.02em;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.featured-title-block {
  margin-top: var(--fc-meta-to-title);
  max-width: min(520px, 100%);
  flex-shrink: 0;
}

.featured-title {
  margin: 0;
  color: var(--text);
  line-height: var(--fc-title-lh);
  font-size: clamp(20px, 2.2vw, 26px);
  font-weight: 600;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;
}

.featured-summary-block {
  margin-top: var(--fc-title-to-summary);
  flex: 1 1 0;
  min-height: 0;
  overflow: hidden;
}

.featured-link :deep(.summary-preview) {
  margin: 0;
  font-size: 14px;
  line-height: var(--fc-summary-lh);
  color: var(--text-soft);
  width: 100%;
}

.featured-link :deep(.summary-preview[data-summary-type="ai"]) {
  color: var(--text-soft);
}

.featured-link :deep(.summary-badge) {
  font-size: 10px;
  padding: 1px 5px 1px;
  vertical-align: baseline;
}

@media (max-width: 768px) {
  .featured-card {
    border-radius: 14px;
    --fc-px: 20px;
    --fc-py: 20px;
  }

  .featured-link {
    min-height: 260px;
    max-height: 340px;
  }

  .featured-title {
    font-size: 19px;
  }

  .featured-title-block {
    max-width: 100%;
  }

  .featured-link :deep(.summary-preview) {
    font-size: 13px;
    line-height: 1.6;
  }
}
</style>
