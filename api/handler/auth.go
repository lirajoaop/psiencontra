package handler

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/repository"
	"github.com/joaop/psiencontra/api/service"
)

const (
	authCookieName  = "psiencontra_auth"
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
		default:
			sendError(c, http.StatusBadRequest, err.Error())
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
	token := readTokenFromRequest(c)
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
	writeCookie(c, authCookieName, token, maxAge)
}

func clearAuthCookie(c *gin.Context) {
	writeCookie(c, authCookieName, "", -1)
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
func writeCookie(c *gin.Context, name, value string, maxAge int) {
	secure := config.GetEnv("COOKIE_SECURE", "true") == "true"
	sameSite := http.SameSiteNoneMode
	if !secure {
		// Local dev over plain HTTP: browsers reject SameSite=None without Secure.
		sameSite = http.SameSiteLaxMode
	}
	c.SetSameSite(sameSite)
	domain := config.GetEnv("COOKIE_DOMAIN", "")
	c.SetCookie(name, value, maxAge, "/", domain, secure, true)
}

func readTokenFromRequest(c *gin.Context) string {
	if cookie, err := c.Cookie(authCookieName); err == nil && cookie != "" {
		return cookie
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
