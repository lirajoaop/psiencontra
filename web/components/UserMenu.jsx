"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { useAuth } from "./AuthProvider";

export default function UserMenu() {
  const router = useRouter();
  const { user, loading, logout } = useAuth();
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    function onClick(e) {
      if (ref.current && !ref.current.contains(e.target)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, []);

  // Render the "Entrar" button optimistically while the initial /auth/me
  // round-trip is in flight. Anonymous visitors (the common case on a
  // public landing page) see the CTA instantly instead of waiting for a
  // cold-started backend. Authenticated returning users briefly see
  // "Entrar" and then the avatar swaps in — a one-off micro-flash that's
  // preferable to showing nothing for several seconds.
  if (loading || !user) {
    return (
      <button
        onClick={() => router.push("/entrar")}
        className="px-4 py-2 rounded-lg bg-violet-100 dark:bg-violet-900/50 text-violet-700 dark:text-violet-300 hover:bg-violet-200 dark:hover:bg-violet-800 transition-colors text-sm font-semibold cursor-pointer"
      >
        Entrar
      </button>
    );
  }

  const initial = (user.name || user.email || "?").trim().charAt(0).toUpperCase();

  return (
    <div className="relative" ref={ref}>
      <motion.button
        whileTap={{ scale: 0.95 }}
        onClick={() => setOpen((o) => !o)}
        className="flex items-center gap-2 p-1 pr-3 rounded-full bg-violet-100 dark:bg-violet-900/50 hover:bg-violet-200 dark:hover:bg-violet-800 transition-colors cursor-pointer"
        aria-label="Menu do usuário"
      >
        {user.avatar_url ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={user.avatar_url} alt="" className="w-8 h-8 rounded-full" />
        ) : (
          <div className="w-8 h-8 rounded-full bg-violet-600 text-white flex items-center justify-center text-sm font-bold">
            {initial}
          </div>
        )}
        <span className="text-sm font-semibold text-violet-700 dark:text-violet-200 max-w-[120px] truncate">
          {user.name || user.email}
        </span>
      </motion.button>

      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: -8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -8 }}
            className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden z-30"
          >
            <button
              onClick={() => {
                setOpen(false);
                router.push("/historico");
              }}
              className="w-full text-left px-4 py-3 text-sm text-gray-700 dark:text-gray-200 hover:bg-violet-50 dark:hover:bg-gray-700 cursor-pointer"
            >
              Meu histórico
            </button>
            <div className="h-px bg-gray-100 dark:bg-gray-700" />
            <button
              onClick={async () => {
                await logout();
                setOpen(false);
                router.push("/");
              }}
              className="w-full text-left px-4 py-3 text-sm text-gray-700 dark:text-gray-200 hover:bg-violet-50 dark:hover:bg-gray-700 cursor-pointer"
            >
              Sair
            </button>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
