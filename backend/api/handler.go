package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
)

type Handler interface {
	LoginAdmin(*gin.Context)
	IsLoggedInAsAdmin(*gin.Context)

	CreateExam(*gin.Context)
}

type handler struct {
	adminAuthService adminauth.Service
	examService      exam.Service
}

func NewHandler(
	adminAuthService adminauth.Service,
	examService exam.Service,
) Handler {
	return &handler{
		adminAuthService: adminAuthService,
		examService:      examService,
	}
}
