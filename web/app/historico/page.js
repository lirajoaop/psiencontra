"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { motion } from "framer-motion";
import { ArrowLeft, Calendar, ChevronRight } from "lucide-react";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import ThemeToggle from "@/components/ThemeToggle";
import { useAuth } from "@/components/AuthProvider";
import { getUserHistory } from "@/lib/api";
import {
  APPROACH_LABELS,
  APPROACH_COLORS,
  FIELD_LABELS,
  FIELD_COLORS,
} from "@/lib/constants";

const DATE_FORMATTER = new Intl.DateTimeFormat("pt-BR", {
  day: "2-digit",
  month: "long",
  year: "numeric",
});

function parseScores(value) {
  if (!value) return {};
  return typeof value === "string" ? JSON.parse(value) : value;
}

function topEntry(scores, labels, colors) {
  const entries = Object.entries(scores);
  if (entries.length === 0) return null;
  const [key, score] = entries.reduce((best, cur) => (cur[1] > best[1] ? cur : best));
  return { key, score, label: labels[key] || key, color: colors[key] || "#7c3aed" };
}

export default function HistoricoPage() {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();
  const [sessions, setSessions] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (authLoading) return;
    if (!user) {
      router.replace("/login");
      return;
    }
    getUserHistory()
      .then((data) => setSessions(data ?? []))
      .catch((err) => setError(err.message));
  }, [authLoading, user, router]);

  if (authLoading || (user && sessions === null && !error)) {
    return <LoadingSpinner message="Carregando seu histórico..." />;
  }

  return (
    <main className="flex-1 pb-12 bg-gradient-to-br from-violet-50 to-purple-50 dark:from-gray-900 dark:to-gray-950 min-h-[60vh]">
      <div className="bg-gradient-to-r from-violet-600 to-purple-700 text-white py-10 px-6 relative">
        <div className="absolute top-4 left-4 z-20">
          <Link
            href="/"
            className="flex items-center gap-2 px-3 py-2 rounded-lg bg-white/10 hover:bg-white/20 transition-colors text-sm font-medium"
          >
            <ArrowLeft size={18} aria-hidden="true" />
            Início
          </Link>
        </div>
        <div className="absolute top-4 right-4 z-20">
          <ThemeToggle />
        </div>
        <div className="max-w-5xl mx-auto text-center">
          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-3xl md:text-4xl font-bold mb-3"
          >
            Meu Histórico
          </motion.h1>
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="text-violet-100 max-w-2xl mx-auto"
          >
            Acompanhe como seus interesses se reorganizam ao longo do tempo.
          </motion.p>
        </div>
      </div>

      <div className="max-w-3xl mx-auto px-6 mt-8">
        {error && (
          <div className="rounded-xl border border-red-200 dark:border-red-900 bg-red-50 dark:bg-red-950/40 text-red-700 dark:text-red-300 p-4 mb-6 text-sm">
            Não foi possível carregar seu histórico: {error}
          </div>
        )}

        {sessions && sessions.length === 0 && !error && (
          <div className="rounded-2xl border border-violet-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-8 text-center">
            <p className="text-gray-700 dark:text-gray-300 mb-6">
              Você ainda não concluiu nenhum questionário.
            </p>
            <Button onClick={() => router.push("/questionario")}>
              Fazer meu primeiro questionário
            </Button>
          </div>
        )}

        {sessions && sessions.length > 0 && (
          <ul className="space-y-4">
            {sessions.map((s, i) => {
              const approaches = parseScores(s.approach_scores);
              const fields = parseScores(s.field_scores);
              const topApproach = topEntry(approaches, APPROACH_LABELS, APPROACH_COLORS);
              const topField = topEntry(fields, FIELD_LABELS, FIELD_COLORS);
              const date = DATE_FORMATTER.format(new Date(s.created_at));
              const typeLabel =
                s.questionnaire_type === "detailed"
                  ? "Detalhado · 76 itens"
                  : "Rápido · 15 perguntas";

              return (
                <motion.li
                  key={s.id}
                  initial={{ opacity: 0, y: 12 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: i * 0.04 }}
                >
                  <Link
                    href={`/resultado/${s.id}`}
                    className="block rounded-2xl border border-violet-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-5 hover:shadow-md hover:border-violet-300 dark:hover:border-violet-600 transition-all"
                  >
                    <div className="flex items-center justify-between gap-4">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400 mb-2">
                          <Calendar size={14} aria-hidden="true" />
                          <span>{date}</span>
                          <span aria-hidden="true">·</span>
                          <span>{typeLabel}</span>
                        </div>
                        <div className="flex flex-wrap gap-2">
                          {topApproach && (
                            <Badge label={topApproach.label} color={topApproach.color} prefix="Abordagem" />
                          )}
                          {topField && (
                            <Badge label={topField.label} color={topField.color} prefix="Campo" />
                          )}
                        </div>
                      </div>
                      <ChevronRight
                        size={20}
                        className="text-gray-400 dark:text-gray-500 shrink-0"
                        aria-hidden="true"
                      />
                    </div>
                  </Link>
                </motion.li>
              );
            })}
          </ul>
        )}
      </div>
    </main>
  );
}

function Badge({ label, color, prefix }) {
  return (
    <span
      className="inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-medium"
      style={{ backgroundColor: `${color}1a`, color }}
    >
      <span className="w-2 h-2 rounded-full" style={{ backgroundColor: color }} />
      <span className="text-gray-500 dark:text-gray-400">{prefix}:</span>
      <span className="text-gray-800 dark:text-gray-100 truncate max-w-[200px]">{label}</span>
    </span>
  );
}
