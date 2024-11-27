package api

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
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
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToDecodeContent.Error(),
		})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", example.ExamZipExampleFilename))
	c.Data(http.StatusOK, "application/zip", fileContent)
}

func (h *handler) UploadExam(c *gin.Context) {
	file, err := c.FormFile(constants.File)
	if err != nil {
		log.Println("[exam][UploadExam] failed to get form file:", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		log.Println("[exam][UploadExam] failed to open uploaded file:", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	defer uploadedFile.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, uploadedFile)
	if err != nil {
		log.Println("[exam][UploadExam] failed to read uploaded file:", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), file.Size)
	if err != nil {
		log.Println("[exam][UploadExam] failed to read uploaded file as zip:", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	fileMap := map[string]*zip.File{}
	for _, f := range zipReader.File {
		fileMap[f.Name] = f
	}

	examFile, ok := fileMap[constants.UjianCSV]
	if !ok {
		log.Println("[exam][UploadExam] ujian.csv not found in zip")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	questionsFile, ok := fileMap[constants.SoalCSV]
	if !ok {
		log.Println("[exam][UploadExam] soal.csv not found in zip")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	mcqOptionsFile, ok := fileMap[constants.KunciCSV]
	if !ok {
		log.Println("[exam][UploadExam] kunci.csv not found in zip")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	participantsFile, ok := fileMap[constants.PesertaCSV]
	if !ok {
		log.Println("[exam][UploadExam] peserta.csv not found in zip")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	openedExamFile, err := examFile.Open()
	if err != nil {
		log.Println("[exam][UploadExam] failed to open ujian.csv:", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	defer openedExamFile.Close()
	examFileHeader, examFileContent, err := lib.ReadCSV(openedExamFile)
	if err != nil {
		log.Println("[exam][UploadExam] failed to read ujian.csv:", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	if len(examFileHeader) != 2 || examFileHeader[0] != constants.Nama || examFileHeader[1] != constants.Durasi || len(examFileContent) != 1 {
		log.Println("[exam][UploadExam] ujian.csv is not formatted as expected")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	openedQuestionsFile, err := questionsFile.Open()
	if err != nil {
		log.Println("[exam][UploadExam] failed to open soal.csv:", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	defer openedQuestionsFile.Close()
	questionsFileHeader, questionsFileContent, err := lib.ReadCSV(openedQuestionsFile)
	if err != nil {
		log.Println("[exam][UploadExam] failed to read soal.csv:", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	if len(questionsFileHeader) != 2 || questionsFileHeader[0] != constants.Nomor || questionsFileHeader[1] != constants.Gambar {
		log.Println("[exam][UploadExam] soal.csv is not formatted as expected")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	for _, questionContent := range questionsFileContent {
		if _, ok := fileMap[questionContent[constants.Gambar]]; !ok {
			log.Println("[exam][UploadExam] file declared in soal.csv not found:", questionContent[constants.Gambar])
			c.JSON(http.StatusBadRequest, lib.BaseResponse{
				Message: lib.ErrFailedToProcessUploadedFile.Error(),
			})
		}
	}

	openedMcqOptionsFile, err := mcqOptionsFile.Open()
	if err != nil {
		log.Println("[exam][UploadExam] failed to open kunci.csv:", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	defer openedMcqOptionsFile.Close()
	mcqOptionsFileHeader, mcqOptionsFileContent, err := lib.ReadCSV(openedMcqOptionsFile)
	if err != nil {
		log.Println("[exam][UploadExam] failed to read kunci.csv:", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	if len(mcqOptionsFileHeader) != 3 || mcqOptionsFileHeader[0] != constants.Soal || mcqOptionsFileHeader[1] != constants.Deskripsi || mcqOptionsFileHeader[2] != constants.Poin {
		log.Println("[exam][UploadExam] kunci.csv is not formatted as expected")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	openedParticipantsFile, err := participantsFile.Open()
	if err != nil {
		log.Println("[exam][UploadExam] failed to open kunci.csv:", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	defer openedParticipantsFile.Close()
	participantsFileHeader, participantsFileContent, err := lib.ReadCSV(openedParticipantsFile)
	if err != nil {
		log.Println("[exam][UploadExam] failed to read kunci.csv:", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}
	if len(participantsFileHeader) != 1 || participantsFileHeader[0] != constants.Kode {
		log.Println("[exam][UploadExam] peserta.csv is not formatted as expected")
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToProcessUploadedFile.Error(),
		})
		return
	}

	log.Println("examFileHeader", examFileHeader)
	log.Println("examFileContent", examFileContent)
	log.Println("questionsFileHeader", questionsFileHeader)
	log.Println("questionsFileContent", questionsFileContent)
	log.Println("mcqOptionsFileHeader", mcqOptionsFileHeader)
	log.Println("mcqOptionsFileContent", mcqOptionsFileContent)
	log.Println("participantsFileHeader", participantsFileHeader)
	log.Println("participantsFileContent", participantsFileContent)

	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
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
