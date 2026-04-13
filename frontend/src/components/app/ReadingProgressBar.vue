<template>
  <div class="reading-progress" :style="progressStyle" aria-hidden="true">
    <div class="reading-progress__bar" :style="{ width: `${progress}%` }"></div>
  </div>
</template>

<script>
export default {
  name: "ReadingProgressBar",
  props: {
    progress: {
      type: Number,
      default: 0,
    },
    topOffset: {
      type: Number,
      default: null,
    },
  },
  computed: {
    progressStyle() {
      if (typeof this.topOffset !== "number") {
        return null;
      }
      return { top: `${this.topOffset}px` };
    },
  },
};
</script>

<style scoped>
.reading-progress {
  position: fixed;
  top: var(--header-h);
  left: 0;
  width: 100%;
  height: 2px;
  background: transparent;
  z-index: 1100;
}

.reading-progress__bar {
  height: 100%;
  background: var(--accent);
  transition: width 120ms linear;
}
</style>
