package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participant"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participantsession"
	"gorm.io/gorm"
)

/***
	entity
***/

type CreateParticipantsRequest struct {
	ExamSerial             string   `json:"exam_serial" binding:"required"`
	ExamID                 uint     `json:"-"`
	Names                  []string `json:"names" binding:"required"`
	AllowedDurationMinutes uint     `json:"allowed_duration_minutes" binding:"required"`
}

type ParticipantData struct {
	ID                     uint       `json:"id"`
	Name                   string     `json:"name"`
	StartedAt              *time.Time `json:"started_at"`
	EndedAt                *time.Time `json:"ended_at"`
	AllowedDurationMinutes uint       `json:"allowed_duration_minutes"`
	TotalPoint             int        `json:"total_point"`
}

type UpdateParticipantRequest struct {
	ID                     uint   `json:"-"`
	Name                   string `json:"name" binding:"required"`
	AllowedDurationMinutes uint   `json:"allowed_duration_minutes" binding:"required"`
}

type StartExamRequest struct {
	ExamID uint   `json:"-"`
	Name   string `json:"name" binding:"required"`
}

type StartExamResponse struct {
	Token         string `json:"token"`
	SessionSerial string `json:"session"`
}

/***
	handler
***/

func (h *handler) CreateParticipant(c *gin.Context) {
	var req CreateParticipantsRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	exam, err := h.examService.GetExamBySerial(req.ExamSerial)
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

	req.ExamID = exam.ID
	svcReq := h.MapCreateParticipantsRequestToParticipantEntityList(&req)

	svcRes, err := h.participantService.CreateParticipants(svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapParticipantEntityListToParticipantDataList(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetParticipantsByExamSerial(c *gin.Context) {
	examSerial := c.Param(constants.Serial)
	exam, err := h.examService.GetExamBySerial(examSerial)
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

	svcRes, err := h.participantService.GetParticipantsByExamID(exam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	totalPoints, err := h.participantService.GetParticipantTotalPointsByExamID(exam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapGetParticipantsByExamSerialResponse(svcRes, totalPoints)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetParticipantByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)

	svcRes, err := h.participantService.GetParticipantByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapParticipantEntityToParticipantData(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) UpdateParticipant(c *gin.Context) {
	var req UpdateParticipantRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)
	req.ID = uint(id)
	svcReq := h.MapUpdateParticipantRequestToParticipantEntity(&req)

	err := h.participantService.UpdateParticipant(svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

func (h *handler) DeleteParticipantByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)

	err := h.participantService.DeleteParticipantByID(uint(id))
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
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

func (h *handler) StartExam(c *gin.Context) {
	// TODO: rate limit per participant
	var req StartExamRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	exam, err := h.examService.GetExamBySerial(c.Param(constants.Serial))
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
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrExamNotFound.Error(),
		})
		return
	}

	req.ExamID = exam.ID
	participant, err := h.participantService.GetParticipantByExamIDAndName(req.ExamID, req.Name)
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

	if participant.EndedAt != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamAlreadySubmitted.Error(),
		})
		return
	}

	participantSession, err := h.participantSessionService.CreateParticipantSession(&participantsession.ParticipantSession{
		ParticipantID: participant.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	// generate exam token
	examToken := h.participantService.GenerateToken(exam.Serial, participant.ID, participantSession.Serial)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data: StartExamResponse{
			Token:         examToken,
			SessionSerial: participantSession.Serial,
		},
	})
}

func (h *handler) SubmitExam(c *gin.Context) {
	jwtClaims, err := lib.GetExamTokenJWTClaimsFromContext(c)
	if err != nil {
		log.Printf("[handler][participant][SubmitExam] error when get jwt: %s", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrUnknownError.Error(),
		})
		return
	}

	participant, err := h.participantService.GetParticipantByID(jwtClaims.ParticipantID)
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
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamNotStarted.Error(),
		})
		return
	}

	if participant.EndedAt != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamAlreadySubmitted.Error(),
		})
		return
	}

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
	if !exam.IsOpen || exam.Serial != c.Param(constants.Serial) {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrExamNotFound.Error(),
		})
		return
	}
	participantSession, err := h.participantSessionService.GetLatestAuthorizedParticipantSessionByParticipantID(participant.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrSessionNotFound.Error(),
		})
		return
	}
	if participantSession.Serial != jwtClaims.SessionSerial {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrSessionNotFound.Error(),
		})
		return
	}

	currentTime := time.Now()
	participant.EndedAt = &currentTime
	err = h.participantService.UpdateParticipant(participant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

func (h *handler) IsSessionAuthorized(c *gin.Context) {
	// TODO: return success, 404, 401, or 500 depending on if the request exam token session serial is authorized or not
	jwtClaims, err := lib.GetExamTokenJWTClaimsFromContext(c)
	if err != nil {
		log.Printf("[handler][participant][IsSessionAuthorized] error when get jwt: %s", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrUnknownError.Error(),
		})
		return
	}

	participant, err := h.participantService.GetParticipantByID(jwtClaims.ParticipantID)
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

	if participant.EndedAt != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamAlreadySubmitted.Error(),
		})
		return
	}

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
	if !exam.IsOpen || exam.Serial != c.Param(constants.Serial) {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrExamNotFound.Error(),
		})
		return
	}
	participantSession, err := h.participantSessionService.GetLatestAuthorizedParticipantSessionByParticipantID(participant.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrSessionNotFound.Error(),
		})
		return
	}
	if participantSession.Serial != jwtClaims.SessionSerial {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrSessionNotFound.Error(),
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

func (h *handler) MapCreateParticipantsRequestToParticipantEntityList(req *CreateParticipantsRequest) []*participant.Participant {
	res := []*participant.Participant{}
	for _, name := range req.Names {
		res = append(res, &participant.Participant{
			ExamID:                 req.ExamID,
			Name:                   name,
			AllowedDurationMinutes: req.AllowedDurationMinutes,
		})
	}
	return res
}

func (h *handler) MapParticipantEntityToParticipantData(svcRes *participant.Participant) *ParticipantData {
	return &ParticipantData{
		ID:                     svcRes.ID,
		Name:                   svcRes.Name,
		StartedAt:              svcRes.StartedAt,
		EndedAt:                svcRes.EndedAt,
		AllowedDurationMinutes: svcRes.AllowedDurationMinutes,
	}
}

func (h *handler) MapParticipantEntityListToParticipantDataList(svcRes []*participant.Participant) []*ParticipantData {
	res := []*ParticipantData{}
	for _, obj := range svcRes {
		res = append(res, h.MapParticipantEntityToParticipantData(obj))
	}
	return res
}

func (h *handler) MapUpdateParticipantRequestToParticipantEntity(req *UpdateParticipantRequest) *participant.Participant {
	return &participant.Participant{
		BaseModel: lib.BaseModel{
			Model: gorm.Model{
				ID: req.ID,
			},
		},
		Name:                   req.Name,
		AllowedDurationMinutes: req.AllowedDurationMinutes,
	}
}

func (h *handler) MapGetParticipantsByExamSerialResponse(svcRes []*participant.Participant, totalPoints []*participant.ParticipantTotalPoint) []*ParticipantData {
	participantData := h.MapParticipantEntityListToParticipantDataList(svcRes)
	totalPointsMap := map[uint]int{}
	for i := range totalPoints {
		totalPointsMap[totalPoints[i].ParticipantID] = totalPoints[i].TotalPoint
	}
	for i := range participantData {
		participantData[i].TotalPoint = totalPointsMap[participantData[i].ID]
	}
	return participantData
}
