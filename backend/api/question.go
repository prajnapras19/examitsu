package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prajnapras19/project-form-exam-sman2/backend/client/storage"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"github.com/prajnapras19/project-form-exam-sman2/backend/mcqoption"
	"github.com/prajnapras19/project-form-exam-sman2/backend/question"
	"github.com/prajnapras19/project-form-exam-sman2/backend/submission"
	"gorm.io/gorm"
)

/***
	entity
***/

type CreateQuestionRequest struct {
	ExamSerial string `json:"exam_serial" binding:"required"`
	ExamID     uint   `json:"-"`
	Data       string `json:"data"`
}

type QuestionDataIDOnly struct {
	ID uint `json:"id"`
}

type QuestionData struct {
	ID   uint   `json:"id"`
	Data string `json:"data"`
}

type UpdateQuestionRequest struct {
	ID   uint   `json:"-"`
	Data string `json:"data"`
}

type ExamSessionQuestionData struct {
	Question *QuestionData                `json:"question"`
	Options  []*McqOptionWithoutPointData `json:"options"`
	AnswerID uint                         `json:"answer"`
}

type SubmitAnswerRequest struct {
	McqOptionID uint `json:"mcq_option_id" binding:"required"`
}

type GetUploadQuestionBlobURLRequest struct {
	FileType string `json:"file_type" binding:"required"`
}

type GetUploadQuestionBlobURLResponse struct {
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
}

type GetExamSessionDetail struct {
	QuestionsIDList []*QuestionDataIDOnly `json:"questions_id_list"`
	StartTime       time.Time             `json:"start_time"`
	Duration        uint                  `json:"duration"`
}

/***
	handler
***/

func (h *handler) CreateQuestion(c *gin.Context) {
	var req CreateQuestionRequest

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
	svcReq := h.MapCreateQuestionRequestToQuestionEntity(&req)

	svcRes, err := h.questionService.CreateQuestion(svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	for _, mcqOption := range h.cfg.InitialMcqOptions {
		h.mcqOptionService.CreateMcqOption(&mcqoption.McqOption{
			QuestionID:  svcRes.ID,
			Description: mcqOption,
		})
	}

	res := h.MapQuestionEntityToQuestionData(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetUploadQuestionBlobURL(c *gin.Context) {
	var req GetUploadQuestionBlobURLRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	fileName, err := lib.GenerateRandomString(constants.DefaultRandomQuestionBlobFilenameLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrUnknownError.Error(),
		})
		return
	}

	svcRes, err := h.storageService.GetUploadURL(&storage.GetUploadURLRequest{
		FileName: fileName,
		FileType: req.FileType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapGetUploadURLResponseEntityToGetUploadQuestionBlobURLResponse(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetQuestions(c *gin.Context) {
	var filter question.GetQuestionsFilter

	if err := c.ShouldBind(&filter); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	if filter.IDEqualsTo == nil && filter.ExamSerialEqualsTo == nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	if filter.ExamSerialEqualsTo != nil {
		exam, err := h.examService.GetExamBySerial(filter.ExamSerialEqualsTo.Value)
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
		filter.ExamIDEqualsTo = &lib.QueryFiltersEqualToUint{
			Value: exam.ID,
		}
	}

	pagination, err := lib.GetQueryPaginationFromContext(c)
	if err != nil {
		log.Printf("[handler][question][GetQuestionsIDOnly] get query pagination error: %s", err.Error())
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}
	pagination.Sort = "order_number ASC"

	svcRes, err := h.questionService.GetQuestionsIDOnly(pagination, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := h.MapQuestionEntityListToQuestionDataIDOnlyList(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetQuestionByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)

	svcRes, err := h.questionService.GetQuestionByID(uint(id))
	if err != nil {
		if errors.Is(err, lib.ErrQuestionNotFound) {
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

	res := h.MapQuestionEntityToQuestionData(svcRes)
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetQuestionsIDByExamSerial(c *gin.Context) {
	jwtClaims, err := lib.GetExamTokenJWTClaimsFromContext(c)
	if err != nil {
		log.Printf("[handler][question][GetQuestionsIDByExamSerial] error when get jwt: %s", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrUnknownError.Error(),
		})
		return
	}

	participant, err := h.participantService.GetParticipantByID(jwtClaims.ParticipantID)
	if err != nil {
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

	if participant.EndedAt != nil || participant.StartedAt.Add(time.Duration(participant.AllowedDurationMinutes)*time.Minute).Before(time.Now()) {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamAlreadySubmitted.Error(),
		})
		return
	}

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

	if participant.ExamID != exam.ID {
		c.JSON(http.StatusUnauthorized, lib.BaseResponse{
			Message: lib.ErrUnauthorizedRequest.Error(),
		})
		return
	}

	if !exam.IsOpen {
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

	svcRes, err := h.questionService.GetQuestionsIDByExamID(exam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	res := GetExamSessionDetail{
		QuestionsIDList: h.MapQuestionEntityListToQuestionDataIDOnlyList(svcRes),
		StartTime:       *participant.StartedAt,
		Duration:        participant.AllowedDurationMinutes,
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) GetQuestionWithOptions(c *gin.Context) {
	jwtClaims, err := lib.GetExamTokenJWTClaimsFromContext(c)
	if err != nil {
		log.Printf("[handler][question][GetQuestionsIDByExamSerial] error when get jwt: %s", err.Error())
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

	if participant.EndedAt != nil || participant.StartedAt.Add(time.Duration(participant.AllowedDurationMinutes)*time.Minute).Before(time.Now()) {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamAlreadySubmitted.Error(),
		})
		return
	}

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

	if participant.ExamID != exam.ID {
		c.JSON(http.StatusUnauthorized, lib.BaseResponse{
			Message: lib.ErrUnauthorizedRequest.Error(),
		})
		return
	}

	if !exam.IsOpen {
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

	questionID, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)
	question, err := h.questionService.GetQuestionByID(uint(questionID))
	if err != nil {
		if errors.Is(err, lib.ErrQuestionNotFound) {
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
	if question.ExamID != exam.ID {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrQuestionNotFound.Error(),
		})
		return
	}

	mcqOptions, err := h.mcqOptionService.GetMcqOptionsByQuestionID(question.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}

	answer, err := h.submissionService.GetAnswer(participant.ID, question.ID)
	if err != nil {
		if !errors.Is(err, lib.ErrAnswerNotFound) {
			c.JSON(http.StatusInternalServerError, lib.BaseResponse{
				Message: err.Error(),
			})
			return
		}
	}

	answerID := uint(0)
	if answer != nil {
		answerID = answer.McqOptionID
	}
	res := ExamSessionQuestionData{
		Question: h.MapQuestionEntityToQuestionData(question),
		Options:  h.MapMcqOptionEntityListToMcqOptionWithoutPointDataList(mcqOptions),
		AnswerID: answerID,
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
		Data:    res,
	})
}

func (h *handler) SubmitAnswer(c *gin.Context) {
	jwtClaims, err := lib.GetExamTokenJWTClaimsFromContext(c)
	if err != nil {
		log.Printf("[handler][question][SubmitAnswer] error when get jwt: %s", err.Error())
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: lib.ErrUnknownError.Error(),
		})
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
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

	if participant.EndedAt != nil || participant.StartedAt.Add(time.Duration(participant.AllowedDurationMinutes)*time.Minute).Before(time.Now()) {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrExamAlreadySubmitted.Error(),
		})
		return
	}

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

	if participant.ExamID != exam.ID {
		c.JSON(http.StatusUnauthorized, lib.BaseResponse{
			Message: lib.ErrUnauthorizedRequest.Error(),
		})
		return
	}

	if !exam.IsOpen {
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

	questionID, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)
	question, err := h.questionService.GetQuestionByID(uint(questionID))
	if err != nil {
		if errors.Is(err, lib.ErrQuestionNotFound) {
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
	if question.ExamID != exam.ID {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrQuestionNotFound.Error(),
		})
		return
	}

	mcqOption, err := h.mcqOptionService.GetMcqOptionByID(req.McqOptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
		return
	}
	if mcqOption.QuestionID != question.ID {
		c.JSON(http.StatusNotFound, lib.BaseResponse{
			Message: lib.ErrMcqOptionNotFound.Error(),
		})
		return
	}

	err = h.submissionService.Answer(&submission.ExamSessionSubmissionCacheObject{
		ParticipantID: participant.ID,
		QuestionID:    question.ID,
		McqOptionID:   mcqOption.ID,
		Timestamp:     time.Now().Truncate(time.Second),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, lib.BaseResponse{
			Message: err.Error(),
		})
	}
	c.JSON(http.StatusOK, lib.BaseResponse{
		Message: constants.Success,
	})
}

func (h *handler) UpdateQuestion(c *gin.Context) {
	var req UpdateQuestionRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, lib.BaseResponse{
			Message: lib.ErrFailedToParseRequest.Error(),
		})
		return
	}

	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)
	req.ID = uint(id)
	svcReq := h.MapUpdateQuestionRequestToQuestionEntity(&req)

	err := h.questionService.UpdateQuestion(svcReq)
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

func (h *handler) DeleteQuestionBySerial(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param(constants.ID), 10, 64)

	err := h.questionService.DeleteQuestionByID(uint(id))
	if err != nil {
		if errors.Is(err, lib.ErrQuestionNotFound) {
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

func (h *handler) MapCreateQuestionRequestToQuestionEntity(req *CreateQuestionRequest) *question.Question {
	return &question.Question{
		ExamID: req.ExamID,
		Data:   req.Data,
	}
}

func (h *handler) MapQuestionEntityToQuestionData(svcRes *question.Question) *QuestionData {
	return &QuestionData{
		ID:   svcRes.ID,
		Data: svcRes.Data,
	}
}

func (h *handler) MapQuestionEntityListToQuestionDataList(svcRes []*question.Question) []*QuestionData {
	res := []*QuestionData{}
	for _, obj := range svcRes {
		res = append(res, h.MapQuestionEntityToQuestionData(obj))
	}
	return res
}

func (h *handler) MapUpdateQuestionRequestToQuestionEntity(req *UpdateQuestionRequest) *question.Question {
	return &question.Question{
		BaseModel: lib.BaseModel{
			Model: gorm.Model{
				ID: req.ID,
			},
		},
		Data: req.Data,
	}
}

func (h *handler) MapQuestionEntityToQuestionDataIDOnly(svcRes *question.Question) *QuestionDataIDOnly {
	return &QuestionDataIDOnly{
		ID: svcRes.ID,
	}
}

func (h *handler) MapQuestionEntityListToQuestionDataIDOnlyList(svcRes []*question.Question) []*QuestionDataIDOnly {
	res := []*QuestionDataIDOnly{}
	for _, obj := range svcRes {
		res = append(res, h.MapQuestionEntityToQuestionDataIDOnly(obj))
	}
	return res
}

func (h *handler) MapGetUploadURLResponseEntityToGetUploadQuestionBlobURLResponse(svcRes *storage.GetUploadURLResponse) *GetUploadQuestionBlobURLResponse {
	return &GetUploadQuestionBlobURLResponse{
		UploadURL: svcRes.UploadURL,
		PublicURL: svcRes.PublicURL,
	}
}
