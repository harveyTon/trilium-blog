<template>
  <p
    v-if="content"
    class="summary-preview"
    :class="[{ 'summary-preview--flush': flush, 'summary-preview--indent': indent, 'summary-preview--card': card }]"
    :style="previewStyle"
    :data-summary-type="type || 'fallback'"
  >
    <span v-if="type === 'ai' && !hideBadge" class="summary-badge">{{ badgeText }}</span>
    {{ content }}
  </p>
</template>

<script>
import { computed } from "vue";

export default {
  name: "SummaryPreview",
  props: {
    content: {
      type: String,
      default: "",
    },
    type: {
      type: String,
      default: "",
    },
    lines: {
      type: Number,
      default: 2,
    },
    flush: {
      type: Boolean,
      default: false,
    },
    tone: {
      type: String,
      default: "default",
    },
    hideBadge: {
      type: Boolean,
      default: false,
    },
    badgeText: {
      type: String,
      default: "AI",
    },
    indent: {
      type: Boolean,
      default: false,
    },
    card: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const previewStyle = computed(() => ({
      WebkitLineClamp: String(props.lines),
    }));

    return {
      previewStyle,
    };
  },
};
</script>

<style scoped>
.summary-badge {
  display: inline-block;
  margin-right: 6px;
  padding: 1px 6px 2px;
  vertical-align: 1px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--accent) 12%, transparent);
  border: 1px solid color-mix(in srgb, var(--accent) 28%, transparent);
  color: white;
  color: var(--accent);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  line-height: 1.1;
}

.summary-preview {
  margin: 10px 0 0;
  font-size: 15px;
  line-height: 1.75;
  color: var(--text-soft);
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.summary-preview--flush {
  margin: 0;
}

.summary-preview--card {
  margin: 8px 0 0;
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-faint);
}

.summary-preview--indent {
  text-indent: 2em;
}

@media (max-width: 768px) {
  .summary-badge {
    font-size: 10px;
    margin-right: 5px;
  }

  .summary-preview {
    font-size: 14px;
    line-height: 1.7;
  }

  .summary-preview--card {
    font-size: 13px;
    line-height: 1.65;
  }
}
</style>
