"use client";

import { motion } from "framer-motion";

export default function ScoreBar({ score, color = "#7c3aed" }) {
  return (
    <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
      <motion.div
        className="h-full rounded-full"
        style={{ backgroundColor: color }}
        initial={{ width: 0 }}
        animate={{ width: `${score}%` }}
        transition={{ duration: 0.8, ease: "easeOut" }}
      />
    </div>
  );
}
