"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import { useAuth } from "@/components/AuthProvider";

const ERROR_MESSAGES = {
  invalid_state: "Sessão de login expirada ou inválida. Tente novamente.",
  exchange_failed: "Não foi possível concluir o login com o Google.",
  user_upsert_failed: "Erro ao criar sua conta. Tente novamente.",
  google_disabled: "Login com Google não está configurado.",
  missing_code: "Resposta do Google incompleta.",
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
