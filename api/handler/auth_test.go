package handler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/service"
)

func newDiscardLogger() *log.Logger { return log.New(io.Discard, "", 0) }

func authCookie(w *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == AuthCookieName {
			return c
		}
	}
	return nil
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"success", `{"email":"Ana@Example.com","password":"supersecret","name":"Ana"}`, http.StatusOK},
		{"invalid email", `{"email":"nope","password":"supersecret"}`, http.StatusBadRequest},
		{"short password", `{"email":"ana@example.com","password":"short"}`, http.StatusBadRequest},
		{"malformed JSON", `{`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setup(t)
			w := env.do(t, req{method: http.MethodPost, path: "/auth/register", body: tt.body})
			assertStatus(t, w, tt.status)
			if tt.status != http.StatusOK {
				if authCookie(w) != nil {
					t.Error("auth cookie set on failed registration")
				}
				return
			}
			data := decode(t, w)["data"].(map[string]any)
			token, _ := data["token"].(string)
			if token == "" {
				t.Fatal("token missing from body")
			}
			c := authCookie(w)
			if c == nil || c.Value != token || !c.HttpOnly {
				t.Errorf("auth cookie = %+v, want HttpOnly cookie with the token", c)
			}
			user := data["user"].(map[string]any)
			if user["email"] != "ana@example.com" {
				t.Errorf("email not normalized: %v", user["email"])
			}
			if _, leaked := user["password_hash"]; leaked {
				t.Error("password hash leaked in response")
			}
		})
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	env := setup(t)
	body := `{"email":"ana@example.com","password":"supersecret"}`
	assertStatus(t, env.do(t, req{method: http.MethodPost, path: "/auth/register", body: body}), http.StatusOK)
	w := env.do(t, req{method: http.MethodPost, path: "/auth/register", body: body})
	assertStatus(t, w, http.StatusConflict)
}

func TestLogin(t *testing.T) {
	env := setup(t)
	if _, _, err := AuthSvc.Register("ana@example.com", "supersecret", "Ana"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"success", `{"email":"ana@example.com","password":"supersecret"}`, http.StatusOK},
		{"wrong password", `{"email":"ana@example.com","password":"wrong-password"}`, http.StatusUnauthorized},
		{"unknown user", `{"email":"bob@example.com","password":"supersecret"}`, http.StatusUnauthorized},
		{"missing password", `{"email":"ana@example.com"}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := env.do(t, req{method: http.MethodPost, path: "/auth/login", body: tt.body})
			assertStatus(t, w, tt.status)
			if (authCookie(w) != nil) != (tt.status == http.StatusOK) {
				t.Errorf("auth cookie presence mismatch for status %d", tt.status)
			}
		})
	}
}

func TestLogout_ClearsCookie(t *testing.T) {
	env := setup(t)
	w := env.do(t, req{method: http.MethodPost, path: "/auth/logout"})
	assertStatus(t, w, http.StatusOK)
	c := authCookie(w)
	if c == nil || c.MaxAge >= 0 || c.Value != "" {
		t.Errorf("cookie = %+v, want an expired empty cookie", c)
	}
}

func TestMe(t *testing.T) {
	env := setup(t)
	user, token, err := AuthSvc.Register("ana@example.com", "supersecret", "Ana")
	if err != nil {
		t.Fatal(err)
	}
	ghostToken, _ := AuthSvc.IssueToken(uuid.New()) // valid JWT, deleted user
	foreignToken, _ := service.NewAuthService(nil, "other-secret").IssueToken(user.ID)

	tests := []struct {
		name        string
		cookie      string
		bearer      string
		wantUser    bool
		wantCleared bool
	}{
		{name: "no token", wantUser: false},
		{name: "cookie token", cookie: token, wantUser: true},
		{name: "bearer token", bearer: token, wantUser: true},
		{name: "token signed with another secret", bearer: foreignToken, wantCleared: true},
		{name: "garbage token", cookie: "garbage", wantCleared: true},
		{name: "user no longer exists", bearer: ghostToken, wantCleared: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := req{method: http.MethodGet, path: "/auth/me"}
			if tt.cookie != "" {
				r.cookies = []*http.Cookie{{Name: AuthCookieName, Value: tt.cookie}}
			}
			if tt.bearer != "" {
				r.headers = map[string]string{"Authorization": "Bearer " + tt.bearer}
			}
			w := env.do(t, r)
			assertStatus(t, w, http.StatusOK)
			data := decode(t, w)["data"]
			if tt.wantUser {
				u, ok := data.(map[string]any)
				if !ok || u["id"] != user.ID.String() {
					t.Fatalf("data = %v, want user %s", data, user.ID)
				}
			} else if data != nil {
				t.Fatalf("data = %v, want null", data)
			}
			if tt.wantCleared {
				if c := authCookie(w); c == nil || c.MaxAge >= 0 {
					t.Errorf("invalid token should clear the cookie, got %+v", c)
				}
			}
		})
	}
}

func TestReadToken(t *testing.T) {
	tests := []struct {
		name   string
		cookie string
		header string
		want   string
	}{
		{"none", "", "", ""},
		{"cookie only", "from-cookie", "", "from-cookie"},
		{"bearer only", "", "Bearer from-header", "from-header"},
		{"cookie wins over bearer", "from-cookie", "Bearer from-header", "from-cookie"},
		{"non-bearer scheme ignored", "", "Basic dXNlcjpwYXNz", ""},
		{"lowercase bearer ignored", "", "bearer from-header", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != "" {
				r.AddCookie(&http.Cookie{Name: AuthCookieName, Value: tt.cookie})
			}
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
			}
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = r
			if got := ReadToken(c); got != tt.want {
				t.Errorf("ReadToken = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriteCookie_SecurityFlags(t *testing.T) {
	tests := []struct {
		name         string
		appEnv       string
		cookieSecure string
		wantSecure   bool
		wantSameSite http.SameSite
	}{
		{"production default", "", "", true, http.SameSiteNoneMode},
		{"development", "development", "", false, http.SameSiteLaxMode},
		{"explicit override wins", "development", "true", true, http.SameSiteNoneMode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.appEnv)
			t.Setenv("COOKIE_SECURE", tt.cookieSecure)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			writeCookie(c, "x", "v", 60)

			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("got %d cookies", len(cookies))
			}
			got := cookies[0]
			if got.Secure != tt.wantSecure || got.SameSite != tt.wantSameSite || !got.HttpOnly {
				t.Errorf("cookie = secure:%v samesite:%v httponly:%v, want secure:%v samesite:%v httponly:true",
					got.Secure, got.SameSite, got.HttpOnly, tt.wantSecure, tt.wantSameSite)
			}
		})
	}
}

func TestGoogleOAuth_DisabledRedirectsWithError(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://front.example")
	env := setup(t)
	for _, path := range []string{"/auth/google", "/auth/google/callback?state=x&code=y"} {
		t.Run(path, func(t *testing.T) {
			w := env.do(t, req{method: http.MethodGet, path: path})
			assertStatus(t, w, http.StatusTemporaryRedirect)
			if loc := w.Header().Get("Location"); loc != "https://front.example/auth/callback?error=google_disabled" {
				t.Errorf("Location = %q", loc)
			}
		})
	}
}

func TestGoogleOAuth_StateValidation(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://front.example")
	env := setup(t)
	GoogleOAuthSvc = service.NewGoogleOAuthService("client-id", "client-secret", "https://api.example/callback")

	t.Run("start sets state cookie and redirects to Google", func(t *testing.T) {
		w := env.do(t, req{method: http.MethodGet, path: "/auth/google"})
		assertStatus(t, w, http.StatusTemporaryRedirect)
		var state string
		for _, c := range w.Result().Cookies() {
			if c.Name == stateCookieName {
				state = c.Value
			}
		}
		if state == "" {
			t.Fatal("state cookie not set")
		}
		loc, err := url.Parse(w.Header().Get("Location"))
		if err != nil || !strings.Contains(loc.Host, "google") {
			t.Fatalf("Location = %q, want Google", w.Header().Get("Location"))
		}
		if loc.Query().Get("state") != state {
			t.Errorf("state in URL %q != cookie %q", loc.Query().Get("state"), state)
		}
	})

	tests := []struct {
		name      string
		query     string
		cookie    string
		wantError string
	}{
		{"missing state cookie", "?state=abc&code=c", "", "invalid_state"},
		{"state mismatch", "?state=abc&code=c", "xyz", "invalid_state"},
		{"provider error", "?state=abc&error=access_denied", "abc", "access_denied"},
		{"missing code", "?state=abc", "abc", "missing_code"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := req{method: http.MethodGet, path: "/auth/google/callback" + tt.query}
			if tt.cookie != "" {
				r.cookies = []*http.Cookie{{Name: stateCookieName, Value: tt.cookie}}
			}
			w := env.do(t, r)
			assertStatus(t, w, http.StatusTemporaryRedirect)
			want := "https://front.example/auth/callback?error=" + tt.wantError
			if loc := w.Header().Get("Location"); loc != want {
				t.Errorf("Location = %q, want %q", loc, want)
			}
			if authCookie(w) != nil {
				t.Error("auth cookie must not be set on a failed callback")
			}
		})
	}
}
