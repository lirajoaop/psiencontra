package handler

import (
	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/repository"
	"github.com/joaop/psiencontra/api/service"
)

var (
	SessionSvc  *service.SessionService
	QuestionSvc *service.QuestionService
	PDFSvc      *service.PDFService
)

func Init() {
	sessionRepo := repository.NewSessionRepo(config.DB)
	responseRepo := repository.NewResponseRepo(config.DB)
	resultRepo := repository.NewResultRepo(config.DB)

	QuestionSvc = service.NewQuestionService()
	aiSvc := service.NewAIService()
	PDFSvc = service.NewPDFService()

	SessionSvc = service.NewSessionService(sessionRepo, responseRepo, resultRepo, QuestionSvc, aiSvc)
}
