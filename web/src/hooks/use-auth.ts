import { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "@/stores/auth-store";
import * as authClient from "@/api/auth-client";

export function useAuth() {
  const { user, accessToken, isLoading, setAuth, logout: clearAuth } = useAuthStore();
  const navigate = useNavigate();

  const login = useCallback(
    async (email: string, password: string) => {
      const res = await authClient.login(email, password);
      setAuth(res.user, res.accessToken);
      navigate("/");
    },
    [setAuth, navigate],
  );

  const register = useCallback(
    async (email: string, password: string, name: string) => {
      const res = await authClient.register(email, password, name);
      setAuth(res.user, res.accessToken);
      navigate("/");
    },
    [setAuth, navigate],
  );

  const logout = useCallback(async () => {
    await authClient.logout().catch(() => {});
    clearAuth();
    navigate("/login");
  }, [clearAuth, navigate]);

  return { user, accessToken, isLoading, login, register, logout };
}
