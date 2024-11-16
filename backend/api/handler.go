package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/question"
)

type Handler interface {
	LoginAdmin(*gin.Context)
	IsLoggedInAsAdmin(*gin.Context)

	CreateExam(*gin.Context)
	GetExamBySerial(*gin.Context)
	GetExams(c *gin.Context)
	UpdateExam(*gin.Context)
	DeleteExamBySerial(*gin.Context)

	CreateQuestion(*gin.Context)
	GetQuestions(c *gin.Context)
	UpdateQuestion(*gin.Context)
	DeleteQuestionBySerial(*gin.Context)
}

type handler struct {
	adminAuthService adminauth.Service
	examService      exam.Service
	questionService  question.Service
}

func NewHandler(
	adminAuthService adminauth.Service,
	examService exam.Service,
	questionService question.Service,
) Handler {
	return &handler{
		adminAuthService: adminAuthService,
		examService:      examService,
		questionService:  questionService,
	}
}
