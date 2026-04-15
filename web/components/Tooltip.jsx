"use client";

import { useEffect, useRef, useState } from "react";
import { Info } from "lucide-react";

export default function Tooltip({ label, side = "top", className = "" }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    if (!open) return;
    const onDown = (e) => {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false);
    };
    const onKey = (e) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onDown);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onDown);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const positions = {
    top: "bottom-full left-1/2 -translate-x-1/2 mb-2 origin-bottom",
    bottom: "top-full left-1/2 -translate-x-1/2 mt-2 origin-top",
    left: "right-full top-1/2 -translate-y-1/2 mr-2 origin-right",
    right: "left-full top-1/2 -translate-y-1/2 ml-2 origin-left",
  };

  return (
    <span ref={ref} className={`relative inline-flex group ${className}`}>
      <button
        type="button"
        aria-label="Mais informações"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
        className="inline-flex items-center justify-center rounded-full text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-gray-500 rounded-full"
      >
        <Info className="h-4 w-4" strokeWidth={2} />
      </button>
      <span
        role="tooltip"
        className={`pointer-events-none absolute z-50 w-max max-w-xs px-2 py-1 rounded bg-gray-700/95 dark:bg-gray-600/95 text-white text-[11px] font-medium leading-snug shadow-md transition-all duration-150 scale-90 opacity-0 group-hover:opacity-100 group-hover:scale-100 group-focus-within:opacity-100 group-focus-within:scale-100 ${open ? "opacity-100 scale-100" : ""} ${positions[side]}`}
      >
        {label}
      </span>
    </span>
  );
}
