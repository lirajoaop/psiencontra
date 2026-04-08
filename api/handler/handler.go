package handler

import (
	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/repository"
	"github.com/joaop/psiencontra/api/service"
)

var (
	SessionSvc     *service.SessionService
	QuestionSvc    *service.QuestionService
	PDFSvc         *service.PDFService
	AuthSvc        *service.AuthService
	GoogleOAuthSvc *service.GoogleOAuthService
)

func Init() {
	sessionRepo := repository.NewSessionRepo(config.DB)
	responseRepo := repository.NewResponseRepo(config.DB)
	resultRepo := repository.NewResultRepo(config.DB)
	userRepo := repository.NewUserRepo(config.DB)

	QuestionSvc = service.NewQuestionService()
	aiSvc := service.NewAIService()
	PDFSvc = service.NewPDFService()

	SessionSvc = service.NewSessionService(sessionRepo, responseRepo, resultRepo, QuestionSvc, aiSvc)

	jwtSecret := config.GetEnv("JWT_SECRET", "dev-only-insecure-secret-change-me")
	AuthSvc = service.NewAuthService(userRepo, jwtSecret)

	GoogleOAuthSvc = service.NewGoogleOAuthService(
		config.GetEnv("GOOGLE_CLIENT_ID", ""),
		config.GetEnv("GOOGLE_CLIENT_SECRET", ""),
		config.GetEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
	)
}
