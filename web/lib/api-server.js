import { cache } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export const getResultServer = cache(async (id) => {
  try {
    const res = await fetch(`${API_URL}/sessions/${id}/result`, {
      cache: "no-store",
    });
    if (!res.ok) return null;
    const body = await res.json().catch(() => null);
    return body?.data ?? null;
  } catch {
    return null;
  }
});
