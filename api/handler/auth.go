package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/repository"
	"github.com/joaop/psiencontra/api/service"
)

const (
	// AuthCookieName is the name of the cookie carrying the JWT. Exported
	// so the router middleware can share the same constant instead of
	// hard-coding a duplicate string.
	AuthCookieName  = "psiencontra_auth"
	stateCookieName = "psiencontra_oauth_state"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := AuthSvc.Register(req.Email, req.Password, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyUsed):
			sendError(c, http.StatusConflict, "email already in use")
		case errors.Is(err, service.ErrEmailRequired):
			sendError(c, http.StatusBadRequest, "email required")
		case errors.Is(err, service.ErrPasswordTooShort):
			sendError(c, http.StatusBadRequest, "password must be at least 8 characters")
		default:
			// Unexpected: DB error, bcrypt failure, JWT signing failure, etc.
			// Log the real cause server-side and return a generic 500 so we
			// don't leak internal details (DB driver messages, stack hints,
			// secret-related errors) to the client.
			config.Log.Error.Printf("register failed: %v", err)
			sendError(c, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	setAuthCookie(c, token)
	sendSuccess(c, gin.H{"user": user})
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := AuthSvc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			sendError(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		config.Log.Error.Printf("login failed: %v", err)
		sendError(c, http.StatusInternalServerError, "login failed")
		return
	}

	setAuthCookie(c, token)
	sendSuccess(c, gin.H{"user": user})
}

func Logout(c *gin.Context) {
	clearAuthCookie(c)
	sendSuccess(c, gin.H{"ok": true})
}

func Me(c *gin.Context) {
	token := ReadToken(c)
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}
	claims, err := AuthSvc.ParseToken(token)
	if err != nil {
		clearAuthCookie(c)
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}
	user, err := AuthSvc.GetUser(claims.UserID)
	if err != nil {
		if repository.IsNotFound(err) {
			clearAuthCookie(c)
			c.JSON(http.StatusOK, gin.H{"data": nil})
			return
		}
		sendError(c, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	sendSuccess(c, user)
}

// GoogleStart begins the OAuth dance: stores a CSRF state cookie and
// redirects to Google. Errors here redirect back to the frontend callback
// with an error code instead of returning JSON, because the browser reaches
// this endpoint via a full-page navigation — a JSON body would render as a
// raw error page with no way back into the app.
func GoogleStart(c *gin.Context) {
	frontendURL := config.GetEnv("FRONTEND_URL", "http://localhost:3000")

	if !GoogleOAuthSvc.Enabled() {
		redirectWithError(c, frontendURL, "google_disabled")
		return
	}

	state, err := service.GenerateState()
	if err != nil {
		redirectWithError(c, frontendURL, "state_generation_failed")
		return
	}
	setStateCookie(c, state)
	c.Redirect(http.StatusTemporaryRedirect, GoogleOAuthSvc.AuthURL(state))
}

// GoogleCallback handles the redirect from Google: validates the state
// cookie, exchanges the code for a profile, upserts the user and sets the
// auth cookie before bouncing the browser back to the frontend.
func GoogleCallback(c *gin.Context) {
	frontendURL := config.GetEnv("FRONTEND_URL", "http://localhost:3000")

	if !GoogleOAuthSvc.Enabled() {
		redirectWithError(c, frontendURL, "google_disabled")
		return
	}

	stateParam := c.Query("state")
	stateCookie, err := c.Cookie(stateCookieName)
	if err != nil || stateParam == "" || stateParam != stateCookie {
		redirectWithError(c, frontendURL, "invalid_state")
		return
	}
	clearStateCookie(c)

	if errParam := c.Query("error"); errParam != "" {
		redirectWithError(c, frontendURL, errParam)
		return
	}

	code := c.Query("code")
	if code == "" {
		redirectWithError(c, frontendURL, "missing_code")
		return
	}

	info, err := GoogleOAuthSvc.Exchange(c.Request.Context(), code)
	if err != nil {
		redirectWithError(c, frontendURL, "exchange_failed")
		return
	}

	_, token, err := AuthSvc.UpsertGoogleUser(info.Sub, info.Email, info.Name, info.Picture)
	if err != nil {
		redirectWithError(c, frontendURL, "user_upsert_failed")
		return
	}

	setAuthCookie(c, token)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/callback?ok=1")
}

// --- helpers ---

func setAuthCookie(c *gin.Context, token string) {
	maxAge := int(AuthSvc.TokenLifetime().Seconds())
	writeCookie(c, AuthCookieName, token, maxAge)
}

func clearAuthCookie(c *gin.Context) {
	writeCookie(c, AuthCookieName, "", -1)
}

func setStateCookie(c *gin.Context, state string) {
	writeCookie(c, stateCookieName, state, 600)
}

func clearStateCookie(c *gin.Context) {
	writeCookie(c, stateCookieName, "", -1)
}

// writeCookie writes a cross-site capable cookie. In production we need
// SameSite=None + Secure so the browser sends the cookie on cross-origin
// requests from the Vercel frontend to the Railway API.
//
// Default policy mirrors the JWT_SECRET pattern: prod-safe by default,
// only relaxed for explicit dev mode. Local HTTP dev must run with
// APP_ENV=development (or COOKIE_SECURE=false), otherwise the browser
// rejects Secure cookies over plain HTTP and login silently fails.
// An explicit COOKIE_SECURE env always wins for edge cases.
func writeCookie(c *gin.Context, name, value string, maxAge int) {
	secure := config.GetEnv("APP_ENV", "") != "development"
	if v := config.GetEnv("COOKIE_SECURE", ""); v != "" {
		secure = v == "true"
	}
	sameSite := http.SameSiteNoneMode
	if !secure {
		// Local dev over plain HTTP: browsers reject SameSite=None without Secure.
		sameSite = http.SameSiteLaxMode
	}
	c.SetSameSite(sameSite)
	domain := config.GetEnv("COOKIE_DOMAIN", "")
	c.SetCookie(name, value, maxAge, "/", domain, secure, true)
}

// ReadToken extracts the JWT from either the auth cookie (preferred for
// browser clients) or the Authorization: Bearer header (for non-browser
// API clients). Centralized so /auth/me and the auth middleware stay
// consistent for all client types.
func ReadToken(c *gin.Context) string {
	if cookie, err := c.Cookie(AuthCookieName); err == nil && cookie != "" {
		return cookie
	}
	if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

func redirectWithError(c *gin.Context, frontendURL, code string) {
	u := frontendURL + "/auth/callback?error=" + url.QueryEscape(code)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

// UserIDFromContext returns the authenticated user ID set by middleware, if any.
func UserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
