"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import { useAuth } from "@/components/AuthProvider";
import { setToken } from "@/lib/auth-token";

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

  useEffect(() => {
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

    if (!token) {
      setError(ERROR_MESSAGES.no_token);
      return;
    }

    setToken(token);
    // Strip the fragment from the URL so the token doesn't sit in the address
    // bar / browser history. replaceState avoids a navigation event.
    window.history.replaceState(null, "", window.location.pathname);

    refresh().then(() => router.push("/"));
  }, [params, refresh, router]);

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4 px-6">
        <p className="text-red-600 dark:text-red-400 text-center">{error}</p>
        <Button onClick={() => router.push("/login")} variant="secondary">
          Voltar para login
        </Button>
      </div>
    );
  }

  return <LoadingSpinner message="Finalizando login..." />;
}
