<template>
  <article class="search-result-card">
    <router-link :to="{ name: 'Article', params: { noteId: item.noteId } }" class="search-result-link">
      <h2 class="search-result-title">
        <SearchHighlightText :text="item.title" :query="query" />
      </h2>
      <p class="search-result-date">{{ item.dateModified }}</p>
      <p v-if="preferredSnippet" class="search-result-snippet">
        <SearchHighlightText :text="preferredSnippet" :query="query" />
      </p>
    </router-link>
  </article>
</template>

<script>
import { computed, toRef } from "vue";
import SearchHighlightText from "./SearchHighlightText.vue";
import { useSummaryStatus } from "../../composables/useSummaryStatus";

export default {
  name: "SearchResultCard",
  components: {
    SearchHighlightText,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
    query: {
      type: String,
      default: "",
    },
  },
  setup(props) {
    const { preferredSummary } = useSummaryStatus(toRef(props, "item"));
    const preferredSnippet = computed(() => props.item.match?.snippet || preferredSummary.value.text);

    return {
      preferredSnippet,
    };
  },
};
</script>

<style scoped>
.search-result-card {
  border: 1px solid var(--border-soft);
  background: var(--surface);
  border-radius: var(--radius-md);
}

.search-result-link {
  display: block;
  padding: 18px 20px;
  text-decoration: none;
}

.search-result-title {
  margin: 0;
  color: var(--text);
  font-size: 20px;
  line-height: 1.4;
}

.search-result-date {
  margin: 10px 0 0;
  color: var(--text-faint);
  font-size: 13px;
}

.search-result-snippet {
  margin: 12px 0 0;
  color: var(--text-soft);
  line-height: 1.7;
}
</style>
