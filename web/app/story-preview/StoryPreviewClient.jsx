"use client";

import { useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import QRCode from "qrcode";
import StoryCard, { CARD_WIDTH, CARD_HEIGHT } from "@/components/StoryCard";
import {
  APPROACH_LABELS_SHORT,
  APPROACH_DESCRIPTIONS,
  APPROACH_HERO,
  APPROACH_TAGLINES,
  FIELD_LABELS,
} from "@/lib/constants";

function buildScores(topKey) {
  const keys = Object.keys(APPROACH_LABELS_SHORT);
  const others = keys.filter((k) => k !== topKey);
  const vals = [78, 72, 60, 50, 45, 38, 28];
  const out = { [topKey]: 87 };
  others.forEach((k, i) => {
    out[k] = vals[i] ?? 30;
  });
  return out;
}

const SHARE_URL = "https://psiencontra.com/resultado/preview";

export default function StoryPreviewClient() {
  const params = useSearchParams();
  const approachKey = params.get("approach") || "humanismo";
  const fieldKey = params.get("field") || "clinica";

  const sample = {
    topApproach: {
      key: approachKey,
      label: APPROACH_HERO[approachKey] ?? APPROACH_LABELS_SHORT[approachKey] ?? approachKey,
      score: 87,
    },
    topField: {
      key: fieldKey,
      label: FIELD_LABELS[fieldKey] ?? fieldKey,
      score: 82,
    },
    tagline: APPROACH_TAGLINES[approachKey] ?? APPROACH_DESCRIPTIONS[approachKey] ?? "",
    approachScores: buildScores(approachKey),
  };

  const [showSafeZones, setShowSafeZones] = useState(false);
  const [scale, setScale] = useState(0.45);
  const [qrDataUrl, setQrDataUrl] = useState(null);

  useEffect(() => {
    let cancelled = false;
    QRCode.toDataURL(SHARE_URL, {
      errorCorrectionLevel: "M",
      margin: 0,
      width: 180,
      color: { dark: "#e9d5ff", light: "#00000000" },
    })
      .then((d) => {
        if (!cancelled) setQrDataUrl(d);
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div
      style={{
        minHeight: "100vh",
        background: "#0a0612",
        padding: "24px",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: 16,
        fontFamily: "var(--font-geist-sans), ui-sans-serif, system-ui, sans-serif",
      }}
    >
      <div
        style={{
          color: "#c4b5fd",
          fontSize: 14,
          display: "flex",
          gap: 24,
          alignItems: "center",
          flexWrap: "wrap",
          justifyContent: "center",
        }}
      >
        <span style={{ opacity: 0.7 }}>preview · {CARD_WIDTH} × {CARD_HEIGHT} · 9:16</span>
        <label style={{ display: "flex", gap: 8, alignItems: "center" }}>
          <input
            type="checkbox"
            checked={showSafeZones}
            onChange={(e) => setShowSafeZones(e.target.checked)}
          />
          safe zones do Instagram
        </label>
        <label style={{ display: "flex", gap: 8, alignItems: "center" }}>
          escala
          <input
            type="range"
            min="0.2"
            max="0.9"
            step="0.05"
            value={scale}
            onChange={(e) => setScale(Number(e.target.value))}
          />
          <span style={{ width: 36 }}>{Math.round(scale * 100)}%</span>
        </label>
      </div>

      <div
        style={{
          width: CARD_WIDTH * scale,
          height: CARD_HEIGHT * scale,
          position: "relative",
          boxShadow: "0 30px 80px rgba(124, 58, 237, 0.25)",
          borderRadius: 12 * scale,
          overflow: "hidden",
        }}
      >
        <div
          style={{
            transform: `scale(${scale})`,
            transformOrigin: "top left",
            width: CARD_WIDTH,
            height: CARD_HEIGHT,
          }}
        >
          <StoryCard data={sample} qrDataUrl={qrDataUrl} showSafeZones={showSafeZones} />
        </div>
      </div>
    </div>
  );
}
