package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/service"
)

func CreateSession(c *gin.Context) {
	session, err := SessionSvc.CreateSession()
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

	var inputs []service.SubmitResponseInput
	for _, r := range req.Responses {
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

	sendSuccess(c, result)
}
