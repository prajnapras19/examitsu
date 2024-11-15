package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/api"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
)

func main() {
	cfg := config.Get()
	initDefault(cfg)
}

func initDefault(cfg *config.Config) {
	// handlers
	handler := api.NewHandler()

	// routes
	router := gin.Default()
	if cfg.AllowCORS {
		router.Use(api.CORSMiddleware())
	}
	router.GET("/_health", func(gc *gin.Context) {
		gc.Status(http.StatusOK)
	})

	apiV1 := router.Group("/api/v1")

	apiV1.POST("/admin/login", handler.LoginAdmin)

	router.Run(fmt.Sprintf(":%d", cfg.RESTPort))
}
