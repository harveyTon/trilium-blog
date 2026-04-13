<template>
  <el-empty :description="description" />
</template>

<script>
export default {
  name: "SearchEmptyState",
  props: {
    description: {
      type: String,
      required: true,
    },
  },
};
</script>
