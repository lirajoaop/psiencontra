package service

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"

	"github.com/joaop/psiencontra/api/schemas"
)

// likertQ builds a synthetic Likert question for scoring tests.
func likertQ(id int, block, mapping string, reversed bool) Question {
	return Question{ID: id, Type: "likert", Block: block, Options: likertOptions(mapping), Reversed: reversed}
}

func likertR(questionID int, value string) schemas.Response {
	return schemas.Response{QuestionID: questionID, AnswerType: "likert", AnswerValue: value}
}

func TestComputeIpsativeScores(t *testing.T) {
	// Two approach mappings with two items each, plus one field mapping pair.
	qs := []Question{
		likertQ(1, "approaches", "psicanalise", false),
		likertQ(2, "approaches", "psicanalise", false),
		likertQ(3, "approaches", "tcc", false),
		likertQ(4, "approaches", "tcc", true),
		likertQ(5, "approaches", "gestalt", false),
		likertQ(10, "fields", "clinica", false),
		likertQ(11, "fields", "escolar", false),
	}

	tests := []struct {
		name          string
		responses     []schemas.Response
		wantApproach  map[string]float64
		wantField     map[string]float64
		wantApproachN int
	}{
		{
			name:         "no responses yields empty maps",
			responses:    nil,
			wantApproach: map[string]float64{},
			wantField:    map[string]float64{},
		},
		{
			name: "flat profile gives every mapping 50",
			responses: []schemas.Response{
				likertR(1, "3"), likertR(2, "3"), likertR(3, "3"), likertR(4, "3"), likertR(5, "3"),
				likertR(10, "4"), likertR(11, "4"),
			},
			wantApproach: map[string]float64{"psicanalise": 50, "tcc": 50, "gestalt": 50},
			wantField:    map[string]float64{"clinica": 50, "escolar": 50},
		},
		{
			name: "highest mean maps to 95, lowest to 5, middle is interpolated",
			responses: []schemas.Response{
				likertR(1, "5"), likertR(2, "5"), // psicanalise mean 5
				likertR(3, "1"), likertR(4, "5"), // tcc: 1 and reversed 5→1, mean 1
				likertR(5, "3"), // gestalt mean 3
			},
			wantApproach: map[string]float64{"psicanalise": 95, "tcc": 5, "gestalt": 50},
			wantField:    map[string]float64{},
		},
		{
			name: "reverse-coded item lowers its mapping when endorsed",
			responses: []schemas.Response{
				likertR(1, "3"), likertR(2, "3"), // psicanalise mean 3
				likertR(3, "3"), likertR(4, "5"), // tcc: 3 and 6-5=1 → mean 2
				likertR(5, "4"), // gestalt mean 4
			},
			wantApproach: map[string]float64{"psicanalise": 50, "tcc": 5, "gestalt": 95},
			wantField:    map[string]float64{},
		},
		{
			name: "reverse-coded item raises its mapping when rejected",
			responses: []schemas.Response{
				likertR(1, "3"), likertR(2, "3"),
				likertR(3, "5"), likertR(4, "1"), // tcc: 5 and 6-1=5 → mean 5
				likertR(5, "1"),
			},
			wantApproach: map[string]float64{"psicanalise": 50, "tcc": 95, "gestalt": 5},
			wantField:    map[string]float64{},
		},
		{
			name: "blocks are normalized independently",
			responses: []schemas.Response{
				likertR(1, "5"), likertR(3, "1"),
				likertR(10, "2"), likertR(11, "4"),
			},
			wantApproach: map[string]float64{"psicanalise": 95, "tcc": 5},
			wantField:    map[string]float64{"clinica": 5, "escolar": 95},
		},
		{
			name: "invalid, unknown and non-likert answers are ignored",
			responses: []schemas.Response{
				likertR(1, "5"),
				likertR(3, "1"),
				likertR(5, "not-a-number"), // unparseable
				likertR(999, "5"),          // unknown question
				{QuestionID: 2, AnswerType: "multiple_choice", AnswerValue: "1"},
			},
			wantApproach: map[string]float64{"psicanalise": 95, "tcc": 5},
			wantField:    map[string]float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotF := ComputeIpsativeScores(tt.responses, qs)
			assertScores(t, "approach", gotA, tt.wantApproach)
			assertScores(t, "field", gotF, tt.wantField)
		})
	}
}

func TestComputeIpsativeScores_IgnoresQuestionsWithoutMapping(t *testing.T) {
	qs := []Question{
		likertQ(1, "approaches", "psicanalise", false),
		{ID: 2, Type: "likert", Block: "approaches"},                               // no options
		{ID: 3, Type: "likert", Block: "approaches", Options: likertOptions("")},   // empty mapping
		{ID: 4, Type: "likert", Block: "vignettes", Options: likertOptions("tcc")}, // unscored block
	}
	resp := []schemas.Response{likertR(1, "4"), likertR(2, "5"), likertR(3, "1"), likertR(4, "1")}

	gotA, gotF := ComputeIpsativeScores(resp, qs)
	assertScores(t, "approach", gotA, map[string]float64{"psicanalise": 50})
	assertScores(t, "field", gotF, map[string]float64{})
}

// Ipsative scoring must cancel acquiescence bias: two respondents with the
// same relative preferences but different overall levels get the same
// profile.
func TestComputeIpsativeScores_AcquiescenceInvariant(t *testing.T) {
	qs := []Question{
		likertQ(1, "approaches", "psicanalise", false),
		likertQ(2, "approaches", "tcc", false),
		likertQ(3, "approaches", "gestalt", false),
		likertQ(4, "approaches", "humanismo", false),
	}
	low := []schemas.Response{likertR(1, "1"), likertR(2, "2"), likertR(3, "3"), likertR(4, "2")}
	high := []schemas.Response{likertR(1, "3"), likertR(2, "4"), likertR(3, "5"), likertR(4, "4")}

	lowA, _ := ComputeIpsativeScores(low, qs)
	highA, _ := ComputeIpsativeScores(high, qs)
	assertScores(t, "approach", highA, lowA)
}

func TestComputeIpsativeScores_Deterministic(t *testing.T) {
	qs := detailedQuestions
	resp := answerDetailed(t, map[string]int{"gestalt": 5, "tcc": 2, "clinica": 5, "juridica": 1}, 3)

	firstA, firstF := ComputeIpsativeScores(resp, qs)
	for i := 0; i < 20; i++ {
		a, f := ComputeIpsativeScores(resp, qs)
		assertScores(t, "approach", a, firstA)
		assertScores(t, "field", f, firstF)
	}
}

// End-to-end over the real 76-item bank: a respondent who strongly endorses
// one approach and one field (including answering their reverse-coded items
// consistently) must see them ranked first, and every score must stay in
// the documented [5, 95] range.
func TestComputeIpsativeScores_RealQuestionnaireRanking(t *testing.T) {
	tests := []struct {
		name          string
		prefs         map[string]int
		wantTopA      string
		wantBottomA   string
		wantTopF      string
		wantBottomF   string
		wantAllFifty  bool
		defaultRating int
	}{
		{
			name:          "psicanalise and clinica preferred, tcc and psicometria rejected",
			prefs:         map[string]int{"psicanalise": 5, "tcc": 1, "clinica": 5, "psicometria": 1},
			defaultRating: 3,
			wantTopA:      "psicanalise", wantBottomA: "tcc",
			wantTopF: "clinica", wantBottomF: "psicometria",
		},
		{
			name:          "humanismo and esporte preferred over a high baseline",
			prefs:         map[string]int{"humanismo": 5, "socio_historica": 2, "esporte": 5, "organizacional": 2},
			defaultRating: 4,
			wantTopA:      "humanismo", wantBottomA: "socio_historica",
			wantTopF: "esporte", wantBottomF: "organizacional",
		},
		{
			name:          "neutral answers everywhere give a flat profile",
			prefs:         map[string]int{},
			defaultRating: 3,
			wantAllFifty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := answerDetailed(t, tt.prefs, tt.defaultRating)
			a, f := ComputeIpsativeScores(resp, detailedQuestions)

			if len(a) != len(approachOrder) {
				t.Fatalf("approach scores: got %d keys, want %d", len(a), len(approachOrder))
			}
			if len(f) != len(fieldOrder) {
				t.Fatalf("field scores: got %d keys, want %d", len(f), len(fieldOrder))
			}
			for _, scores := range []map[string]float64{a, f} {
				for k, v := range scores {
					if v < 5 || v > 95 {
						t.Errorf("score %s = %v out of [5, 95]", k, v)
					}
					if v != float64(int(v)) {
						t.Errorf("score %s = %v is not an integer", k, v)
					}
				}
			}

			if tt.wantAllFifty {
				for k, v := range a {
					if v != 50 {
						t.Errorf("approach %s = %v, want 50", k, v)
					}
				}
				for k, v := range f {
					if v != 50 {
						t.Errorf("field %s = %v, want 50", k, v)
					}
				}
				return
			}

			assertRank(t, a, tt.wantTopA, tt.wantBottomA)
			assertRank(t, f, tt.wantTopF, tt.wantBottomF)
		})
	}
}

func TestOverrideDetailScores(t *testing.T) {
	canonical := map[string]float64{"psicanalise": 95, "tcc": 5}

	tests := []struct {
		name      string
		details   string
		canonical map[string]float64
		want      map[string]detail
		wantRaw   string // when set, output must equal this verbatim
	}{
		{
			name:      "replaces AI scores and keeps descriptions",
			details:   `{"psicanalise":{"score":40,"description":"forte"},"tcc":{"score":80,"description":"fraca"}}`,
			canonical: canonical,
			want:      map[string]detail{"psicanalise": {95, "forte"}, "tcc": {5, "fraca"}},
		},
		{
			name:      "adds missing keys with empty description and drops extra keys",
			details:   `{"psicanalise":{"score":40,"description":"forte"},"gestalt":{"score":70,"description":"x"}}`,
			canonical: canonical,
			want:      map[string]detail{"psicanalise": {95, "forte"}, "tcc": {5, ""}},
		},
		{
			name:      "malformed AI JSON is rebuilt from canonical scores",
			details:   `{"psicanalise": [oops`,
			canonical: canonical,
			want:      map[string]detail{"psicanalise": {95, ""}, "tcc": {5, ""}},
		},
		{
			name:      "empty details are returned untouched",
			details:   ``,
			canonical: canonical,
			wantRaw:   ``,
		},
		{
			name:      "empty canonical scores leave details untouched",
			details:   `{"psicanalise":{"score":40,"description":"forte"}}`,
			canonical: map[string]float64{},
			wantRaw:   `{"psicanalise":{"score":40,"description":"forte"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := overrideDetailScores(json.RawMessage(tt.details), tt.canonical)
			if tt.want == nil {
				if string(got) != tt.wantRaw {
					t.Fatalf("got %s, want %s", got, tt.wantRaw)
				}
				return
			}
			var parsed map[string]detail
			if err := json.Unmarshal(got, &parsed); err != nil {
				t.Fatalf("output is not valid JSON: %v (%s)", err, got)
			}
			if len(parsed) != len(tt.want) {
				t.Fatalf("got %d keys, want %d: %s", len(parsed), len(tt.want), got)
			}
			for k, w := range tt.want {
				if parsed[k] != w {
					t.Errorf("%s: got %+v, want %+v", k, parsed[k], w)
				}
			}
		})
	}
}

type detail struct {
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

// answerDetailed answers the full detailed questionnaire. Mappings listed
// in prefs get that rating (the rating is mirrored on reverse-coded items so
// the respondent is consistent); everything else gets defaultRating.
func answerDetailed(t *testing.T, prefs map[string]int, defaultRating int) []schemas.Response {
	t.Helper()
	var out []schemas.Response
	for _, q := range detailedQuestions {
		rating := defaultRating
		if p, ok := prefs[q.Options[0].Mapping]; ok {
			rating = p
		}
		if q.Reversed {
			rating = 6 - rating
		}
		out = append(out, schemas.Response{
			QuestionID:  q.ID,
			AnswerType:  q.Type,
			AnswerValue: strconv.Itoa(rating),
		})
	}
	return out
}

func assertScores(t *testing.T, label string, got, want map[string]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s scores: got %v, want %v", label, got, want)
	}
	for k, w := range want {
		if g, ok := got[k]; !ok || g != w {
			t.Errorf("%s score %q: got %v, want %v", label, k, got[k], w)
		}
	}
}

func assertRank(t *testing.T, scores map[string]float64, top, bottom string) {
	t.Helper()
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return scores[keys[i]] > scores[keys[j]] })
	if keys[0] != top || scores[top] != 95 {
		t.Errorf("top = %s (%v), want %s with 95; scores %v", keys[0], scores[keys[0]], top, scores)
	}
	if keys[len(keys)-1] != bottom || scores[bottom] != 5 {
		t.Errorf("bottom = %s (%v), want %s with 5; scores %v", keys[len(keys)-1], scores[keys[len(keys)-1]], bottom, scores)
	}
}
