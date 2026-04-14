"use client";

import { motion } from "framer-motion";
import ScoreBar from "./ScoreBar";

export default function ApproachCard({ name, score, description, color, rank }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: rank * 0.1 }}
      className="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-md border border-violet-100 dark:border-gray-700 hover:shadow-lg transition-shadow"
    >
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-3">
          <span
            className="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold"
            style={{ backgroundColor: color }}
          >
            {rank}
          </span>
          <h4 className="font-bold text-violet-900 dark:text-violet-200">{name}</h4>
        </div>
        <span className="text-2xl font-bold" style={{ color }}>{Math.round(score)}%</span>
      </div>
      <ScoreBar score={Math.round(score)} color={color} />
      {description && (
        <p className="text-sm text-gray-600 dark:text-gray-300 mt-3 leading-relaxed">{description}</p>
      )}
    </motion.div>
  );
}
