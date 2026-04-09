package handler

import (
	"log"

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

	// JWT_SECRET must be set explicitly. We refuse to fall back to a hard-coded
	// value in production because that would make tokens trivially forgeable.
	// Local development can opt in via APP_ENV=development, which uses an
	// obviously-insecure placeholder so the dev never confuses it with a real
	// secret.
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		if config.GetEnv("APP_ENV", "") == "development" {
			config.Log.Info.Println("WARNING: APP_ENV=development — using insecure placeholder JWT_SECRET. NEVER set APP_ENV=development in production.")
			jwtSecret = "insecure-development-only-do-not-use-in-production"
		} else {
			log.Fatal("JWT_SECRET environment variable is required (set APP_ENV=development to use an insecure dev placeholder)")
		}
	}
	AuthSvc = service.NewAuthService(userRepo, jwtSecret)

	// All three Google OAuth env vars must be set together. No localhost
	// default for the redirect URL: a silent default would let production
	// boot without the real callback URL and only fail at Google's
	// redirect_uri_mismatch check, after the user clicked the login button.
	// If any of the three is missing, GoogleOAuthService.Enabled() returns
	// false and the handlers degrade to "google login not configured".
	GoogleOAuthSvc = service.NewGoogleOAuthService(
		config.GetEnv("GOOGLE_CLIENT_ID", ""),
		config.GetEnv("GOOGLE_CLIENT_SECRET", ""),
		config.GetEnv("GOOGLE_REDIRECT_URL", ""),
	)
}
