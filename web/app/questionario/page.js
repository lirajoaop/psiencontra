"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import Link from "next/link";
import ProgressBar from "@/components/ProgressBar";
import QuestionCard from "@/components/QuestionCard";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import ThemeToggle from "@/components/ThemeToggle";
import { getQuestions, createSession, submitResponses } from "@/lib/api";

const BLOCK_LABELS = {
  approaches: "Abordagens",
  fields: "Campos de Atuação",
  vignettes: "Situações Práticas",
  reflection: "Reflexão Pessoal",
};

export default function Questionario() {
  const router = useRouter();
  const [questionnaireType, setQuestionnaireType] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [current, setCurrent] = useState(0);
  const [answers, setAnswers] = useState({});
  const [direction, setDirection] = useState(1);
  const [loading, setLoading] = useState(false);

  function selectType(type) {
    setQuestionnaireType(type);
    setLoading(true);
  }
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!questionnaireType) return;
    let cancelled = false;
    getQuestions(questionnaireType)
      .then((data) => {
        if (cancelled) return;
        setQuestions(data);
        setLoading(false);
      })
      .catch((err) => {
        if (cancelled) return;
        setError(err.message);
        setLoading(false);
      });
    return () => { cancelled = true; };
  }, [questionnaireType]);

  const question = questions[current];
  const isFirst = current === 0;
  const isLast = current === questions.length - 1;
  const hasAnswer = answers[question?.id] != null && answers[question?.id] !== "";

  function handleAnswer(value) {
    setAnswers((prev) => ({ ...prev, [question.id]: value }));
  }

  function handleNext() {
    if (!isLast) {
      setDirection(1);
      setCurrent((c) => c + 1);
    }
  }

  function handlePrev() {
    if (!isFirst) {
      setDirection(-1);
      setCurrent((c) => c - 1);
    }
  }

  async function handleSubmit() {
    setSubmitting(true);
    setError(null);

    try {
      const session = await createSession(questionnaireType);
      const sessionId = session.id;

      const responsePayload = questions.map((q) => ({
        question_id: q.id,
        answer_value: answers[q.id] || "",
      }));

      await submitResponses(sessionId, responsePayload);
      router.push(`/resultado/${sessionId}`);
    } catch (err) {
      setError(err.message);
      setSubmitting(false);
    }
  }

  // --- Selection screen ---
  if (!questionnaireType) {
    return (
      <main className="flex-1 flex items-center justify-center px-6 py-16 bg-gradient-to-br from-violet-50 to-purple-50 dark:from-gray-900 dark:to-gray-950">
        <div className="absolute top-4 right-4 z-20">
          <ThemeToggle />
        </div>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="w-full max-w-2xl"
        >
          <h1 className="text-3xl md:text-4xl font-bold text-violet-900 dark:text-violet-200 text-center mb-3">
            Escolha o Questionário
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-center mb-10">
            Dois formatos, o mesmo objetivo: descobrir seu perfil em Psicologia.
          </p>

          <div className="grid md:grid-cols-2 gap-6">
            {/* Simple */}
            <motion.button
              onClick={() => selectType("simple")}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className="cursor-pointer text-left p-6 rounded-2xl border-2 border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 hover:border-violet-400 dark:hover:border-violet-500 transition-all shadow-sm hover:shadow-lg"
            >
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-full bg-violet-100 dark:bg-violet-900/50 text-violet-600 dark:text-violet-300 flex items-center justify-center text-lg font-bold">
                  R
                </div>
                <h2 className="text-lg font-bold text-violet-900 dark:text-violet-200">Rápido</h2>
              </div>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                15 perguntas entre objetivas e dissertativas. Análise interpretativa feita por IA — leve e exploratório.
              </p>
              <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
                <span className="flex items-center gap-1">
                  <ClockIcon /> ~5 min
                </span>
                <span className="flex items-center gap-1">
                  <ListIcon /> 15 perguntas
                </span>
              </div>
            </motion.button>

            {/* Detailed */}
            <motion.button
              onClick={() => selectType("detailed")}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className="cursor-pointer text-left p-6 rounded-2xl border-2 border-violet-300 dark:border-violet-600 bg-gradient-to-br from-violet-50 to-white dark:from-violet-900/30 dark:to-gray-800 hover:border-violet-500 dark:hover:border-violet-400 transition-all shadow-sm hover:shadow-lg relative"
            >
              <div className="absolute top-3 right-3">
                <span className="text-[10px] font-bold bg-violet-600 text-white px-2 py-0.5 rounded-full">
                  RECOMENDADO
                </span>
              </div>
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-full bg-violet-600 text-white flex items-center justify-center text-lg font-bold">
                  D
                </div>
                <h2 className="text-lg font-bold text-violet-900 dark:text-violet-200">Detalhado</h2>
              </div>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                76 afirmações em escala Likert. Scoring determinístico e reprodutível, sem interferência da IA nos números.
              </p>
              <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
                <span className="flex items-center gap-1">
                  <ClockIcon /> ~10–12 min
                </span>
                <span className="flex items-center gap-1">
                  <ListIcon /> 76 itens
                </span>
              </div>
            </motion.button>
          </div>
        </motion.div>
      </main>
    );
  }

  // --- Loading / Submitting / Error ---
  if (loading) {
    return <LoadingSpinner message="Carregando perguntas..." />;
  }

  if (submitting) {
    return <LoadingSpinner message="Analisando suas respostas com IA..." />;
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4 px-6">
        <p className="text-red-600 dark:text-red-400 text-center">{error}</p>
        <Button onClick={() => setError(null)} variant="secondary">Tentar novamente</Button>
      </div>
    );
  }

  if (!question) return null;

  const allAnswered = questions.every((q) => answers[q.id] != null && answers[q.id] !== "");
  const blockLabel = BLOCK_LABELS[question.block] || question.block;

  return (
    <main className="flex-1 flex flex-col">
      {/* Header */}
      <div className="bg-white dark:bg-gray-900 border-b border-violet-100 dark:border-gray-700 px-6 py-4">
        <div className="max-w-2xl mx-auto">
          <div className="flex items-center justify-between mb-3">
            <h1 className="text-lg font-bold text-violet-900 dark:text-violet-200">PsiEncontra</h1>
            <div className="flex items-center gap-3">
              <span className="text-xs text-violet-400 bg-violet-50 dark:bg-violet-900/50 dark:text-violet-300 px-2 py-1 rounded-full">
                {blockLabel}
              </span>
              <ThemeToggle />
            </div>
          </div>
          <ProgressBar current={current} total={questions.length} />
        </div>
      </div>

      {/* Question */}
      <div className="flex-1 flex items-center justify-center px-6 py-8">
        <div className="w-full max-w-2xl">
          <QuestionCard
            question={question}
            answer={answers[question.id]}
            onAnswer={handleAnswer}
            direction={direction}
          />
        </div>
      </div>

      {/* Navigation */}
      <div className="bg-white dark:bg-gray-900 border-t border-violet-100 dark:border-gray-700 px-6 py-4">
        {isLast && (
          <p className="max-w-2xl mx-auto text-xs text-gray-500 dark:text-gray-400 text-center mb-3">
            Ao clicar em <strong>Enviar</strong>, você concorda que suas respostas sejam
            usadas para gerar o seu resultado, conforme a{" "}
            <Link
              href="/privacidade"
              target="_blank"
              className="underline text-violet-700 dark:text-violet-300"
            >
              Política de Privacidade
            </Link>
            .
          </p>
        )}
        <div className="max-w-2xl mx-auto flex justify-between items-center">
          <Button
            onClick={handlePrev}
            variant="ghost"
            disabled={isFirst}
          >
            Anterior
          </Button>

          <div className="flex gap-1">
            {questions.length <= 30 ? (
              questions.map((_, i) => (
                <button
                  key={i}
                  onClick={() => { setDirection(i > current ? 1 : -1); setCurrent(i); }}
                  className={`w-2 h-2 rounded-full transition-colors cursor-pointer ${
                    i === current ? "bg-violet-600" : answers[questions[i]?.id] ? "bg-violet-300 dark:bg-violet-500" : "bg-gray-200 dark:bg-gray-600"
                  }`}
                />
              ))
            ) : (
              <span className="text-xs text-gray-400 dark:text-gray-500">
                {current + 1} / {questions.length}
              </span>
            )}
          </div>

          {isLast ? (
            <Button
              onClick={handleSubmit}
              disabled={!allAnswered}
            >
              Enviar
            </Button>
          ) : (
            <Button
              onClick={handleNext}
              disabled={!hasAnswer}
            >
              Próxima
            </Button>
          )}
        </div>
      </div>
    </main>
  );
}

function ClockIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <polyline points="12 6 12 12 16 14" />
    </svg>
  );
}

function ListIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="8" y1="6" x2="21" y2="6" />
      <line x1="8" y1="12" x2="21" y2="12" />
      <line x1="8" y1="18" x2="21" y2="18" />
      <line x1="3" y1="6" x2="3.01" y2="6" />
      <line x1="3" y1="12" x2="3.01" y2="12" />
      <line x1="3" y1="18" x2="3.01" y2="18" />
    </svg>
  );
}
