package adminauth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	Login(req *LoginRequest) (*LoginResponse, error)
	GenerateToken(username string) string
	ValidateToken(tokenString string) (*lib.JWTClaims, error)
	VerifyToken(token *jwt.Token) (interface{}, error)
}
type service struct {
	cfg *config.Config
}

func NewService(
	cfg *config.Config,
) Service {
	return &service{
		cfg: cfg,
	}
}

func (s *service) Login(req *LoginRequest) (*LoginResponse, error) {
	if req.Password != s.cfg.SystemPassword {
		return nil, lib.ErrIncorrectPassword
	}
	return &LoginResponse{
		Token: s.GenerateToken(constants.SystemUser),
	}, nil
}

func (s *service) GenerateToken(username string) string {
	claims := lib.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    s.cfg.AuthConfig.ApplicationName,
			ExpiresAt: time.Now().Add(s.cfg.AuthConfig.LoginTokenExpirationDuration).Unix(),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	signedToken, _ := token.SignedString(s.cfg.AuthConfig.SignatureKey)
	return signedToken
}

func (s *service) VerifyToken(token *jwt.Token) (interface{}, error) {
	if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, lib.ErrSigningMethodInvalid
	} else if method != jwt.SigningMethodHS256 {
		return nil, lib.ErrSigningMethodInvalid
	}
	return s.cfg.AuthConfig.SignatureKey, nil
}

func (s *service) ValidateToken(tokenString string) (*lib.JWTClaims, error) {
	claims := &lib.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, s.VerifyToken)
	if err != nil {
		return nil, lib.ErrUnauthorizedRequest
	}
	if !token.Valid {
		return nil, lib.ErrUnauthorizedRequest
	}
	claims, ok := token.Claims.(*lib.JWTClaims)
	if !ok {
		return nil, lib.ErrUnauthorizedRequest
	}
	if claims.Username != constants.SystemUser {
		return nil, lib.ErrUnauthorizedRequest
	}

	return claims, nil
}
