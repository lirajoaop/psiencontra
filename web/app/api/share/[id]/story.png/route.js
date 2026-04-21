import { ImageResponse } from "next/og";
import QRCode from "qrcode";
import StoryCard from "@/components/StoryCard";
import {
  APPROACH_LABELS_SHORT,
  APPROACH_DESCRIPTIONS,
  APPROACH_HERO,
  APPROACH_TAGLINES,
  FIELD_LABELS,
} from "@/lib/constants";

export const runtime = "nodejs";

const FONT_URLS = {
  frauncesRegular:
    "https://cdn.jsdelivr.net/fontsource/fonts/fraunces@latest/latin-400-normal.ttf",
  frauncesItalic:
    "https://cdn.jsdelivr.net/fontsource/fonts/fraunces@latest/latin-400-italic.ttf",
  geistRegular:
    "https://cdn.jsdelivr.net/fontsource/fonts/geist-sans@latest/latin-400-normal.ttf",
  geistMedium:
    "https://cdn.jsdelivr.net/fontsource/fonts/geist-sans@latest/latin-500-normal.ttf",
  geistSemibold:
    "https://cdn.jsdelivr.net/fontsource/fonts/geist-sans@latest/latin-600-normal.ttf",
};

let fontCache = null;
async function loadFonts() {
  if (fontCache) return fontCache;
  const entries = await Promise.all(
    Object.entries(FONT_URLS).map(async ([k, url]) => [
      k,
      await fetch(url).then((r) => {
        if (!r.ok) throw new Error(`Font ${k} returned ${r.status}`);
        return r.arrayBuffer();
      }),
    ]),
  );
  fontCache = Object.fromEntries(entries);
  return fontCache;
}

function parseScores(v) {
  return typeof v === "string" ? JSON.parse(v) : v;
}

function topKey(scores) {
  return Object.entries(scores).sort((a, b) => b[1] - a[1])[0][0];
}

export async function GET(request, { params }) {
  const { id } = await params;

  const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
  const siteBase =
    process.env.NEXT_PUBLIC_SITE_URL || new URL(request.url).origin;

  let result;
  try {
    const res = await fetch(`${apiBase}/sessions/${id}/result`);
    if (!res.ok) {
      return new Response(`Result not found (${res.status})`, { status: res.status });
    }
    const json = await res.json();
    result = json.data;
  } catch (err) {
    return new Response(`Failed to fetch result: ${err.message}`, { status: 502 });
  }

  const approachScores = parseScores(result.approach_scores);
  const fieldScores = parseScores(result.field_scores);
  const topApproachKey = topKey(approachScores);
  const topFieldKey = topKey(fieldScores);

  const shareUrl = `${siteBase}/resultado/${id}`;
  const qrDataUrl = await QRCode.toDataURL(shareUrl, {
    errorCorrectionLevel: "M",
    margin: 0,
    width: 180,
    color: { dark: "#e9d5ff", light: "#00000000" },
  });

  const data = {
    topApproach: {
      key: topApproachKey,
      label:
        APPROACH_HERO[topApproachKey] ??
        APPROACH_LABELS_SHORT[topApproachKey] ??
        topApproachKey,
      score: Math.round(approachScores[topApproachKey]),
    },
    topField: {
      key: topFieldKey,
      label: FIELD_LABELS[topFieldKey] ?? topFieldKey,
      score: Math.round(fieldScores[topFieldKey]),
    },
    tagline:
      APPROACH_TAGLINES[topApproachKey] ??
      APPROACH_DESCRIPTIONS[topApproachKey] ??
      "",
    approachScores,
  };

  const fonts = await loadFonts();

  return new ImageResponse(
    (
      <StoryCard
        data={data}
        qrDataUrl={qrDataUrl}
        fontSerif="Fraunces"
        fontSans="Geist"
      />
    ),
    {
      width: 1080,
      height: 1920,
      fonts: [
        { name: "Fraunces", data: fonts.frauncesRegular, weight: 400, style: "normal" },
        { name: "Fraunces", data: fonts.frauncesItalic, weight: 400, style: "italic" },
        { name: "Geist", data: fonts.geistRegular, weight: 400, style: "normal" },
        { name: "Geist", data: fonts.geistMedium, weight: 500, style: "normal" },
        { name: "Geist", data: fonts.geistSemibold, weight: 600, style: "normal" },
      ],
      headers: {
        "cache-control": "public, immutable, no-transform, max-age=31536000",
      },
    },
  );
}
