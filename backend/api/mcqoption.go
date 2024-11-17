package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"github.com/prajnapras19/project-form-exam-sman2/backend/mcqoption"
	"gorm.io/gorm"
)

/***
	entity
***/

type CreateMcqOptionRequest struct {
	McqOptionID uint   `json:"question_id" binding:"required"`
	Description string `json:"description"`
	Point       int    `json:"point"`
}

type McqOptionData struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Point       int    `json:"point"`
}

type UpdateMcqOptionRequest struct {
	ID          uint   `json:"-"`
	Description string `json:"description"`
	Point       int    `json:"point"`
}

type McqOptionWithoutPointData struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
}

/***
	handler
***/

func (h *handler) CreateMcqOption(c *gin.Context) {
	var req CreateMcqOptionRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	svcReq := h.MapCreateMcqOptionRequestToMcqOptionEntity(&req)

	svcRes, err := h.mcqOptionService.CreateMcqOption(svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapMcqOptionEntityToMcqOptionData(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetMcqOptionsByQuestionID(c *gin.Context) {
	questionID, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)
	svcRes, err := h.mcqOptionService.GetMcqOptionsByQuestionID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapMcqOptionEntityListToMcqOptionDataList(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) UpdateMcqOption(c *gin.Context) {
	var req UpdateMcqOptionRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)
	req.ID = uint(id)
	svcReq := h.MapUpdateMcqOptionRequestToMcqOptionEntity(&req)

	err := h.mcqOptionService.UpdateMcqOption(svcReq)
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

func (h *handler) DeleteMcqOptionByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)

	err := h.mcqOptionService.DeleteMcqOptionByID(uint(id))
	if err != nil {
		if errors.Is(err, lib.ErrMcqOptionNotFound) {
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

/***
	mapping
***/

func (h *handler) MapCreateMcqOptionRequestToMcqOptionEntity(req *CreateMcqOptionRequest) *mcqoption.McqOption {
	return &mcqoption.McqOption{
		QuestionID:  req.McqOptionID,
		Description: req.Description,
		Point:       req.Point,
	}
}

func (h *handler) MapMcqOptionEntityToMcqOptionData(svcRes *mcqoption.McqOption) *McqOptionData {
	return &McqOptionData{
		ID:          svcRes.ID,
		Description: svcRes.Description,
		Point:       svcRes.Point,
	}
}

func (h *handler) MapMcqOptionEntityListToMcqOptionDataList(svcRes []*mcqoption.McqOption) []*McqOptionData {
	res := []*McqOptionData{}
	for _, obj := range svcRes {
		res = append(res, h.MapMcqOptionEntityToMcqOptionData(obj))
	}
	return res
}

func (h *handler) MapMcqOptionEntityToMcqOptionWithoutPointData(svcRes *mcqoption.McqOption) *McqOptionWithoutPointData {
	return &McqOptionWithoutPointData{
		ID:          svcRes.ID,
		Description: svcRes.Description,
	}
}

func (h *handler) MapMcqOptionEntityListToMcqOptionWithoutPointDataList(svcRes []*mcqoption.McqOption) []*McqOptionWithoutPointData {
	res := []*McqOptionWithoutPointData{}
	for _, obj := range svcRes {
		res = append(res, h.MapMcqOptionEntityToMcqOptionWithoutPointData(obj))
	}
	return res
}

func (h *handler) MapUpdateMcqOptionRequestToMcqOptionEntity(req *UpdateMcqOptionRequest) *mcqoption.McqOption {
	return &mcqoption.McqOption{
		BaseModel: lib.BaseModel{
			Model: gorm.Model{
				ID: req.ID,
			},
		},
		Description: req.Description,
		Point:       req.Point,
	}
}
