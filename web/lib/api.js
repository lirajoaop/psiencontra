const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function request(path, options = {}) {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    ...options,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: "Request failed" }));
    const error = new Error(body.error || `HTTP ${res.status}`);
    // Attach the HTTP status so callers can branch on a stable contract
    // instead of matching backend display strings (which aren't a contract).
    error.status = res.status;
    throw error;
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

// --- Auth ---

export function register({ email, password, name }) {
  return request("/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password, name }),
  });
}

export function login({ email, password }) {
  return request("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export function logout() {
  return request("/auth/logout", { method: "POST" });
}

export function getMe() {
  return request("/auth/me");
}

export function getGoogleLoginURL() {
  return `${API_URL}/auth/google`;
}
