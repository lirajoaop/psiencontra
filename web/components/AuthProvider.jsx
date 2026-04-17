"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { getMe, login as apiLogin, logout as apiLogout, register as apiRegister, claimSession } from "@/lib/api";
import { clearToken, getToken, setToken } from "@/lib/auth-token";

const PENDING_CLAIM_KEY = "psiencontra_pending_claim";

export function setPendingClaim(sessionId) {
  localStorage.setItem(PENDING_CLAIM_KEY, sessionId);
}

export function getPendingClaim() {
  return localStorage.getItem(PENDING_CLAIM_KEY);
}

export function clearPendingClaim() {
  localStorage.removeItem(PENDING_CLAIM_KEY);
}

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
      return null;
    }
    try {
      const me = await getMe();
      setUser(me);
      // Claim any pending anonymous session (Google OAuth lands here via
      // refresh() rather than login()).
      const sessionId = getPendingClaim();
      if (sessionId) {
        try { await claimSession(sessionId); } catch { /* best-effort */ }
        clearPendingClaim();
      }
      return sessionId;
    } catch {
      clearToken();
      setUser(null);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  // After login/register, if there's a pending anonymous session, claim it
  // so the result gets linked to the new account. Best-effort: if the claim
  // fails (already owned, network error) we still complete the login.
  const tryClaimPending = useCallback(async () => {
    const sessionId = getPendingClaim();
    if (!sessionId) return null;
    try {
      await claimSession(sessionId);
    } catch {
      // ignore — session may already be claimed or deleted
    }
    clearPendingClaim();
    return sessionId;
  }, []);

  const login = useCallback(async (credentials) => {
    const data = await apiLogin(credentials);
    if (data.token) setToken(data.token);
    setUser(data.user);
    await tryClaimPending();
    return data.user;
  }, [tryClaimPending]);

  const register = useCallback(async (credentials) => {
    const data = await apiRegister(credentials);
    if (data.token) setToken(data.token);
    setUser(data.user);
    await tryClaimPending();
    return data.user;
  }, [tryClaimPending]);

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
