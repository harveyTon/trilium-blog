<template>
  <article class="post-item">
    <div class="post-date">
      <div class="post-day">{{ day }}</div>
      <div class="post-month">{{ month }}</div>
    </div>
    <div class="post-body">
      <h2 class="post-title">
        <router-link :to="{ name: 'Article', params: { noteId: post.noteId } }">{{ post.title }}</router-link>
      </h2>
      <div class="post-meta">{{ fullDate }}</div>
      <SummaryPreview :content="resolvedSummary" :type="summaryType" hide-badge />
    </div>
  </article>
</template>

<script>
import { computed, toRef } from "vue";
import SummaryPreview from "./SummaryPreview.vue";
import { useSummaryStatus } from "../../composables/useSummaryStatus";

export default {
  name: "PostCard",
  components: {
    SummaryPreview,
  },
  props: {
    post: {
      type: Object,
      required: true,
    },
    day: {
      type: String,
      required: true,
    },
    month: {
      type: String,
      required: true,
    },
    fullDate: {
      type: String,
      required: true,
    },
    summary: {
      type: String,
      default: "",
    },
  },
  setup(props) {
    const { preferredSummary } = useSummaryStatus(toRef(props, "post"));
    const resolvedSummary = computed(() => preferredSummary.value.text || props.summary);
    const summaryType = computed(() => preferredSummary.value.type);

    return {
      resolvedSummary,
      summaryType,
    };
  },
};
</script>

<style scoped>
.post-item {
  display: flex;
  align-items: flex-start;
  gap: 28px;
  padding: 18px 0 22px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  transition: transform 0.18s ease, border-color 0.18s ease;
}

.post-date {
  width: 92px;
  min-width: 92px;
  height: 92px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-soft);
  background: var(--surface-muted);
  flex-shrink: 0;
}

.post-day {
  font-size: 52px;
  line-height: 1;
  font-weight: 600;
  color: var(--text);
  font-variant-numeric: tabular-nums;
}

.post-month {
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-faint);
  line-height: 1.2;
}

.post-body {
  flex: 1;
  min-width: 0;
}

.post-title {
  margin: 0;
  font-size: 22px;
  line-height: 1.32;
  font-weight: 700;
}

.post-title a {
  color: var(--text);
  text-decoration: none;
  position: relative;
  display: inline-block;
  transition: color 0.18s ease;
}

.post-title a::after {
  content: "";
  display: block;
  width: 0;
  max-width: 100%;
  height: 2px;
  margin-top: 12px;
  background: color-mix(in srgb, var(--accent) 88%, transparent);
  transition: width 0.18s ease, background 0.18s ease;
}

.post-item:hover {
  border-color: color-mix(in srgb, var(--border) 72%, var(--accent) 28%);
}

.post-item:hover .post-title a,
.post-title a:hover {
  color: var(--link-hover);
}

.post-item:hover .post-title a::after {
  width: min(180px, 100%);
}

.post-title a:focus-visible {
  outline: none;
  color: var(--link-hover);
}

.post-title a:focus-visible::after {
  width: min(180px, 100%);
}

.post-meta {
  display: none;
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-faint);
  line-height: 1.4;
}

html.dark .post-item {
  border-color: var(--border-soft);
}

@media (max-width: 768px) {
  html.dark .post-item {
    border-color: var(--border);
  }
}

html.dark .post-date {
  border-color: var(--border);
  background: var(--surface-muted);
}

@media (max-width: 768px) {
  .post-item {
    display: block;
    padding: 16px;
    margin-bottom: 12px;
    border: 1px solid var(--border-soft);
    border-radius: 12px;
    background: var(--surface-muted);
  }

  .post-date {
    display: none;
  }

  .post-title {
    font-size: 20px;
    line-height: 1.45;
  }

  .post-title a::after {
    margin-top: 12px;
  }

  .post-item:hover .post-title a::after,
  .post-title a:focus-visible::after {
    width: min(132px, 100%);
  }

  .post-meta {
    display: block;
  }
}
</style>
