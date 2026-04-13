import axios from "axios";

export const api = axios.create({
  baseURL: "/api",
  timeout: 30000,
});

api.interceptors.response.use(
  response => response,
  error => {
    console.error("API Error:", error);
    return Promise.reject(error);
  }
);

export async function fetchSite() {
  const response = await api.get("/site");
  return response.data;
}

export async function fetchPosts(page = 1) {
  const response = await api.get("/posts", {
    params: { page },
  });
  return response.data;
}

export async function fetchPost(noteId) {
  const response = await api.get(`/posts/${noteId}`);
  return response.data;
}

export async function fetchAsset(attachmentId) {
  const response = await api.get(`/assets/${attachmentId}`, {
    responseType: "arraybuffer",
  });
  return response.data;
}
