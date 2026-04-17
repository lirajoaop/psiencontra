"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import { useAuth, getPendingClaim } from "@/components/AuthProvider";
import { getToken, setToken } from "@/lib/auth-token";

const ERROR_MESSAGES = {
  invalid_state: "Sessão de login expirada ou inválida. Tente novamente.",
  exchange_failed: "Não foi possível concluir o login com o Google.",
  user_upsert_failed: "Erro ao criar sua conta. Tente novamente.",
  google_disabled: "Login com Google não está configurado.",
  missing_code: "Resposta do Google incompleta.",
  state_generation_failed: "Não foi possível iniciar o login. Tente novamente.",
  no_token: "Resposta de autenticação incompleta. Tente novamente.",
};

export default function AuthCallback() {
  const router = useRouter();
  const params = useSearchParams();
  const { refresh } = useAuth();
  const [error, setError] = useState(null);
  // React 19 + Strict Mode runs effects twice on mount in development. Our
  // effect mutates global state (URL fragment, localStorage), so without
  // this guard the second run sees an already-cleaned hash, falls into the
  // "no_token" branch, and flashes an error before the first run's
  // redirect lands. The ref survives across the double-invoke.
  const handled = useRef(false);

  useEffect(() => {
    if (handled.current) return;
    handled.current = true;

    const errCode = params.get("error");
    if (errCode) {
      setError(ERROR_MESSAGES[errCode] || "Falha ao entrar com Google.");
      return;
    }

    // The backend returns the JWT in the URL fragment (#token=...). Hash
    // fragments aren't sent to servers, so the token never lands in any
    // access log along the redirect chain — only this browser sees it.
    const hash = window.location.hash.startsWith("#")
      ? window.location.hash.slice(1)
      : window.location.hash;
    const hashParams = new URLSearchParams(hash);
    const token = hashParams.get("token");

    // Capture before refresh() clears it during the claim flow.
    const pendingSession = getPendingClaim();
    const redirectAfterAuth = () => {
      router.push(pendingSession ? `/resultado/${pendingSession}` : "/");
    };

    // If there's no token in the fragment, we may already be authenticated
    // from a prior tab (Strict Mode double-mount, browser back/forward, etc).
    // Trust the existing localStorage token instead of flashing an error.
    if (!token) {
      if (getToken()) {
        refresh().then(redirectAfterAuth);
        return;
      }
      setError(ERROR_MESSAGES.no_token);
      return;
    }

    setToken(token);
    // Strip the fragment from the URL so the token doesn't sit in the address
    // bar / browser history. replaceState avoids a navigation event.
    window.history.replaceState(null, "", window.location.pathname);

    refresh().then(redirectAfterAuth);
  }, [params, refresh, router]);

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4 px-6">
        <p className="text-red-600 dark:text-red-400 text-center">{error}</p>
        <Button onClick={() => router.push("/entrar")} variant="secondary">
          Voltar para login
        </Button>
      </div>
    );
  }

  return <LoadingSpinner message="Finalizando login..." />;
}
