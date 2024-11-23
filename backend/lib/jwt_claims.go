package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
)

// TODO: add more claims
type JWTClaims struct {
	jwt.StandardClaims
	Username string `json:"username" binding:"required"`
}

func GetJWTClaimsFromContext(c *gin.Context) (*JWTClaims, error) {
	if val, exists := c.Get(constants.JWTClaims); exists {
		if res, ok := val.(*JWTClaims); ok {
			return res, nil
		}
		return nil, ErrFailedToParseJWTClaimsInContext
	}
	return nil, ErrJWTClaimsNotFoundInContext
}

type ExamTokenJWTClaims struct {
	jwt.StandardClaims

	ParticipantID uint   `json:"pid"`
	SessionSerial string `json:"sess"`
}

func GetExamTokenJWTClaimsFromContext(c *gin.Context) (*ExamTokenJWTClaims, error) {
	if val, exists := c.Get(constants.JWTClaims); exists {
		if res, ok := val.(*ExamTokenJWTClaims); ok {
			return res, nil
		}
		return nil, ErrFailedToParseJWTClaimsInContext
	}
	return nil, ErrJWTClaimsNotFoundInContext
}
