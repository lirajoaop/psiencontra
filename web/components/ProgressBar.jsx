"use client";

import { motion } from "framer-motion";

export default function ProgressBar({ current, total, blockLabel }) {
  const pct = ((current + 1) / total) * 100;

  return (
    <div className="w-full">
      <div className="flex justify-between items-center gap-3 text-sm text-violet-600 dark:text-violet-300 mb-2">
        <div className="flex items-center gap-2 min-w-0">
          <span className="shrink-0">Pergunta {current + 1} de {total}</span>
          {blockLabel && (
            <span className="text-xs text-violet-400 bg-violet-50 dark:bg-violet-900/50 dark:text-violet-300 px-2 py-0.5 rounded-full truncate">
              {blockLabel}
            </span>
          )}
        </div>
        <span className="shrink-0 tabular-nums">{Math.round(pct)}%</span>
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
