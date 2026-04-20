"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import RadarChartResult from "@/components/RadarChartResult";
import ApproachCard from "@/components/ApproachCard";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import ThemeToggle from "@/components/ThemeToggle";
import ShareButtons from "@/components/ShareButtons";
import Link from "next/link";
import { getResult, getPDFUrl } from "@/lib/api";
import { useAuth, setPendingClaim } from "@/components/AuthProvider";
import {
  APPROACH_LABELS,
  APPROACH_COLORS,
  APPROACH_AUTHORS,
  FIELD_LABELS,
  FIELD_COLORS,
} from "@/lib/constants";

export default function ResultadoClient({ id, initialResult = null }) {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();
  const [result, setResult] = useState(initialResult);
  const [loading, setLoading] = useState(!initialResult);
  const [error, setError] = useState(null);
  const [bannerDismissed, setBannerDismissed] = useState(false);
  const [shareUrl, setShareUrl] = useState("");

  useEffect(() => {
    if (typeof window !== "undefined") {
      setShareUrl(window.location.href);
    }
  }, []);

  useEffect(() => {
    if (initialResult) return;
    getResult(id)
      .then((data) => {
        setResult(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id, initialResult]);

  if (loading) return <LoadingSpinner message="Carregando resultado..." />;

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4 px-6">
        <p className="text-red-600 dark:text-red-400">{error}</p>
        <Button onClick={() => router.push("/")} variant="secondary">Voltar ao início</Button>
      </div>
    );
  }

  if (!result) return null;

  const approachScores = typeof result.approach_scores === "string"
    ? JSON.parse(result.approach_scores)
    : result.approach_scores;

  const fieldScores = typeof result.field_scores === "string"
    ? JSON.parse(result.field_scores)
    : result.field_scores;

  const approachDetails = typeof result.approach_details === "string"
    ? JSON.parse(result.approach_details)
    : result.approach_details;

  const fieldDetails = typeof result.field_details === "string"
    ? JSON.parse(result.field_details)
    : result.field_details;

  const sortedApproaches = Object.entries(approachScores)
    .sort(([, a], [, b]) => b - a)
    .map(([key, score], i) => ({
      key,
      score,
      label: APPROACH_LABELS[key] || key,
      color: APPROACH_COLORS[key] || "#7c3aed",
      description: approachDetails?.[key]?.description || "",
      rank: i + 1,
    }));

  const sortedFields = Object.entries(fieldScores)
    .sort(([, a], [, b]) => b - a)
    .map(([key, score], i) => ({
      key,
      score,
      label: FIELD_LABELS[key] || key,
      color: FIELD_COLORS[key] || "#7c3aed",
      description: fieldDetails?.[key]?.description || "",
      rank: i + 1,
    }));

  return (
    <main className="flex-1 pb-12">
      {/* Header */}
      <div className="bg-gradient-to-r from-violet-600 to-purple-700 text-white py-10 px-6 relative">
        <div className="absolute top-4 right-4 z-20">
          <ThemeToggle />
        </div>
        <div className="max-w-5xl mx-auto text-center">
          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-3xl md:text-4xl font-bold mb-3"
          >
            Seu Resultado
          </motion.h1>
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="text-violet-100 max-w-2xl mx-auto"
          >
            {result.explanation}
          </motion.p>
        </div>
      </div>

      <div className="max-w-5xl mx-auto px-6 mt-8">
        {/* Save results banner — shown only for anonymous users */}
        {!authLoading && !user && !bannerDismissed && (
          <motion.section
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            className="rounded-2xl border border-violet-200 dark:border-violet-800 bg-violet-50 dark:bg-violet-950/40 p-6 md:p-8 mb-8 text-center"
          >
            <h3 className="text-lg font-bold text-violet-900 dark:text-violet-200 mb-2">
              Salve seus resultados
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-5 max-w-lg mx-auto">
              Crie uma conta ou faça login para guardar este resultado no seu histórico
              e acessá-lo a qualquer momento.
            </p>
            <div className="flex flex-wrap gap-3 justify-center">
              <Link
                href={`/cadastro?redirect=/resultado/${id}`}
                onClick={() => setPendingClaim(id)}
                className="px-6 py-2.5 rounded-xl font-semibold bg-violet-600 text-white hover:bg-violet-700 shadow-lg shadow-violet-200 dark:shadow-violet-900/50 transition-all text-sm text-center"
              >
                Criar conta
              </Link>
              <Link
                href={`/entrar?redirect=/resultado/${id}`}
                onClick={() => setPendingClaim(id)}
                className="px-6 py-2.5 rounded-xl font-semibold border-2 border-violet-300 dark:border-violet-700 text-violet-700 dark:text-violet-300 hover:border-violet-500 dark:hover:border-violet-500 transition-colors text-sm text-center"
              >
                Já tenho conta
              </Link>
              <button
                type="button"
                onClick={() => setBannerDismissed(true)}
                className="px-5 py-2.5 rounded-xl text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors cursor-pointer"
              >
                Continuar sem conta
              </button>
            </div>
          </motion.section>
        )}

        {/* Radar Charts */}
        <div className="grid md:grid-cols-2 gap-6 mb-10">
          <RadarChartResult
            title="Abordagens Teóricas"
            scores={approachScores}
            labels={APPROACH_LABELS}
            color="#7c3aed"
          />
          <RadarChartResult
            title="Campos de Atuação"
            scores={fieldScores}
            labels={FIELD_LABELS}
            color="#2563eb"
          />
        </div>

        {/* Approach Cards */}
        <h2 className="text-xl font-bold text-violet-900 dark:text-violet-200 mb-4">Ranking de Abordagens</h2>
        <div className="grid gap-4 mb-10">
          {sortedApproaches.map((a) => (
            <ApproachCard
              key={a.key}
              name={a.label}
              score={a.score}
              description={a.description}
              color={a.color}
              rank={a.rank}
            />
          ))}
        </div>

        {/* Field Cards */}
        <h2 className="text-xl font-bold text-violet-900 dark:text-violet-200 mb-4">Ranking de Campos de Atuação</h2>
        <div className="grid gap-4 mb-10">
          {sortedFields.map((f) => (
            <ApproachCard
              key={f.key}
              name={f.label}
              score={f.score}
              description={f.description}
              color={f.color}
              rank={f.rank}
            />
          ))}
        </div>

        {/* Next Steps */}
        <section className="rounded-2xl border border-violet-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 md:p-8 mb-10">
          <h2 className="text-xl font-bold text-violet-900 dark:text-violet-200 mb-3">
            Próximos passos
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-5">
            Este resultado é um ponto de partida para reflexão — não é um teste
            psicológico nem um veredito sobre sua carreira. Algumas formas de seguir
            explorando:
          </p>
          <ul className="space-y-3 text-sm text-gray-700 dark:text-gray-300">
            <li className="flex gap-3">
              <span className="text-violet-600 dark:text-violet-400 font-bold">1.</span>
              <span>
                Converse com a <strong>coordenação do seu curso</strong> ou com um(a)
                supervisor(a) sobre estes resultados — eles ajudam a situá-los frente à
                sua trajetória.
              </span>
            </li>
            <li className="flex gap-3">
              <span className="text-violet-600 dark:text-violet-400 font-bold">2.</span>
              <span>
                Busque <strong>estágios, eletivas ou grupos de estudo</strong> em{" "}
                <strong>{sortedApproaches[0]?.label}</strong> e em{" "}
                <strong>{sortedFields[0]?.label}</strong>, suas áreas de maior afinidade
                neste questionário.
              </span>
            </li>
            {APPROACH_AUTHORS[sortedApproaches[0]?.key] && (
              <li className="flex gap-3">
                <span className="text-violet-600 dark:text-violet-400 font-bold">
                  3.
                </span>
                <span>
                  Comece pelas leituras-fundadoras de{" "}
                  <strong>{sortedApproaches[0].label}</strong>:{" "}
                  {APPROACH_AUTHORS[sortedApproaches[0].key]}.
                </span>
              </li>
            )}
            <li className="flex gap-3">
              <span className="text-violet-600 dark:text-violet-400 font-bold">
                {APPROACH_AUTHORS[sortedApproaches[0]?.key] ? "4." : "3."}
              </span>
              <span>
                Refaça o questionário <strong>daqui a alguns meses</strong>: interesses
                tendem a se reorganizar com novas experiências práticas e leituras.
              </span>
            </li>
          </ul>
        </section>

        {/* Share */}
        {shareUrl && (
          <ShareButtons
            url={shareUrl}
            topApproach={sortedApproaches[0]?.label}
            topField={sortedFields[0]?.label}
          />
        )}

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mt-8">
          <a
            href={getPDFUrl(id)}
            target="_blank"
            rel="noopener noreferrer"
            className="px-6 py-3 rounded-xl font-semibold bg-violet-600 text-white hover:bg-violet-700 shadow-lg shadow-violet-200 dark:shadow-violet-900/50 transition-all text-center"
          >
            Baixar PDF
          </a>
          <Button onClick={() => router.push("/questionario")} variant="secondary">
            Fazer Novamente
          </Button>
          {user && (
            <Button onClick={() => router.push("/historico")} variant="secondary">
              Ver Histórico
            </Button>
          )}
          <Button onClick={() => router.push("/")} variant="ghost">
            Voltar ao Início
          </Button>
        </div>

        {/* Method Badge */}
        <p className="text-center text-xs text-gray-400 mt-8">
          {result.questionnaire_type === "detailed"
            ? "Questionário Detalhado · Scoring psicométrico determinístico (ipsativo)"
            : `Questionário Rápido · Análise por ${result.ai_provider === "gemini" ? "Google Gemini" : "Groq (Llama 3.3)"}`}
        </p>
      </div>
    </main>
  );
}
