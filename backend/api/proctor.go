package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

/***
	entity
***/

/***
	handler
***/

func (h *handler) LoginProctor(c *gin.Context) {
	var req LoginAdminRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	svcReq := h.MapLoginAdminRequestToAdminAuthLoginRequest(&req)
	svcRes, err := h.adminAuthService.LoginProctor(svcReq)

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

func (h *handler) IsLoggedInAsProctor(c *gin.Context) {
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

func (h *handler) CheckSession(c *gin.Context) {
	// TODO: return data need to be seen by proctor to be able to authorize session
}

func (h *handler) AuthorizeSession(c *gin.Context) {
	// TODO: authorize session as proctor
}

/***
	mapping
***/
