import { useEffect, useCallback } from "react";
import { useAuthStore } from "../store/authStore";

export function useAuth() {
  const store = useAuthStore();

  useEffect(() => {
    store.restore();
  }, []);

  return {
    user: store.user,
    token: store.token,
    isLoading: store.isLoading,
    isAuthenticated: store.isAuthenticated,
    login: useCallback((e: string, p: string) => store.login(e, p), []),
    register: useCallback((u: string, e: string, p: string) => store.register(u, e, p), []),
    logout: useCallback(() => store.logout(), []),
  };
}
