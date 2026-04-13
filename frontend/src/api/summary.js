export function normalizeSummaryPayload(post) {
  if (!post) {
    return {
      ai: null,
      code: null,
      fallback: "",
    };
  }

  if (post.summaries && typeof post.summaries === "object") {
    return {
      ai: post.summaries.ai || null,
      code: post.summaries.code || null,
      fallback: post.summaries.fallback || post.summary || "",
    };
  }

  if (typeof post.summary === "string") {
    return {
      ai: null,
      code: post.summary ? { type: "code", text: post.summary, status: "ready" } : null,
      fallback: post.summary || "",
    };
  }

  if (post.summary && typeof post.summary === "object") {
    return {
      ai: post.summary.ai || null,
      code: post.summary.code || null,
      fallback: post.summary.fallback || "",
    };
  }

  return {
    ai: null,
    code: null,
    fallback: "",
  };
}
