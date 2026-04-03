const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function request(path, options = {}) {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }

  const contentType = res.headers.get("content-type");
  if (contentType?.includes("application/pdf")) {
    return res.blob();
  }

  const json = await res.json();
  return json.data;
}

export function getQuestions() {
  return request("/questions");
}

export function createSession() {
  return request("/sessions", { method: "POST" });
}

export function submitResponses(sessionId, responses) {
  return request(`/sessions/${sessionId}/responses`, {
    method: "POST",
    body: JSON.stringify({ responses }),
  });
}

export function getResult(sessionId) {
  return request(`/sessions/${sessionId}/result`);
}

export function getPDFUrl(sessionId) {
  return `${API_URL}/sessions/${sessionId}/pdf`;
}
