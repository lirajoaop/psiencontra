package service

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/joaop/psiencontra/api/schemas"
)

func TestParseScoresWithDetails(t *testing.T) {
	scores := json.RawMessage(`{"tcc":40,"psicanalise":90,"unknown_key":10}`)
	details := json.RawMessage(`{"psicanalise":{"score":90,"description":"forte afinidade"}}`)

	got := parseScoresWithDetails(scores, details, approachLabels)

	want := []scoreEntry{
		{Key: "Psicanálise", Score: 90, Desc: "forte afinidade"},
		{Key: "Terapia Cognitivo-Comportamental", Score: 40},
		{Key: "unknown_key", Score: 10}, // falls back to the raw key
	}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestParseScoresWithDetails_InvalidJSON(t *testing.T) {
	if got := parseScoresWithDetails(json.RawMessage(`nope`), nil, approachLabels); len(got) != 0 {
		t.Errorf("expected no entries for invalid JSON, got %+v", got)
	}
}

func TestPDFGenerate(t *testing.T) {
	result := &schemas.Result{
		ApproachScores:  json.RawMessage(`{"psicanalise":95,"tcc":5}`),
		FieldScores:     json.RawMessage(`{"clinica":95,"psicometria":5}`),
		ApproachDetails: json.RawMessage(`{"psicanalise":{"score":95,"description":"Afinidade com a escuta do inconsciente."}}`),
		FieldDetails:    json.RawMessage(`{}`),
		Explanation:     "Perfil voltado à clínica.",
	}
	for _, qType := range []string{"simple", "detailed"} {
		t.Run(qType, func(t *testing.T) {
			pdf, err := NewPDFService().Generate(result, qType)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
				t.Errorf("output is not a PDF (starts with %q)", pdf[:min(len(pdf), 8)])
			}
		})
	}
}
