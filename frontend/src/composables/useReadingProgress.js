import { onMounted, onUnmounted, ref } from "vue";

export function useReadingProgress() {
  const progress = ref(0);

  const update = () => {
    const doc = document.documentElement;
    const scrollable = doc.scrollHeight - window.innerHeight;
    if (scrollable <= 0) {
      progress.value = 0;
      return;
    }
    progress.value = Math.min(100, Math.max(0, (window.scrollY / scrollable) * 100));
  };

  onMounted(() => {
    update();
    window.addEventListener("scroll", update, { passive: true });
    window.addEventListener("resize", update);
  });

  onUnmounted(() => {
    window.removeEventListener("scroll", update);
    window.removeEventListener("resize", update);
  });

  return {
    progress,
    updateReadingProgress: update,
  };
}
