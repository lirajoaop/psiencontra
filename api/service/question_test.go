package service

import (
	"sort"
	"testing"
)

// The detailed questionnaire is the input to the deterministic scoring
// engine, so its structure is part of the contract: 8 approaches × 5 items
// (one reverse-coded each) and 9 fields × 4 items.
func TestDetailedQuestionBankStructure(t *testing.T) {
	qs := NewQuestionService().GetAllDetailed()
	if len(qs) != 76 {
		t.Fatalf("detailed bank has %d items, want 76", len(qs))
	}

	ids := map[int]bool{}
	itemsPer := map[string]int{}
	reversedPer := map[string]int{}
	for _, q := range qs {
		if ids[q.ID] {
			t.Errorf("duplicate question ID %d", q.ID)
		}
		ids[q.ID] = true

		if q.Type != "likert" {
			t.Errorf("question %d: type %q, want likert", q.ID, q.Type)
		}
		if len(q.Options) != 5 {
			t.Fatalf("question %d: %d options, want 5", q.ID, len(q.Options))
		}
		mapping := q.Options[0].Mapping
		for i, o := range q.Options {
			if o.Mapping != mapping {
				t.Errorf("question %d: option %d maps to %q, want %q", q.ID, i, o.Mapping, mapping)
			}
		}
		key := q.Block + "/" + mapping
		itemsPer[key]++
		if q.Reversed {
			reversedPer[key]++
		}
	}

	for _, a := range approachOrder {
		key := "approaches/" + a
		if itemsPer[key] != 5 {
			t.Errorf("approach %s has %d items, want 5", a, itemsPer[key])
		}
		if reversedPer[key] != 1 {
			t.Errorf("approach %s has %d reverse-coded items, want 1", a, reversedPer[key])
		}
	}
	for _, f := range fieldOrder {
		key := "fields/" + f
		if itemsPer[key] != 4 {
			t.Errorf("field %s has %d items, want 4", f, itemsPer[key])
		}
	}
	if len(itemsPer) != len(approachOrder)+len(fieldOrder) {
		t.Errorf("unexpected mappings in bank: %v", itemsPer)
	}
}

// Every mapping used by either questionnaire must have a human label in the
// PDF, otherwise the report would print the raw key.
func TestEveryMappingHasLabel(t *testing.T) {
	svc := NewQuestionService()
	for _, qs := range [][]Question{svc.GetAll(), svc.GetAllDetailed()} {
		for _, q := range qs {
			for _, o := range q.Options {
				if o.Mapping == "" {
					continue
				}
				var labels map[string]string
				switch q.Block {
				case "approaches":
					labels = approachLabels
				case "fields":
					labels = fieldLabels
				default:
					continue
				}
				if _, ok := labels[o.Mapping]; !ok {
					t.Errorf("question %d (%s): mapping %q has no label", q.ID, q.Block, o.Mapping)
				}
			}
		}
	}
}

func TestGetByType(t *testing.T) {
	svc := NewQuestionService()
	tests := []struct {
		qType string
		want  int
	}{
		{"detailed", len(detailedQuestions)},
		{"simple", len(questions)},
		{"", len(questions)},
		{"unknown", len(questions)},
	}
	for _, tt := range tests {
		t.Run(tt.qType, func(t *testing.T) {
			if got := len(svc.GetByType(tt.qType)); got != tt.want {
				t.Errorf("GetByType(%q) returned %d questions, want %d", tt.qType, got, tt.want)
			}
		})
	}
}

func TestGetShuffledByType_PreservesContentAndBlockOrder(t *testing.T) {
	svc := NewQuestionService()
	for _, qType := range []string{"simple", "detailed"} {
		t.Run(qType, func(t *testing.T) {
			original := svc.GetByType(qType)
			originalIDs := idsOf(original)
			firstOptionLabel := map[int]string{}
			for _, q := range original {
				if len(q.Options) > 0 {
					firstOptionLabel[q.ID] = q.Options[0].Label
				}
			}

			shuffled := svc.GetShuffledByType(qType)

			if got := sortedInts(idsOf(shuffled)); !equalInts(got, sortedInts(originalIDs)) {
				t.Fatalf("shuffled set of IDs differs from original")
			}
			if got, want := blockSequence(shuffled), blockSequence(original); !equalStrings(got, want) {
				t.Errorf("block order changed: got %v, want %v", got, want)
			}

			byID := map[int]Question{}
			for _, q := range original {
				byID[q.ID] = q
			}
			for _, q := range shuffled {
				orig := byID[q.ID]
				if len(q.Options) != len(orig.Options) {
					t.Fatalf("question %d lost options", q.ID)
				}
				if q.Type == "likert" {
					for i := range q.Options {
						if q.Options[i] != orig.Options[i] {
							t.Errorf("question %d: Likert scale must not be shuffled", q.ID)
							break
						}
					}
				}
			}

			// The source slice must not be mutated by shuffling.
			if !equalInts(idsOf(svc.GetByType(qType)), originalIDs) {
				t.Error("shuffling mutated the original question order")
			}
			for _, q := range svc.GetByType(qType) {
				if len(q.Options) > 0 && q.Options[0].Label != firstOptionLabel[q.ID] {
					t.Errorf("shuffling mutated options of question %d", q.ID)
				}
			}
		})
	}
}

func idsOf(qs []Question) []int {
	out := make([]int, len(qs))
	for i, q := range qs {
		out[i] = q.ID
	}
	return out
}

// blockSequence returns the distinct blocks in order of first appearance.
func blockSequence(qs []Question) []string {
	var out []string
	for _, q := range qs {
		if len(out) == 0 || out[len(out)-1] != q.Block {
			out = append(out, q.Block)
		}
	}
	return out
}

func sortedInts(in []int) []int {
	out := append([]int(nil), in...)
	sort.Ints(out)
	return out
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
