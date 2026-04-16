package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joaop/psiencontra/api/handler"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	// Global rate limit: 60 req/min per IP (burst of 10).
	v1.Use(RateLimit(60, 10))

	{
		v1.GET("/health", handler.HealthCheck)
		v1.GET("/questions", handler.GetQuestions)

		// Auth — stricter limit: 10 req/min per IP.
		auth := v1.Group("/auth")
		auth.Use(RateLimit(10, 5))
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/logout", handler.Logout)
			auth.GET("/me", handler.Me)
			auth.GET("/google", handler.GoogleStart)
			auth.GET("/google/callback", handler.GoogleCallback)
		}

		// Sessions: anonymous-friendly. OptionalAuth attaches the user_id
		// when a JWT cookie is present, otherwise the session is created
		// without an owner (anonymous flow).
		v1.POST("/sessions", OptionalAuth(), handler.CreateSession)

		// Submit responses triggers AI calls — stricter limit: 5 req/min.
		v1.POST("/sessions/:id/responses", RateLimit(5, 2), handler.SubmitResponses)

		v1.GET("/sessions/:id/result", handler.GetResult)
		v1.GET("/sessions/:id/pdf", handler.DownloadPDF)

		// User-scoped history: only completed sessions belonging to the caller.
		v1.GET("/user/sessions", RequireAuth(), handler.GetUserHistory)
	}
}
