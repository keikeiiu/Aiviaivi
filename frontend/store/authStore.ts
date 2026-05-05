import { create } from "zustand";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { auth as authApi, users as usersApi } from "../services/api";

interface User {
  id: string;
  username: string;
  email: string;
  avatar_url: string;
  bio: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;

  login: (email: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  restore: () => Promise<void>;
  refreshProfile: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: null,
  isLoading: true,
  isAuthenticated: false,

  login: async (email, password) => {
    const { data } = await authApi.login(email, password);
    const { user, token, refresh_token } = data.data;
    await AsyncStorage.multiSet([
      ["token", token],
      ["refresh_token", refresh_token],
      ["user", JSON.stringify(user)],
    ]);
    set({ user, token, isAuthenticated: true });
  },

  register: async (username, email, password) => {
    const { data } = await authApi.register(username, email, password);
    const { user, token, refresh_token } = data.data;
    await AsyncStorage.multiSet([
      ["token", token],
      ["refresh_token", refresh_token],
      ["user", JSON.stringify(user)],
    ]);
    set({ user, token, isAuthenticated: true });
  },

  logout: async () => {
    await AsyncStorage.multiRemove(["token", "refresh_token", "user"]);
    set({ user: null, token: null, isAuthenticated: false });
  },

  restore: async () => {
    try {
      const [[, token], [, userStr]] = await AsyncStorage.multiGet(["token", "user"]);
      if (token && userStr) {
        const user = JSON.parse(userStr);
        set({ user, token, isAuthenticated: true, isLoading: false });
        // Refresh profile in background
        get().refreshProfile();
        return;
      }
    } catch {}
    set({ isLoading: false });
  },

  refreshProfile: async () => {
    try {
      const { data } = await usersApi.me();
      const user = data.data;
      await AsyncStorage.setItem("user", JSON.stringify(user));
      set({ user });
    } catch {}
  },
}));
