package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/api"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/mysql"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/redis"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/mcqoption"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participant"
	"github.com/prajnapras19/project-form-exam-sman2/backend/question"
)

func main() {
	cfg := config.Get()
	initDefault(cfg)
}

func initDefault(cfg *config.Config) {
	// clients
	dbmysql := mysql.NewService(cfg.MySQLConfig)
	dbredis := redis.NewService(cfg.RedisConfig)

	// repositories
	examRepository := exam.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	questionRepository := question.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	mcqOptionRepository := mcqoption.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	participantRepository := participant.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())

	// services
	adminAuthService := adminauth.NewService(cfg)
	examService := exam.NewService(examRepository)
	questionService := question.NewService(questionRepository)
	mcqOptionService := mcqoption.NewService(mcqOptionRepository)
	participantService := participant.NewService(cfg, participantRepository, examService)

	// handlers
	handler := api.NewHandler(
		cfg,
		adminAuthService,
		examService,
		questionService,
		mcqOptionService,
		participantService,
	)

	// routes
	router := gin.Default()
	if cfg.AllowCORS {
		router.Use(api.CORSMiddleware())
	}
	router.GET("/_health", func(gc *gin.Context) {
		gc.Status(http.StatusOK)
	})

	apiV1 := router.Group("/api/v1")

	adminGroup := apiV1.Group("/admin")
	adminGroup.POST("/login", handler.LoginAdmin)
	adminGroup.Use(api.JWTAdminMiddleware(adminAuthService))
	adminGroup.GET("/is-logged-in", handler.IsLoggedInAsAdmin)

	adminGroup.PUT("/exams", handler.CreateExam)
	adminGroup.POST("/exams", handler.GetExams)
	adminGroup.POST("/exams/:serial", handler.GetExamBySerial)
	adminGroup.PATCH("/exams/:serial", handler.UpdateExam)
	adminGroup.DELETE("/exams/:serial", handler.DeleteExamBySerial)

	adminGroup.PUT("/questions", handler.CreateQuestion)
	adminGroup.POST("/questions", handler.GetQuestions)
	adminGroup.POST("/questions/:id", handler.GetQuestionByID)
	adminGroup.PATCH("/questions/:id", handler.UpdateQuestion)
	adminGroup.DELETE("/questions/:id", handler.DeleteQuestionBySerial)

	adminGroup.PUT("/mcq-options", handler.CreateMcqOption)
	adminGroup.POST("/mcq-options/question-id/:id", handler.GetMcqOptionsByQuestionID)
	adminGroup.PATCH("/mcq-options/:id", handler.UpdateMcqOption)
	adminGroup.DELETE("/mcq-options/:id", handler.DeleteMcqOptionByID)

	adminGroup.PUT("/participants", handler.CreateParticipant)
	adminGroup.POST("/participants/exam-serial/:serial", handler.GetParticipantsByExamSerial)
	adminGroup.POST("/participants/id/:id", handler.GetParticipantByID)
	adminGroup.PATCH("/participants/:id", handler.UpdateParticipant)
	adminGroup.DELETE("/participants/:id", handler.DeleteParticipantByID)

	apiV1.GET("/exams/:serial", handler.GetOpenedExam)
	apiV1.POST("/exams/:serial/start", handler.StartExam)

	examSessionGroup := apiV1.Group("/exam-session")
	examSessionGroup.Use(api.JWTExamTokenMiddleware(participantService))
	examSessionGroup.GET("/:serial/questions", handler.GetQuestionsIDByExamSerial)

	router.Run(fmt.Sprintf(":%d", cfg.RESTPort))
}
