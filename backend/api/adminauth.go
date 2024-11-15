package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/adminauth"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

/***
	entity
***/

type LoginAdminRequest struct {
	Password string `json:"password" binding:"required"`
}

type LoginAdminResponse struct {
	Token string `json:"token"`
}

/***
	handler
***/

func (h *handler) LoginAdmin(c *gin.Context) {
	var req LoginAdminRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	svcReq := h.MapLoginAdminRequestToAdminAuthLoginRequest(&req)
	svcRes, err := h.adminAuthService.Login(svcReq)

	if err != nil {
		if errors.Is(err, lib.ErrIncorrectPassword) {
			c.JSON(http.StatusBadRequest, lib.BaseResponse{
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, lib.BaseResponse{
				Message: err.Error(),
			})
		}
		return
	}

	res := h.MapAdminAuthLoginResponseToLoginResponse(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

/***
	mapping
***/

func (h *handler) MapLoginAdminRequestToAdminAuthLoginRequest(req *LoginAdminRequest) *adminauth.LoginRequest {
	return &adminauth.LoginRequest{
		Password: req.Password,
	}
}

func (h *handler) MapAdminAuthLoginResponseToLoginResponse(svcRes *adminauth.LoginResponse) *LoginAdminResponse {
	return &LoginAdminResponse{
		Token: svcRes.Token,
	}
}
