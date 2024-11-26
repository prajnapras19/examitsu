package api

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/example"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

/***
	entity
***/

type CreateExamRequest struct {
	Name                   string `json:"name" binding:"required"`
	IsOpen                 bool   `json:"is_open"`
	AllowedDurationMinutes uint   `json:"allowed_duration_minutes"`
}

type ExamData struct {
	Serial                 string `json:"serial"`
	Name                   string `json:"name"`
	IsOpen                 bool   `json:"is_open"`
	AllowedDurationMinutes uint   `json:"allowed_duration_minutes"`
}

type UpdateExamRequest struct {
	Serial                 string `json:"-"`
	Name                   string `json:"name" binding:"required"`
	IsOpen                 bool   `json:"is_open"`
	AllowedDurationMinutes uint   `json:"allowed_duration_minutes"`
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

func (h *handler) GetExams(c *gin.Context) {
	var filter exam.GetExamsFilter

	if err := c.ShouldBind(&filter); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	pagination, err := lib.GetQueryPaginationFromContext(c)
	if err != nil {
		log.Printf("[handler][exam][GetExams] get query pagination error: %s", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	svcRes, err := h.examService.GetExams(pagination, &filter)
	if err != nil {
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

func (h *handler) GetOpenedExam(c *gin.Context) {
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
	if !res.IsOpen {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrExamNotFound.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetAllOpenedExams(c *gin.Context) {
	svcRes, err := h.examService.GetAllOpenedExams()
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

func (h *handler) GetExamTemplate(c *gin.Context) {
	fileContent, err := base64.StdEncoding.DecodeString(example.ExamZipExample)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode Base64 content"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", example.ExamZipExampleFilename))
	c.Data(http.StatusOK, "application/zip", fileContent)
}

func (h *handler) UploadExam(c *gin.Context) {
	// Retrieve the uploaded ZIP file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: "Failed to get uploaded file",
		})
		return
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: "Failed to get uploaded file",
		})
		return
	}
	defer uploadedFile.Close()

	// Read the file into memory
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, uploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: "Failed to read file",
		})
		return
	}

	// Open the ZIP archive
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), file.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: "Failed to read ZIP file",
		})
		return
	}

	// Process each file in the ZIP archive
	for _, zipFile := range zipReader.File {
		log.Printf("Processing file: %s\n", zipFile.Name)

		// Check if the file is a CSV
		if !isCSV(zipFile.Name) {
			log.Printf("Skipping non-CSV file: %s\n", zipFile.Name)
			continue
		}

		// Open the CSV file
		fileInZip, err := zipFile.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v\n", zipFile.Name, err)
			continue
		}
		defer fileInZip.Close()

		// Read and print the contents of the CSV file
		reader := csv.NewReader(fileInZip)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV file %s: %v\n", zipFile.Name, err)
				break
			}
			log.Printf("Record from %s: %v\n", zipFile.Name, record)
		}
	}

	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

func isCSV(filename string) bool {
	// Check if the file has a .csv extension
	return len(filename) > 4 && filename[len(filename)-4:] == ".csv"
}

/***
	mapping
***/

func (h *handler) MapCreateExamRequestToExamEntity(req *CreateExamRequest) *exam.Exam {
	return &exam.Exam{
		Name:                   req.Name,
		IsOpen:                 req.IsOpen,
		AllowedDurationMinutes: req.AllowedDurationMinutes,
	}
}

func (h *handler) MapExamEntityToExamData(svcRes *exam.Exam) *ExamData {
	return &ExamData{
		Serial:                 svcRes.Serial,
		Name:                   svcRes.Name,
		IsOpen:                 svcRes.IsOpen,
		AllowedDurationMinutes: svcRes.AllowedDurationMinutes,
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
		Serial:                 req.Serial,
		Name:                   req.Name,
		IsOpen:                 req.IsOpen,
		AllowedDurationMinutes: req.AllowedDurationMinutes,
	}
}
