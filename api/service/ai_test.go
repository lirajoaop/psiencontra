package service

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeProvider is an LLMProvider that returns a canned result or error and
// records how many times it was called.
type fakeProvider struct {
	name   string
	result *AIResult
	err    error
	calls  int
}

func (p *fakeProvider) Name() string { return p.name }

func (p *fakeProvider) Generate(string) (*AIResult, error) {
	p.calls++
	return p.result, p.err
}

func TestAIService_Analyze(t *testing.T) {
	okGemini := &AIResult{Summary: "from gemini"}
	okGroq := &AIResult{Summary: "from groq"}

	tests := []struct {
		name         string
		gemini       *fakeProvider // nil = not configured
		groq         *fakeProvider
		wantProvider string
		wantSummary  string
		wantErr      string
		wantGroqCall int
	}{
		{
			name:         "gemini succeeds and groq is never called",
			gemini:       &fakeProvider{name: "gemini", result: okGemini},
			groq:         &fakeProvider{name: "groq", result: okGroq},
			wantProvider: "gemini",
			wantSummary:  "from gemini",
			wantGroqCall: 0,
		},
		{
			name:         "gemini fails and groq is used as fallback",
			gemini:       &fakeProvider{name: "gemini", err: errors.New("gemini API error 503")},
			groq:         &fakeProvider{name: "groq", result: okGroq},
			wantProvider: "groq",
			wantSummary:  "from groq",
			wantGroqCall: 1,
		},
		{
			name:         "only groq configured",
			groq:         &fakeProvider{name: "groq", result: okGroq},
			wantProvider: "groq",
			wantSummary:  "from groq",
			wantGroqCall: 1,
		},
		{
			name:         "both fail returns an error wrapping the last failure",
			gemini:       &fakeProvider{name: "gemini", err: errors.New("gemini down")},
			groq:         &fakeProvider{name: "groq", err: errors.New("groq rate limited")},
			wantErr:      "groq rate limited",
			wantGroqCall: 1,
		},
		{
			name:    "only gemini configured and it fails",
			gemini:  &fakeProvider{name: "gemini", err: errors.New("gemini down")},
			wantErr: "gemini down",
		},
		{
			name:    "no provider configured",
			wantErr: "no AI provider configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var providers []LLMProvider
			if tt.gemini != nil {
				providers = append(providers, tt.gemini)
			}
			if tt.groq != nil {
				providers = append(providers, tt.groq)
			}
			svc := NewAIServiceWithProviders(providers...)

			result, provider, err := svc.Analyze("prompt")

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
				if result != nil || provider != "" {
					t.Errorf("expected no result on error, got %v / %q", result, provider)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if provider != tt.wantProvider {
					t.Errorf("provider = %q, want %q", provider, tt.wantProvider)
				}
				if result.Summary != tt.wantSummary {
					t.Errorf("summary = %q, want %q", result.Summary, tt.wantSummary)
				}
			}
			if tt.groq != nil && tt.groq.calls != tt.wantGroqCall {
				t.Errorf("groq called %d times, want %d", tt.groq.calls, tt.wantGroqCall)
			}
		})
	}
}

const validAIJSON = `{"approach_scores":{"psicanalise":80},"field_scores":{"clinica":70},"approach_details":{},"field_details":{},"summary":"ok"}`

func TestParseAIResponse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantSummary string
		wantErr     bool
	}{
		{name: "plain JSON", input: validAIJSON, wantSummary: "ok"},
		{name: "markdown json fence", input: "```json\n" + validAIJSON + "\n```", wantSummary: "ok"},
		{name: "bare fence and whitespace", input: "  ```\n" + validAIJSON + "```  ", wantSummary: "ok"},
		{
			name:        "repairs object closed with bracket",
			input:       `{"approach_details":{"psicanalise":{"description":"x"]},"summary":"fixed"}`,
			wantSummary: "fixed",
		},
		{name: "garbage", input: "I cannot help with that", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAIResponse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Summary != tt.wantSummary {
				t.Errorf("summary = %q, want %q", got.Summary, tt.wantSummary)
			}
		})
	}
}

func TestGeminiProvider_Generate(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{
			name:   "success",
			status: http.StatusOK,
			body:   geminiBody(t, validAIJSON),
		},
		{name: "http error", status: http.StatusServiceUnavailable, body: "overloaded", wantErr: "gemini API error 503"},
		{name: "no candidates", status: http.StatusOK, body: `{"candidates":[]}`, wantErr: "empty gemini response"},
		{name: "invalid payload", status: http.StatusOK, body: geminiBody(t, "not json"), wantErr: "failed to parse AI JSON"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotKey, gotPrompt string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotKey = r.URL.Query().Get("key")
				var req struct {
					Contents []struct {
						Parts []struct{ Text string } `json:"parts"`
					} `json:"contents"`
				}
				_ = json.NewDecoder(r.Body).Decode(&req)
				if len(req.Contents) > 0 && len(req.Contents[0].Parts) > 0 {
					gotPrompt = req.Contents[0].Parts[0].Text
				}
				w.WriteHeader(tt.status)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer srv.Close()

			p := &geminiProvider{apiKey: "test-key", baseURL: srv.URL, httpClient: srv.Client()}
			got, err := p.Generate("hello")

			if gotKey != "test-key" {
				t.Errorf("api key sent = %q, want test-key", gotKey)
			}
			if gotPrompt != "hello" {
				t.Errorf("prompt sent = %q, want hello", gotPrompt)
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Summary != "ok" {
				t.Errorf("summary = %q, want ok", got.Summary)
			}
		})
	}
}

func TestGroqProvider_Generate(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{name: "success", status: http.StatusOK, body: groqBody(t, validAIJSON)},
		{name: "http error", status: http.StatusTooManyRequests, body: "slow down", wantErr: "groq API error 429"},
		{name: "no choices", status: http.StatusOK, body: `{"choices":[]}`, wantErr: "empty groq response"},
		{name: "undecodable body", status: http.StatusOK, body: `<html>`, wantErr: "invalid character"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotAuth, gotModel string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotAuth = r.Header.Get("Authorization")
				var req struct{ Model string }
				_ = json.NewDecoder(r.Body).Decode(&req)
				gotModel = req.Model
				w.WriteHeader(tt.status)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer srv.Close()

			p := &groqProvider{apiKey: "groq-key", url: srv.URL, httpClient: srv.Client()}
			got, err := p.Generate("hello")

			if gotAuth != "Bearer groq-key" {
				t.Errorf("Authorization = %q, want Bearer groq-key", gotAuth)
			}
			if gotModel == "" {
				t.Error("request did not specify a model")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Summary != "ok" {
				t.Errorf("summary = %q, want ok", got.Summary)
			}
		})
	}
}

func TestNewAIService_ProvidersFromEnv(t *testing.T) {
	tests := []struct {
		name   string
		gemini string
		groq   string
		want   []string
	}{
		{"both keys", "g", "q", []string{"gemini", "groq"}},
		{"gemini only", "g", "", []string{"gemini"}},
		{"groq only", "", "q", []string{"groq"}},
		{"none", "", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GEMINI_API_KEY", tt.gemini)
			t.Setenv("GROQ_API_KEY", tt.groq)
			svc := NewAIService()
			var got []string
			for _, p := range svc.providers {
				got = append(got, p.Name())
			}
			if !equalStrings(got, tt.want) {
				t.Errorf("providers = %v, want %v", got, tt.want)
			}
		})
	}
}

func geminiBody(t *testing.T, text string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"candidates": []any{map[string]any{
			"content": map[string]any{"parts": []any{map[string]string{"text": text}}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func groqBody(t *testing.T, content string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"choices": []any{map[string]any{"message": map[string]string{"content": content}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestGeminiURL(t *testing.T) {
	got := geminiURL("gemini-flash-lite-latest")
	want := "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-lite-latest:generateContent"
	if got != want {
		t.Fatalf("geminiURL() = %q, want %q", got, want)
	}
}

func TestNewAIServiceGeminiModel(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  string
	}{
		{"default alias", "", geminiURL(defaultGeminiModel)},
		{"pinned by env", "gemini-3.6-flash", geminiURL("gemini-3.6-flash")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GEMINI_API_KEY", "test-key")
			t.Setenv("GROQ_API_KEY", "")
			t.Setenv("GEMINI_MODEL", tt.model)

			svc := NewAIService()
			if len(svc.providers) != 1 {
				t.Fatalf("got %d providers, want 1", len(svc.providers))
			}
			p, ok := svc.providers[0].(*geminiProvider)
			if !ok {
				t.Fatalf("provider is %T, want *geminiProvider", svc.providers[0])
			}
			if p.baseURL != tt.want {
				t.Fatalf("baseURL = %q, want %q", p.baseURL, tt.want)
			}
		})
	}
}
