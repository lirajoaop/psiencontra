package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/joaop/psiencontra/api/config"
)

type AIService struct {
	geminiKey  string
	groqKey    string
	httpClient *http.Client
}

type AIResult struct {
	ApproachScores  json.RawMessage `json:"approach_scores"`
	FieldScores     json.RawMessage `json:"field_scores"`
	ApproachDetails json.RawMessage `json:"approach_details"`
	FieldDetails    json.RawMessage `json:"field_details"`
	Summary         string          `json:"summary"`
}

func NewAIService() *AIService {
	return &AIService{
		geminiKey: config.GetEnv("GEMINI_API_KEY", ""),
		groqKey:   config.GetEnv("GROQ_API_KEY", ""),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *AIService) Analyze(prompt string) (*AIResult, string, error) {
	if s.geminiKey != "" {
		result, err := s.callGemini(prompt)
		if err == nil {
			return result, "gemini", nil
		}
		config.Log.Error.Printf("Gemini failed, falling back to Groq: %v", err)
	}

	if s.groqKey != "" {
		result, err := s.callGroq(prompt)
		if err == nil {
			return result, "groq", nil
		}
		return nil, "", fmt.Errorf("both AI providers failed: %v", err)
	}

	return nil, "", fmt.Errorf("no AI provider configured")
}

func (s *AIService) callGemini(prompt string) (*AIResult, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", s.geminiKey)

	body := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":      0.7,
			"maxOutputTokens":  8192,
			"responseMimeType": "application/json",
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini API error %d: %s", resp.StatusCode, string(respBody))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty gemini response")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text
	return parseAIResponse(text)
}

func (s *AIService) callGroq(prompt string) (*AIResult, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	body := map[string]any{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a JSON generator. You MUST return only valid JSON. Every opening { must have a matching }. Every opening [ must have a matching ]. Never use ] to close an object that was opened with {."},
			{"role": "user", "content": prompt},
		},
		"temperature":     0.7,
		"max_tokens":      8192,
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.groqKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq API error %d: %s", resp.StatusCode, string(respBody))
	}

	var groqResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, err
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("empty groq response")
	}

	return parseAIResponse(groqResp.Choices[0].Message.Content)
}

func parseAIResponse(text string) (*AIResult, error) {
	text = strings.TrimSpace(text)
	// Remove markdown code fences if present
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result AIResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		// Try to repair common JSON issues from LLMs
		repaired := repairJSON(text)
		if err2 := json.Unmarshal([]byte(repaired), &result); err2 != nil {
			return nil, fmt.Errorf("failed to parse AI JSON: %w\nraw: %s", err, text[:min(len(text), 500)])
		}
	}

	return &result, nil
}

func repairJSON(text string) string {
	// Fix common LLM mistake: "] instead of "} to close objects
	text = strings.ReplaceAll(text, "\"]", "\"}")
	// Fix "] with whitespace variations
	text = strings.ReplaceAll(text, "\" ]", "\" }")
	return text
}
