package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/repository"
	"github.com/joaop/psiencontra/api/schemas"
)

type SessionService struct {
	sessionRepo  *repository.SessionRepo
	responseRepo *repository.ResponseRepo
	resultRepo   *repository.ResultRepo
	questionSvc  *QuestionService
	aiSvc        *AIService
}

func NewSessionService(
	sessionRepo *repository.SessionRepo,
	responseRepo *repository.ResponseRepo,
	resultRepo *repository.ResultRepo,
	questionSvc *QuestionService,
	aiSvc *AIService,
) *SessionService {
	return &SessionService{
		sessionRepo:  sessionRepo,
		responseRepo: responseRepo,
		resultRepo:   resultRepo,
		questionSvc:  questionSvc,
		aiSvc:        aiSvc,
	}
}

func (s *SessionService) CreateSession(userID *uuid.UUID) (*schemas.Session, error) {
	session := &schemas.Session{UserID: userID}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

type SubmitResponseInput struct {
	QuestionID  int    `json:"question_id"`
	AnswerValue string `json:"answer_value"`
}

func (s *SessionService) SubmitResponses(sessionID uuid.UUID, inputs []SubmitResponseInput) (*schemas.Result, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	questions := s.questionSvc.GetAll()
	questionMap := make(map[int]Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	var responses []schemas.Response
	for _, input := range inputs {
		q, ok := questionMap[input.QuestionID]
		if !ok {
			return nil, fmt.Errorf("invalid question_id: %d", input.QuestionID)
		}

		responses = append(responses, schemas.Response{
			SessionID:    session.ID,
			QuestionID:   input.QuestionID,
			QuestionText: q.Text,
			AnswerType:   q.Type,
			AnswerValue:  input.AnswerValue,
		})
	}

	if err := s.responseRepo.CreateBatch(responses); err != nil {
		return nil, fmt.Errorf("failed to save responses: %w", err)
	}

	savedResponses, err := s.responseRepo.FindBySessionID(session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch responses: %w", err)
	}

	prompt := BuildPrompt(savedResponses)
	aiResult, provider, err := s.aiSvc.Analyze(prompt)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	approachDetails, _ := json.Marshal(aiResult.ApproachDetails)
	fieldDetails, _ := json.Marshal(aiResult.FieldDetails)

	result := &schemas.Result{
		SessionID:       session.ID,
		ApproachScores:  aiResult.ApproachScores,
		FieldScores:     aiResult.FieldScores,
		Explanation:     aiResult.Summary,
		ApproachDetails: approachDetails,
		FieldDetails:    fieldDetails,
		AIProvider:      provider,
	}

	if err := s.resultRepo.Create(result); err != nil {
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	return result, nil
}

func (s *SessionService) GetResult(sessionID uuid.UUID) (*schemas.Result, error) {
	return s.resultRepo.FindBySessionID(sessionID)
}
