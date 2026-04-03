"use client";

import { motion } from "framer-motion";

export default function LoadingSpinner({ message = "Analisando suas respostas..." }) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] gap-6">
      <motion.div
        className="w-16 h-16 border-4 border-violet-200 border-t-violet-600 rounded-full"
        animate={{ rotate: 360 }}
        transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
      />
      <motion.p
        className="text-violet-700 text-lg font-medium"
        animate={{ opacity: [0.5, 1, 0.5] }}
        transition={{ duration: 2, repeat: Infinity }}
      >
        {message}
      </motion.p>
    </div>
  );
}
