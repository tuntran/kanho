import { useEffect } from "react";
import { useAuthStore } from "@/stores/auth-store";
import * as authClient from "@/api/auth-client";

export function useSilentRefresh() {
  const { setAuth, setLoading, logout } = useAuthStore();

  useEffect(() => {
    let mounted = true;

    async function tryRefresh() {
      try {
        const res = await authClient.refresh();
        if (mounted) {
          setAuth(res.user, res.accessToken);
        }
      } catch {
        if (mounted) {
          logout();
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    }

    tryRefresh();

    return () => {
      mounted = false;
    };
  }, [setAuth, setLoading, logout]);
}
