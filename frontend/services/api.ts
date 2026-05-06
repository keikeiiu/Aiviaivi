import axios, { AxiosError } from "axios";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { API_BASE_URL, PAGE_SIZE } from "../utils/constants";

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
  headers: { "Content-Type": "application/json" },
});

// Attach token to requests
api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auto-refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      const refreshToken = await AsyncStorage.getItem("refresh_token");
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          });
          await AsyncStorage.setItem("token", data.data.token);
          if (error.config) {
            error.config.headers.Authorization = `Bearer ${data.data.token}`;
            return axios(error.config);
          }
        } catch {
          await AsyncStorage.multiRemove(["token", "refresh_token", "user"]);
        }
      }
    }
    return Promise.reject(error);
  }
);

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  pagination?: { page: number; size: number; total: number };
}

// --- Auth ---
export const auth = {
  register: (username: string, email: string, password: string) =>
    api.post("/auth/register", { username, email, password }),
  login: (email: string, password: string) =>
    api.post("/auth/login", { email, password }),
  refresh: (refreshToken: string) =>
    api.post("/auth/refresh", { refresh_token: refreshToken }),
};

// --- Videos ---
export const videos = {
  list: (params?: { page?: number; size?: number; sort?: string; category?: number }) =>
    api.get("/videos", { params: { size: PAGE_SIZE, ...params } }),
  detail: (id: string) => api.get(`/videos/${id}`),
  related: (id: string) => api.get(`/videos/${id}/related`),
  upload: (formData: FormData) =>
    api.post("/videos/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
      transformRequest: [(data: any) => data],
    }),
  update: (id: string, data: { title?: string; description?: string; tags?: string[] }) =>
    api.put(`/videos/${id}`, data),
  delete: (id: string) => api.delete(`/videos/${id}`),
};

// --- Danmaku ---
export const danmaku = {
  list: (videoId: string, tStart: number, tEnd: number) =>
    api.get(`/videos/${videoId}/danmaku`, { params: { t_start: tStart, t_end: tEnd } }),
  send: (videoId: string, data: { content: string; video_time: number; color?: string; font_size?: string; mode?: string }) =>
    api.post(`/videos/${videoId}/danmaku`, data),
};

// --- Comments ---
export const comments = {
  list: (videoId: string, params?: { page?: number; size?: number; sort?: string }) =>
    api.get(`/videos/${videoId}/comments`, { params: { size: PAGE_SIZE, ...params } }),
  create: (videoId: string, content: string, parentId?: string) =>
    api.post(`/videos/${videoId}/comments`, { content, parent_id: parentId }),
  delete: (commentId: string) => api.delete(`/comments/${commentId}`),
};

// --- Social ---
export const social = {
  like: (videoId: string) => api.post(`/videos/${videoId}/like`),
  unlike: (videoId: string) => api.delete(`/videos/${videoId}/like`),
  favorite: (videoId: string) => api.post(`/videos/${videoId}/favorite`),
  unfavorite: (videoId: string) => api.delete(`/videos/${videoId}/favorite`),
  subscribe: (userId: string) => api.post(`/users/${userId}/subscribe`),
  unsubscribe: (userId: string) => api.delete(`/users/${userId}/subscribe`),
};

// --- Users ---
export const users = {
  me: () => api.get("/users/me"),
  profile: (id: string) => api.get(`/users/${id}`),
  update: (id: string, data: { bio?: string; avatar_url?: string }) =>
    api.put(`/users/${id}`, data),
  videos: (id: string, params?: { page?: number }) =>
    api.get(`/users/${id}/videos`, { params: { size: PAGE_SIZE, ...params } }),
  favorites: (id: string, params?: { page?: number }) =>
    api.get(`/users/${id}/favorites`, { params: { size: PAGE_SIZE, ...params } }),
  history: (params?: { page?: number }) =>
    api.get("/users/me/history", { params: { size: PAGE_SIZE, ...params } }),
};

// --- Other ---
export const search = (q: string, type?: "video" | "user", params?: { page?: number }) =>
  api.get("/search", { params: { q, type, size: PAGE_SIZE, ...params } });

export const trending = (params?: { page?: number }) =>
  api.get("/feed/trending", { params: { size: PAGE_SIZE, ...params } });

export const categories = () => api.get("/categories");

export const playlists = {
  list: () => api.get("/playlists"),
  create: (data: { name: string; description?: string; is_public?: boolean }) =>
    api.post("/playlists", data),
  detail: (id: string) => api.get(`/playlists/${id}`),
  update: (id: string, data: { name: string; description?: string; is_public?: boolean }) =>
    api.put(`/playlists/${id}`, data),
  delete: (id: string) => api.delete(`/playlists/${id}`),
  addVideo: (id: string, videoId: string) =>
    api.post(`/playlists/${id}/videos`, { video_id: videoId }),
  removeVideo: (id: string, videoId: string) =>
    api.delete(`/playlists/${id}/videos/${videoId}`),
};

export const analytics = {
  overview: () => api.get("/analytics/overview"),
  videos: () => api.get("/analytics/videos"),
  videoDetail: (id: string) => api.get(`/analytics/videos/${id}`),
};

export const watch = {
  record: (videoId: string, progress: number, duration: number) =>
    api.post(`/videos/${videoId}/watch`, { progress, duration }),
};

export default api;
