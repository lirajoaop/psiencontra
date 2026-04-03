package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joaop/psiencontra/api/handler"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", handler.HealthCheck)
		v1.GET("/questions", handler.GetQuestions)
		v1.POST("/sessions", handler.CreateSession)
		v1.POST("/sessions/:id/responses", handler.SubmitResponses)
		v1.GET("/sessions/:id/result", handler.GetResult)
		v1.GET("/sessions/:id/pdf", handler.DownloadPDF)
	}
}
