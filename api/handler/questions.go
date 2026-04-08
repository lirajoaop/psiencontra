package handler

import "github.com/gin-gonic/gin"

func GetQuestions(c *gin.Context) {
	questions := QuestionSvc.GetAllShuffled()
	sendSuccess(c, questions)
}
