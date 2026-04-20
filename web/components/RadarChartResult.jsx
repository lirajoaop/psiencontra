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
import { useTheme } from "./ThemeProvider";

export default function RadarChartResult({ title, scores, labels, color = "#7c3aed" }) {
  const { theme } = useTheme();
  const isDark = theme === "dark";

  const data = Object.entries(scores).map(([key, value]) => ({
    subject: labels[key] || key,
    value: value,
    fullMark: 100,
  }));

  return (
    <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg border border-violet-100 dark:border-gray-700">
      <h3 className="text-lg font-bold text-violet-900 dark:text-violet-200 mb-4 text-center">{title}</h3>
      <ResponsiveContainer width="100%" height={350}>
        <RadarChart data={data} cx="50%" cy="50%" outerRadius="62%">
          <PolarGrid stroke={isDark ? "#374151" : "#e5e7eb"} />
          <PolarAngleAxis
            dataKey="subject"
            tick={{ fontSize: 11, fill: isDark ? "#d1d5db" : "#4b5563" }}
            className="text-xs"
          />
          <PolarRadiusAxis
            angle={90}
            domain={[0, 100]}
            tick={{ fontSize: 9, fill: isDark ? "#6b7280" : "#9ca3af" }}
          />
          <Radar
            dataKey="value"
            stroke={color}
            fill={color}
            fillOpacity={isDark ? 0.35 : 0.25}
            strokeWidth={2}
          />
          <Tooltip
            formatter={(value) => [`${value}%`, "Afinidade"]}
            contentStyle={{
              backgroundColor: isDark ? "#1f2937" : "white",
              border: `1px solid ${isDark ? "#374151" : "#e5e7eb"}`,
              borderRadius: "8px",
              fontSize: "12px",
              color: isDark ? "#f3f4f6" : undefined,
            }}
          />
        </RadarChart>
      </ResponsiveContainer>
    </div>
  );
}
