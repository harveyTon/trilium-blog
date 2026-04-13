import { api } from "./blog";

export async function fetchSearchResults(query, options = {}) {
  const response = await api.get("/search", {
    params: {
      q: query,
      ...options,
    },
  });
  return response.data;
}

export async function fetchSearchPreview(query) {
  const response = await api.get("/search", {
    params: {
      q: query,
      preview: true,
      limit: 5,
    },
  });
  return response.data;
}
