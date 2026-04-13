<template>
  <article class="featured-card">
    <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }" class="featured-link">
      <div class="featured-copy">
        <div class="featured-top">
          <p class="featured-kicker">精选</p>
          <p class="featured-date">{{ formattedDate }}</p>
        </div>
        <div class="featured-title-wrap">
          <h2 class="featured-title">{{ post.title }}</h2>
        </div>
        <div class="featured-summary-wrap">
          <SummaryPreview :content="resolvedSummary" :type="summaryType" :lines="4" :flush="true" />
        </div>
      </div>
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
    const formattedDate = computed(() => {
      const date = new Date(props.post.dateModified);
      if (Number.isNaN(date.getTime())) {
        return props.post.dateModified;
      }

      return new Intl.DateTimeFormat("zh-CN", {
        year: "numeric",
        month: "short",
        day: "numeric",
      }).format(date);
    });

    return {
      formattedDate,
      resolvedSummary,
      summaryType,
    };
  },
};
</script>

<style scoped>
.featured-card {
  --featured-title-lines: 2;
  --featured-title-line-height: 1.24;
  --featured-summary-line-height: 1.78;
  border: 1px solid var(--border-soft);
  background:
    linear-gradient(180deg, color-mix(in srgb, var(--surface) 98%, white 2%), color-mix(in srgb, var(--surface) 94%, var(--bg) 6%));
  border-radius: 20px;
  height: 100%;
  box-shadow: 0 8px 20px rgba(20, 30, 40, 0.06);
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease, background 0.18s ease;
  cursor: pointer;
}

.featured-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 18px 34px rgba(20, 30, 40, 0.12);
  border-color: color-mix(in srgb, var(--border-soft) 58%, var(--accent) 42%);
}

.featured-link {
  display: block;
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
  text-decoration: none;
  min-height: 312px;
  cursor: pointer;
  border-radius: inherit;
  overflow: hidden;
}

.featured-link:focus-visible {
  outline: none;
}

.featured-link:focus-visible .featured-copy {
  border-color: color-mix(in srgb, var(--accent) 42%, var(--border-soft) 58%);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent) 16%, transparent);
}

.featured-card:hover .featured-title,
.featured-link:hover .featured-title {
  color: var(--link-hover);
}

.featured-copy {
  display: flex;
  flex-direction: column;
  gap: 14px;
  height: 100%;
  min-width: 0;
  padding: 2px;
  border-radius: 14px;
}

.featured-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 22px;
}

.featured-kicker {
  margin: 0;
  font-size: 12px;
  color: var(--text-soft);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-weight: 600;
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.featured-kicker::before {
  content: "";
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--accent) 82%, transparent);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--accent) 10%, transparent);
}

.featured-title-wrap {
  min-height: calc(2 * 1em * var(--featured-title-line-height) + 8px);
  padding-bottom: 4px;
}

.featured-title {
  margin: 0;
  color: var(--text);
  line-height: var(--featured-title-line-height);
  font-size: clamp(20px, 1.9vw, 24px);
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: var(--featured-title-lines);
  overflow: hidden;
  text-wrap: pretty;
}

.featured-summary-wrap {
  flex: 1 1 auto;
  min-height: 8.8em;
  display: flex;
  align-items: flex-start;
  overflow: hidden;
}

.featured-date {
  margin: 0;
  color: var(--text-faint);
  font-size: 12px;
  letter-spacing: 0.02em;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  flex-shrink: 0;
}

.featured-link :deep(.summary-preview) {
  font-size: 15px;
  line-height: var(--featured-summary-line-height);
  color: var(--text);
  width: 100%;
}

.featured-link :deep(.summary-preview[data-summary-type="fallback"]) {
  color: var(--text-soft);
}

.featured-link :deep(.summary-badge) {
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .featured-card {
    border-radius: 18px;
  }

  .featured-link {
    min-height: 292px;
    padding: 18px;
  }

  .featured-title {
    font-size: 19px;
  }

  .featured-copy {
    gap: 12px;
  }

  .featured-title-wrap {
    min-height: calc(2 * 1em * var(--featured-title-line-height) + 6px);
    padding-bottom: 2px;
  }

  .featured-summary-wrap {
    min-height: 8.6em;
  }

  .featured-link :deep(.summary-preview) {
    font-size: 14px;
    line-height: 1.72;
  }
}
</style>
