import { APPROACH_LABELS_SHORT } from "@/lib/constants";

export const CARD_WIDTH = 1080;
export const CARD_HEIGHT = 1920;

const DEFAULT_SERIF = "var(--font-fraunces), ui-serif, Georgia, serif";
const DEFAULT_SANS = "var(--font-geist-sans), ui-sans-serif, system-ui, sans-serif";

function heroLines(label) {
  if (label.length > 20 && label.includes("-")) {
    const idx = label.indexOf("-");
    return [label.slice(0, idx + 1), label.slice(idx + 1)];
  }
  return [label];
}

function heroFontSize(maxLineLen) {
  if (maxLineLen <= 10) return 152;
  if (maxLineLen <= 14) return 132;
  if (maxLineLen <= 18) return 110;
  if (maxLineLen <= 22) return 92;
  return 76;
}

export default function StoryCard({
  data,
  qrDataUrl,
  showSafeZones = false,
  fontSerif = DEFAULT_SERIF,
  fontSans = DEFAULT_SANS,
}) {
  return (
    <div
      style={{
        width: CARD_WIDTH,
        height: CARD_HEIGHT,
        position: "relative",
        overflow: "hidden",
        color: "#f5f3ff",
        fontFamily: fontSans,
        display: "flex",
        background: "linear-gradient(180deg, #1a0b2e 0%, #0a0612 100%)",
      }}
    >
      <div
        style={{
          position: "absolute",
          top: 0,
          left: 0,
          width: CARD_WIDTH,
          height: CARD_HEIGHT,
          background:
            "radial-gradient(circle at 18% 12%, #7c3aed 0%, rgba(124,58,237,0) 55%)",
        }}
      />
      <div
        style={{
          position: "absolute",
          top: 0,
          left: 0,
          width: CARD_WIDTH,
          height: CARD_HEIGHT,
          background:
            "radial-gradient(circle at 88% 92%, #ec4899 0%, rgba(236,72,153,0) 55%)",
        }}
      />
      <div
        style={{
          position: "absolute",
          top: 0,
          left: 0,
          width: CARD_WIDTH,
          height: CARD_HEIGHT,
          background:
            "radial-gradient(circle at 50% 55%, #4c1d95 0%, rgba(76,29,149,0) 70%)",
        }}
      />

      <div
        style={{
          position: "absolute",
          inset: 0,
          padding: "260px 90px",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
          <div
            style={{
              width: 64,
              height: 64,
              borderRadius: "50%",
              border: "1.5px solid rgba(196, 181, 253, 0.7)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontFamily: fontSerif,
              fontSize: 40,
              color: "#f0abfc",
              lineHeight: 1,
              paddingBottom: 4,
            }}
          >
            ψ
          </div>
          <div style={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <div style={{ fontSize: 26, fontWeight: 600, letterSpacing: "-0.01em" }}>
              psiencontra
            </div>
            <div
              style={{
                fontSize: 16,
                color: "#c4b5fd",
                letterSpacing: "0.14em",
                textTransform: "uppercase",
              }}
            >
              meu perfil em psicologia
            </div>
          </div>
        </div>

        <div style={{ height: 80 }} />

        {(() => {
          const lines = heroLines(data.topApproach.label);
          const maxLen = Math.max(...lines.map((l) => l.length));
          return (
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                fontFamily: fontSerif,
                fontWeight: 400,
                fontSize: heroFontSize(maxLen),
                lineHeight: 0.95,
                letterSpacing: "-0.028em",
                backgroundImage:
                  "linear-gradient(135deg, #ffffff 0%, #ede9fe 35%, #f0abfc 100%)",
                backgroundClip: "text",
                WebkitBackgroundClip: "text",
                color: "transparent",
                paddingRight: 12,
              }}
            >
              {lines.map((l, i) => (
                <div key={i}>{l}</div>
              ))}
            </div>
          );
        })()}

        <div style={{ height: 28 }} />

        <div style={{ fontSize: 44, fontWeight: 500, color: "#ddd6fe" }}>
          {`+ ${data.topField.label}`}
        </div>

        <div style={{ height: 16 }} />

        <div style={{ display: "flex", alignItems: "center", gap: 14, fontSize: 22, color: "#a78bfa" }}>
          <Dot />
          <span>{data.topApproach.score}% afinidade · abordagem</span>
          <span style={{ opacity: 0.5 }}>·</span>
          <span>{data.topField.score}% · campo</span>
        </div>

        <div style={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center" }}>
          <RadarSVG scores={data.approachScores} labels={APPROACH_LABELS_SHORT} fontSans={fontSans} />
        </div>

        <div
          style={{
            fontFamily: fontSerif,
            fontStyle: "italic",
            fontWeight: 400,
            fontSize: 42,
            lineHeight: 1.3,
            color: "#ede9fe",
            maxWidth: 820,
          }}
        >
          {`“${data.tagline}”`}
        </div>

        <div style={{ height: 64 }} />

        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
          <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
            <div style={{ fontSize: 20, color: "#c4b5fd", letterSpacing: "0.08em", textTransform: "uppercase" }}>
              faça o seu
            </div>
            <div style={{ fontSize: 36, fontWeight: 600, letterSpacing: "-0.01em" }}>
              psiencontra.com
            </div>
          </div>
          <QRBadge dataUrl={qrDataUrl} />
        </div>
      </div>

      {showSafeZones && (
        <>
          <div
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              right: 0,
              height: 250,
              background: "rgba(239, 68, 68, 0.12)",
              borderBottom: "2px dashed rgba(252, 165, 165, 0.7)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              color: "#fecaca",
              fontSize: 20,
              letterSpacing: "0.1em",
              textTransform: "uppercase",
            }}
          >
            safe zone topo · UI do IG cobre
          </div>
          <div
            style={{
              position: "absolute",
              bottom: 0,
              left: 0,
              right: 0,
              height: 250,
              background: "rgba(239, 68, 68, 0.12)",
              borderTop: "2px dashed rgba(252, 165, 165, 0.7)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              color: "#fecaca",
              fontSize: 20,
              letterSpacing: "0.1em",
              textTransform: "uppercase",
            }}
          >
            safe zone rodapé · UI do IG cobre
          </div>
        </>
      )}
    </div>
  );
}

function Dot() {
  return (
    <div
      style={{
        width: 8,
        height: 8,
        borderRadius: "50%",
        background: "linear-gradient(135deg, #a78bfa, #f0abfc)",
      }}
    />
  );
}

function QRBadge({ dataUrl, size = 180 }) {
  return (
    <div
      style={{
        width: size,
        height: size,
        background: "#140821",
        border: "2px solid rgba(196, 181, 253, 0.6)",
        borderRadius: 22,
        padding: 14,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        boxSizing: "border-box",
      }}
    >
      {dataUrl && (
        <img
          src={dataUrl}
          alt="QR code psiencontra"
          style={{ width: "100%", height: "100%", display: "block" }}
        />
      )}
    </div>
  );
}

function RadarSVG({ scores, labels, fontSans }) {
  const width = 900;
  const height = 600;
  const cx = width / 2;
  const cy = height / 2;
  const r = 230;
  const keys = Object.keys(labels);
  const n = keys.length;
  const topKey = Object.keys(scores).sort((a, b) => scores[b] - scores[a])[0];

  const angleFor = (i) => (i * 2 * Math.PI) / n - Math.PI / 2;
  const pointFor = (i, pct) => {
    const a = angleFor(i);
    return [cx + Math.cos(a) * r * pct, cy + Math.sin(a) * r * pct];
  };

  const dataPolygon = keys
    .map((k, i) => pointFor(i, (scores[k] ?? 0) / 100).join(","))
    .join(" ");
  const gridLevels = [0.25, 0.5, 0.75, 1];

  const labelBoxW = 320;
  const labelBoxH = 40;

  return (
    <div style={{ position: "relative", width, height, display: "flex" }}>
      <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`}>
        <defs>
          <linearGradient id="radarFill" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#a78bfa" stopOpacity="0.55" />
            <stop offset="100%" stopColor="#f0abfc" stopOpacity="0.25" />
          </linearGradient>
          <linearGradient id="radarStroke" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#ddd6fe" />
            <stop offset="100%" stopColor="#f0abfc" />
          </linearGradient>
        </defs>

        {gridLevels.map((lvl, idx) => {
          const pts = keys.map((_, i) => pointFor(i, lvl).join(",")).join(" ");
          return (
            <polygon
              key={idx}
              points={pts}
              fill="none"
              stroke="#6d28d9"
              strokeWidth={1}
              strokeOpacity={lvl === 1 ? 0.5 : 0.25}
            />
          );
        })}

        {keys.map((_, i) => {
          const [x, y] = pointFor(i, 1);
          return (
            <line
              key={i}
              x1={cx}
              y1={cy}
              x2={x}
              y2={y}
              stroke="#6d28d9"
              strokeWidth={1}
              strokeOpacity={0.25}
            />
          );
        })}

        <polygon
          points={dataPolygon}
          fill="url(#radarFill)"
          stroke="url(#radarStroke)"
          strokeWidth={2.5}
          strokeLinejoin="round"
        />

        {keys.map((k, i) => {
          const [x, y] = pointFor(i, (scores[k] ?? 0) / 100);
          return <circle key={k} cx={x} cy={y} r={4.5} fill="#f0abfc" />;
        })}
      </svg>

      {keys.map((k, i) => {
        const a = angleFor(i);
        const cos = Math.cos(a);
        const sin = Math.sin(a);
        const [x, y] = pointFor(i, 1.08);
        const anchor = Math.abs(cos) < 0.2 ? "center" : cos > 0 ? "flex-start" : "flex-end";
        const left =
          anchor === "flex-end" ? x - labelBoxW : anchor === "center" ? x - labelBoxW / 2 : x;
        const dy = sin > 0.5 ? 10 : sin < -0.5 ? -6 : 0;
        const isTop = k === topKey;
        return (
          <div
            key={k}
            style={{
              position: "absolute",
              left,
              top: y - labelBoxH / 2 + dy,
              width: labelBoxW,
              height: labelBoxH,
              display: "flex",
              alignItems: "center",
              justifyContent: anchor,
              fontFamily: fontSans,
              fontSize: 18,
              fontWeight: isTop ? 600 : 500,
              color: isTop ? "#f0abfc" : "#c4b5fd",
            }}
          >
            {labels[k]}
          </div>
        );
      })}
    </div>
  );
}
