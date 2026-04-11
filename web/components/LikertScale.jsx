"use client";

import { motion } from "framer-motion";

const LABELS = [
  { value: "1", label: "Discordo totalmente" },
  { value: "2", label: "Discordo" },
  { value: "3", label: "Neutro" },
  { value: "4", label: "Concordo" },
  { value: "5", label: "Concordo totalmente" },
];

export default function LikertScale({ selected, onSelect }) {
  return (
    <div className="flex flex-col gap-4 mt-4">
      {/* Mobile: vertical buttons */}
      <div className="flex flex-col gap-2 sm:hidden">
        {LABELS.map((item) => (
          <motion.button
            key={item.value}
            onClick={() => onSelect(item.value)}
            whileTap={{ scale: 0.97 }}
            className={`w-full text-left px-4 py-3 rounded-xl border-2 transition-all duration-200 cursor-pointer text-sm ${
              selected === item.value
                ? "border-violet-500 bg-violet-50 text-violet-900 shadow-md shadow-violet-100 dark:bg-violet-900/50 dark:text-violet-100 dark:border-violet-400 dark:shadow-violet-900/30"
                : "border-gray-200 bg-white text-gray-700 hover:border-violet-300 hover:bg-violet-50/50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-300 dark:hover:border-violet-500 dark:hover:bg-violet-900/30"
            }`}
          >
            {item.label}
          </motion.button>
        ))}
      </div>

      {/* Desktop: horizontal scale */}
      <div className="hidden sm:block">
        <div className="flex justify-between mb-2">
          <span className="text-xs text-gray-400 dark:text-gray-500">Discordo totalmente</span>
          <span className="text-xs text-gray-400 dark:text-gray-500">Concordo totalmente</span>
        </div>
        <div className="flex justify-between gap-2">
          {LABELS.map((item) => (
            <motion.button
              key={item.value}
              onClick={() => onSelect(item.value)}
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.95 }}
              className={`flex-1 py-3 rounded-xl border-2 transition-all duration-200 cursor-pointer font-bold text-lg ${
                selected === item.value
                  ? "border-violet-500 bg-violet-600 text-white shadow-lg shadow-violet-200 dark:shadow-violet-900/50"
                  : "border-gray-200 bg-white text-gray-500 hover:border-violet-300 hover:bg-violet-50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:border-violet-500 dark:hover:bg-violet-900/30"
              }`}
            >
              {item.value}
            </motion.button>
          ))}
        </div>
      </div>
    </div>
  );
}
