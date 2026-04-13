<template>
  <nav class="post-list" aria-label="文章列表">
    <PostCard
      v-for="post in items"
      :key="post.noteId"
      :post="post"
      :day="getDay(post.dateModified)"
      :month="getMonth(post.dateModified)"
      :full-date="formatFullDate(post.dateModified)"
      :summary="sanitizeSummary(post.summary)"
    />
  </nav>
</template>

<script>
import PostCard from "./PostCard.vue";

export default {
  name: "PostFeed",
  components: {
    PostCard,
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    getDay: {
      type: Function,
      required: true,
    },
    getMonth: {
      type: Function,
      required: true,
    },
    formatFullDate: {
      type: Function,
      required: true,
    },
    sanitizeSummary: {
      type: Function,
      required: true,
    },
  },
};
</script>

<style scoped>
.post-list {
  display: flex;
  flex-direction: column;
}

.post-list:focus-within {
  outline: none;
}
</style>
