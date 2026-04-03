"use client";

import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import Button from "@/components/Button";
import ThemeToggle from "@/components/ThemeToggle";
import { APPROACH_LABELS, APPROACH_DESCRIPTIONS, APPROACH_COLORS } from "@/lib/constants";

export default function Home() {
  const router = useRouter();

  const approaches = Object.entries(APPROACH_LABELS).map(([key, label]) => ({
    key,
    label,
    description: APPROACH_DESCRIPTIONS[key],
    color: APPROACH_COLORS[key],
  }));

  return (
    <main className="flex-1">
      {/* Hero */}
      <section className="relative overflow-hidden bg-gradient-to-br from-violet-600 via-violet-700 to-purple-800 text-white">
        <div className="absolute top-4 right-4 z-20">
          <ThemeToggle />
        </div>
        <div className="max-w-5xl mx-auto px-6 py-24 md:py-32 text-center relative z-10">
          <motion.h1
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-4xl md:text-6xl font-bold mb-6"
          >
            PsiEncontra
          </motion.h1>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.15 }}
            className="text-lg md:text-xl text-violet-100 max-w-2xl mx-auto mb-4"
          >
            Descubra qual abordagem teórica e campo de atuação da Psicologia mais combinam com você
          </motion.p>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.25 }}
            className="text-violet-200 max-w-xl mx-auto mb-10"
          >
            Responda 15 perguntas e nossa IA analisará seu perfil, gerando um ranking de afinidade com 8 abordagens e 9 campos de atuação.
          </motion.p>
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.35 }}
          >
            <Button
              onClick={() => router.push("/questionario")}
              className="!bg-white !text-violet-700 hover:!bg-violet-50 !text-lg !px-10 !py-4"
            >
              Começar Questionário
            </Button>
          </motion.div>
        </div>
      </section>

      {/* How it works */}
      <section className="max-w-5xl mx-auto px-6 py-16">
        <h2 className="text-2xl md:text-3xl font-bold text-violet-900 dark:text-violet-200 text-center mb-12">
          Como funciona?
        </h2>
        <div className="grid md:grid-cols-3 gap-8">
          {[
            { step: "1", title: "Responda", desc: "15 perguntas sobre sua visão da Psicologia, entre objetivas e dissertativas." },
            { step: "2", title: "IA Analisa", desc: "Nossa inteligência artificial analisa suas respostas e identifica padrões." },
            { step: "3", title: "Descubra", desc: "Veja seu ranking de afinidade com gráficos radar e explicações detalhadas." },
          ].map((item) => (
            <motion.div
              key={item.step}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-center"
            >
              <div className="w-14 h-14 rounded-full bg-violet-100 dark:bg-violet-900/50 text-violet-700 dark:text-violet-300 flex items-center justify-center text-xl font-bold mx-auto mb-4">
                {item.step}
              </div>
              <h3 className="font-bold text-violet-900 dark:text-violet-200 text-lg mb-2">{item.title}</h3>
              <p className="text-gray-600 dark:text-gray-400">{item.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>

      {/* Approaches Grid */}
      <section className="bg-white dark:bg-gray-900 py-16">
        <div className="max-w-5xl mx-auto px-6">
          <h2 className="text-2xl md:text-3xl font-bold text-violet-900 dark:text-violet-200 text-center mb-12">
            As 8 Abordagens Teóricas
          </h2>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {approaches.map((a, i) => (
              <motion.div
                key={a.key}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
                className="p-4 rounded-xl border border-violet-100 dark:border-gray-700 hover:shadow-md transition-shadow dark:bg-gray-800"
              >
                <div
                  className="w-3 h-3 rounded-full mb-3"
                  style={{ backgroundColor: a.color }}
                />
                <h4 className="font-bold text-violet-900 dark:text-violet-200 text-sm mb-1">{a.label}</h4>
                <p className="text-xs text-gray-500 dark:text-gray-400">{a.description}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="max-w-5xl mx-auto px-6 py-16 text-center">
        <h2 className="text-2xl font-bold text-violet-900 dark:text-violet-200 mb-4">
          Pronto para se descobrir?
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-8">
          O questionário leva cerca de 10 minutos. Suas respostas são anônimas.
        </p>
        <Button onClick={() => router.push("/questionario")}>
          Iniciar Questionário
        </Button>
      </section>

      {/* Footer */}
      <footer className="bg-violet-900 dark:bg-gray-900 dark:border-t dark:border-gray-800 text-violet-200 py-6 text-center text-sm">
        <p>PsiEncontra — Resultado orientativo, não substitui acompanhamento profissional.</p>
      </footer>
    </main>
  );
}
