import type { AuthProvider } from "@refinedev/core";
import { authApi, unwrap } from "../api/client";

// JWT auth provider. Login goes through the generated AuthApi
// (POST /api/v1/login, swagger-documented); everything else carries the Bearer
// token via the shared axios instance (see api/client.ts).
export const authProvider: AuthProvider = {
  login: async (payload: any) => {
    try {
      const data = await unwrap<any>(
        authApi.loginPost({
          body: {
            username: payload?.username,
            password: payload?.password,
            code: payload?.code,
            uuid: payload?.uuid,
          },
        })
      );
      localStorage.setItem("token", data.token);
      localStorage.setItem("user", JSON.stringify(data.user));
      // Permission strings derived from the role's granted F-menu nodes,
      // "*" means everything (admin).
      localStorage.setItem("perms", JSON.stringify(data.perms ?? []));
      return { success: true };
    } catch (e: any) {
      return {
        success: false,
        error: { message: e?.message || "登录失败", name: "LoginError" },
      };
    }
  },
  logout: async () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    localStorage.removeItem("perms");
    return { success: true };
  },
  check: async () => {
    return localStorage.getItem("token")
      ? { authenticated: true }
      : { authenticated: false, redirectTo: "/login" };
  },
  getPermissions: async () => {
    return JSON.parse(localStorage.getItem("perms") || "[]");
  },
  getIdentity: async () => {
    const u = JSON.parse(localStorage.getItem("user") || "{}");
    return { id: u.id, name: u.username, avatar: u.avatar };
  },
  onError: async (error: any) => {
    if (error?.statusCode === 401 || error?.code === 401 || error?.response?.status === 401) {
      return { logout: true, redirectTo: "/login" };
    }
    return { error };
  },
};
