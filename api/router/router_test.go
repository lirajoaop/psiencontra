package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/handler"
	"github.com/joaop/psiencontra/api/internal/memrepo"
	"github.com/joaop/psiencontra/api/schemas"
	"github.com/joaop/psiencontra/api/service"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	discard := log.New(io.Discard, "", 0)
	config.Log = &config.Logger{Info: discard, Error: discard}
	os.Exit(m.Run())
}

type noopAnalyzer struct{}

func (noopAnalyzer) Analyze(string) (*service.AIResult, string, error) {
	return &service.AIResult{}, "gemini", nil
}

// newTestRouter wires the real routes and middleware to in-memory services.
func newTestRouter(t *testing.T) (*gin.Engine, *memrepo.Store) {
	t.Helper()
	store := memrepo.New()
	handler.QuestionSvc = service.NewQuestionService()
	handler.PDFSvc = service.NewPDFService()
	handler.SessionSvc = service.NewSessionService(store.Sessions(), store.Responses(), store.Results(), handler.QuestionSvc, noopAnalyzer{})
	handler.AuthSvc = service.NewAuthService(store.Users(), "router-test-secret")
	handler.GoogleOAuthSvc = service.NewGoogleOAuthService("", "", "")

	r := gin.New()
	SetupRoutes(r)
	return r, store
}

func serve(r *gin.Engine, method, path, remoteAddr string, mutate func(*http.Request)) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if remoteAddr != "" {
		req.RemoteAddr = remoteAddr
	}
	if mutate != nil {
		mutate(req)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Each scope has its own burst; once exhausted the API answers 429 with a
// JSON error and the handler is not reached.
func TestRateLimit_PerScope(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		burst  int
	}{
		{"global scope", http.MethodGet, "/api/v1/health", 10},
		{"auth scope", http.MethodPost, "/api/v1/auth/logout", 5},
		{"answer submission scope", http.MethodPost, "/api/v1/sessions/not-a-uuid/responses", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := newTestRouter(t)
			for i := 0; i < tt.burst; i++ {
				if w := serve(r, tt.method, tt.path, "", nil); w.Code == http.StatusTooManyRequests {
					t.Fatalf("request %d was rate limited before the burst of %d", i+1, tt.burst)
				}
			}
			w := serve(r, tt.method, tt.path, "", nil)
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("request %d: status = %d, want 429", tt.burst+1, w.Code)
			}
			var body map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body["error"] == "" {
				t.Errorf("429 body = %s, want JSON error", w.Body.String())
			}

			// Another client is unaffected.
			if w := serve(r, tt.method, tt.path, "198.51.100.20:4321", nil); w.Code == http.StatusTooManyRequests {
				t.Error("a different IP was rate limited")
			}
		})
	}
}

// The stricter auth limit must not consume the global budget of other
// routes beyond the requests actually made.
func TestRateLimit_ScopesAreIndependent(t *testing.T) {
	r, _ := newTestRouter(t)
	for i := 0; i < 5; i++ {
		serve(r, http.MethodPost, "/api/v1/auth/logout", "", nil)
	}
	if w := serve(r, http.MethodPost, "/api/v1/auth/logout", "", nil); w.Code != http.StatusTooManyRequests {
		t.Fatalf("auth scope not exhausted: %d", w.Code)
	}
	if w := serve(r, http.MethodGet, "/api/v1/health", "", nil); w.Code != http.StatusOK {
		t.Errorf("health blocked after auth limit: %d", w.Code)
	}
}

func TestRequireAuth(t *testing.T) {
	r, _ := newTestRouter(t)
	user := uuid.New()
	token, err := handler.AuthSvc.IssueToken(user)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*http.Request)
		status int
	}{
		{"no credentials", nil, http.StatusUnauthorized},
		{"invalid bearer", func(r *http.Request) { r.Header.Set("Authorization", "Bearer nope") }, http.StatusUnauthorized},
		{"bearer token", func(r *http.Request) { r.Header.Set("Authorization", "Bearer "+token) }, http.StatusOK},
		{"cookie token", func(r *http.Request) {
			r.AddCookie(&http.Cookie{Name: handler.AuthCookieName, Value: token})
		}, http.StatusOK},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Distinct IPs keep the global limiter out of the way.
			addr := "192.0.2." + itoa(i+10) + ":1000"
			w := serve(r, http.MethodGet, "/api/v1/user/sessions", addr, tt.mutate)
			if w.Code != tt.status {
				t.Fatalf("status = %d, want %d; body %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

// The anonymous flow: a session created without credentials has no owner;
// after logging in, the user claims it through RequireAuth.
func TestAnonymousSessionClaimFlow(t *testing.T) {
	r, store := newTestRouter(t)

	w := serve(r, http.MethodPost, "/api/v1/sessions", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("create session: %d %s", w.Code, w.Body.String())
	}
	var created struct {
		Data struct {
			ID uuid.UUID `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	sess, _ := store.Sessions().FindByID(created.Data.ID)
	if sess.UserID != nil {
		t.Fatalf("anonymous session has owner %v", sess.UserID)
	}

	claimPath := "/api/v1/sessions/" + created.Data.ID.String() + "/claim"
	if w := serve(r, http.MethodPost, claimPath, "", nil); w.Code != http.StatusUnauthorized {
		t.Fatalf("claim without auth: %d, want 401", w.Code)
	}

	user := &schemas.User{Email: "ana@example.com"}
	if err := store.Users().Create(user); err != nil {
		t.Fatal(err)
	}
	token, _ := handler.AuthSvc.IssueToken(user.ID)
	bearer := func(r *http.Request) { r.Header.Set("Authorization", "Bearer "+token) }

	if w := serve(r, http.MethodPost, claimPath, "", bearer); w.Code != http.StatusOK {
		t.Fatalf("claim: %d %s", w.Code, w.Body.String())
	}
	sess, _ = store.Sessions().FindByID(created.Data.ID)
	if sess.UserID == nil || *sess.UserID != user.ID {
		t.Fatalf("owner = %v, want %v", sess.UserID, user.ID)
	}

	if w := serve(r, http.MethodPost, claimPath, "", bearer); w.Code != http.StatusConflict {
		t.Errorf("second claim: %d, want 409", w.Code)
	}
}

// OptionalAuth attaches the owner when a valid token is sent and silently
// ignores invalid tokens instead of rejecting the request.
func TestOptionalAuth_CreateSession(t *testing.T) {
	user := uuid.New()
	tests := []struct {
		name      string
		token     func() string
		wantOwner bool
	}{
		{"valid token sets owner", func() string { tok, _ := handler.AuthSvc.IssueToken(user); return tok }, true},
		{"invalid token is ignored", func() string { return "garbage" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, store := newTestRouter(t)
			tok := tt.token()
			w := serve(r, http.MethodPost, "/api/v1/sessions", "", func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+tok)
			})
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d", w.Code)
			}
			var created struct {
				Data struct {
					ID uuid.UUID `json:"id"`
				} `json:"data"`
			}
			_ = json.Unmarshal(w.Body.Bytes(), &created)
			sess, _ := store.Sessions().FindByID(created.Data.ID)
			if (sess.UserID != nil) != tt.wantOwner {
				t.Errorf("owner = %v, wantOwner %v", sess.UserID, tt.wantOwner)
			}
			if tt.wantOwner && *sess.UserID != user {
				t.Errorf("owner = %v, want %v", *sess.UserID, user)
			}
		})
	}
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}
