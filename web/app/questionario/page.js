"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import ProgressBar from "@/components/ProgressBar";
import QuestionCard from "@/components/QuestionCard";
import LoadingSpinner from "@/components/LoadingSpinner";
import Button from "@/components/Button";
import ThemeToggle from "@/components/ThemeToggle";
import { getQuestions, createSession, submitResponses } from "@/lib/api";

export default function Questionario() {
  const router = useRouter();
  const [questions, setQuestions] = useState([]);
  const [current, setCurrent] = useState(0);
  const [answers, setAnswers] = useState({});
  const [direction, setDirection] = useState(1);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    getQuestions()
      .then((data) => {
        setQuestions(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

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
      const session = await createSession();
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

  return (
    <main className="flex-1 flex flex-col">
      {/* Header */}
      <div className="bg-white dark:bg-gray-900 border-b border-violet-100 dark:border-gray-700 px-6 py-4">
        <div className="max-w-2xl mx-auto">
          <div className="flex items-center justify-between mb-3">
            <h1 className="text-lg font-bold text-violet-900 dark:text-violet-200">PsiEncontra</h1>
            <div className="flex items-center gap-3">
              <span className="text-xs text-violet-400 bg-violet-50 dark:bg-violet-900/50 dark:text-violet-300 px-2 py-1 rounded-full">
                {question.block === "approaches" ? "Abordagens" : "Campos de Atuação"}
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
        <div className="max-w-2xl mx-auto flex justify-between items-center">
          <Button
            onClick={handlePrev}
            variant="ghost"
            disabled={isFirst}
          >
            Anterior
          </Button>

          <div className="flex gap-1">
            {questions.map((_, i) => (
              <button
                key={i}
                onClick={() => { setDirection(i > current ? 1 : -1); setCurrent(i); }}
                className={`w-2 h-2 rounded-full transition-colors cursor-pointer ${
                  i === current ? "bg-violet-600" : answers[questions[i]?.id] ? "bg-violet-300 dark:bg-violet-500" : "bg-gray-200 dark:bg-gray-600"
                }`}
              />
            ))}
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
