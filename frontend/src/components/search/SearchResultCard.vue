<template>
  <article class="search-result-card">
    <router-link :to="{ name: 'Article', params: { noteId: item.noteId } }" class="search-result-link">
      <h2 class="search-result-title">
        <SearchHighlightText :text="item.title" :query="query" />
      </h2>
      <p class="search-result-date">{{ item.dateModified }}</p>
      <p v-if="item.match?.snippet" class="search-result-snippet">
        <SearchHighlightText :text="item.match.snippet" :query="query" />
      </p>
      <p v-else-if="item.summary" class="search-result-snippet">
        {{ item.summary }}
      </p>
    </router-link>
  </article>
</template>

<script>
import SearchHighlightText from "./SearchHighlightText.vue";

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
