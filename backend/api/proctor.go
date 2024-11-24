package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

/***
	entity
***/

type CheckSessionResponse struct {
	IsStartExam bool             `json:"is_start_exam"`
	IsSubmitted bool             `json:"is_submitted"`
	Participant *ParticipantData `json:"participant"`
	Exam        *ExamData        `json:"exam"`
}

type AuthorizeSessionRequest struct {
	AllowedDurationMinutes uint `json:"allowed_duration_minutes" binding:"required"`
}

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
	// return data need to be seen by proctor to be able to authorize session

	participantSession, err := h.participantSessionService.GetParticipantSessionBySerial(c.Param(constants.Serial))
	if err != nil {
		if errors.Is(err, lib.ErrParticipantSessionNotFound) {
			c.JSON(http.StatusNotFound, lib.BaseResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := &CheckSessionResponse{}

	participant, err := h.participantService.GetParticipantByID(participantSession.ParticipantID)
	if err != nil {
		if errors.Is(err, lib.ErrParticipantNotFound) {
			c.JSON(http.StatusNotFound, lib.BaseResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}
	if participant.StartedAt == nil {
		res.IsStartExam = true
	} else {
		if participant.EndedAt != nil || participant.StartedAt.Add(time.Duration(participant.AllowedDurationMinutes)*time.Minute).Before(time.Now()) {
			res.IsSubmitted = true
		}
	}
	res.Participant = h.MapParticipantEntityToParticipantData(participant)

	exam, err := h.examService.GetExamByID(participant.ExamID)
	if err != nil {
		if errors.Is(err, lib.ErrExamNotFound) {
			c.JSON(http.StatusNotFound, lib.BaseResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}
	if !exam.IsOpen {
		if !exam.IsOpen {
			c.JSON(http.StatusNotFound, lib.BaseResponse{
				Message: lib.ErrExamNotFound.Error(),
			})
			return
		}
	}
	res.Exam = h.MapExamEntityToExamData(exam)

	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) AuthorizeSession(c *gin.Context) {
	// authorize session as proctor
	var req AuthorizeSessionRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	// TODO: validation, but it seems like even if an invalid session is authorized, it won't affect anything because it's being validated in other endpoints
	err := h.participantSessionService.AuthorizeSession(c.Param(constants.Serial), req.AllowedDurationMinutes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: constants.Success,
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

/***
	mapping
***/
