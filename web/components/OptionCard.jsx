"use client";

import { motion } from "framer-motion";
import GlossaryText from "./GlossaryText";

export default function OptionCard({ label, selected, onClick }) {
  return (
    <motion.button
      onClick={onClick}
      whileHover={{ scale: 1.01 }}
      whileTap={{ scale: 0.99 }}
      className={`w-full text-left p-4 rounded-xl border-2 transition-all duration-200 cursor-pointer ${
        selected
          ? "border-violet-500 bg-violet-50 text-violet-900 shadow-md shadow-violet-100 dark:bg-violet-900/50 dark:text-violet-100 dark:border-violet-400 dark:shadow-violet-900/30"
          : "border-gray-200 bg-white text-gray-700 hover:border-violet-300 hover:bg-violet-50/50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-300 dark:hover:border-violet-500 dark:hover:bg-violet-900/30"
      }`}
    >
      <span className="text-sm md:text-base"><GlossaryText>{label}</GlossaryText></span>
    </motion.button>
  );
}
