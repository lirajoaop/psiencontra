"use client";

import { motion } from "framer-motion";

export default function ProgressBar({ current, total }) {
  const pct = ((current + 1) / total) * 100;

  return (
    <div className="w-full">
      <div className="flex justify-between text-sm text-violet-600 dark:text-violet-300 mb-2">
        <span>Pergunta {current + 1} de {total}</span>
        <span>{Math.round(pct)}%</span>
      </div>
      <div className="w-full h-2 bg-violet-100 dark:bg-violet-900 rounded-full overflow-hidden">
        <motion.div
          className="h-full bg-gradient-to-r from-violet-500 to-violet-600 rounded-full"
          initial={{ width: 0 }}
          animate={{ width: `${pct}%` }}
          transition={{ duration: 0.4, ease: "easeOut" }}
        />
      </div>
    </div>
  );
}
