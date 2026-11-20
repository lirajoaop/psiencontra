package service

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/internal/memrepo"
	"github.com/joaop/psiencontra/api/repository"
)

// Compile-time checks: both the GORM repositories and the in-memory fakes
// satisfy the interfaces SessionService depends on.
var (
	_ SessionRepository  = (*repository.SessionRepo)(nil)
	_ ResponseRepository = (*repository.ResponseRepo)(nil)
	_ ResultRepository   = (*repository.ResultRepo)(nil)
	_ UserRepository     = (*repository.UserRepo)(nil)
	_ SessionRepository  = (*memrepo.Sessions)(nil)
	_ ResponseRepository = (*memrepo.Responses)(nil)
	_ ResultRepository   = (*memrepo.Results)(nil)
	_ UserRepository     = (*memrepo.Users)(nil)
	_ Analyzer           = (*AIService)(nil)
)

// stubAnalyzer records the prompt it receives and returns a canned result.
type stubAnalyzer struct {
	result     *AIResult
	provider   string
	err        error
	lastPrompt string
}

func (a *stubAnalyzer) Analyze(prompt string) (*AIResult, string, error) {
	a.lastPrompt = prompt
	return a.result, a.provider, a.err
}

func newSessionService(ai Analyzer) (*SessionService, *memrepo.Store) {
	store := memrepo.New()
	svc := NewSessionService(store.Sessions(), store.Responses(), store.Results(), NewQuestionService(), ai)
	return svc, store
}

func TestCreateSession(t *testing.T) {
	svc, _ := newSessionService(&stubAnalyzer{})
	uid := uuid.New()

	tests := []struct {
		name     string
		userID   *uuid.UUID
		qType    string
		wantType string
	}{
		{"anonymous defaults to simple", nil, "", "simple"},
		{"owned detailed", &uid, "detailed", "detailed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sess, err := svc.CreateSession(tt.userID, tt.qType)
			if err != nil {
				t.Fatal(err)
			}
			if sess.ID == uuid.Nil {
				t.Error("session ID was not generated")
			}
			if sess.QuestionnaireType != tt.wantType {
				t.Errorf("type = %q, want %q", sess.QuestionnaireType, tt.wantType)
			}
			stored, err := svc.GetSession(sess.ID)
			if err != nil {
				t.Fatalf("session not persisted: %v", err)
			}
			if (stored.UserID == nil) != (tt.userID == nil) {
				t.Errorf("owner mismatch: got %v, want %v", stored.UserID, tt.userID)
			}
		})
	}
}

func TestCreateSession_RepoError(t *testing.T) {
	svc, store := newSessionService(&stubAnalyzer{})
	store.CreateSessionErr = errors.New("db down")
	if _, err := svc.CreateSession(nil, ""); err == nil {
		t.Fatal("expected error")
	}
}

// For the detailed questionnaire the numeric scores come from the ipsative
// engine, never from the LLM, even if the LLM returns different numbers.
func TestSubmitResponses_DetailedUsesDeterministicScores(t *testing.T) {
	ai := &stubAnalyzer{
		provider: "groq",
		result: &AIResult{
			ApproachScores:  json.RawMessage(`{"psicanalise":1,"tcc":99}`),
			FieldScores:     json.RawMessage(`{"clinica":1}`),
			ApproachDetails: json.RawMessage(`{"psicanalise":{"score":1,"description":"desc psicanalise"}}`),
			FieldDetails:    json.RawMessage(`{"clinica":{"score":1,"description":"desc clinica"}}`),
			Summary:         "summary",
		},
	}
	svc, _ := newSessionService(ai)
	sess, _ := svc.CreateSession(nil, "detailed")

	var inputs []SubmitResponseInput
	for _, r := range answerDetailed(t, map[string]int{"psicanalise": 5, "tcc": 1, "clinica": 5, "psicometria": 1}, 3) {
		inputs = append(inputs, SubmitResponseInput{QuestionID: r.QuestionID, AnswerValue: r.AnswerValue})
	}

	result, err := svc.SubmitResponses(sess.ID, inputs)
	if err != nil {
		t.Fatalf("SubmitResponses: %v", err)
	}

	var approach, field map[string]float64
	mustUnmarshal(t, result.ApproachScores, &approach)
	mustUnmarshal(t, result.FieldScores, &field)
	if approach["psicanalise"] != 95 || approach["tcc"] != 5 {
		t.Errorf("approach scores not deterministic: %v", approach)
	}
	if field["clinica"] != 95 || field["psicometria"] != 5 {
		t.Errorf("field scores not deterministic: %v", field)
	}

	var details map[string]detail
	mustUnmarshal(t, result.ApproachDetails, &details)
	if got := details["psicanalise"]; got.Score != 95 || got.Description != "desc psicanalise" {
		t.Errorf("approach detail = %+v, want score 95 with AI description", got)
	}
	if len(details) != len(approachOrder) {
		t.Errorf("approach details have %d keys, want %d", len(details), len(approachOrder))
	}

	if result.AIProvider != "groq" || result.Explanation != "summary" {
		t.Errorf("provider/summary not propagated: %q / %q", result.AIProvider, result.Explanation)
	}
	if !strings.Contains(ai.lastPrompt, "- psicanalise: 95") {
		t.Error("detailed prompt does not carry the precomputed scores")
	}

	stored, err := svc.GetResult(sess.ID)
	if err != nil || stored.SessionID != sess.ID {
		t.Fatalf("result not persisted: %v", err)
	}
}

// For the simple questionnaire the LLM scores are stored as-is.
func TestSubmitResponses_SimpleUsesAIScores(t *testing.T) {
	ai := &stubAnalyzer{
		provider: "gemini",
		result: &AIResult{
			ApproachScores: json.RawMessage(`{"psicanalise":72}`),
			FieldScores:    json.RawMessage(`{"clinica":64}`),
			Summary:        "s",
		},
	}
	svc, _ := newSessionService(ai)
	sess, _ := svc.CreateSession(nil, "simple")
	q := NewQuestionService().GetAll()[0]

	result, err := svc.SubmitResponses(sess.ID, []SubmitResponseInput{{QuestionID: q.ID, AnswerValue: q.Options[0].Value}})
	if err != nil {
		t.Fatal(err)
	}
	if string(result.ApproachScores) != `{"psicanalise":72}` {
		t.Errorf("approach scores = %s, want AI output", result.ApproachScores)
	}
	if !strings.Contains(ai.lastPrompt, q.Text) {
		t.Error("simple prompt does not include the question text")
	}
}

func TestSubmitResponses_Errors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*SessionService, *memrepo.Store, *stubAnalyzer) uuid.UUID
		inputs  []SubmitResponseInput
		wantErr string
	}{
		{
			name:    "unknown session",
			setup:   func(*SessionService, *memrepo.Store, *stubAnalyzer) uuid.UUID { return uuid.New() },
			wantErr: "session not found",
		},
		{
			name: "question from another questionnaire",
			setup: func(svc *SessionService, _ *memrepo.Store, _ *stubAnalyzer) uuid.UUID {
				s, _ := svc.CreateSession(nil, "simple")
				return s.ID
			},
			inputs:  []SubmitResponseInput{{QuestionID: detailedQuestions[0].ID, AnswerValue: "5"}},
			wantErr: "invalid question_id",
		},
		{
			name: "AI failure",
			setup: func(svc *SessionService, _ *memrepo.Store, ai *stubAnalyzer) uuid.UUID {
				ai.err = errors.New("all AI providers failed")
				s, _ := svc.CreateSession(nil, "detailed")
				return s.ID
			},
			inputs:  []SubmitResponseInput{{QuestionID: detailedQuestions[0].ID, AnswerValue: "5"}},
			wantErr: "AI analysis failed",
		},
		{
			name: "result cannot be saved",
			setup: func(svc *SessionService, store *memrepo.Store, ai *stubAnalyzer) uuid.UUID {
				ai.result = &AIResult{}
				store.CreateResultErr = errors.New("db down")
				s, _ := svc.CreateSession(nil, "detailed")
				return s.ID
			},
			inputs:  []SubmitResponseInput{{QuestionID: detailedQuestions[0].ID, AnswerValue: "5"}},
			wantErr: "failed to save result",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := &stubAnalyzer{}
			svc, store := newSessionService(ai)
			id := tt.setup(svc, store, ai)
			_, err := svc.SubmitResponses(id, tt.inputs)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestClaimSession(t *testing.T) {
	svc, _ := newSessionService(&stubAnalyzer{})
	alice, bob := uuid.New(), uuid.New()

	anon, _ := svc.CreateSession(nil, "simple")
	owned, _ := svc.CreateSession(&alice, "simple")

	tests := []struct {
		name    string
		session uuid.UUID
		user    uuid.UUID
		wantErr bool
	}{
		{"anonymous session is claimed", anon.ID, bob, false},
		{"second claim is rejected", anon.ID, alice, true},
		{"owned session cannot be claimed", owned.ID, bob, true},
		{"unknown session", uuid.New(), bob, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ClaimSession(tt.session, tt.user)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	got, _ := svc.GetSession(anon.ID)
	if got.UserID == nil || *got.UserID != bob {
		t.Errorf("anonymous session owner = %v, want %v", got.UserID, bob)
	}
}

func TestGetUserHistory_OnlyCompletedOwnSessions(t *testing.T) {
	ai := &stubAnalyzer{result: &AIResult{}, provider: "gemini"}
	svc, _ := newSessionService(ai)
	alice, bob := uuid.New(), uuid.New()
	answer := []SubmitResponseInput{{QuestionID: detailedQuestions[0].ID, AnswerValue: "4"}}

	done, _ := svc.CreateSession(&alice, "detailed")
	if _, err := svc.SubmitResponses(done.ID, answer); err != nil {
		t.Fatal(err)
	}
	_, _ = svc.CreateSession(&alice, "detailed") // abandoned, no result
	other, _ := svc.CreateSession(&bob, "detailed")
	if _, err := svc.SubmitResponses(other.ID, answer); err != nil {
		t.Fatal(err)
	}

	history, err := svc.GetUserHistory(alice)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].ID != done.ID {
		t.Fatalf("history = %+v, want only session %s", history, done.ID)
	}
}

func mustUnmarshal(t *testing.T, raw json.RawMessage, v any) {
	t.Helper()
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatalf("unmarshal %s: %v", raw, err)
	}
}
