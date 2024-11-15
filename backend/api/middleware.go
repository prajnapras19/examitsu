package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

var (
	ErrUnauthorizedRequest = errors.New("unauthorized request")
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func JWTAdminMiddleware(adminAuthService adminauth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, lib.BaseResponse{
				Message: ErrUnauthorizedRequest.Error(),
			})
			c.Abort()
			return
		}
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		claims, err := adminAuthService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, lib.BaseResponse{
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set(constants.JWTClaims, claims)

		c.Next()
	}
}
