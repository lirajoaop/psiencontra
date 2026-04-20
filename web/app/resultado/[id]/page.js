import ResultadoClient from "./ResultadoClient";
import { getResultServer } from "@/lib/api-server";
import { APPROACH_LABELS, FIELD_LABELS } from "@/lib/constants";

const SITE_URL =
  process.env.NEXT_PUBLIC_SITE_URL ||
  (process.env.VERCEL_PROJECT_PRODUCTION_URL && `https://${process.env.VERCEL_PROJECT_PRODUCTION_URL}`) ||
  (process.env.VERCEL_URL && `https://${process.env.VERCEL_URL}`) ||
  "http://localhost:3000";

function parseScores(raw) {
  if (!raw) return null;
  return typeof raw === "string" ? JSON.parse(raw) : raw;
}

function topKey(scores) {
  if (!scores) return null;
  const entries = Object.entries(scores);
  if (!entries.length) return null;
  entries.sort(([, a], [, b]) => b - a);
  return entries[0][0];
}

export async function generateMetadata({ params }) {
  const { id } = await params;
  const baseTitle = "Meu resultado no PsiEncontra";
  const baseDescription = "Descubra qual abordagem teórica e campo de atuação da Psicologia mais combinam com você.";
  const url = `${SITE_URL}/resultado/${id}`;

  const result = await getResultServer(id);

  let description = baseDescription;
  if (result) {
    const approachKey = topKey(parseScores(result.approach_scores));
    const fieldKey = topKey(parseScores(result.field_scores));
    const parts = [];
    if (approachKey) parts.push(`Abordagem: ${APPROACH_LABELS[approachKey] || approachKey}`);
    if (fieldKey) parts.push(`Campo: ${FIELD_LABELS[fieldKey] || fieldKey}`);
    if (parts.length) description = parts.join(" · ");
  }

  return {
    metadataBase: new URL(SITE_URL),
    title: baseTitle,
    description,
    alternates: { canonical: url },
    openGraph: {
      title: baseTitle,
      description,
      url,
      siteName: "PsiEncontra",
      locale: "pt_BR",
      type: "article",
    },
    twitter: {
      card: "summary",
      title: baseTitle,
      description,
    },
  };
}

export default async function Page({ params }) {
  const { id } = await params;
  const initialResult = await getResultServer(id);
  return <ResultadoClient id={id} initialResult={initialResult} />;
}
