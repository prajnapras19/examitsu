package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
)

type Handler interface {
	LoginAdmin(*gin.Context)
	IsLoggedInAsAdmin(*gin.Context)
}

type handler struct {
	adminAuthService adminauth.Service
}

func NewHandler(
	adminAuthService adminauth.Service,
) Handler {
	return &handler{
		adminAuthService: adminAuthService,
	}
}
