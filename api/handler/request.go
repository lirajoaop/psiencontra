package handler

type SubmitResponsesRequest struct {
	Responses []ResponseInput `json:"responses" binding:"required"`
}

type ResponseInput struct {
	QuestionID  int    `json:"question_id" binding:"required"`
	AnswerValue string `json:"answer_value" binding:"required"`
}
