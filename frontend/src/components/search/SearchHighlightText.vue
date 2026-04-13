<template>
  <span>
    <template v-for="(segment, index) in segments" :key="index">
      <mark v-if="segment.highlight" class="search-highlight-mark">{{ segment.text }}</mark>
      <span v-else>{{ segment.text }}</span>
    </template>
  </span>
</template>

<script>
import { computed } from "vue";

export default {
  name: "SearchHighlightText",
  props: {
    text: {
      type: String,
      default: "",
    },
    query: {
      type: String,
      default: "",
    },
  },
  setup(props) {
    const segments = computed(() => {
      if (!props.text || !props.query) {
        return [{ text: props.text, highlight: false }];
      }

      const lowerText = props.text.toLowerCase();
      const lowerQuery = props.query.toLowerCase();
      const parts = [];
      let cursor = 0;

      while (cursor < props.text.length) {
        const matchIndex = lowerText.indexOf(lowerQuery, cursor);
        if (matchIndex < 0) {
          parts.push({
            text: props.text.slice(cursor),
            highlight: false,
          });
          break;
        }

        if (matchIndex > cursor) {
          parts.push({
            text: props.text.slice(cursor, matchIndex),
            highlight: false,
          });
        }

        parts.push({
          text: props.text.slice(matchIndex, matchIndex + props.query.length),
          highlight: true,
        });

        cursor = matchIndex + props.query.length;
      }

      return parts.filter((part) => part.text);
    });

    return {
      segments,
    };
  },
};
</script>

<style scoped>
.search-highlight-mark {
  background: color-mix(in srgb, var(--accent) 24%, transparent);
  color: inherit;
  padding: 0 0.08em;
  border-radius: 4px;
}
</style>
