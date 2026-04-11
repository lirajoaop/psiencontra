"use client";

import { motion, AnimatePresence } from "framer-motion";
import OptionCard from "./OptionCard";
import LikertScale from "./LikertScale";

export default function QuestionCard({ question, answer, onAnswer, direction }) {
  return (
    <AnimatePresence mode="wait" custom={direction}>
      <motion.div
        key={question.id}
        custom={direction}
        initial={{ x: direction > 0 ? 300 : -300, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        exit={{ x: direction > 0 ? -300 : 300, opacity: 0 }}
        transition={{ duration: 0.3, ease: "easeInOut" }}
        className="w-full"
      >
        <h2 className="text-xl md:text-2xl font-bold text-violet-900 dark:text-violet-100 mb-6">
          {question.text}
        </h2>

        {question.type === "multiple_choice" ? (
          <div className="grid gap-3">
            {question.options.map((opt) => (
              <OptionCard
                key={opt.value}
                label={opt.label}
                selected={answer === opt.value}
                onClick={() => onAnswer(opt.value)}
              />
            ))}
          </div>
        ) : question.type === "likert" ? (
          <LikertScale
            selected={answer}
            onSelect={onAnswer}
          />
        ) : (
          <textarea
            value={answer || ""}
            onChange={(e) => onAnswer(e.target.value)}
            placeholder="Escreva sua resposta aqui..."
            className="w-full h-36 p-4 rounded-xl border-2 border-violet-200 focus:border-violet-500 focus:outline-none resize-none text-gray-700 bg-white transition-colors dark:bg-gray-800 dark:border-violet-700 dark:text-gray-200 dark:focus:border-violet-400 dark:placeholder-gray-500"
          />
        )}
      </motion.div>
    </AnimatePresence>
  );
}
