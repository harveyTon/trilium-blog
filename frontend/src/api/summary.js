import { api } from "./blog";

export async function fetchPostSummary(noteId) {
  const response = await api.get(`/posts/${noteId}/summary`);
  return response.data;
}

export function normalizeSummaryPayload(post) {
  if (!post) {
    return {
      noteId: "",
      aiEnabled: false,
      ai: null,
      code: null,
      fallback: "",
    };
  }

  if (post.summaries && typeof post.summaries === "object") {
    return {
      noteId: post.summaries.noteId || post.noteId || "",
      aiEnabled: Boolean(post.summaries.aiEnabled),
      ai: post.summaries.ai || null,
      code: post.summaries.code || null,
      fallback: post.summary || post.summaries.code?.text || "",
    };
  }

  if (typeof post.summary === "string") {
    return {
      noteId: post.noteId || "",
      aiEnabled: false,
      ai: null,
      code: post.summary ? { type: "code", text: post.summary, status: "ready" } : null,
      fallback: post.summary || "",
    };
  }

  if (post.summary && typeof post.summary === "object") {
    return {
      noteId: post.summary.noteId || post.noteId || "",
      aiEnabled: Boolean(post.summary.aiEnabled),
      ai: post.summary.ai || null,
      code: post.summary.code || null,
      fallback: post.summary.code?.text || "",
    };
  }

  return {
    noteId: post.noteId || "",
    aiEnabled: false,
    ai: null,
    code: null,
    fallback: "",
  };
}
