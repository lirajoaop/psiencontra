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

		// Auth
		v1.POST("/auth/register", handler.Register)
		v1.POST("/auth/login", handler.Login)
		v1.POST("/auth/logout", handler.Logout)
		v1.GET("/auth/me", handler.Me)
		v1.GET("/auth/google", handler.GoogleStart)
		v1.GET("/auth/google/callback", handler.GoogleCallback)

		// Sessions: anonymous-friendly. OptionalAuth attaches the user_id
		// when a JWT cookie is present, otherwise the session is created
		// without an owner (anonymous flow).
		v1.POST("/sessions", OptionalAuth(), handler.CreateSession)
		v1.POST("/sessions/:id/responses", handler.SubmitResponses)
		v1.GET("/sessions/:id/result", handler.GetResult)
		v1.GET("/sessions/:id/pdf", handler.DownloadPDF)

		// User-scoped history: only completed sessions belonging to the caller.
		v1.GET("/user/sessions", RequireAuth(), handler.GetUserHistory)
	}
}
