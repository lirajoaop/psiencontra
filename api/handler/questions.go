package handler

import "github.com/gin-gonic/gin"

func GetQuestions(c *gin.Context) {
	questions := QuestionSvc.GetAll()
	sendSuccess(c, questions)
}
