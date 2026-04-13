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
      <SummaryPreview :content="resolvedSummary" :type="summaryType" />
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
  font-size: 20px;
  line-height: 1.4;
  font-weight: 700;
}

.post-title a {
  color: var(--text);
  text-decoration: none;
  position: relative;
  display: inline-block;
}

.post-title a::after {
  content: "";
  display: block;
  width: 120px;
  height: 1px;
  margin-top: 14px;
  background: linear-gradient(
    to right,
    var(--border) 0%,
    var(--border) 35%,
    var(--border) 65%,
    transparent 100%
  );
  transition: background 0.2s ease;
}

.post-item:hover .post-title a::after {
  background: linear-gradient(
    to right,
    var(--accent) 0%,
    var(--accent) 35%,
    var(--accent) 65%,
    transparent 100%
  );
}

.post-meta {
  display: none;
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-faint);
  line-height: 1.4;
}

html.dark .post-item {
  border-color: rgba(255, 255, 255, 0.04);
}

html.dark .post-date {
  border-color: var(--border);
  background: var(--surface-muted);
}

@media (max-width: 768px) {
  .post-item {
    display: block;
    padding: 18px 0;
  }

  .post-date {
    display: none;
  }

  .post-title {
    font-size: 18px;
    line-height: 1.45;
  }

  .post-title a::after {
    width: 96px;
    margin-top: 12px;
  }

  .post-meta {
    display: block;
  }
}
</style>
