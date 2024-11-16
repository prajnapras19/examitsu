package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/mcqoption"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participant"
	"github.com/prajnapras19/project-form-exam-sman2/backend/question"
)

type Handler interface {
	LoginAdmin(*gin.Context)
	IsLoggedInAsAdmin(*gin.Context)

	CreateExam(*gin.Context)
	GetExamBySerial(*gin.Context)
	GetExams(*gin.Context)
	UpdateExam(*gin.Context)
	DeleteExamBySerial(*gin.Context)

	CreateQuestion(*gin.Context)
	GetQuestions(*gin.Context)
	GetQuestionByID(*gin.Context)
	UpdateQuestion(*gin.Context)
	DeleteQuestionBySerial(*gin.Context)

	CreateMcqOption(*gin.Context)
	GetMcqOptionsByQuestionID(*gin.Context)
	UpdateMcqOption(*gin.Context)
	DeleteMcqOptionByID(*gin.Context)

	CreateParticipant(*gin.Context)
	GetParticipantsByExamSerial(*gin.Context)
	UpdateParticipant(*gin.Context)
	DeleteParticipantByID(*gin.Context)
}

type handler struct {
	adminAuthService   adminauth.Service
	examService        exam.Service
	questionService    question.Service
	mcqOptionService   mcqoption.Service
	participantService participant.Service
}

func NewHandler(
	adminAuthService adminauth.Service,
	examService exam.Service,
	questionService question.Service,
	mcqOptionService mcqoption.Service,
	participantService participant.Service,
) Handler {
	return &handler{
		adminAuthService:   adminAuthService,
		examService:        examService,
		questionService:    questionService,
		mcqOptionService:   mcqOptionService,
		participantService: participantService,
	}
}
