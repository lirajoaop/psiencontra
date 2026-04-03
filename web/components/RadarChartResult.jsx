"use client";

import {
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

export default function RadarChartResult({ title, scores, labels, color = "#7c3aed" }) {
  const data = Object.entries(scores).map(([key, value]) => ({
    subject: labels[key] || key,
    value: value,
    fullMark: 100,
  }));

  return (
    <div className="bg-white rounded-2xl p-6 shadow-lg border border-violet-100">
      <h3 className="text-lg font-bold text-violet-900 mb-4 text-center">{title}</h3>
      <ResponsiveContainer width="100%" height={350}>
        <RadarChart data={data} cx="50%" cy="50%" outerRadius="70%">
          <PolarGrid stroke="#e5e7eb" />
          <PolarAngleAxis
            dataKey="subject"
            tick={{ fontSize: 10, fill: "#4b5563" }}
            className="text-xs"
          />
          <PolarRadiusAxis
            angle={90}
            domain={[0, 100]}
            tick={{ fontSize: 9, fill: "#9ca3af" }}
          />
          <Radar
            dataKey="value"
            stroke={color}
            fill={color}
            fillOpacity={0.25}
            strokeWidth={2}
          />
          <Tooltip
            formatter={(value) => [`${value}%`, "Afinidade"]}
            contentStyle={{
              backgroundColor: "white",
              border: "1px solid #e5e7eb",
              borderRadius: "8px",
              fontSize: "12px",
            }}
          />
        </RadarChart>
      </ResponsiveContainer>
    </div>
  );
}
