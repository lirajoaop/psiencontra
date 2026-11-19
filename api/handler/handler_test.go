package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/internal/memrepo"
	"github.com/joaop/psiencontra/api/schemas"
	"github.com/joaop/psiencontra/api/service"
)

const testJWTSecret = "handler-test-secret"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config.Log = &config.Logger{
		Info:  newDiscardLogger(),
		Error: newDiscardLogger(),
	}
	os.Exit(m.Run())
}

// stubAnalyzer stands in for the Gemini/Groq client.
type stubAnalyzer struct {
	err error
}

func (a *stubAnalyzer) Analyze(string) (*service.AIResult, string, error) {
	if a.err != nil {
		return nil, "", a.err
	}
	return &service.AIResult{
		ApproachScores: json.RawMessage(`{"psicanalise":80}`),
		FieldScores:    json.RawMessage(`{"clinica":70}`),
		Summary:        "resumo",
	}, "gemini", nil
}

type testEnv struct {
	store  *memrepo.Store
	ai     *stubAnalyzer
	engine *gin.Engine
}

// setup wires the package-level services to in-memory fakes and returns an
// engine exposing the handlers. The X-Test-User header simulates what the
// auth middleware (tested in the router package) puts in the context.
func setup(t *testing.T) *testEnv {
	t.Helper()
	store := memrepo.New()
	ai := &stubAnalyzer{}
	QuestionSvc = service.NewQuestionService()
	PDFSvc = service.NewPDFService()
	SessionSvc = service.NewSessionService(store.Sessions(), store.Responses(), store.Results(), QuestionSvc, ai)
	AuthSvc = service.NewAuthService(store.Users(), testJWTSecret)
	GoogleOAuthSvc = service.NewGoogleOAuthService("", "", "")

	r := gin.New()
	r.Use(func(c *gin.Context) {
		if h := c.GetHeader("X-Test-User"); h != "" {
			c.Set("user_id", uuid.MustParse(h))
		}
		c.Next()
	})
	r.GET("/health", HealthCheck)
	r.GET("/questions", GetQuestions)
	r.POST("/sessions", CreateSession)
	r.POST("/sessions/:id/responses", SubmitResponses)
	r.POST("/sessions/:id/claim", ClaimSession)
	r.GET("/sessions/:id/result", GetResult)
	r.GET("/sessions/:id/pdf", DownloadPDF)
	r.GET("/user/sessions", GetUserHistory)
	r.POST("/auth/register", Register)
	r.POST("/auth/login", Login)
	r.POST("/auth/logout", Logout)
	r.GET("/auth/me", Me)
	r.GET("/auth/google", GoogleStart)
	r.GET("/auth/google/callback", GoogleCallback)

	return &testEnv{store: store, ai: ai, engine: r}
}

type req struct {
	method  string
	path    string
	body    string
	user    *uuid.UUID
	headers map[string]string
	cookies []*http.Cookie
}

func (e *testEnv) do(t *testing.T, r req) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if r.body != "" {
		body = strings.NewReader(r.body)
	}
	httpReq := httptest.NewRequest(r.method, r.path, body)
	if r.body != "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	if r.user != nil {
		httpReq.Header.Set("X-Test-User", r.user.String())
	}
	for k, v := range r.headers {
		httpReq.Header.Set(k, v)
	}
	for _, c := range r.cookies {
		httpReq.AddCookie(c)
	}
	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, httpReq)
	return w
}

// newSession creates a session directly in the store.
func (e *testEnv) newSession(t *testing.T, owner *uuid.UUID, qType string) uuid.UUID {
	t.Helper()
	s := &schemas.Session{UserID: owner, QuestionnaireType: qType}
	if err := e.store.Sessions().Create(s); err != nil {
		t.Fatal(err)
	}
	return s.ID
}

// newResult stores a finished result for the given session.
func (e *testEnv) newResult(t *testing.T, sessionID uuid.UUID) {
	t.Helper()
	err := e.store.Results().Create(&schemas.Result{
		SessionID:      sessionID,
		ApproachScores: json.RawMessage(`{"psicanalise":95,"tcc":5}`),
		FieldScores:    json.RawMessage(`{"clinica":95}`),
		Explanation:    "resumo",
		AIProvider:     "groq",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func decode(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v (%s)", err, w.Body.String())
	}
	return out
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, want, w.Body.String())
	}
}

func TestHealthCheck(t *testing.T) {
	env := setup(t)
	w := env.do(t, req{method: http.MethodGet, path: "/health"})
	assertStatus(t, w, http.StatusOK)
	if decode(t, w)["status"] != "ok" {
		t.Errorf("body = %s", w.Body.String())
	}
}

func TestGetQuestions(t *testing.T) {
	env := setup(t)
	tests := []struct {
		query string
		want  int
	}{
		{"", len(QuestionSvc.GetAll())},
		{"?type=simple", len(QuestionSvc.GetAll())},
		{"?type=detailed", 76},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			w := env.do(t, req{method: http.MethodGet, path: "/questions" + tt.query})
			assertStatus(t, w, http.StatusOK)
			data := decode(t, w)["data"].([]any)
			if len(data) != tt.want {
				t.Errorf("got %d questions, want %d", len(data), tt.want)
			}
		})
	}
}

func TestCreateSession(t *testing.T) {
	user := uuid.New()
	tests := []struct {
		name      string
		body      string
		user      *uuid.UUID
		repoErr   error
		status    int
		wantType  string
		wantOwner bool
	}{
		{name: "anonymous, no body", status: http.StatusOK, wantType: "simple"},
		{name: "detailed", body: `{"questionnaire_type":"detailed"}`, status: http.StatusOK, wantType: "detailed"},
		{name: "malformed body falls back to simple", body: `{oops`, status: http.StatusOK, wantType: "simple"},
		{name: "logged-in user owns the session", user: &user, status: http.StatusOK, wantType: "simple", wantOwner: true},
		{name: "repository failure", repoErr: errors.New("db down"), status: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setup(t)
			env.store.CreateSessionErr = tt.repoErr
			w := env.do(t, req{method: http.MethodPost, path: "/sessions", body: tt.body, user: tt.user})
			assertStatus(t, w, tt.status)
			if tt.status != http.StatusOK {
				return
			}
			id := uuid.MustParse(decode(t, w)["data"].(map[string]any)["id"].(string))
			sess, err := SessionSvc.GetSession(id)
			if err != nil {
				t.Fatal(err)
			}
			if sess.QuestionnaireType != tt.wantType {
				t.Errorf("type = %q, want %q", sess.QuestionnaireType, tt.wantType)
			}
			if (sess.UserID != nil) != tt.wantOwner {
				t.Errorf("owner = %v, wantOwner %v", sess.UserID, tt.wantOwner)
			}
		})
	}
}

func TestSubmitResponses(t *testing.T) {
	detailedQ := service.NewQuestionService().GetAllDetailed()[0]
	validBody := `{"responses":[{"question_id":` + itoa(detailedQ.ID) + `,"answer_value":"5"}]}`

	tests := []struct {
		name    string
		path    func(sessionID uuid.UUID) string
		body    string
		aiErr   error
		status  int
		errPart string
	}{
		{
			name:   "success",
			path:   func(id uuid.UUID) string { return "/sessions/" + id.String() + "/responses" },
			body:   validBody,
			status: http.StatusOK,
		},
		{
			name:    "invalid session id",
			path:    func(uuid.UUID) string { return "/sessions/not-a-uuid/responses" },
			body:    validBody,
			status:  http.StatusBadRequest,
			errPart: "invalid session id",
		},
		{
			name:    "missing responses field",
			path:    func(id uuid.UUID) string { return "/sessions/" + id.String() + "/responses" },
			body:    `{}`,
			status:  http.StatusBadRequest,
			errPart: "invalid request body",
		},
		{
			name:    "answer over 500 characters",
			path:    func(id uuid.UUID) string { return "/sessions/" + id.String() + "/responses" },
			body:    `{"responses":[{"question_id":1,"answer_value":"` + strings.Repeat("á", 501) + `"}]}`,
			status:  http.StatusBadRequest,
			errPart: "500 character limit",
		},
		{
			name:    "AI providers down",
			path:    func(id uuid.UUID) string { return "/sessions/" + id.String() + "/responses" },
			body:    validBody,
			aiErr:   errors.New("all AI providers failed"),
			status:  http.StatusInternalServerError,
			errPart: "AI analysis failed",
		},
		{
			name:    "unknown session",
			path:    func(uuid.UUID) string { return "/sessions/" + uuid.NewString() + "/responses" },
			body:    validBody,
			status:  http.StatusInternalServerError,
			errPart: "session not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setup(t)
			env.ai.err = tt.aiErr
			id := env.newSession(t, nil, "detailed")

			w := env.do(t, req{method: http.MethodPost, path: tt.path(id), body: tt.body})
			assertStatus(t, w, tt.status)
			if tt.errPart != "" {
				if msg, _ := decode(t, w)["error"].(string); !strings.Contains(msg, tt.errPart) {
					t.Errorf("error = %q, want containing %q", msg, tt.errPart)
				}
				return
			}
			data := decode(t, w)["data"].(map[string]any)
			if data["ai_provider"] != "gemini" {
				t.Errorf("ai_provider = %v", data["ai_provider"])
			}
		})
	}
}

func TestClaimSession(t *testing.T) {
	alice, bob := uuid.New(), uuid.New()
	tests := []struct {
		name   string
		owner  *uuid.UUID
		path   func(uuid.UUID) string
		user   *uuid.UUID
		status int
	}{
		{"claims anonymous session", nil, sessionPath("claim"), &alice, http.StatusOK},
		{"unauthenticated", nil, sessionPath("claim"), nil, http.StatusUnauthorized},
		{"invalid id", nil, func(uuid.UUID) string { return "/sessions/x/claim" }, &alice, http.StatusBadRequest},
		{"already owned", &bob, sessionPath("claim"), &alice, http.StatusConflict},
		{"unknown session", nil, func(uuid.UUID) string { return "/sessions/" + uuid.NewString() + "/claim" }, &alice, http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setup(t)
			id := env.newSession(t, tt.owner, "simple")
			w := env.do(t, req{method: http.MethodPost, path: tt.path(id), user: tt.user})
			assertStatus(t, w, tt.status)
			if tt.status == http.StatusOK {
				sess, _ := SessionSvc.GetSession(id)
				if sess.UserID == nil || *sess.UserID != *tt.user {
					t.Errorf("owner = %v, want %v", sess.UserID, *tt.user)
				}
			}
		})
	}
}

func TestGetResult(t *testing.T) {
	tests := []struct {
		name       string
		withResult bool
		path       func(uuid.UUID) string
		status     int
	}{
		{"found", true, sessionPath("result"), http.StatusOK},
		{"no result yet", false, sessionPath("result"), http.StatusNotFound},
		{"invalid id", false, func(uuid.UUID) string { return "/sessions/123/result" }, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setup(t)
			id := env.newSession(t, nil, "detailed")
			if tt.withResult {
				env.newResult(t, id)
			}
			w := env.do(t, req{method: http.MethodGet, path: tt.path(id)})
			assertStatus(t, w, tt.status)
			if tt.status != http.StatusOK {
				return
			}
			data := decode(t, w)["data"].(map[string]any)
			if data["questionnaire_type"] != "detailed" || data["session_id"] != id.String() {
				t.Errorf("unexpected payload: %v", data)
			}
		})
	}
}

func TestDownloadPDF(t *testing.T) {
	tests := []struct {
		name       string
		withResult bool
		path       func(uuid.UUID) string
		status     int
	}{
		{"generates pdf", true, sessionPath("pdf"), http.StatusOK},
		{"no result", false, sessionPath("pdf"), http.StatusNotFound},
		{"invalid id", false, func(uuid.UUID) string { return "/sessions/zzz/pdf" }, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setup(t)
			id := env.newSession(t, nil, "simple")
			if tt.withResult {
				env.newResult(t, id)
			}
			w := env.do(t, req{method: http.MethodGet, path: tt.path(id)})
			assertStatus(t, w, tt.status)
			if tt.status != http.StatusOK {
				return
			}
			if ct := w.Header().Get("Content-Type"); ct != "application/pdf" {
				t.Errorf("Content-Type = %q", ct)
			}
			wantName := "psiencontra-" + id.String()[:8] + ".pdf"
			if cd := w.Header().Get("Content-Disposition"); !strings.Contains(cd, wantName) {
				t.Errorf("Content-Disposition = %q, want filename %s", cd, wantName)
			}
			if !strings.HasPrefix(w.Body.String(), "%PDF-") {
				t.Error("body is not a PDF")
			}
		})
	}
}

func TestGetUserHistory(t *testing.T) {
	env := setup(t)
	alice, bob := uuid.New(), uuid.New()
	done := env.newSession(t, &alice, "detailed")
	env.newResult(t, done)
	env.newSession(t, &alice, "detailed") // abandoned
	other := env.newSession(t, &bob, "simple")
	env.newResult(t, other)

	t.Run("unauthenticated", func(t *testing.T) {
		w := env.do(t, req{method: http.MethodGet, path: "/user/sessions"})
		assertStatus(t, w, http.StatusUnauthorized)
	})

	t.Run("lists only own completed sessions", func(t *testing.T) {
		w := env.do(t, req{method: http.MethodGet, path: "/user/sessions", user: &alice})
		assertStatus(t, w, http.StatusOK)
		data := decode(t, w)["data"].([]any)
		if len(data) != 1 {
			t.Fatalf("got %d sessions, want 1: %v", len(data), data)
		}
		item := data[0].(map[string]any)
		if item["id"] != done.String() {
			t.Errorf("id = %v, want %s", item["id"], done)
		}
		if _, leaked := item["explanation"]; leaked {
			t.Error("history DTO should not include the explanation")
		}
	})

	t.Run("repository failure", func(t *testing.T) {
		env.store.HistoryErr = errors.New("db down")
		defer func() { env.store.HistoryErr = nil }()
		w := env.do(t, req{method: http.MethodGet, path: "/user/sessions", user: &alice})
		assertStatus(t, w, http.StatusInternalServerError)
	})
}

func sessionPath(suffix string) func(uuid.UUID) string {
	return func(id uuid.UUID) string { return "/sessions/" + id.String() + "/" + suffix }
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}
