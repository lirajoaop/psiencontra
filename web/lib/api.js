import { getToken } from "@/lib/auth-token";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function request(path, options = {}) {
  // Build headers in a single object so callers passing custom headers in
  // `options` still get the auth token automatically. The Bearer token is the
  // primary auth path; we keep `credentials: include` only as a redundant
  // path for browsers that still allow cross-site cookies.
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };
  const token = getToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_URL}${path}`, {
    credentials: "include",
    ...options,
    headers,
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

export function getQuestions(type = "simple") {
  return request(`/questions?type=${type}`);
}

export function createSession(questionnaireType = "simple") {
  return request("/sessions", {
    method: "POST",
    body: JSON.stringify({ questionnaire_type: questionnaireType }),
  });
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
