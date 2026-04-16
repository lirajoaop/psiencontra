package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/service"
)

func CreateSession(c *gin.Context) {
	var userID *uuid.UUID
	if id, ok := UserIDFromContext(c); ok {
		userID = &id
	}

	var req struct {
		QuestionnaireType string `json:"questionnaire_type"`
	}
	// Body is optional — default to "simple" if missing or unparseable.
	_ = c.ShouldBindJSON(&req)

	session, err := SessionSvc.CreateSession(userID, req.QuestionnaireType)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "failed to create session")
		return
	}
	sendSuccess(c, gin.H{"id": session.ID})
}

func SubmitResponses(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "invalid session id")
		return
	}

	var req SubmitResponsesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	const maxAnswerLen = 500
	var inputs []service.SubmitResponseInput
	for _, r := range req.Responses {
		if len([]rune(r.AnswerValue)) > maxAnswerLen {
			sendError(c, http.StatusBadRequest, "answer exceeds 500 character limit")
			return
		}
		inputs = append(inputs, service.SubmitResponseInput{
			QuestionID:  r.QuestionID,
			AnswerValue: r.AnswerValue,
		})
	}

	result, err := SessionSvc.SubmitResponses(id, inputs)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

func GetUserHistory(c *gin.Context) {
	userID, ok := UserIDFromContext(c)
	if !ok {
		sendError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessions, err := SessionSvc.GetUserHistory(userID)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "failed to fetch history")
		return
	}

	// Slim DTO: preview cards only need the top-ranking approach/field, which
	// the frontend derives from the score maps. Skipping explanation + details
	// keeps the payload light even for users with many sessions.
	out := make([]gin.H, 0, len(sessions))
	for _, s := range sessions {
		if s.Result == nil {
			continue
		}
		out = append(out, gin.H{
			"id":                 s.ID,
			"created_at":         s.CreatedAt,
			"questionnaire_type": s.QuestionnaireType,
			"approach_scores":    s.Result.ApproachScores,
			"field_scores":       s.Result.FieldScores,
		})
	}

	sendSuccess(c, out)
}

func GetResult(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendError(c, http.StatusBadRequest, "invalid session id")
		return
	}

	result, err := SessionSvc.GetResult(id)
	if err != nil {
		sendError(c, http.StatusNotFound, "result not found")
		return
	}

	questionnaireType := "simple"
	if session, err := SessionSvc.GetSession(id); err == nil && session != nil {
		questionnaireType = session.QuestionnaireType
	}

	sendSuccess(c, gin.H{
		"id":                 result.ID,
		"session_id":         result.SessionID,
		"approach_scores":    result.ApproachScores,
		"field_scores":       result.FieldScores,
		"explanation":        result.Explanation,
		"approach_details":   result.ApproachDetails,
		"field_details":      result.FieldDetails,
		"ai_provider":        result.AIProvider,
		"questionnaire_type": questionnaireType,
	})
}
