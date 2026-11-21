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

const (
	geminiModelsURL = "https://generativelanguage.googleapis.com/v1beta/models/"
	// defaultGeminiModel follows Google's current Flash release. Override with GEMINI_MODEL.
	defaultGeminiModel = "gemini-flash-latest"
	defaultGroqURL     = "https://api.groq.com/openai/v1/chat/completions"
)

// geminiURL returns the generateContent endpoint for a Gemini model.
func geminiURL(model string) string {
	return geminiModelsURL + model + ":generateContent"
}

// LLMProvider is a single LLM backend able to turn a prompt into an AIResult.
// Gemini and Groq implement it; tests substitute fakes.
type LLMProvider interface {
	Name() string
	Generate(prompt string) (*AIResult, error)
}

// AIService tries each configured provider in order and returns the first
// successful result (Gemini first, Groq as fallback).
type AIService struct {
	providers []LLMProvider
}

type AIResult struct {
	ApproachScores  json.RawMessage `json:"approach_scores"`
	FieldScores     json.RawMessage `json:"field_scores"`
	ApproachDetails json.RawMessage `json:"approach_details"`
	FieldDetails    json.RawMessage `json:"field_details"`
	Summary         string          `json:"summary"`
}

func NewAIService() *AIService {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	var providers []LLMProvider
	if key := config.GetEnv("GEMINI_API_KEY", ""); key != "" {
		model := config.GetEnv("GEMINI_MODEL", defaultGeminiModel)
		providers = append(providers, &geminiProvider{apiKey: key, baseURL: geminiURL(model), httpClient: httpClient})
	}
	if key := config.GetEnv("GROQ_API_KEY", ""); key != "" {
		providers = append(providers, &groqProvider{apiKey: key, url: defaultGroqURL, httpClient: httpClient})
	}
	return NewAIServiceWithProviders(providers...)
}

// NewAIServiceWithProviders builds an AIService from an explicit, ordered
// list of providers. Mainly useful for tests.
func NewAIServiceWithProviders(providers ...LLMProvider) *AIService {
	return &AIService{providers: providers}
}

// Analyze sends the prompt to each provider in order, falling back to the
// next one on failure. It returns the result and the name of the provider
// that produced it.
func (s *AIService) Analyze(prompt string) (*AIResult, string, error) {
	if len(s.providers) == 0 {
		return nil, "", fmt.Errorf("no AI provider configured")
	}

	var lastErr error
	for i, p := range s.providers {
		result, err := p.Generate(prompt)
		if err == nil {
			return result, p.Name(), nil
		}
		lastErr = err
		if i < len(s.providers)-1 && config.Log != nil {
			config.Log.Error.Printf("%s failed, falling back to %s: %v", p.Name(), s.providers[i+1].Name(), err)
		}
	}
	return nil, "", fmt.Errorf("all AI providers failed: %w", lastErr)
}

type geminiProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func (p *geminiProvider) Name() string { return "gemini" }

type groqProvider struct {
	apiKey     string
	url        string
	httpClient *http.Client
}

func (p *groqProvider) Name() string { return "groq" }

func (p *geminiProvider) Generate(prompt string) (*AIResult, error) {
	url := p.baseURL + "?key=" + p.apiKey

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

	resp, err := p.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
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

func (p *groqProvider) Generate(prompt string) (*AIResult, error) {

	body := map[string]any{
		"model": "openai/gpt-oss-20b",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a JSON generator. You MUST return only valid JSON. Every opening { must have a matching }. Every opening [ must have a matching ]. Never use ] to close an object that was opened with {."},
			{"role": "user", "content": prompt},
		},
		"temperature":     0.7,
		"max_tokens":      4096,
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", p.url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(req)
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
