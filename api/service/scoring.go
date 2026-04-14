package service

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/joaop/psiencontra/api/schemas"
)

// ComputeIpsativeScores computes deterministic, reproducible scores for a
// Likert-only questionnaire using ipsative normalization (deviation from the
// respondent's personal mean). This cancels the acquiescence bias: two
// respondents with different response styles but the same relative preference
// produce the same profile.
//
// Algorithm (per block):
//  1. Group ratings by mapping (e.g. "psicanalise", "clinica").
//  2. Compute mean rating per mapping.
//  3. Compute the grand mean across all ratings in the block.
//  4. Compute each mapping's deviation from the grand mean.
//  5. Linearly normalize deviations into [5, 95]:
//     - The largest positive deviation becomes ~95.
//     - The largest negative deviation becomes ~5.
//     - If all deviations are equal (flat profile), every mapping gets 50.
//
// References:
//   - Holland (1997), Making Vocational Choices (RIASEC ipsative scoring)
//   - Harmon et al. (1994), Strong Interest Inventory technical manual
//   - Cattell & Brennan (1994), on ipsative vs. normative scoring
func ComputeIpsativeScores(responses []schemas.Response, questions []Question) (approachScores, fieldScores map[string]float64) {
	// Index questions by ID so we can look up their mapping and block.
	qMap := make(map[int]Question, len(questions))
	for _, q := range questions {
		qMap[q.ID] = q
	}

	approachRatings := map[string][]float64{}
	fieldRatings := map[string][]float64{}

	for _, r := range responses {
		if r.AnswerType != "likert" {
			continue
		}
		q, ok := qMap[r.QuestionID]
		if !ok || len(q.Options) == 0 {
			continue
		}
		rating, err := strconv.ParseFloat(r.AnswerValue, 64)
		if err != nil {
			continue
		}
		if q.Reversed {
			// 1-5 Likert → flip so that endorsement of an anti-construct item
			// lowers the mapped dimension.
			rating = 6 - rating
		}
		mapping := q.Options[0].Mapping
		if mapping == "" {
			continue
		}
		switch q.Block {
		case "approaches":
			approachRatings[mapping] = append(approachRatings[mapping], rating)
		case "fields":
			fieldRatings[mapping] = append(fieldRatings[mapping], rating)
		}
	}

	return normalizeIpsative(approachRatings), normalizeIpsative(fieldRatings)
}

func normalizeIpsative(ratingsByKey map[string][]float64) map[string]float64 {
	if len(ratingsByKey) == 0 {
		return map[string]float64{}
	}

	// Step 1: mean per mapping.
	means := make(map[string]float64, len(ratingsByKey))
	var totalSum float64
	var totalCount int
	for key, ratings := range ratingsByKey {
		if len(ratings) == 0 {
			continue
		}
		var sum float64
		for _, r := range ratings {
			sum += r
		}
		means[key] = sum / float64(len(ratings))
		totalSum += sum
		totalCount += len(ratings)
	}

	if totalCount == 0 {
		return map[string]float64{}
	}

	// Step 2: grand mean (across all items in the block).
	grandMean := totalSum / float64(totalCount)

	// Step 3: deviations from grand mean.
	deviations := make(map[string]float64, len(means))
	var minDev, maxDev float64
	first := true
	for key, m := range means {
		dev := m - grandMean
		deviations[key] = dev
		if first {
			minDev, maxDev = dev, dev
			first = false
			continue
		}
		if dev < minDev {
			minDev = dev
		}
		if dev > maxDev {
			maxDev = dev
		}
	}

	// Step 4: linear normalization to [5, 95]. If all deviations are equal
	// (the respondent rated every block item the same), return a flat 50 —
	// there's no signal to differentiate preferences. Scores are rounded to
	// integers: sub-integer precision is not psychometrically meaningful and
	// float residue (e.g. 95.00000000000001) leaks into the AI prompt, PDF,
	// and UI.
	scores := make(map[string]float64, len(deviations))
	span := maxDev - minDev
	if span < 1e-9 {
		for key := range deviations {
			scores[key] = 50
		}
		return scores
	}
	for key, dev := range deviations {
		scores[key] = math.Round(5 + 90*(dev-minDev)/span)
	}
	return scores
}

// overrideDetailScores rewrites the "score" field inside each entry of a
// details JSON blob ({"key": {"score": N, "description": "..."}}) using the
// canonical scores. This keeps the AI-generated descriptions but guarantees
// that the numeric scores shown in the UI match the deterministic ones.
func overrideDetailScores(detailsJSON json.RawMessage, canonicalScores map[string]float64) json.RawMessage {
	if len(detailsJSON) == 0 || len(canonicalScores) == 0 {
		return detailsJSON
	}

	var details map[string]struct {
		Score       float64 `json:"score"`
		Description string  `json:"description"`
	}
	if err := json.Unmarshal(detailsJSON, &details); err != nil {
		// If the AI returned malformed JSON, rebuild the blob from scratch so
		// at least the scores match — descriptions will be empty.
		rebuilt := make(map[string]map[string]any, len(canonicalScores))
		for key, score := range canonicalScores {
			rebuilt[key] = map[string]any{"score": score, "description": ""}
		}
		if b, err := json.Marshal(rebuilt); err == nil {
			return b
		}
		return detailsJSON
	}

	rebuilt := make(map[string]map[string]any, len(canonicalScores))
	for key, score := range canonicalScores {
		desc := ""
		if d, ok := details[key]; ok {
			desc = d.Description
		}
		rebuilt[key] = map[string]any{"score": score, "description": desc}
	}
	if b, err := json.Marshal(rebuilt); err == nil {
		return b
	}
	return detailsJSON
}
