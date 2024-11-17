package main

import (
	"fmt"
	"net/http"

	rmq "github.com/adjust/rmq/v5"
	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/api"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/mysql"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/redis"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/mcqoption"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participant"
	"github.com/prajnapras19/project-form-exam-sman2/backend/question"
	"github.com/prajnapras19/project-form-exam-sman2/backend/submission"
	"github.com/prajnapras19/project-form-exam-sman2/backend/worker"
)

func main() {
	cfg := config.Get()
	if cfg.Role == constants.Worker {
		initWorker(cfg)
	} else {
		initDefault(cfg)
	}
}

func initDefault(cfg *config.Config) {
	// clients
	dbmysql := mysql.NewService(cfg.MySQLConfig)
	dbredis := redis.NewService(cfg.RedisConfig)
	redisMQConnection, err := rmq.OpenConnectionWithRedisClient(constants.Examitsu, dbredis.GetClient(), nil)
	if err != nil {
		panic(err)
	}
	updateAnswerQueue, err := redisMQConnection.OpenQueue(constants.UpdateAnswerQueueName)
	if err != nil {
		panic(err)
	}

	// repositories
	examRepository := exam.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	questionRepository := question.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	mcqOptionRepository := mcqoption.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	participantRepository := participant.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())
	submissionRepository := submission.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())

	// services
	adminAuthService := adminauth.NewService(cfg)
	examService := exam.NewService(examRepository)
	questionService := question.NewService(questionRepository)
	mcqOptionService := mcqoption.NewService(mcqOptionRepository)
	participantService := participant.NewService(cfg, participantRepository, examService)
	submissionService := submission.NewService(submissionRepository, dbredis.GetClient(), updateAnswerQueue)

	// handlers
	handler := api.NewHandler(
		cfg,
		adminAuthService,
		examService,
		questionService,
		mcqOptionService,
		participantService,
		submissionService,
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
	examSessionGroup.GET("/:serial/questions/:id", handler.GetQuestionWithOptions)
	examSessionGroup.POST("/:serial/questions/:id", handler.SubmitAnswer)

	router.Run(fmt.Sprintf(":%d", cfg.RESTPort))
}

func initWorker(cfg *config.Config) {
	// clients
	dbmysql := mysql.NewService(cfg.MySQLConfig)
	dbredis := redis.NewService(cfg.RedisConfig)
	redisMQConnection, err := rmq.OpenConnectionWithRedisClient(constants.Examitsu, dbredis.GetClient(), nil)
	if err != nil {
		panic(err)
	}
	updateAnswerQueue, err := redisMQConnection.OpenQueue(constants.UpdateAnswerQueueName)
	if err != nil {
		panic(err)
	}

	// repositories
	submissionRepository := submission.NewRepository(cfg, dbmysql.GetDB(), dbredis.GetClient())

	// services
	submissionService := submission.NewService(submissionRepository, dbredis.GetClient(), updateAnswerQueue)

	// consumers
	updateAnswerConsumer := worker.NewUpdateAnswerQueueConsumer(submissionService)

	// routes
	router := gin.Default()
	if cfg.AllowCORS {
		router.Use(api.CORSMiddleware())
	}
	router.GET("/_health", func(gc *gin.Context) {
		gc.Status(http.StatusOK)
	})

	// workers
	workerService := worker.NewService(
		cfg,
		updateAnswerQueue,
		updateAnswerConsumer,
	)
	workerService.InitConsumers()

	router.Run(fmt.Sprintf(":%d", cfg.RESTPort))
}
