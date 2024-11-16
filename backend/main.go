package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/api"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/mysql"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
)

func main() {
	cfg := config.Get()
	initDefault(cfg)
}

func initDefault(cfg *config.Config) {
	// clients
	dbmysql := mysql.NewService(cfg.MySQLConfig)

	// repositories
	examRepository := exam.NewRepository(dbmysql.GetDB())

	// services
	adminAuthService := adminauth.NewService(cfg)
	examService := exam.NewService(examRepository)

	// handlers
	handler := api.NewHandler(
		adminAuthService,
		examService,
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
	adminGroup.PATCH("/exams/:serial", handler.UpdateExam)
	adminGroup.DELETE("/exams/:serial", handler.DeleteExamBySerial)

	router.Run(fmt.Sprintf(":%d", cfg.RESTPort))
}
