"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import RadarChartResult from "@/components/RadarChartResult";
import ApproachCard from "@/components/ApproachCard";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import ThemeToggle from "@/components/ThemeToggle";
import { getResult, getPDFUrl } from "@/lib/api";
import {
  APPROACH_LABELS,
  APPROACH_COLORS,
  FIELD_LABELS,
  FIELD_COLORS,
} from "@/lib/constants";

export default function Resultado({ params }) {
  const { id } = use(params);
  const router = useRouter();
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    getResult(id)
      .then((data) => {
        setResult(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id]);

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
          <Button onClick={() => router.push("/")} variant="ghost">
            Voltar ao Início
          </Button>
        </div>

        {/* AI Provider Badge */}
        <p className="text-center text-xs text-gray-400 mt-8">
          Análise gerada por {result.ai_provider === "gemini" ? "Google Gemini" : "Groq (Llama 3.3)"}
        </p>
      </div>
    </main>
  );
}
