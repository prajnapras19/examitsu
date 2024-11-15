package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

/***
	entity
***/

type CreateExamRequest struct {
	Name   string `json:"name" binding:"required"`
	IsOpen bool   `json:"is_open"`
}

type ExamData struct {
	Serial string `json:"serial"`
	Name   string `json:"name"`
	IsOpen bool   `json:"is_open"`
}

type UpdateExamRequest struct {
	Serial string `json:"-"`
	Name   string `json:"name" binding:"required"`
	IsOpen bool   `json:"is_open"`
}

/***
	handler
***/

func (h *handler) CreateExam(c *gin.Context) {
	var req CreateExamRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	svcReq := h.MapCreateExamRequestToExamEntity(&req)

	svcRes, err := h.examService.CreateExam(svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapExamEntityToExamData(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetExamBySerial(c *gin.Context) {
	svcRes, err := h.examService.GetExamBySerial(c.Param(constants.Serial))
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

	res := h.MapExamEntityToExamData(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetAllExams(c *gin.Context) {
	svcRes, err := h.examService.GetAllExams()
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

	res := h.MapExamEntityListToExamDataList(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) UpdateExam(c *gin.Context) {
	var req UpdateExamRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}
	req.Serial = c.Param(constants.Serial)

	svcReq := h.MapUpdateExamRequestToExamEntity(&req)

	err := h.examService.UpdateExam(svcReq)
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

func (h *handler) DeleteExamBySerial(c *gin.Context) {
	err := h.examService.DeleteExamBySerial(c.Param(constants.Serial))
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
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

/***
	mapping
***/

func (h *handler) MapCreateExamRequestToExamEntity(req *CreateExamRequest) *exam.Exam {
	return &exam.Exam{
		Name:   req.Name,
		IsOpen: req.IsOpen,
	}
}

func (h *handler) MapExamEntityToExamData(svcRes *exam.Exam) *ExamData {
	return &ExamData{
		Serial: svcRes.Serial,
		Name:   svcRes.Name,
		IsOpen: svcRes.IsOpen,
	}
}

func (h *handler) MapExamEntityListToExamDataList(svcRes []*exam.Exam) []*ExamData {
	res := []*ExamData{}
	for _, obj := range svcRes {
		res = append(res, h.MapExamEntityToExamData(obj))
	}
	return res
}

func (h *handler) MapUpdateExamRequestToExamEntity(req *UpdateExamRequest) *exam.Exam {
	return &exam.Exam{
		Serial: req.Serial,
		Name:   req.Name,
		IsOpen: req.IsOpen,
	}
}
