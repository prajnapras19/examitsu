package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/storage"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/mcqoption"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participant"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participantsession"
	"github.com/prajnapras19/project-form-exam-sman2/backend/question"
	"github.com/prajnapras19/project-form-exam-sman2/backend/submission"
)

type Handler interface {
	LoginAdmin(*gin.Context)
	IsLoggedInAsAdmin(*gin.Context)

	CreateExam(*gin.Context)
	GetExamBySerial(*gin.Context)
	GetExams(*gin.Context)
	UpdateExam(*gin.Context)
	DeleteExamBySerial(*gin.Context)
	GetOpenedExam(*gin.Context)
	GetAllOpenedExams(*gin.Context)
	GetExamTemplate(*gin.Context)
	UploadExam(*gin.Context)

	CreateQuestion(*gin.Context)
	GetUploadQuestionBlobURL(*gin.Context)
	GetQuestions(*gin.Context)
	GetQuestionByID(*gin.Context)
	UpdateQuestion(*gin.Context)
	DeleteQuestionBySerial(*gin.Context)

	// exam session auth
	StartExam(*gin.Context)
	IsSessionAuthorized(*gin.Context)
	CheckSession(*gin.Context)
	AuthorizeSession(*gin.Context)

	// exam session
	GetQuestionsIDByExamSerial(*gin.Context)
	GetQuestionWithOptions(*gin.Context)
	SubmitAnswer(*gin.Context)
	SubmitExam(*gin.Context)

	CreateMcqOption(*gin.Context)
	GetMcqOptionsByQuestionID(*gin.Context)
	UpdateMcqOption(*gin.Context)
	DeleteMcqOptionByID(*gin.Context)

	CreateParticipant(*gin.Context)
	GetParticipantsByExamSerial(*gin.Context)
	GetParticipantByID(*gin.Context)
	UpdateParticipant(*gin.Context)
	DeleteParticipantByID(*gin.Context)
	GetParticipantsReport(*gin.Context)

	LoginProctor(*gin.Context)
	IsLoggedInAsProctor(*gin.Context)
}

type handler struct {
	cfg                       *config.Config
	adminAuthService          adminauth.Service
	examService               exam.Service
	questionService           question.Service
	mcqOptionService          mcqoption.Service
	participantService        participant.Service
	submissionService         submission.Service
	storageService            storage.Service
	participantSessionService participantsession.Service
}

func NewHandler(
	cfg *config.Config,
	adminAuthService adminauth.Service,
	examService exam.Service,
	questionService question.Service,
	mcqOptionService mcqoption.Service,
	participantService participant.Service,
	submissionService submission.Service,
	storageService storage.Service,
	participantSessionService participantsession.Service,
) Handler {
	return &handler{
		cfg:                       cfg,
		adminAuthService:          adminAuthService,
		examService:               examService,
		questionService:           questionService,
		mcqOptionService:          mcqOptionService,
		participantService:        participantService,
		submissionService:         submissionService,
		storageService:            storageService,
		participantSessionService: participantSessionService,
	}
}
