"use client";

import { motion } from "framer-motion";

export default function OptionCard({ label, selected, onClick }) {
  return (
    <motion.button
      onClick={onClick}
      whileHover={{ scale: 1.01 }}
      whileTap={{ scale: 0.99 }}
      className={`w-full text-left p-4 rounded-xl border-2 transition-all duration-200 cursor-pointer ${
        selected
          ? "border-violet-500 bg-violet-50 text-violet-900 shadow-md shadow-violet-100"
          : "border-gray-200 bg-white text-gray-700 hover:border-violet-300 hover:bg-violet-50/50"
      }`}
    >
      <span className="text-sm md:text-base">{label}</span>
    </motion.button>
  );
}
