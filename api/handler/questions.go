package handler

import "github.com/gin-gonic/gin"

func GetQuestions(c *gin.Context) {
	qType := c.DefaultQuery("type", "simple")
	questions := QuestionSvc.GetShuffledByType(qType)
	sendSuccess(c, questions)
}
