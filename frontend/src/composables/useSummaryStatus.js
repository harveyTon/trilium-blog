import { computed } from "vue";
import { normalizeSummaryPayload } from "../api/summary";

export function useSummaryStatus(summarySource) {
  const summaryState = computed(() => normalizeSummaryPayload(summarySource.value));

  const preferredSummary = computed(() => {
    if (summaryState.value.ai?.status === "ready" && summaryState.value.ai?.text) {
      return { type: "ai", text: summaryState.value.ai.text };
    }
    if (summaryState.value.code?.status === "ready" && summaryState.value.code?.text) {
      return { type: "code", text: summaryState.value.code.text };
    }
    if (summaryState.value.fallback) {
      return { type: "fallback", text: summaryState.value.fallback };
    }
    return { type: "", text: "" };
  });

  const aiSummary = computed(() => summaryState.value.ai || null);

  return {
    summaryState,
    preferredSummary,
    aiSummary,
  };
}
