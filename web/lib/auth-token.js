// Centralized accessors for the auth token in localStorage. SSR-safe: every
// helper checks for `window` because Next.js renders this code on the server
// during the initial pass, where `localStorage` doesn't exist.
//
// Why localStorage instead of an httpOnly cookie: cross-site cookies are
// blocked by Safari, mobile WebKit and incognito modes when the frontend
// (Vercel) and backend (Railway) live on different domains. localStorage works
// everywhere. The trade-off is XSS exposure — acceptable here because the app
// renders no user-generated HTML and has no third-party script soup.

const TOKEN_KEY = "psiencontra_auth";

export function getToken() {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setToken(token) {
  if (typeof window === "undefined") return;
  if (!token) {
    window.localStorage.removeItem(TOKEN_KEY);
    return;
  }
  window.localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(TOKEN_KEY);
}
