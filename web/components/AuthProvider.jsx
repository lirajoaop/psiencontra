"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { getMe, login as apiLogin, logout as apiLogout, register as apiRegister } from "@/lib/api";
import { clearToken, getToken, setToken } from "@/lib/auth-token";

const AuthContext = createContext(null);

export default function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    // No token = definitely anonymous. Skip the network round-trip so the
    // landing page renders the "Entrar" button instantly for first-time
    // visitors instead of waiting on a cold-started backend.
    if (!getToken()) {
      setUser(null);
      setLoading(false);
      return;
    }
    try {
      const me = await getMe();
      setUser(me);
    } catch {
      // Token is present but rejected by the server (expired, signed with a
      // different secret, deleted user, etc.). Drop it so we don't keep
      // sending an invalid Authorization header on every subsequent request.
      clearToken();
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const login = useCallback(async (credentials) => {
    const data = await apiLogin(credentials);
    if (data.token) setToken(data.token);
    setUser(data.user);
    return data.user;
  }, []);

  const register = useCallback(async (credentials) => {
    const data = await apiRegister(credentials);
    if (data.token) setToken(data.token);
    setUser(data.user);
    return data.user;
  }, []);

  const logout = useCallback(async () => {
    // Best-effort: ask the backend to clear its cookie too. If the network
    // request fails (offline, backend down) we still want to log the user out
    // locally, so we swallow the error.
    try {
      await apiLogout();
    } catch {
      // ignore
    }
    clearToken();
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout, refresh }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside <AuthProvider>");
  }
  return ctx;
}
